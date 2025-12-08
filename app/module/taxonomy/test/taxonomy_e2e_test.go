package test

import (
	"fmt"
	"net/http"
	"testing"

	"go-fiber-starter/app/database/schema"
)

// =============================================================================
// INDEX TESTS - GET /v1/business/:businessID/taxonomies
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

	// Create test taxonomies
	ta.CreateTestTaxonomy(t, "Category One", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Tag One", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), nil, token)

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
		t.Errorf("expected 2 taxonomies, got %d", len(data))
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/taxonomies", nil, "")

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
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), nil, token)

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

	// Create multiple taxonomies
	for i := 1; i <= 5; i++ {
		ta.CreateTestTaxonomy(t, fmt.Sprintf("Category %d", i), schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 taxonomies in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["total"] != float64(5) {
		t.Errorf("expected total 5 taxonomies, got: %v", meta["total"])
	}
}

func TestIndex_FilterByType(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create taxonomies of different types
	ta.CreateTestTaxonomy(t, "Category One", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Category Two", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Tag One", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request filtering by category type
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies?Type=category", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only categories
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 categories, got %d", len(data))
	}
}

func TestIndex_FilterByDomain(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create taxonomies for different domains
	ta.CreateTestTaxonomy(t, "Post Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Product Category", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request filtering by post domain
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies?Domain=post", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only post domain taxonomies
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 1 {
		t.Errorf("expected 1 post domain taxonomy, got %d", len(data))
	}
}

// =============================================================================
// SEARCH TESTS - GET /v1/business/:businessID/taxonomies/search
// =============================================================================

func TestSearch_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test taxonomies
	ta.CreateTestTaxonomy(t, "Electronics", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Clothing", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Electronic Accessories", schema.TaxonomyTypeTag, schema.PostTypeProduct, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with keyword search
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/search?Keyword=Electronic", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains matching taxonomies
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 matching taxonomies, got %d", len(data))
	}
}

func TestSearch_WithParentID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create parent taxonomy
	parent := ta.CreateTestTaxonomy(t, "Parent Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Create child taxonomies
	ta.CreateTestTaxonomy(t, "Child One", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, &parent.ID)
	ta.CreateTestTaxonomy(t, "Child Two", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, &parent.ID)

	// Create another root taxonomy
	ta.CreateTestTaxonomy(t, "Another Root", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request filtering by parent ID
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/search?ParentID=%d", business.ID, parent.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only children of the parent
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 child taxonomies, got %d", len(data))
	}
}

func TestSearch_RootTaxonomiesOnly(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create parent taxonomy
	parent := ta.CreateTestTaxonomy(t, "Parent Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Create child taxonomy
	ta.CreateTestTaxonomy(t, "Child One", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, &parent.ID)

	// Create another root taxonomy
	ta.CreateTestTaxonomy(t, "Another Root", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request filtering by ParentID = -1 (root only)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/search?ParentID=-1", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only root taxonomies
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 root taxonomies, got %d", len(data))
	}
}

// =============================================================================
// SHOW TESTS - GET /v1/business/:businessID/taxonomies/:id
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

	// Create test taxonomy
	taxonomy := ta.CreateTestTaxonomy(t, "Test Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, taxonomy.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains taxonomy data
	if result["Title"] != "Test Category" {
		t.Errorf("expected title 'Test Category', got: %v", result["Title"])
	}

	if result["Type"] != string(schema.TaxonomyTypeCategory) {
		t.Errorf("expected type 'category', got: %v", result["Type"])
	}

	if result["Domain"] != string(schema.PostTypePost) {
		t.Errorf("expected domain 'post', got: %v", result["Domain"])
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

	// Make request for non-existent taxonomy
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// STORE TESTS - POST /v1/business/:businessID/taxonomies
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

	// Create taxonomy request
	storeReq := map[string]interface{}{
		"Title":       "New Category",
		"Type":        string(schema.TaxonomyTypeCategory),
		"Domain":      string(schema.PostTypePost),
		"Description": "A new test category",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}

	// Verify taxonomy was created in database
	var count int64
	ta.DB.Model(&schema.Taxonomy{}).Where("title = ?", "New Category").Count(&count)
	if count != 1 {
		t.Errorf("expected taxonomy to be created in database, count: %d", count)
	}
}

func TestStore_WithParent(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create parent taxonomy
	parent := ta.CreateTestTaxonomy(t, "Parent Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Create child taxonomy request
	storeReq := map[string]interface{}{
		"Title":       "Child Category",
		"Type":        string(schema.TaxonomyTypeCategory),
		"Domain":      string(schema.PostTypePost),
		"ParentID":    parent.ID,
		"Description": "A child category",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify child taxonomy was created with parent reference
	var child schema.Taxonomy
	ta.DB.Where("title = ?", "Child Category").First(&child)
	if child.ParentID == nil || *child.ParentID != parent.ID {
		t.Errorf("expected child taxonomy to have parent ID %d, got: %v", parent.ID, child.ParentID)
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

	// Create taxonomy request without required fields
	storeReq := map[string]interface{}{
		"Title": "T", // Too short (min=2)
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_ValidationError_InvalidType(t *testing.T) {
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

	// Create taxonomy request with invalid type
	storeReq := map[string]interface{}{
		"Title":  "Test Category",
		"Type":   "invalidType", // Invalid type
		"Domain": string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_ValidationError_InvalidDomain(t *testing.T) {
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

	// Create taxonomy request with invalid domain
	storeReq := map[string]interface{}{
		"Title":  "Test Category",
		"Type":   string(schema.TaxonomyTypeCategory),
		"Domain": "invalidDomain", // Invalid domain
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_AllTaxonomyTypes(t *testing.T) {
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

	taxonomyTypes := []schema.TaxonomyType{
		schema.TaxonomyTypeCategory,
		schema.TaxonomyTypeTag,
		schema.TaxonomyTypeProductAttributes,
	}

	for _, taxonomyType := range taxonomyTypes {
		// Create taxonomy request
		storeReq := map[string]interface{}{
			"Title":  fmt.Sprintf("Test %s", taxonomyType),
			"Type":   string(taxonomyType),
			"Domain": string(schema.PostTypeProduct),
		}

		// Make request
		resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

		if resp.StatusCode != http.StatusOK {
			result := ParseResponse(t, resp)
			t.Errorf("expected status 200 for type %s, got %d, response: %v", taxonomyType, resp.StatusCode, result)
		}
	}

	// Verify all types were created
	var count int64
	ta.DB.Model(&schema.Taxonomy{}).Where("business_id = ?", business.ID).Count(&count)
	if count != 3 {
		t.Errorf("expected 3 taxonomies, got %d", count)
	}
}

// =============================================================================
// UPDATE TESTS - PUT /v1/business/:businessID/taxonomies/:id
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

	// Create test taxonomy
	taxonomy := ta.CreateTestTaxonomy(t, "Original Title", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update taxonomy request
	updateReq := map[string]interface{}{
		"Title":       "Updated Title",
		"Type":        string(schema.TaxonomyTypeTag),
		"Domain":      string(schema.PostTypeProduct),
		"Description": "Updated description",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, taxonomy.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify taxonomy was updated in database
	var updatedTaxonomy schema.Taxonomy
	ta.DB.First(&updatedTaxonomy, taxonomy.ID)

	if updatedTaxonomy.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got: %s", updatedTaxonomy.Title)
	}

	if updatedTaxonomy.Type != schema.TaxonomyTypeTag {
		t.Errorf("expected type 'tag', got: %s", updatedTaxonomy.Type)
	}

	if updatedTaxonomy.Domain != schema.PostTypeProduct {
		t.Errorf("expected domain 'product', got: %s", updatedTaxonomy.Domain)
	}

	if updatedTaxonomy.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got: %s", updatedTaxonomy.Description)
	}
}

func TestUpdate_ChangeParent(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create parent taxonomies
	parent1 := ta.CreateTestTaxonomy(t, "Parent One", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	parent2 := ta.CreateTestTaxonomy(t, "Parent Two", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Create child taxonomy under parent1
	child := ta.CreateTestTaxonomy(t, "Child Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, &parent1.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update taxonomy to move to parent2
	updateReq := map[string]interface{}{
		"Title":    "Child Category",
		"Type":     string(schema.TaxonomyTypeCategory),
		"Domain":   string(schema.PostTypePost),
		"ParentID": parent2.ID,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, child.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify parent was changed
	var updatedChild schema.Taxonomy
	ta.DB.First(&updatedChild, child.ID)

	if updatedChild.ParentID == nil || *updatedChild.ParentID != parent2.ID {
		t.Errorf("expected parent ID %d, got: %v", parent2.ID, updatedChild.ParentID)
	}
}

func TestUpdate_ValidationError(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test taxonomy
	taxonomy := ta.CreateTestTaxonomy(t, "Test Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update request with invalid data
	updateReq := map[string]interface{}{
		"Title":  "T",                            // Too short
		"Type":   string(schema.TaxonomyTypeTag), // Valid but title is not
		"Domain": string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, taxonomy.ID), updateReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

// =============================================================================
// DELETE TESTS - DELETE /v1/business/:businessID/taxonomies/:id
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

	// Create test taxonomy
	taxonomy := ta.CreateTestTaxonomy(t, "To Delete", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, taxonomy.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify taxonomy was deleted (soft delete)
	var count int64
	ta.DB.Unscoped().Model(&schema.Taxonomy{}).Where("id = ? AND deleted_at IS NOT NULL", taxonomy.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected taxonomy to be soft deleted")
	}
}

func TestDelete_WithChildren_ShouldFail(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create parent taxonomy
	parent := ta.CreateTestTaxonomy(t, "Parent Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)

	// Create child taxonomy
	ta.CreateTestTaxonomy(t, "Child Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, &parent.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Try to delete parent (should fail due to child)
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, parent.ID), nil, token)

	// Should return FailedDependency (424)
	if resp.StatusCode != http.StatusFailedDependency {
		t.Errorf("expected status 424 for failed dependency, got %d", resp.StatusCode)
	}

	// Verify parent was not deleted
	var count int64
	ta.DB.Model(&schema.Taxonomy{}).Where("id = ?", parent.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected parent taxonomy to still exist")
	}
}

func TestDelete_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodDelete, "/v1/business/1/taxonomies/1", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

// =============================================================================
// USER ROUTES TESTS - GET /v1/user/business/:businessID/taxonomies
// =============================================================================

func TestUserSearch_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user (no special permissions needed for user routes)
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test taxonomies
	ta.CreateTestTaxonomy(t, "User Category", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID, nil)
	ta.CreateTestTaxonomy(t, "User Tag", schema.TaxonomyTypeTag, schema.PostTypeProduct, business.ID, nil)

	// Note: User routes (/v1/user/...) are typically unprotected or have mdl.ForUser middleware
	// Make request without token (public route)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/taxonomies", business.ID), nil, "")

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
		t.Errorf("expected 2 taxonomies, got %d", len(data))
	}
}

func TestUserSearch_WithFilters(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test taxonomies
	ta.CreateTestTaxonomy(t, "Category One", schema.TaxonomyTypeCategory, schema.PostTypeProduct, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Category Two", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID, nil)
	ta.CreateTestTaxonomy(t, "Tag One", schema.TaxonomyTypeTag, schema.PostTypeProduct, business.ID, nil)

	// Make request with filters
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/taxonomies?Type=category&Domain=product", business.ID), nil, "")

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only matching taxonomies
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 1 {
		t.Errorf("expected 1 taxonomy with type=category and domain=product, got %d", len(data))
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestCRUD_Integration(t *testing.T) {
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
	createReq := map[string]interface{}{
		"Title":       "Integration Category",
		"Type":        string(schema.TaxonomyTypeCategory),
		"Domain":      string(schema.PostTypePost),
		"Description": "Integration test category",
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE failed: %d", createResp.StatusCode)
	}

	// Get created taxonomy ID
	var createdTaxonomy schema.Taxonomy
	ta.DB.Where("title = ?", "Integration Category").First(&createdTaxonomy)

	// 2. READ (Show)
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, createdTaxonomy.ID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("READ failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	if showResult["Title"] != "Integration Category" {
		t.Errorf("READ: expected title 'Integration Category', got: %v", showResult["Title"])
	}

	// 3. UPDATE
	updateReq := map[string]interface{}{
		"Title":       "Updated Integration Category",
		"Type":        string(schema.TaxonomyTypeTag),
		"Domain":      string(schema.PostTypeProduct),
		"Description": "Updated integration test",
	}

	updateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, createdTaxonomy.ID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("UPDATE failed: %d", updateResp.StatusCode)
	}

	// Verify update
	var updatedTaxonomy schema.Taxonomy
	ta.DB.First(&updatedTaxonomy, createdTaxonomy.ID)
	if updatedTaxonomy.Title != "Updated Integration Category" {
		t.Errorf("UPDATE: expected title 'Updated Integration Category', got: %s", updatedTaxonomy.Title)
	}

	// 4. DELETE
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, createdTaxonomy.ID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed: %d", deleteResp.StatusCode)
	}

	// Verify deletion
	var count int64
	ta.DB.Model(&schema.Taxonomy{}).Where("id = ?", createdTaxonomy.ID).Count(&count)
	if count != 0 {
		t.Errorf("DELETE: expected taxonomy to be deleted, count: %d", count)
	}
}

func TestHierarchicalTaxonomy_Integration(t *testing.T) {
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

	// 1. Create root category
	rootReq := map[string]interface{}{
		"Title":  "Electronics",
		"Type":   string(schema.TaxonomyTypeCategory),
		"Domain": string(schema.PostTypeProduct),
	}

	rootResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), rootReq, token)
	if rootResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE root failed: %d", rootResp.StatusCode)
	}

	var rootTaxonomy schema.Taxonomy
	ta.DB.Where("title = ?", "Electronics").First(&rootTaxonomy)

	// 2. Create child categories
	childReq := map[string]interface{}{
		"Title":    "Smartphones",
		"Type":     string(schema.TaxonomyTypeCategory),
		"Domain":   string(schema.PostTypeProduct),
		"ParentID": rootTaxonomy.ID,
	}

	childResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), childReq, token)
	if childResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE child failed: %d", childResp.StatusCode)
	}

	var childTaxonomy schema.Taxonomy
	ta.DB.Where("title = ?", "Smartphones").First(&childTaxonomy)

	// 3. Create grandchild
	grandchildReq := map[string]interface{}{
		"Title":    "Android Phones",
		"Type":     string(schema.TaxonomyTypeCategory),
		"Domain":   string(schema.PostTypeProduct),
		"ParentID": childTaxonomy.ID,
	}

	grandchildResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), grandchildReq, token)
	if grandchildResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE grandchild failed: %d", grandchildResp.StatusCode)
	}

	// 4. Verify hierarchy by searching with ParentID
	searchResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies/search?ParentID=%d", business.ID, rootTaxonomy.ID), nil, token)
	if searchResp.StatusCode != http.StatusOK {
		t.Fatalf("SEARCH children failed: %d", searchResp.StatusCode)
	}

	searchResult := ParseResponse(t, searchResp)
	data, ok := searchResult["Data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("expected 1 direct child, got: %v", data)
	}

	// 5. Try to delete parent with children (should fail)
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, rootTaxonomy.ID), nil, token)
	if deleteResp.StatusCode != http.StatusFailedDependency {
		t.Errorf("expected status 424 when deleting parent with children, got %d", deleteResp.StatusCode)
	}

	// 6. Delete in correct order (grandchild -> child -> root)
	var grandchildTaxonomy schema.Taxonomy
	ta.DB.Where("title = ?", "Android Phones").First(&grandchildTaxonomy)

	// Delete grandchild
	ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, grandchildTaxonomy.ID), nil, token)
	// Delete child
	ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, childTaxonomy.ID), nil, token)
	// Delete root
	finalDeleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/taxonomies/%d", business.ID, rootTaxonomy.ID), nil, token)

	if finalDeleteResp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 when deleting root after children removed, got %d", finalDeleteResp.StatusCode)
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
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), nil, token)

	// Should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

func TestLongTitle(t *testing.T) {
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

	// Create taxonomy with title exceeding max length (100 chars)
	longTitle := ""
	for i := 0; i < 150; i++ {
		longTitle += "a"
	}

	storeReq := map[string]interface{}{
		"Title":  longTitle,
		"Type":   string(schema.TaxonomyTypeCategory),
		"Domain": string(schema.PostTypePost),
	}

	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	// Should return validation error for exceeding max length
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for title exceeding max length, got %d", resp.StatusCode)
	}
}

func TestLongDescription(t *testing.T) {
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

	// Create taxonomy with description exceeding max length (500 chars)
	longDesc := ""
	for i := 0; i < 600; i++ {
		longDesc += "a"
	}

	storeReq := map[string]interface{}{
		"Title":       "Valid Title",
		"Type":        string(schema.TaxonomyTypeCategory),
		"Domain":      string(schema.PostTypePost),
		"Description": longDesc,
	}

	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

	// Should return validation error for exceeding max length
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for description exceeding max length, got %d", resp.StatusCode)
	}
}

func TestMultipleDomains(t *testing.T) {
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

	domains := []schema.PostType{
		schema.PostTypePost,
		schema.PostTypePage,
		schema.PostTypeProduct,
	}

	for _, domain := range domains {
		storeReq := map[string]interface{}{
			"Title":  fmt.Sprintf("Category for %s", domain),
			"Type":   string(schema.TaxonomyTypeCategory),
			"Domain": string(domain),
		}

		resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/taxonomies", business.ID), storeReq, token)

		if resp.StatusCode != http.StatusOK {
			result := ParseResponse(t, resp)
			t.Errorf("expected status 200 for domain %s, got %d, response: %v", domain, resp.StatusCode, result)
		}
	}

	// Verify all domains were created
	var count int64
	ta.DB.Model(&schema.Taxonomy{}).Where("business_id = ?", business.ID).Count(&count)
	if count != 3 {
		t.Errorf("expected 3 taxonomies, got %d", count)
	}
}

func TestCrossBusinessIsolation(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create two users with their own businesses
	user1 := ta.CreateTestUser(t, 9123456781, "testPassword123", "User", "One", 0, nil)
	business1 := ta.CreateTestBusiness(t, "Business One", schema.BTypeGymManager, user1.ID)
	user1.Permissions[business1.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user1)

	user2 := ta.CreateTestUser(t, 9123456782, "testPassword123", "User", "Two", 0, nil)
	business2 := ta.CreateTestBusiness(t, "Business Two", schema.BTypeGymManager, user2.ID)
	user2.Permissions[business2.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user2)

	// Create taxonomies for each business
	ta.CreateTestTaxonomy(t, "Business1 Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business1.ID, nil)
	ta.CreateTestTaxonomy(t, "Business2 Category", schema.TaxonomyTypeCategory, schema.PostTypePost, business2.ID, nil)

	// Generate token for user1
	token1 := ta.GenerateTestToken(t, user1)

	// User1 should only see business1 taxonomies
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies", business1.ID), nil, token1)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Fatalf("expected Data array in response")
	}

	if len(data) != 1 {
		t.Errorf("expected 1 taxonomy for business1, got %d", len(data))
	}

	// User1 should not have access to business2 taxonomies
	resp2 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/taxonomies", business2.ID), nil, token1)

	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for cross-business access, got %d", resp2.StatusCode)
	}
}

