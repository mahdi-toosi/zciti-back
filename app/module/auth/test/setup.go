package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/auth"
	"go-fiber-starter/app/module/auth/controller"
	"go-fiber-starter/app/module/auth/service"
	userRepo "go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/internal"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/helpers"

	MessageWay "github.com/MessageWay/MessageWayGolang"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TestApp holds the test application components
type TestApp struct {
	App        *fiber.App
	DB         *gorm.DB
	Config     *config.Config
	Cleanup    func()
	UserRepo   userRepo.IRepository
	AuthRouter *auth.Router
}

// MockMessageWayService mocks the MessageWay service for testing
type MockMessageWayService struct {
	SendCalled   bool
	VerifyCalled bool
	SentOTP      string
	VerifyOTP    string
}

func (m *MockMessageWayService) Send(req MessageWay.Message) (*MessageWay.SendResponse, error) {
	m.SendCalled = true
	return &MessageWay.SendResponse{
		Status:      "success",
		ReferenceID: "test-reference-id",
	}, nil
}

func (m *MockMessageWayService) Verify(req MessageWay.OTPVerifyRequest) (*MessageWay.OTPVerifyResponse, error) {
	m.VerifyCalled = true
	m.VerifyOTP = req.OTP
	return &MessageWay.OTPVerifyResponse{
		Status: "success",
	}, nil
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/auth/test/ to project root
	return filepath.Join(filepath.Dir(b), "..", "..", "..", "..")
}

// createTestErrorHandler creates an error handler that properly handles validation errors
func createTestErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		// Handle validation errors
		if _, ok := err.(validator.ValidationErrors); ok {
			code = fiber.StatusUnprocessableEntity
		} else if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return c.Status(code).JSON(fiber.Map{
			"Code":     code,
			"Messages": []string{err.Error()},
		})
	}
}

// migrateTestModels creates the users table matching schema.User structure
// Note: Cannot use db.AutoMigrate(&schema.User{}) because GORM follows relationships
// and tries to create/check Business, Taxonomy tables due to FK definitions in schema
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	db.Exec("DROP TABLE IF EXISTS business_users CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Create users table - structure matches schema.User fields
	// Using raw SQL because GORM AutoMigrate follows relationship definitions
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			mobile BIGINT NOT NULL UNIQUE,
			mobile_confirmed BOOLEAN DEFAULT FALSE,
			show_mobile BOOLEAN,
			is_suspended BOOLEAN DEFAULT FALSE,
			suspense_reason VARCHAR(500),
			permissions JSONB NOT NULL DEFAULT '{}',
			password VARCHAR(255) NOT NULL,
			city_id BIGINT,
			workspace_id BIGINT,
			dormitory_id BIGINT,
			reservation_count BIGINT DEFAULT 0,
			meta JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create indexes matching schema.User GORM tags
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)")

	return nil
}

// SetupTestApp initializes the test application with a test database
func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()

	// Load test config from project root
	configPath := filepath.Join(getProjectRoot(), "config", "zciti-test.toml")
	cfg, err := config.ParseConfig(configPath, true)
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}

	// Create logger
	logger := zerolog.Nop()

	// Create database wrapper using the bootstrap database module
	dbWrapper := database.NewDatabase(cfg, logger)
	dbWrapper.ConnectDatabase()

	if dbWrapper.Main == nil {
		t.Fatalf("failed to connect to test database")
	}

	// Migrate test models (creates users table)
	if err := migrateTestModels(dbWrapper.Main); err != nil {
		t.Fatalf("failed to migrate test models: %v", err)
	}

	// Create Fiber app with proper error handling
	app := fiber.New(fiber.Config{
		ErrorHandler: createTestErrorHandler(),
	})

	// Create user repository
	repo := userRepo.Repository(dbWrapper)

	// Create mock MessageWay service with initialized App
	mockMW := internal.NewMessageWay(cfg, logger)

	// Create auth service
	authService := service.Service(repo, mockMW)

	// Create auth controller
	authController := controller.Controllers(authService, cfg)

	// Create auth router
	authRouter := auth.NewRouter(app, authController)
	authRouter.RegisterRoutes()

	// Cleanup function
	cleanup := func() {
		// Clean up test data - delete in reverse order of dependencies
		dbWrapper.Main.Exec("DELETE FROM users")
		dbWrapper.ShutdownDatabase()
	}

	return &TestApp{
		App:        app,
		DB:         dbWrapper.Main,
		Config:     cfg,
		Cleanup:    cleanup,
		UserRepo:   repo,
		AuthRouter: authRouter,
	}
}

// CreateTestUser creates a test user in the database
func (ta *TestApp) CreateTestUser(t *testing.T, mobile uint64, password string, firstName, lastName string) *schema.User {
	t.Helper()

	isSuspended := false
	user := &schema.User{
		Mobile:      mobile,
		FirstName:   firstName,
		LastName:    lastName,
		Password:    helpers.Hash([]byte(password)),
		Permissions: schema.UserPermissions{},
		IsSuspended: &isSuspended,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

// MakeRequest makes an HTTP request to the test server
func (ta *TestApp) MakeRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ta.App.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	return resp
}

// ParseResponse parses the response body into a map
func ParseResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Try parsing as string
		var strResult string
		if err := json.Unmarshal(body, &strResult); err != nil {
			t.Fatalf("failed to parse response: %v, body: %s", err, string(body))
		}
		return map[string]interface{}{"result": strResult}
	}

	return result
}

// ParseResponseTo parses the response body into the provided struct
func ParseResponseTo(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to parse response: %v, body: %s", err, string(body))
	}
}

// CleanupUsers removes all users from the test database
func (ta *TestApp) CleanupUsers(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM users")
}
