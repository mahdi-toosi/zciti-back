package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"go-fiber-starter/app/database/schema"
	productRepo "go-fiber-starter/app/module/product/repository"
	"go-fiber-starter/app/module/reservation"
	"go-fiber-starter/app/module/reservation/controller"
	"go-fiber-starter/app/module/reservation/repository"
	"go-fiber-starter/app/module/reservation/service"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// TestApp holds the test application components
type TestApp struct {
	App               *fiber.App
	DB                *gorm.DB
	Config            *config.Config
	Cleanup           func()
	ReservationRepo   repository.IRepository
	ReservationRouter *reservation.Router
}

// getProjectRoot returns the project root directory
func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// Navigate from app/module/reservation/test/ to project root
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

// migrateTestModels creates the necessary tables for reservation testing
func migrateTestModels(db *gorm.DB) error {
	// Drop existing tables to ensure clean state
	tablesToDrop := []string{
		"posts_taxonomies", "taxonomies", "reservations", "products",
		"posts", "business_users", "businesses", "users",
	}
	for _, table := range tablesToDrop {
		db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE")
	}

	// Create tables
	if err := createUsersTable(db); err != nil {
		return err
	}
	if err := createBusinessesTable(db); err != nil {
		return err
	}
	if err := createBusinessUsersTable(db); err != nil {
		return err
	}
	if err := createPostsTable(db); err != nil {
		return err
	}
	if err := createProductsTable(db); err != nil {
		return err
	}
	if err := createReservationsTable(db); err != nil {
		return err
	}
	if err := createTaxonomiesTable(db); err != nil {
		return err
	}
	if err := createPostsTaxonomiesTable(db); err != nil {
		return err
	}

	// Create indexes
	createIndexes(db)

	return nil
}

func createUsersTable(db *gorm.DB) error {
	return db.Exec(`
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
	`).Error
}

func createBusinessesTable(db *gorm.DB) error {
	return db.Exec(`
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
	`).Error
}

func createBusinessUsersTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS business_users (
			business_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			PRIMARY KEY (business_id, user_id)
		)
	`).Error
}

func createPostsTable(db *gorm.DB) error {
	return db.Exec(`
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
	`).Error
}

func createProductsTable(db *gorm.DB) error {
	return db.Exec(`
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
	`).Error
}

func createReservationsTable(db *gorm.DB) error {
	return db.Exec(`
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
	`).Error
}

func createTaxonomiesTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS taxonomies (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			domain VARCHAR(100) NOT NULL,
			parent_id BIGINT,
			business_id BIGINT NOT NULL,
			description VARCHAR(500),
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)
	`).Error
}

func createPostsTaxonomiesTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS posts_taxonomies (
			post_id BIGINT NOT NULL,
			taxonomy_id BIGINT NOT NULL,
			PRIMARY KEY (post_id, taxonomy_id)
		)
	`).Error
}

func createIndexes(db *gorm.DB) {
	indexes := []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile ON users(mobile)",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_business_id ON reservations(business_id)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_product_id ON reservations(product_id)",
		"CREATE INDEX IF NOT EXISTS idx_reservations_deleted_at ON reservations(deleted_at)",
		"CREATE INDEX IF NOT EXISTS idx_products_post_id ON products(post_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_business_id ON products(business_id)",
		"CREATE INDEX IF NOT EXISTS idx_posts_business_id ON posts(business_id)",
		"CREATE INDEX IF NOT EXISTS idx_taxonomies_business_id ON taxonomies(business_id)",
		"CREATE INDEX IF NOT EXISTS idx_posts_taxonomies_post_id ON posts_taxonomies(post_id)",
		"CREATE INDEX IF NOT EXISTS idx_posts_taxonomies_taxonomy_id ON posts_taxonomies(taxonomy_id)",
	}
	for _, idx := range indexes {
		db.Exec(idx)
	}
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
	reservationRepo := repository.Repository(dbWrapper)
	productRepository := productRepo.Repository(dbWrapper)

	// Create reservation service
	reservationSvc := service.Service(reservationRepo, productRepository)

	// Create reservation controller
	reservationController := controller.Controllers(reservationSvc)

	// Create reservation router manually
	reservationRouter := &reservation.Router{
		App:        app,
		Controller: reservationController,
	}
	reservationRouter.RegisterRoutes(cfg)

	// Cleanup function
	cleanup := func() {
		cleanupTables := []string{
			"posts_taxonomies", "taxonomies", "reservations", "products",
			"posts", "business_users", "businesses", "users",
		}
		for _, table := range cleanupTables {
			dbWrapper.Main.Exec("DELETE FROM " + table)
		}
		dbWrapper.ShutdownDatabase()
	}

	return &TestApp{
		App:               app,
		DB:                dbWrapper.Main,
		Config:            cfg,
		Cleanup:           cleanup,
		ReservationRepo:   reservationRepo,
		ReservationRouter: reservationRouter,
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
	tables := []string{
		"posts_taxonomies", "taxonomies", "reservations", "products",
		"posts", "business_users", "businesses", "users",
	}
	for _, table := range tables {
		ta.DB.Exec("DELETE FROM " + table)
	}
}

// SetupTestUser creates a user and business with permissions for common test scenarios
func (ta *TestApp) SetupTestUser(t *testing.T, opts ...SetupOption) *TestSetup {
	t.Helper()

	cfg := &setupConfig{
		mobile:    9123456789,
		password:  "testPassword123",
		firstName: "Test",
		lastName:  "User",
		roles:     []schema.UserRole{schema.URBusinessOwner},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	user := ta.CreateTestUser(t, cfg.mobile, cfg.password, cfg.firstName, cfg.lastName, 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)
	user.Permissions[business.ID] = cfg.roles
	ta.DB.Save(user)

	token := ta.GenerateTestToken(t, user)

	return &TestSetup{
		User:     user,
		Business: business,
		Token:    token,
	}
}

// TestSetup holds commonly used test entities
type TestSetup struct {
	User     *schema.User
	Business *schema.Business
	Token    string
}

// setupConfig holds configuration for test setup
type setupConfig struct {
	mobile    uint64
	password  string
	firstName string
	lastName  string
	roles     []schema.UserRole
}

// SetupOption is a functional option for SetupTestUser
type SetupOption func(*setupConfig)

// WithMobile sets the mobile number for the test user
func WithMobile(mobile uint64) SetupOption {
	return func(c *setupConfig) {
		c.mobile = mobile
	}
}

// WithRoles sets the roles for the test user
func WithRoles(roles []schema.UserRole) SetupOption {
	return func(c *setupConfig) {
		c.roles = roles
	}
}
