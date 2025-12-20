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
	"time"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/middleware"
	brequest "go-fiber-starter/app/module/business/request"
	bresponse "go-fiber-starter/app/module/business/response"
	orequest "go-fiber-starter/app/module/order/request"
	oresponse "go-fiber-starter/app/module/order/response"
	"go-fiber-starter/app/module/user"
	"go-fiber-starter/app/module/user/controller"
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/service"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/helpers"
	"go-fiber-starter/utils/paginator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TestApp holds the test application components
type TestApp struct {
	App        *fiber.App
	DB         *gorm.DB
	Config     *config.Config
	Cleanup    func()
	UserRepo   repository.IRepository
	UserRouter *user.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/user/test/ to project root
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

// migrateTestModels creates the users table and related tables for testing
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	db.Exec("DROP TABLE IF EXISTS business_users CASCADE")
	db.Exec("DROP TABLE IF EXISTS businesses CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Create users table
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

	// Create businesses table for permission tests
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS businesses (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255),
			owner_id BIGINT,
			account VARCHAR(50) DEFAULT 'default',
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create business_users junction table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS business_users (
			user_id BIGINT NOT NULL,
			business_id BIGINT NOT NULL,
			PRIMARY KEY (user_id, business_id)
		)
	`).Error; err != nil {
		return err
	}

	// Create indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)")

	return nil
}

// MockBusinessService mocks the business service for testing
type MockBusinessService struct {
	OwnerID uint64
}

func (m *MockBusinessService) Index(req brequest.Businesses) ([]*bresponse.Business, paginator.Pagination, error) {
	return nil, paginator.Pagination{}, nil
}

func (m *MockBusinessService) Show(id uint64, role schema.UserRole) (*bresponse.Business, error) {
	ownerID := m.OwnerID
	if ownerID == 0 {
		ownerID = 1 // Default owner ID for tests
	}
	return &bresponse.Business{
		ID:      id,
		Title:   "Test Business",
		OwnerID: ownerID,
	}, nil
}

func (m *MockBusinessService) Store(req brequest.Business) error {
	return nil
}

func (m *MockBusinessService) Update(id uint64, req brequest.Business) error {
	return nil
}

func (m *MockBusinessService) Destroy(id uint64) error {
	return nil
}

func (m *MockBusinessService) RoleMenuItems(businessID uint64, user schema.User) ([]bresponse.MenuItem, error) {
	return nil, nil
}

// MockOrderService mocks the order service for testing
type MockOrderService struct{}

func (m *MockOrderService) Index(req orequest.Orders) ([]*oresponse.Order, uint64, paginator.Pagination, error) {
	return nil, 0, paginator.Pagination{}, nil
}

func (m *MockOrderService) Show(userID uint64, id uint64) (*oresponse.Order, error) {
	return nil, nil
}

func (m *MockOrderService) Store(req orequest.Order) (uint64, string, error) {
	return 0, "", nil
}

func (m *MockOrderService) Status(userID, orderID uint64, refNum string) (string, error) {
	return "OK", nil
}

func (m *MockOrderService) Update(id uint64, req orequest.Order) error {
	return nil
}

func (m *MockOrderService) Destroy(id uint64) error {
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
	repo := repository.Repository(dbWrapper)

	// Create user service
	userService := service.Service(repo)

	// Create mock services
	mockBusinessService := &MockBusinessService{}
	mockOrderService := &MockOrderService{}

	// Create user controller with mock dependencies
	userController := &controller.Controller{
		RestController: controller.RestController(
			userService,
			mockBusinessService,
			mockOrderService,
			cfg,
		),
	}

	// Create user router
	userRouter := &user.Router{
		App:        app,
		Controller: userController,
	}
	userRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		dbWrapper.Main.Exec("DELETE FROM business_users")
		dbWrapper.Main.Exec("DELETE FROM businesses")
		dbWrapper.Main.Exec("DELETE FROM users")
		dbWrapper.ShutdownDatabase()
	}

	return &TestApp{
		App:        app,
		DB:         dbWrapper.Main,
		Config:     cfg,
		Cleanup:    cleanup,
		UserRepo:   repo,
		UserRouter: userRouter,
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

// CreateAdminUser creates a test user with admin permissions
func (ta *TestApp) CreateAdminUser(t *testing.T, mobile uint64, password string, firstName, lastName string) *schema.User {
	t.Helper()

	isSuspended := false
	user := &schema.User{
		Mobile:    mobile,
		FirstName: firstName,
		LastName:  lastName,
		Password:  helpers.Hash([]byte(password)),
		Permissions: schema.UserPermissions{
			1: []schema.UserRole{schema.URAdmin}, // Admin for business ID 1
		},
		IsSuspended: &isSuspended,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create admin test user: %v", err)
	}

	return user
}

// CreateBusinessOwnerUser creates a test user with business owner permissions
func (ta *TestApp) CreateBusinessOwnerUser(t *testing.T, mobile uint64, password string, firstName, lastName string, businessID uint64) *schema.User {
	t.Helper()

	isSuspended := false
	user := &schema.User{
		Mobile:    mobile,
		FirstName: firstName,
		LastName:  lastName,
		Password:  helpers.Hash([]byte(password)),
		Permissions: schema.UserPermissions{
			businessID: []schema.UserRole{schema.URBusinessOwner},
		},
		IsSuspended: &isSuspended,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create business owner test user: %v", err)
	}

	return user
}

// CreateTestBusiness creates a test business in the database
func (ta *TestApp) CreateTestBusiness(t *testing.T, ownerID uint64, title string) *schema.Business {
	t.Helper()

	business := &schema.Business{
		Title:   title,
		OwnerID: ownerID,
	}

	if err := ta.DB.Create(business).Error; err != nil {
		t.Fatalf("failed to create test business: %v", err)
	}

	return business
}

// GenerateToken generates a JWT token for the given user using the test config secret
func (ta *TestApp) GenerateToken(t *testing.T, user *schema.User) string {
	t.Helper()

	// Generate token directly using test config's secret to avoid config mismatch
	expiresAt := jwt.NewNumericDate(time.Now().Add(ta.Config.Middleware.Jwt.Expiration * time.Second))

	jwtCustomClaim := middleware.JWTCustomClaim{
		User: schema.User{
			ID:              user.ID,
			Meta:            user.Meta,
			Mobile:          user.Mobile,
			LastName:        user.LastName,
			FirstName:       user.FirstName,
			Permissions:     user.Permissions,
			MobileConfirmed: user.MobileConfirmed,
			IsSuspended:     user.IsSuspended,
		},
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expiresAt},
	}

	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaim)
	token, err := unSignedToken.SignedString([]byte(ta.Config.Middleware.Jwt.Secret))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return token
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

// MakeAuthenticatedRequest makes an authenticated HTTP request to the test server
func (ta *TestApp) MakeAuthenticatedRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
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
	req.Header.Set("Authorization", "Bearer "+token)

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
	ta.DB.Exec("DELETE FROM business_users")
	ta.DB.Exec("DELETE FROM users")
}

// CleanupAll removes all test data from the database
func (ta *TestApp) CleanupAll(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM business_users")
	ta.DB.Exec("DELETE FROM businesses")
	ta.DB.Exec("DELETE FROM users")
}

// WaitForDB waits for database operations to complete
func (ta *TestApp) WaitForDB() {
	time.Sleep(10 * time.Millisecond)
}
