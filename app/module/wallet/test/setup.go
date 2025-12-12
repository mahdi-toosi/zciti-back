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
	"go-fiber-starter/app/module/wallet"
	"go-fiber-starter/app/module/wallet/controller"
	"go-fiber-starter/app/module/wallet/repository"
	"go-fiber-starter/app/module/wallet/service"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TestApp holds the test application components
type TestApp struct {
	App          *fiber.App
	DB           *gorm.DB
	Config       *config.Config
	Cleanup      func()
	WalletRepo   repository.IRepository
	WalletRouter *wallet.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/wallet/test/ to project root
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

// migrateTestModels creates the required tables for wallet tests
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state (in reverse dependency order)
	db.Exec("DROP TABLE IF EXISTS wallets CASCADE")
	db.Exec("DROP TABLE IF EXISTS business_users CASCADE")
	db.Exec("DROP TABLE IF EXISTS businesses CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Create users table first (businesses depends on it)
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

	// Create businesses table (depends on users)
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS businesses (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			owner_id BIGINT NOT NULL REFERENCES users(id),
			account VARCHAR(100) DEFAULT 'default',
			meta JSONB,
			description VARCHAR(500),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create wallets table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id BIGSERIAL PRIMARY KEY,
			amount FLOAT DEFAULT 0,
			user_id BIGINT REFERENCES users(id),
			business_id BIGINT REFERENCES businesses(id),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_wallet ON wallets(user_id, business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_wallets_deleted_at ON wallets(deleted_at)")

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

	// Migrate test models
	if err := migrateTestModels(dbWrapper.Main); err != nil {
		t.Fatalf("failed to migrate test models: %v", err)
	}

	// Create Fiber app with proper error handling
	app := fiber.New(fiber.Config{
		ErrorHandler: createTestErrorHandler(),
	})

	// Create wallet repository
	walletRepo := repository.Repository(dbWrapper)

	// Create wallet service
	walletService := service.Service(walletRepo)

	// Create wallet controller
	walletController := controller.Controllers(walletService)

	// Create wallet router
	walletRouter := wallet.NewRouter(app, walletController)
	walletRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		// Clean up test data (in reverse dependency order)
		dbWrapper.Main.Exec("DELETE FROM wallets")
		dbWrapper.Main.Exec("DELETE FROM businesses")
		dbWrapper.Main.Exec("DELETE FROM users")
		dbWrapper.ShutdownDatabase()
	}

	return &TestApp{
		App:          app,
		DB:           dbWrapper.Main,
		Config:       cfg,
		Cleanup:      cleanup,
		WalletRepo:   walletRepo,
		WalletRouter: walletRouter,
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

// CreateTestBusiness creates a test business in the database
func (ta *TestApp) CreateTestBusiness(t *testing.T, title string, ownerID uint64) *schema.Business {
	t.Helper()

	business := &schema.Business{
		Title:   title,
		Type:    schema.BTypeGymManager,
		OwnerID: ownerID,
		Account: schema.BusinessAccountDefault,
	}

	if err := ta.DB.Create(business).Error; err != nil {
		t.Fatalf("failed to create test business: %v", err)
	}

	return business
}

// CreateTestUserWithBusinessOwner creates a user who owns a business
func (ta *TestApp) CreateTestUserWithBusinessOwner(t *testing.T, mobile uint64, password string, firstName, lastName string, businessID uint64) *schema.User {
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
		t.Fatalf("failed to create test user with business owner role: %v", err)
	}

	return user
}

// CreateTestUserWithBusinessObserver creates a user with business observer role
func (ta *TestApp) CreateTestUserWithBusinessObserver(t *testing.T, mobile uint64, password string, firstName, lastName string, businessID uint64) *schema.User {
	t.Helper()

	isSuspended := false
	user := &schema.User{
		Mobile:    mobile,
		FirstName: firstName,
		LastName:  lastName,
		Password:  helpers.Hash([]byte(password)),
		Permissions: schema.UserPermissions{
			businessID: []schema.UserRole{schema.URBusinessObserver},
		},
		IsSuspended: &isSuspended,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user with business observer role: %v", err)
	}

	return user
}

// CreateTestWallet creates a test wallet in the database
func (ta *TestApp) CreateTestWallet(t *testing.T, userID *uint64, businessID *uint64, amount float64) *schema.Wallet {
	t.Helper()

	wallet := &schema.Wallet{
		UserID:     userID,
		BusinessID: businessID,
		Amount:     amount,
	}

	if err := ta.DB.Create(wallet).Error; err != nil {
		t.Fatalf("failed to create test wallet: %v", err)
	}

	return wallet
}

// GenerateTestToken generates a JWT token for a test user
func (ta *TestApp) GenerateTestToken(t *testing.T, user *schema.User) string {
	t.Helper()

	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour))

	jwtCustomClaim := middleware.JWTCustomClaim{
		User: schema.User{
			ID:              user.ID,
			Meta:            user.Meta,
			Mobile:          user.Mobile,
			LastName:        user.LastName,
			FirstName:       user.FirstName,
			Permissions:     user.Permissions,
			MobileConfirmed: user.MobileConfirmed,
		},
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expiresAt},
	}

	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaim)
	token, err := unSignedToken.SignedString([]byte(ta.Config.Middleware.Jwt.Secret))
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
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

// CleanupData removes all test data from the database
func (ta *TestApp) CleanupData(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM wallets")
	ta.DB.Exec("DELETE FROM businesses")
	ta.DB.Exec("DELETE FROM users")
}

