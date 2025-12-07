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
	"go-fiber-starter/app/module/order"
	"go-fiber-starter/app/module/order/controller"
	"go-fiber-starter/app/module/order/repository"
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
	App         *fiber.App
	DB          *gorm.DB
	Config      *config.Config
	Cleanup     func()
	OrderRepo   repository.IRepository
	OrderRouter *order.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/order/test/ to project root
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

// migrateTestModels creates the necessary tables for order testing
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	db.Exec("DROP TABLE IF EXISTS order_items CASCADE")
	db.Exec("DROP TABLE IF EXISTS orders CASCADE")
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS wallets CASCADE")
	db.Exec("DROP TABLE IF EXISTS reservations CASCADE")
	db.Exec("DROP TABLE IF EXISTS coupons CASCADE")
	db.Exec("DROP TABLE IF EXISTS products CASCADE")
	db.Exec("DROP TABLE IF EXISTS posts CASCADE")
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
			title VARCHAR(255) NOT NULL,
			excerpt VARCHAR(255),
			content TEXT NOT NULL,
			status VARCHAR(50) DEFAULT 'published',
			type VARCHAR(50) NOT NULL,
			parent_id BIGINT,
			slug VARCHAR(600),
			author_id BIGINT NOT NULL,
			business_id BIGINT NOT NULL,
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
			post_id BIGINT NOT NULL,
			is_root BOOLEAN DEFAULT FALSE,
			type VARCHAR(50) NOT NULL,
			variant_type VARCHAR(50),
			price FLOAT NOT NULL,
			min_price FLOAT DEFAULT 0,
			max_price FLOAT DEFAULT 0,
			on_sale BOOLEAN DEFAULT FALSE,
			stock_status VARCHAR(40) NOT NULL DEFAULT 'inStock',
			total_sales FLOAT DEFAULT 0,
			meta JSONB,
			business_id BIGINT NOT NULL,
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

	// Create wallets table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS wallets (
			id BIGSERIAL PRIMARY KEY,
			amount FLOAT DEFAULT 0,
			user_id BIGINT,
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
			status VARCHAR(50) NOT NULL DEFAULT 'available',
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			product_id BIGINT NOT NULL,
			user_id BIGINT,
			business_id BIGINT NOT NULL,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create orders table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			total_amt FLOAT NOT NULL DEFAULT 0,
			payment_method VARCHAR(20) NOT NULL DEFAULT 'online',
			meta JSONB,
			user_id BIGINT NOT NULL,
			business_id BIGINT NOT NULL,
			coupon_id BIGINT,
			parent_id BIGINT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create order_items table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			type VARCHAR(50) NOT NULL DEFAULT 'lineItem',
			quantity INT NOT NULL DEFAULT 1,
			price FLOAT NOT NULL,
			subtotal FLOAT NOT NULL,
			tax_amt FLOAT DEFAULT 0,
			reservation_id BIGINT,
			post_id BIGINT NOT NULL,
			order_id BIGINT NOT NULL,
			meta JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create transactions table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			amount FLOAT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			description VARCHAR(255) NOT NULL,
			order_payment_method VARCHAR(20) NOT NULL DEFAULT 'online',
			gateway_transaction_id VARCHAR(255),
			wallet_id BIGINT,
			order_id BIGINT,
			user_id BIGINT NOT NULL,
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
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_business_id ON orders(business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)")
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

	// Create order repository
	orderRepo := repository.Repository(dbWrapper)

	// Create mock service for testing
	// Note: For E2E tests, we test the controller directly without the full service dependencies
	// since Store operation requires many external services (payment gateway, etc.)
	mockService := &MockOrderService{
		repo: orderRepo,
		db:   dbWrapper.Main,
	}

	// Create order controller
	orderController := controller.Controllers(mockService)

	// Create order router manually
	orderRouter := &order.Router{
		App:        app,
		Controller: orderController,
	}
	orderRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		dbWrapper.Main.Exec("DELETE FROM order_items")
		dbWrapper.Main.Exec("DELETE FROM orders")
		dbWrapper.Main.Exec("DELETE FROM transactions")
		dbWrapper.Main.Exec("DELETE FROM wallets")
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
		App:         app,
		DB:          dbWrapper.Main,
		Config:      cfg,
		Cleanup:     cleanup,
		OrderRepo:   orderRepo,
		OrderRouter: orderRouter,
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
func (ta *TestApp) CreateTestPost(t *testing.T, title string, businessID, authorID uint64) *schema.Post {
	t.Helper()

	post := &schema.Post{
		Title:      title,
		Excerpt:    "Test excerpt",
		Content:    "Test content",
		Status:     schema.PostStatusPublished,
		Type:       schema.PostTypeProduct,
		AuthorID:   authorID,
		BusinessID: businessID,
		Meta: schema.PostMeta{
			CommentsStatus: schema.PostCommentStatusOpen,
		},
	}

	if err := ta.DB.Create(post).Error; err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return post
}

// CreateTestProduct creates a test product in the database
func (ta *TestApp) CreateTestProduct(t *testing.T, postID, businessID uint64, price float64, productType schema.ProductType, variantType *schema.ProductVariantType) *schema.Product {
	t.Helper()

	product := &schema.Product{
		PostID:      postID,
		BusinessID:  businessID,
		IsRoot:      true,
		Type:        productType,
		VariantType: variantType,
		Price:       price,
		MinPrice:    price,
		MaxPrice:    price,
		StockStatus: schema.ProductStockStatusInStock,
		Meta:        schema.ProductMeta{},
	}

	if err := ta.DB.Create(product).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

// CreateTestOrder creates a test order in the database
func (ta *TestApp) CreateTestOrder(t *testing.T, userID, businessID uint64, totalAmt float64, status schema.OrderStatus, paymentMethod schema.OrderPaymentMethod) *schema.Order {
	t.Helper()

	order := &schema.Order{
		UserID:        userID,
		BusinessID:    businessID,
		TotalAmt:      totalAmt,
		Status:        status,
		PaymentMethod: paymentMethod,
		Meta:          schema.OrderMeta{},
	}

	if err := ta.DB.Create(order).Error; err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}

	return order
}

// CreateTestOrderItem creates a test order item in the database
func (ta *TestApp) CreateTestOrderItem(t *testing.T, orderID, postID uint64, quantity int, price float64) *schema.OrderItem {
	t.Helper()

	orderItem := &schema.OrderItem{
		OrderID:  orderID,
		PostID:   postID,
		Quantity: quantity,
		Price:    price,
		Subtotal: float64(quantity) * price,
		Type:     schema.OrderItemTypeLineItem,
		Meta:     schema.OrderItemMeta{},
	}

	if err := ta.DB.Create(orderItem).Error; err != nil {
		t.Fatalf("failed to create test order item: %v", err)
	}

	return orderItem
}

// CreateTestCoupon creates a test coupon in the database
func (ta *TestApp) CreateTestCoupon(t *testing.T, code, title string, value float64, couponType schema.CouponType, businessID uint64, startTime, endTime time.Time) *schema.Coupon {
	t.Helper()

	coupon := &schema.Coupon{
		Code:       code,
		Title:      title,
		Value:      value,
		Type:       couponType,
		BusinessID: businessID,
		StartTime:  startTime,
		EndTime:    endTime,
		TimesUsed:  0,
		Meta:       schema.CouponMeta{MaxUsage: 100},
	}

	if err := ta.DB.Create(coupon).Error; err != nil {
		t.Fatalf("failed to create test coupon: %v", err)
	}

	return coupon
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

// CleanupOrders removes all orders from the test database
func (ta *TestApp) CleanupOrders(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM order_items")
	ta.DB.Exec("DELETE FROM orders")
}

// CleanupAll removes all test data from the database
func (ta *TestApp) CleanupAll(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM order_items")
	ta.DB.Exec("DELETE FROM orders")
	ta.DB.Exec("DELETE FROM transactions")
	ta.DB.Exec("DELETE FROM wallets")
	ta.DB.Exec("DELETE FROM reservations")
	ta.DB.Exec("DELETE FROM coupons")
	ta.DB.Exec("DELETE FROM products")
	ta.DB.Exec("DELETE FROM posts")
	ta.DB.Exec("DELETE FROM business_users")
	ta.DB.Exec("DELETE FROM businesses")
	ta.DB.Exec("DELETE FROM users")
}

