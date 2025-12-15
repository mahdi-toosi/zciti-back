package test

import (
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/utils/helpers"

	"github.com/golang-jwt/jwt/v4"
)

// =============================================================================
// User Factory
// =============================================================================

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

// CreateTestUserWithMeta creates a test user with meta data for observer permissions
func (ta *TestApp) CreateTestUserWithMeta(t *testing.T, mobile uint64, password string, firstName, lastName string, businessID uint64, roles []schema.UserRole, meta *schema.UserMeta) *schema.User {
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
		Meta:        meta,
	}

	if err := ta.DB.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return user
}

// =============================================================================
// Business Factory
// =============================================================================

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

// =============================================================================
// Post Factory
// =============================================================================

// CreateTestPost creates a test post in the database
func (ta *TestApp) CreateTestPost(t *testing.T, title string, postType schema.PostType, businessID uint64, authorID uint64) *schema.Post {
	t.Helper()

	post := &schema.Post{
		Title:      title,
		Content:    "Test content",
		Status:     schema.PostStatusPublished,
		Type:       postType,
		BusinessID: businessID,
		AuthorID:   authorID,
	}

	if err := ta.DB.Create(post).Error; err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return post
}

// =============================================================================
// Product Factory
// =============================================================================

// CreateTestProduct creates a test product in the database
func (ta *TestApp) CreateTestProduct(t *testing.T, postID uint64, businessID uint64, price float64, productType schema.ProductType, variantType *schema.ProductVariantType) *schema.Product {
	t.Helper()

	product := &schema.Product{
		PostID:      postID,
		BusinessID:  businessID,
		Price:       price,
		Type:        productType,
		VariantType: variantType,
		StockStatus: schema.ProductStockStatusInStock,
		Meta: schema.ProductMeta{
			SKU:                  "TEST-SKU-001",
			UniWashMachineStatus: schema.UniWashMachineStatusON,
		},
	}

	if err := ta.DB.Create(product).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

// =============================================================================
// Reservation Factory
// =============================================================================

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

// =============================================================================
// Taxonomy Factory
// =============================================================================

// CreateTestTaxonomy creates a test taxonomy in the database
func (ta *TestApp) CreateTestTaxonomy(t *testing.T, title string, taxonomyType schema.TaxonomyType, businessID uint64, parentID *uint64) *schema.Taxonomy {
	t.Helper()

	taxonomy := &schema.Taxonomy{
		Title:      title,
		Type:       taxonomyType,
		Domain:     schema.PostTypePost,
		Slug:       title + "-slug",
		BusinessID: businessID,
		ParentID:   parentID,
	}

	if err := ta.DB.Create(taxonomy).Error; err != nil {
		t.Fatalf("failed to create test taxonomy: %v", err)
	}

	return taxonomy
}

// AttachTaxonomyToPost attaches a taxonomy to a post
func (ta *TestApp) AttachTaxonomyToPost(t *testing.T, postID, taxonomyID uint64) {
	t.Helper()

	if err := ta.DB.Exec("INSERT INTO posts_taxonomies (post_id, taxonomy_id) VALUES (?, ?)", postID, taxonomyID).Error; err != nil {
		t.Fatalf("failed to attach taxonomy to post: %v", err)
	}
}

// =============================================================================
// Token Generation
// =============================================================================

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

// =============================================================================
// Composite Factory Helpers
// =============================================================================

// ProductTestData holds product and its related entities
type ProductTestData struct {
	Post    *schema.Post
	Product *schema.Product
}

// CreateTestProductWithPost creates a post and product together
func (ta *TestApp) CreateTestProductWithPost(t *testing.T, businessID uint64, authorID uint64, price float64) *ProductTestData {
	t.Helper()

	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, businessID, authorID)
	variantType := schema.ProductVariantTypeWashingMachine
	product := ta.CreateTestProduct(t, post.ID, businessID, price, schema.ProductTypeVariant, &variantType)

	return &ProductTestData{
		Post:    post,
		Product: product,
	}
}

// ReservationTestData holds reservation and its related entities
type ReservationTestData struct {
	User        *schema.User
	Business    *schema.Business
	Post        *schema.Post
	Product     *schema.Product
	Reservation *schema.Reservation
	Token       string
}

// CreateFullTestReservation creates all entities needed for a complete reservation test
func (ta *TestApp) CreateFullTestReservation(t *testing.T) *ReservationTestData {
	t.Helper()

	setup := ta.SetupTestUser(t)
	productData := ta.CreateTestProductWithPost(t, setup.Business.ID, setup.User.ID, 10000)

	now := time.Now()
	reservation := ta.CreateTestReservation(t, setup.User.ID, productData.Product.ID, setup.Business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)

	return &ReservationTestData{
		User:        setup.User,
		Business:    setup.Business,
		Post:        productData.Post,
		Product:     productData.Product,
		Reservation: reservation,
		Token:       setup.Token,
	}
}
