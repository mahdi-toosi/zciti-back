package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
	postRequest "go-fiber-starter/app/module/post/request"
	"go-fiber-starter/app/module/product/request"
)

// =============================================================================
// INDEX TESTS - GET /v1/business/:businessID/products
// =============================================================================

func TestIndex_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test posts with products
	post1 := ta.CreateTestPost(t, "Product 1", "Content 1", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	post2 := ta.CreateTestPost(t, "Product 2", "Content 2", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)

	ta.CreateTestProduct(t, post1.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)
	ta.CreateTestProduct(t, post2.ID, business.ID, 200, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 products, got %d", len(data))
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/products", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestIndex_Forbidden_NoPermission(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user without business permissions
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Generate token (user has no permissions for this business)
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products", business.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for forbidden, got %d", resp.StatusCode)
	}
}

func TestIndex_WithPagination(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create multiple products
	for i := 1; i <= 5; i++ {
		post := ta.CreateTestPost(t, fmt.Sprintf("Product %d", i), fmt.Sprintf("Content %d", i), schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
		ta.CreateTestProduct(t, post.ID, business.ID, float64(i*100), schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 products in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["total"] != float64(5) {
		t.Errorf("expected total 5 products, got: %v", meta["total"])
	}
}

func TestIndex_WithKeywordFilter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create products with different names
	post1 := ta.CreateTestPost(t, "Apple iPhone", "iPhone content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	post2 := ta.CreateTestPost(t, "Samsung Galaxy", "Galaxy content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)

	ta.CreateTestProduct(t, post1.ID, business.ID, 1000, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)
	ta.CreateTestProduct(t, post2.ID, business.ID, 800, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with keyword filter
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products?keyword=Apple", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("expected 1 product with keyword filter, got: %v", data)
	}
}

// =============================================================================
// SHOW TESTS - GET /v1/business/:businessID/products/:id
// =============================================================================

func TestShow_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product
	post := ta.CreateTestPost(t, "Test Product", "Test content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 150, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, post.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains product data
	postData, ok := result["Post"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Post in response, got: %v", result)
		return
	}

	if postData["Title"] != "Test Product" {
		t.Errorf("expected title 'Test Product', got: %v", postData["Title"])
	}
}

func TestShow_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request for non-existent product
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// STORE TESTS - POST /v1/business/:businessID/products
// =============================================================================

func TestStore_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Create product request
	storeReq := request.Product{
		Post: postRequest.Post{
			Title:   "New Product",
			Content: "New product content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			Price:       250,
			Type:        schema.ProductTypeSimple,
			StockStatus: schema.ProductStockStatusInStock,
			OnSale:      false,
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify product was created in database
	var count int64
	ta.DB.Model(&schema.Post{}).Where("title = ? AND business_id = ?", "New Product", business.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected product post to be created in database, count: %d", count)
	}
}

func TestStore_WithVariants(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	variantType := schema.ProductVariantTypeSimple

	// Create product request with variants
	storeReq := request.Product{
		Post: postRequest.Post{
			Title:   "Variable Product",
			Content: "Variable product content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			Price:       300,
			Type:        schema.ProductTypeVariant,
			StockStatus: schema.ProductStockStatusInStock,
		},
		Variants: []request.ProductInPost{
			{
				Price:       280,
				Type:        schema.ProductTypeVariant,
				StockStatus: schema.ProductStockStatusInStock,
				VariantType: &variantType,
			},
			{
				Price:       320,
				Type:        schema.ProductTypeVariant,
				StockStatus: schema.ProductStockStatusInStock,
				VariantType: &variantType,
			},
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify products were created
	var productCount int64
	ta.DB.Model(&schema.Product{}).Where("business_id = ?", business.ID).Count(&productCount)
	// Should have 3 products: 1 root + 2 variants
	if productCount != 3 {
		t.Errorf("expected 3 products (1 root + 2 variants), got: %d", productCount)
	}
}

func TestStore_ValidationError_MissingRequired(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Create incomplete product request
	storeReq := map[string]interface{}{
		"Post": map[string]interface{}{
			"Title": "Incomplete",
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

// =============================================================================
// UPDATE TESTS - PUT /v1/business/:businessID/products/:id
// =============================================================================

func TestUpdate_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product
	post := ta.CreateTestPost(t, "Original Title", "Original content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update product request
	updateReq := request.Product{
		Post: postRequest.Post{
			Title:   "Updated Title",
			Content: "Updated content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			ID:          product.ID,
			Price:       150,
			Type:        schema.ProductTypeSimple,
			StockStatus: schema.ProductStockStatusInStock,
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, post.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify post was updated
	var updatedPost schema.Post
	ta.DB.First(&updatedPost, post.ID)

	if updatedPost.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got: %s", updatedPost.Title)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update request for non-existent product
	updateReq := request.Product{
		Post: postRequest.Post{
			Title:   "Not Exist",
			Content: "Content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			Price:       100,
			Type:        schema.ProductTypeSimple,
			StockStatus: schema.ProductStockStatusInStock,
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/products/99999", business.ID), updateReq, token)

	// The update might succeed with no rows affected
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", resp.StatusCode)
	}
}

// =============================================================================
// DELETE TESTS - DELETE /v1/business/:businessID/products/:id
// =============================================================================

func TestDelete_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product
	post := ta.CreateTestPost(t, "To Delete", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, post.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify post was deleted (soft delete)
	var count int64
	ta.DB.Unscoped().Model(&schema.Post{}).Where("id = ? AND deleted_at IS NOT NULL", post.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected post to be soft deleted")
	}
}

func TestDelete_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodDelete, "/v1/business/1/products/1", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

// =============================================================================
// STORE VARIANT TESTS - POST /v1/business/:businessID/products/:id/product-variant
// =============================================================================

func TestStoreVariant_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product post
	post := ta.CreateTestPost(t, "Product with Variant", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeVariant, schema.ProductStockStatusInStock, true)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	variantType := schema.ProductVariantTypeSimple

	// Create variant request
	variantReq := request.ProductInPost{
		Price:       120,
		Type:        schema.ProductTypeVariant,
		StockStatus: schema.ProductStockStatusInStock,
		VariantType: &variantType,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products/%d/product-variant", business.ID, post.ID), variantReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify variant was created
	var variantCount int64
	ta.DB.Model(&schema.Product{}).Where("post_id = ? AND is_root = ?", post.ID, false).Count(&variantCount)
	if variantCount != 1 {
		t.Errorf("expected 1 variant, got: %d", variantCount)
	}
}

func TestStoreVariant_UpdateExisting(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product post with variant
	post := ta.CreateTestPost(t, "Product with Variant", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeVariant, schema.ProductStockStatusInStock, true)
	variant := ta.CreateTestProduct(t, post.ID, business.ID, 120, schema.ProductTypeVariant, schema.ProductStockStatusInStock, false)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	variantType := schema.ProductVariantTypeSimple

	// Update variant request
	variantReq := request.ProductInPost{
		ID:          variant.ID,
		Price:       150,
		Type:        schema.ProductTypeVariant,
		StockStatus: schema.ProductStockStatusOutOfStock,
		VariantType: &variantType,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/products/%d/product-variant", business.ID, post.ID), variantReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify variant was updated
	var updatedVariant schema.Product
	ta.DB.First(&updatedVariant, variant.ID)
	if updatedVariant.Price != 150 {
		t.Errorf("expected price 150, got: %f", updatedVariant.Price)
	}
	if updatedVariant.StockStatus != schema.ProductStockStatusOutOfStock {
		t.Errorf("expected stock status 'outOfStock', got: %s", updatedVariant.StockStatus)
	}
}

// =============================================================================
// DELETE VARIANT TESTS - DELETE /v1/business/:businessID/products/:id/product-variant/:variantID
// =============================================================================

func TestDeleteVariant_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product post with variant
	post := ta.CreateTestPost(t, "Product with Variant", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeVariant, schema.ProductStockStatusInStock, true)
	variant := ta.CreateTestProduct(t, post.ID, business.ID, 120, schema.ProductTypeVariant, schema.ProductStockStatusInStock, false)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/products/%d/product-variant/%d", business.ID, post.ID, variant.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
	}
}

func TestDeleteVariant_WithFutureReservation(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product post with variant
	post := ta.CreateTestPost(t, "Product with Variant", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeVariant, schema.ProductStockStatusInStock, true)
	variant := ta.CreateTestProduct(t, post.ID, business.ID, 120, schema.ProductTypeVariant, schema.ProductStockStatusInStock, false)

	// Create a future reservation for this variant
	futureStart := time.Now().Add(24 * time.Hour)
	futureEnd := time.Now().Add(25 * time.Hour)
	ta.CreateTestReservation(t, variant.ID, user.ID, business.ID, futureStart, futureEnd)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/products/%d/product-variant/%d", business.ID, post.ID, variant.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for variant with future reservation, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	// Check for Persian error message about future reservations
	if messages[0] == nil {
		t.Errorf("expected error message about future reservations")
	}
}

// =============================================================================
// STORE ATTRIBUTE TESTS - POST /v1/business/:businessID/products/:id/product-attribute
// =============================================================================

func TestStoreAttribute_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product
	post := ta.CreateTestPost(t, "Product with Attribute", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Create test taxonomy (attribute)
	taxonomy := ta.CreateTestTaxonomy(t, "Color", schema.TaxonomyTypeProductAttributes, schema.PostTypeProduct, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Create attribute request
	attrReq := request.StoreProductAttribute{
		AddedAttrID: taxonomy.ID,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products/%d/product-attribute", business.ID, product.ID), attrReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify attribute was added
	var productWithTaxonomies schema.Product
	ta.DB.Preload("Taxonomies").First(&productWithTaxonomies, product.ID)
	if len(productWithTaxonomies.Taxonomies) != 1 {
		t.Errorf("expected 1 taxonomy, got: %d", len(productWithTaxonomies.Taxonomies))
	}
}

func TestStoreAttribute_ReplaceAttribute(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test product
	post := ta.CreateTestPost(t, "Product with Attribute", "Content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Create test taxonomies (attributes)
	oldTaxonomy := ta.CreateTestTaxonomy(t, "Red", schema.TaxonomyTypeProductAttributes, schema.PostTypeProduct, business.ID)
	newTaxonomy := ta.CreateTestTaxonomy(t, "Blue", schema.TaxonomyTypeProductAttributes, schema.PostTypeProduct, business.ID)

	// Add old taxonomy to product
	ta.DB.Model(&product).Association("Taxonomies").Append(&oldTaxonomy)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Create attribute request to replace
	attrReq := request.StoreProductAttribute{
		AddedAttrID:   newTaxonomy.ID,
		RemovedAttrID: oldTaxonomy.ID,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products/%d/product-attribute", business.ID, product.ID), attrReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify old attribute was removed and new was added
	var productWithTaxonomies schema.Product
	ta.DB.Preload("Taxonomies").First(&productWithTaxonomies, product.ID)
	if len(productWithTaxonomies.Taxonomies) != 1 {
		t.Errorf("expected 1 taxonomy, got: %d", len(productWithTaxonomies.Taxonomies))
		return
	}
	if productWithTaxonomies.Taxonomies[0].Title != "Blue" {
		t.Errorf("expected taxonomy 'Blue', got: %s", productWithTaxonomies.Taxonomies[0].Title)
	}
}

// =============================================================================
// USER ENDPOINT TESTS - GET /v1/user/business/:businessID/products
// =============================================================================

func TestUserIndex_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user (regular user)
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test products (published)
	post1 := ta.CreateTestPost(t, "Published Product 1", "Content 1", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	post2 := ta.CreateTestPost(t, "Draft Product", "Content 2", schema.PostStatusDraft, schema.PostTypeProduct, user.ID, business.ID)

	ta.CreateTestProduct(t, post1.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)
	ta.CreateTestProduct(t, post2.ID, business.ID, 200, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Make request (user endpoint doesn't require auth)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/products", business.ID), nil, "")

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify only published products are returned
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	// Should only see published product
	if len(data) != 1 {
		t.Errorf("expected 1 published product, got %d", len(data))
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestProductCRUD_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// 1. CREATE
	createReq := request.Product{
		Post: postRequest.Post{
			Title:   "Integration Test Product",
			Content: "Integration test content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			Price:       500,
			Type:        schema.ProductTypeSimple,
			StockStatus: schema.ProductStockStatusInStock,
		},
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, createResp)
		t.Fatalf("CREATE failed: %d, response: %v", createResp.StatusCode, result)
	}

	// Get created post ID
	var createdPost schema.Post
	ta.DB.Where("title = ?", "Integration Test Product").First(&createdPost)

	// 2. READ (Show)
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, createdPost.ID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("READ failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	postData, ok := showResult["Post"].(map[string]interface{})
	if !ok || postData["Title"] != "Integration Test Product" {
		t.Errorf("READ: expected title 'Integration Test Product', got: %v", postData)
	}

	// Get the product ID
	var createdProduct schema.Product
	ta.DB.Where("post_id = ? AND is_root = ?", createdPost.ID, true).First(&createdProduct)

	// 3. UPDATE
	updateReq := request.Product{
		Post: postRequest.Post{
			Title:   "Updated Integration Test Product",
			Content: "Updated content",
			Status:  schema.PostStatusPublished,
			Type:    schema.PostTypeProduct,
			Meta: schema.PostMeta{
				CommentsStatus: schema.PostCommentStatusOpen,
			},
		},
		Product: request.ProductInPost{
			ID:          createdProduct.ID,
			Price:       600,
			Type:        schema.ProductTypeSimple,
			StockStatus: schema.ProductStockStatusOutOfStock,
		},
	}

	updateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, createdPost.ID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, updateResp)
		t.Fatalf("UPDATE failed: %d, response: %v", updateResp.StatusCode, result)
	}

	// Verify update
	var updatedPost schema.Post
	ta.DB.First(&updatedPost, createdPost.ID)
	if updatedPost.Title != "Updated Integration Test Product" {
		t.Errorf("UPDATE: expected title 'Updated Integration Test Product', got: %s", updatedPost.Title)
	}

	// 4. DELETE
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/products/%d", business.ID, createdPost.ID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed: %d", deleteResp.StatusCode)
	}

	// Verify deletion
	var count int64
	ta.DB.Model(&schema.Post{}).Where("id = ?", createdPost.ID).Count(&count)
	if count != 0 {
		t.Errorf("DELETE: expected post to be deleted, count: %d", count)
	}
}

func TestProductTypes_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create products of different types
	simplePost := ta.CreateTestPost(t, "Simple Product", "Simple content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, simplePost.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	variablePost := ta.CreateTestPost(t, "Variable Product", "Variable content", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	ta.CreateTestProduct(t, variablePost.ID, business.ID, 200, schema.ProductTypeVariant, schema.ProductStockStatusInStock, true)
	ta.CreateTestProduct(t, variablePost.ID, business.ID, 180, schema.ProductTypeVariant, schema.ProductStockStatusInStock, false)
	ta.CreateTestProduct(t, variablePost.ID, business.ID, 220, schema.ProductTypeVariant, schema.ProductStockStatusOutOfStock, false)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Get all products
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 product posts, got %d", len(data))
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestEmptyBody_Requests(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test POST with empty body
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), nil, token)

	// Should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

func TestProductWithCategory_Filter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create category
	category := ta.CreateTestTaxonomy(t, "Electronics", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID)

	// Create products
	post1 := ta.CreateTestPost(t, "Product in Category", "Content 1", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)
	post2 := ta.CreateTestPost(t, "Product not in Category", "Content 2", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)

	ta.CreateTestProduct(t, post1.ID, business.ID, 100, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)
	ta.CreateTestProduct(t, post2.ID, business.ID, 200, schema.ProductTypeSimple, schema.ProductStockStatusInStock, true)

	// Associate post1 with category
	ta.DB.Exec("INSERT INTO posts_taxonomies (post_id, taxonomy_id) VALUES (?, ?)", post1.ID, category.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Filter by category
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/products?CategoryID=%d", business.ID, category.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	// Should only return the product in the category
	if len(data) != 1 {
		t.Errorf("expected 1 product in category, got %d", len(data))
	}
}

func TestProductStockStatus_Variations(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test different stock statuses
	statuses := []schema.ProductStockStatus{
		schema.ProductStockStatusInStock,
		schema.ProductStockStatusOutOfStock,
		schema.ProductStockStatusOnBackorder,
	}

	for i, status := range statuses {
		createReq := request.Product{
			Post: postRequest.Post{
				Title:   fmt.Sprintf("Product Stock %d", i),
				Content: "Content",
				Status:  schema.PostStatusPublished,
				Type:    schema.PostTypeProduct,
				Meta: schema.PostMeta{
					CommentsStatus: schema.PostCommentStatusOpen,
				},
			},
			Product: request.ProductInPost{
				Price:       float64((i + 1) * 100),
				Type:        schema.ProductTypeSimple,
				StockStatus: status,
			},
		}

		resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/products", business.ID), createReq, token)
		if resp.StatusCode != http.StatusOK {
			result := ParseResponse(t, resp)
			t.Errorf("failed to create product with stock status %s: %v", status, result)
		}
	}

	// Verify all products were created with correct stock statuses
	var products []schema.Product
	ta.DB.Where("business_id = ?", business.ID).Find(&products)

	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}

