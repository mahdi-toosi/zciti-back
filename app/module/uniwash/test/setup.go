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
	couponRepo "go-fiber-starter/app/module/coupon/repository"
	couponService "go-fiber-starter/app/module/coupon/service"
	productRepo "go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/uniwash"
	"go-fiber-starter/app/module/uniwash/controller"
	"go-fiber-starter/app/module/uniwash/repository"
	"go-fiber-starter/app/module/uniwash/service"
	userRepo "go-fiber-starter/app/module/user/repository"
	userService "go-fiber-starter/app/module/user/service"
	"go-fiber-starter/internal"
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
	App           *fiber.App
	DB            *gorm.DB
	Config        *config.Config
	Cleanup       func()
	UniWashRouter *uniwash.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/uniwash/test/ to project root
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

// migrateTestModels creates the necessary tables for uniwash testing
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	db.Exec("DROP TABLE IF EXISTS reservations CASCADE")
	db.Exec("DROP TABLE IF EXISTS products CASCADE")
	db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	db.Exec("DROP TABLE IF EXISTS coupons CASCADE")
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

	// Create businesses table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS businesses (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			owner_id BIGINT NOT NULL,
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

	// Create business_users junction table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS business_users (
			business_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			PRIMARY KEY (business_id, user_id)
		)
	`).Error; err != nil {
		return err
	}

	// Create posts table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255),
			excerpt VARCHAR(255),
			content TEXT NOT NULL,
			status VARCHAR(50) DEFAULT 'published',
			type VARCHAR(50) NOT NULL,
			parent_id BIGINT,
			slug VARCHAR(600),
			author_id BIGINT NOT NULL,
			business_id BIGINT,
			meta JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create products table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id BIGSERIAL PRIMARY KEY,
			post_id BIGINT,
			is_root BOOLEAN DEFAULT FALSE,
			type VARCHAR(50) NOT NULL,
			variant_type VARCHAR(50),
			price FLOAT NOT NULL,
			min_price FLOAT NOT NULL DEFAULT 0,
			max_price FLOAT NOT NULL DEFAULT 0,
			on_sale BOOLEAN DEFAULT FALSE,
			stock_status VARCHAR(40) NOT NULL DEFAULT 'inStock',
			total_sales FLOAT DEFAULT 0,
			meta JSONB,
			business_id BIGINT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create reservations table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS reservations (
			id BIGSERIAL PRIMARY KEY,
			status VARCHAR(50) DEFAULT 'reserved',
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			user_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL,
			business_id BIGINT NOT NULL,
			meta JSONB,
			user_usage_count BIGINT DEFAULT 0,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create coupons table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS coupons (
			id BIGSERIAL PRIMARY KEY,
			code VARCHAR(255) NOT NULL,
			title VARCHAR(255) NOT NULL,
			description VARCHAR(500),
			value FLOAT NOT NULL,
			type VARCHAR(50) NOT NULL,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			times_used INT DEFAULT 0,
			business_id BIGINT NOT NULL,
			meta JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			UNIQUE(code, business_id)
		)
	`).Error; err != nil {
		return err
	}

	// Create indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_reservations_business_id ON reservations(business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_reservations_product_id ON reservations(product_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_reservations_deleted_at ON reservations(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_business_id ON products(business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_post_id ON products(post_id)")

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

	// Create repositories
	uniwashRepo := repository.Repository(dbWrapper)
	productRepository := productRepo.Repository(dbWrapper)
	userRepository := userRepo.Repository(dbWrapper)
	couponRepository := couponRepo.Repository(dbWrapper)

	// Create user service
	userSvc := userService.Service(userRepository)

	// Create mock MessageWay service
	mockMW := internal.NewMessageWay(cfg, logger)

	// Create coupon service
	couponSvc := couponService.Service(couponRepository, userSvc, mockMW)

	// Create uniwash service
	uniwashSvc := service.Service(uniwashRepo, couponSvc, productRepository, mockMW)

	// Create uniwash controller
	uniwashController := controller.Controllers(uniwashSvc)

	// Create uniwash router manually
	uniwashRouter := &uniwash.Router{
		App:        app,
		Controller: uniwashController,
	}
	uniwashRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		dbWrapper.Main.Exec("DELETE FROM reservations")
		dbWrapper.Main.Exec("DELETE FROM coupons")
		dbWrapper.Main.Exec("DELETE FROM products")
		dbWrapper.Main.Exec("DELETE FROM posts")
		dbWrapper.Main.Exec("DELETE FROM business_users")
		dbWrapper.Main.Exec("DELETE FROM businesses")
		dbWrapper.Main.Exec("DELETE FROM users")
		dbWrapper.ShutdownDatabase()
	}

	return &TestApp{
		App:           app,
		DB:            dbWrapper.Main,
		Config:        cfg,
		Cleanup:       cleanup,
		UniWashRouter: uniwashRouter,
	}
}

// CreateTestUser creates a test user in the database with business permissions
func (ta *TestApp) CreateTestUser(t *testing.T, mobile uint64, password string, firstName, lastName string, businessID uint64, roles []schema.UserRole) *schema.User {
	t.Helper()

	isSuspended := false
	permissions := schema.UserPermissions{}
	if businessID > 0 && len(roles) > 0 {
		permissions[businessID] = roles
	}

	user := &schema.User{
		Mobile:      mobile,
		FirstName:   firstName,
		LastName:    lastName,
		Password:    helpers.Hash([]byte(password)),
		Permissions: permissions,
		IsSuspended: &isSuspended,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

// CreateTestBusiness creates a test business in the database
func (ta *TestApp) CreateTestBusiness(t *testing.T, title string, businessType schema.BusinessType, ownerID uint64) *schema.Business {
	t.Helper()

	business := &schema.Business{
		Title:   title,
		Type:    businessType,
		OwnerID: ownerID,
		Account: schema.BusinessAccountDefault,
		Meta:    schema.BusinessMeta{},
	}

	if err := ta.DB.Create(business).Error; err != nil {
		t.Fatalf("failed to create test business: %v", err)
	}

	return business
}

// CreateTestPost creates a test post in the database
func (ta *TestApp) CreateTestPost(t *testing.T, title string, businessID uint64, authorID uint64) *schema.Post {
	t.Helper()

	post := &schema.Post{
		Title:      title,
		Content:    "Test product content",
		Status:     schema.PostStatusPublished,
		Type:       schema.PostTypeProduct,
		AuthorID:   authorID,
		BusinessID: businessID,
	}

	if err := ta.DB.Create(post).Error; err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return post
}

// CreateTestProduct creates a test product (washing machine) in the database
func (ta *TestApp) CreateTestProduct(t *testing.T, postID uint64, businessID uint64, price float64, mobileNumber string, machineStatus schema.UniWashMachineStatus) *schema.Product {
	t.Helper()

	variantType := schema.ProductVariantTypeWashingMachine
	product := &schema.Product{
		PostID:      postID,
		Type:        schema.ProductTypeVariant,
		VariantType: &variantType,
		Price:       price,
		MinPrice:    price,
		MaxPrice:    price,
		StockStatus: schema.ProductStockStatusInStock,
		BusinessID:  businessID,
		Meta: schema.ProductMeta{
			SKU:                  "WM-001",
			UniWashMobileNumber:  mobileNumber,
			UniWashMachineStatus: machineStatus,
		},
	}

	if err := ta.DB.Create(product).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

// CreateTestReservation creates a test reservation in the database
func (ta *TestApp) CreateTestReservation(t *testing.T, userID, productID, businessID uint64, startTime, endTime time.Time, status schema.ReservationStatus) *schema.Reservation {
	t.Helper()

	reservation := &schema.Reservation{
		UserID:     userID,
		ProductID:  productID,
		BusinessID: businessID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     status,
		Meta:       schema.ReservationMeta{},
	}

	if err := ta.DB.Create(reservation).Error; err != nil {
		t.Fatalf("failed to create test reservation: %v", err)
	}

	return reservation
}

// GenerateTestToken generates a JWT token for test user
func (ta *TestApp) GenerateTestToken(t *testing.T, user *schema.User) string {
	t.Helper()

	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

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
func (ta *TestApp) MakeRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
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
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

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

// CleanupReservations removes all reservations from the test database
func (ta *TestApp) CleanupReservations(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM reservations")
}

// CleanupAll removes all test data from the database
func (ta *TestApp) CleanupAll(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM reservations")
	ta.DB.Exec("DELETE FROM coupons")
	ta.DB.Exec("DELETE FROM products")
	ta.DB.Exec("DELETE FROM posts")
	ta.DB.Exec("DELETE FROM business_users")
	ta.DB.Exec("DELETE FROM businesses")
	ta.DB.Exec("DELETE FROM users")
}

