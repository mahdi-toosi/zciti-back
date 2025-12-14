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
	postController "go-fiber-starter/app/module/post/controller"
	postRepo "go-fiber-starter/app/module/post/repository"
	postService "go-fiber-starter/app/module/post/service"
	"go-fiber-starter/app/module/product"
	"go-fiber-starter/app/module/product/controller"
	"go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/product/service"
	userRepo "go-fiber-starter/app/module/user/repository"
	userService "go-fiber-starter/app/module/user/service"
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
	ProductRepo   repository.IRepository
	ProductRouter *product.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/product/test/ to project root
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

// migrateTestModels creates the necessary tables for product testing
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	db.Exec("DROP TABLE IF EXISTS products_taxonomies CASCADE")
	db.Exec("DROP TABLE IF EXISTS posts_taxonomies CASCADE")
	db.Exec("DROP TABLE IF EXISTS reservations CASCADE")
	db.Exec("DROP TABLE IF EXISTS products CASCADE")
	db.Exec("DROP TABLE IF EXISTS posts CASCADE")
	db.Exec("DROP TABLE IF EXISTS taxonomies CASCADE")
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

	// Create taxonomies table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS taxonomies (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			type VARCHAR(100) NOT NULL,
			domain VARCHAR(100) NOT NULL,
			slug VARCHAR(200) NOT NULL,
			business_id BIGINT NOT NULL,
			parent_id BIGINT,
			description VARCHAR(500),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
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
			business_id BIGINT NOT NULL,
			meta JSONB,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error; err != nil {
		return err
	}

	// Create posts_taxonomies junction table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts_taxonomies (
			post_id BIGINT NOT NULL,
			taxonomy_id BIGINT NOT NULL,
			PRIMARY KEY (post_id, taxonomy_id)
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
			min_price FLOAT NOT NULL,
			max_price FLOAT NOT NULL,
			on_sale BOOLEAN DEFAULT FALSE,
			stock_status VARCHAR(40) NOT NULL,
			total_sales FLOAT DEFAULT 0,
			meta JSONB,
			business_id BIGINT NOT NULL,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		)
	`).Error; err != nil {
		return err
	}

	// Create products_taxonomies junction table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products_taxonomies (
			product_id BIGINT NOT NULL,
			taxonomy_id BIGINT NOT NULL,
			PRIMARY KEY (product_id, taxonomy_id)
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

	// Create indexes
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_business_id ON posts(business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_post_id ON products(post_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_business_id ON products(business_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_reservations_product_id ON reservations(product_id)")

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
	productRepo := repository.Repository(dbWrapper)
	postRepository := postRepo.Repository(dbWrapper)
	userRepository := userRepo.Repository(dbWrapper)

	// Create services
	userSvc := userService.Service(userRepository)
	postSvc := postService.Service(postRepository)
	productSvc := service.Service(productRepo, postSvc, userSvc)

	// Create controllers
	productCtrl := controller.Controllers(productSvc)
	postCtrl := postController.Controllers(postSvc)

	// Create product router manually
	productRouter := &product.Router{
		App:            app,
		Controller:     productCtrl,
		PostController: postCtrl,
	}
	productRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		dbWrapper.Main.Exec("DELETE FROM products_taxonomies")
		dbWrapper.Main.Exec("DELETE FROM posts_taxonomies")
		dbWrapper.Main.Exec("DELETE FROM reservations")
		dbWrapper.Main.Exec("DELETE FROM products")
		dbWrapper.Main.Exec("DELETE FROM posts")
		dbWrapper.Main.Exec("DELETE FROM taxonomies")
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
		ProductRepo:   productRepo,
		ProductRouter: productRouter,
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
func (ta *TestApp) CreateTestPost(t *testing.T, title, content string, status schema.PostStatus, postType schema.PostType, authorID, businessID uint64) *schema.Post {
	t.Helper()

	post := &schema.Post{
		Title:      title,
		Content:    content,
		Excerpt:    "Test excerpt",
		Status:     status,
		Type:       postType,
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
func (ta *TestApp) CreateTestProduct(t *testing.T, postID, businessID uint64, price float64, productType schema.ProductType, stockStatus schema.ProductStockStatus, isRoot bool) *schema.Product {
	t.Helper()

	product := &schema.Product{
		PostID:      postID,
		BusinessID:  businessID,
		Price:       price,
		MinPrice:    price,
		MaxPrice:    price,
		Type:        productType,
		StockStatus: stockStatus,
		IsRoot:      isRoot,
		OnSale:      false,
		Meta:        schema.ProductMeta{},
	}

	if err := ta.DB.Create(product).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

// CreateTestTaxonomy creates a test taxonomy in the database
func (ta *TestApp) CreateTestTaxonomy(t *testing.T, title string, taxType schema.TaxonomyType, domain schema.PostType, businessID uint64) *schema.Taxonomy {
	t.Helper()

	taxonomy := &schema.Taxonomy{
		Title:      title,
		Type:       taxType,
		Domain:     domain,
		Slug:       title + "-slug",
		BusinessID: businessID,
	}

	if err := ta.DB.Create(taxonomy).Error; err != nil {
		t.Fatalf("failed to create test taxonomy: %v", err)
	}

	return taxonomy
}

// CreateTestReservation creates a test reservation in the database
func (ta *TestApp) CreateTestReservation(t *testing.T, productID, userID, businessID uint64, startTime, endTime time.Time) *schema.Reservation {
	t.Helper()

	reservation := &schema.Reservation{
		ProductID:  productID,
		UserID:     userID,
		BusinessID: businessID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     schema.ReservationStatusReserved,
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

// CleanupProducts removes all products from the test database
func (ta *TestApp) CleanupProducts(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM products_taxonomies")
	ta.DB.Exec("DELETE FROM products")
}

// CleanupPosts removes all posts from the test database
func (ta *TestApp) CleanupPosts(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM posts_taxonomies")
	ta.DB.Exec("DELETE FROM posts")
}

// CleanupAll removes all test data from the database
func (ta *TestApp) CleanupAll(t *testing.T) {
	t.Helper()
	ta.DB.Exec("DELETE FROM products_taxonomies")
	ta.DB.Exec("DELETE FROM posts_taxonomies")
	ta.DB.Exec("DELETE FROM reservations")
	ta.DB.Exec("DELETE FROM products")
	ta.DB.Exec("DELETE FROM posts")
	ta.DB.Exec("DELETE FROM taxonomies")
	ta.DB.Exec("DELETE FROM business_users")
	ta.DB.Exec("DELETE FROM businesses")
	ta.DB.Exec("DELETE FROM users")
}
