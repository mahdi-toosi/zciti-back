package test

import (
	"fmt"
	"net/http"
	"testing"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/request"
)

// ==================== INDEX TESTS ====================

func TestIndex_WithAdminUser_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user with ROOT business permissions
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// Create test businesses
	ta.CreateTestBusiness(t, "Business 1", schema.BTypeBakery, adminUser.ID, "Description 1")
	ta.CreateTestBusiness(t, "Business 2", schema.BTypeGymManager, adminUser.ID, "Description 2")

	// Make request with businessID in path (required by middleware)
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/businesses?Page=1&Limit=10", nil, token)

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
		t.Errorf("expected 2 businesses, got %d", len(data))
	}
}

func TestIndex_WithoutAuth_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	resp := ta.MakeRequest(t, http.MethodGet, "/v1/businesses", nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestIndex_WithKeyword_FilteredResults(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// Create test businesses with different names
	ta.CreateTestBusiness(t, "Bakery Shop", schema.BTypeBakery, adminUser.ID, "Fresh bread")
	ta.CreateTestBusiness(t, "Gym Center", schema.BTypeGymManager, adminUser.ID, "Fitness center")

	// Search with keyword
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/businesses?Keyword=Bakery&Page=1&Limit=10", nil, token)

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

	if len(data) != 1 {
		t.Errorf("expected 1 business matching keyword, got %d", len(data))
	}
}

// ==================== SHOW TESTS ====================

func TestShow_WithValidID_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create a dummy business first to avoid business ID 1 (ROOT_BUSINESS_ID check in controller)
	ta.CreateTestBusiness(t, "Dummy Business", schema.BTypeBakery, ownerUser.ID, "Placeholder")

	// Create the actual test business (will have ID > 1)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeBakery, ownerUser.ID, "Test Description")

	// Update user permissions to include business owner role
	ownerPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ownerUser.Permissions = ownerPermissions
	ta.DB.Save(ownerUser)

	token := ta.GenerateJWT(t, *ownerUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/businesses/%d", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["Title"] != "Test Business" {
		t.Errorf("expected Title 'Test Business', got: %v", result["Title"])
	}
	if result["Type"] != string(schema.BTypeBakery) {
		t.Errorf("expected Type '%s', got: %v", schema.BTypeBakery, result["Type"])
	}
}

func TestShow_WithInvalidID_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create user
	permissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	user := ta.CreateTestUser(t, 9123456789, "userPass123", "Test", "User", permissions)
	token := ta.GenerateJWT(t, *user)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/businesses/99999", nil, token)

	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500 for non-existent business, got %d", resp.StatusCode)
	}
}

// ==================== STORE TESTS ====================

func TestStore_WithValidData_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// Create business request
	businessReq := request.Business{
		Title:       "New Business",
		Type:        schema.BTypeBakery,
		Description: "A new bakery",
		OwnerID:     adminUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPost, "/v1/businesses", businessReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}

	// Verify business was created
	var count int64
	ta.DB.Model(&schema.Business{}).Where("title = ?", "New Business").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 business to be created, found %d", count)
	}
}

func TestStore_WithMissingTitle_ValidationError(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// Create business request without title
	businessReq := map[string]interface{}{
		"Type":        "bakery",
		"Description": "A bakery without title",
		"OwnerID":     adminUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPost, "/v1/businesses", businessReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_WithShortTitle_ValidationError(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// Create business request with short title
	businessReq := request.Business{
		Title:       "A", // min is 2
		Type:        schema.BTypeBakery,
		Description: "A bakery",
		OwnerID:     adminUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPost, "/v1/businesses", businessReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_WithoutAuth_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	businessReq := request.Business{
		Title:       "New Business",
		Type:        schema.BTypeBakery,
		Description: "A new bakery",
		OwnerID:     1,
	}

	resp := ta.MakeRequest(t, http.MethodPost, "/v1/businesses", businessReq)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

// ==================== UPDATE TESTS ====================

func TestUpdate_AsOwner_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user
	ownerPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create business
	business := ta.CreateTestBusiness(t, "Original Title", schema.BTypeBakery, ownerUser.ID, "Original Description")

	// Update owner permissions
	ownerPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ownerUser.Permissions = ownerPermissions
	ta.DB.Save(ownerUser)

	token := ta.GenerateJWT(t, *ownerUser)

	// Update business
	updateReq := request.Business{
		Title:       "Updated Title",
		Type:        schema.BTypeBakery,
		Description: "Updated Description",
		OwnerID:     ownerUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/v1/businesses/%d", business.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify business was updated
	var updatedBusiness schema.Business
	ta.DB.First(&updatedBusiness, business.ID)
	if updatedBusiness.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", updatedBusiness.Title)
	}
	if updatedBusiness.Description != "Updated Description" {
		t.Errorf("expected description 'Updated Description', got '%s'", updatedBusiness.Description)
	}
}

func TestUpdate_AsNonOwner_Forbidden(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create another user (not owner)
	otherPermissions := schema.UserPermissions{}
	otherUser := ta.CreateTestUser(t, 9123456780, "otherPass123", "Other", "User", otherPermissions)

	// Create business owned by ownerUser
	business := ta.CreateTestBusiness(t, "Original Title", schema.BTypeBakery, ownerUser.ID, "Original Description")

	// Give otherUser some permission on business (but not owner)
	otherPermissions[business.ID] = []schema.UserRole{schema.URUser}
	otherUser.Permissions = otherPermissions
	ta.DB.Save(otherUser)

	token := ta.GenerateJWT(t, *otherUser)

	// Try to update business as non-owner
	updateReq := request.Business{
		Title:       "Hacked Title",
		Type:        schema.BTypeBakery,
		Description: "Hacked Description",
		OwnerID:     ownerUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/v1/businesses/%d", business.ID), updateReq, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for non-owner update, got %d", resp.StatusCode)
	}
}

func TestUpdate_NonExistentBusiness_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	updateReq := request.Business{
		Title:       "Updated Title",
		Type:        schema.BTypeBakery,
		Description: "Updated Description",
		OwnerID:     adminUser.ID,
	}

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPut, "/v1/businesses/99999", updateReq, token)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent business, got %d", resp.StatusCode)
	}
}

// ==================== DELETE TESTS ====================

func TestDelete_AsOwnerWithoutAdmin_Forbidden(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user (NOT admin - only business owner role)
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create business
	business := ta.CreateTestBusiness(t, "To Delete", schema.BTypeBakery, ownerUser.ID, "Will NOT be deleted")

	// Update owner permissions - only business owner, NOT admin
	// Note: URBusinessOwner doesn't have PDelete permission for DBusiness
	ownerPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ownerUser.Permissions = ownerPermissions
	ta.DB.Save(ownerUser)

	token := ta.GenerateJWT(t, *ownerUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/v1/businesses/%d", business.ID), nil, token)

	// Business owner without admin cannot delete (PDelete not in URBusinessOwner permissions)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 (business owner doesn't have delete permission), got %d", resp.StatusCode)
	}

	// Verify business was NOT deleted
	var existingBusiness schema.Business
	result := ta.DB.First(&existingBusiness, business.ID)
	if result.Error != nil {
		t.Errorf("business should still exist: %v", result.Error)
	}
}

func TestDelete_AsNonOwner_Forbidden(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create another user (not owner)
	otherPermissions := schema.UserPermissions{}
	otherUser := ta.CreateTestUser(t, 9123456780, "otherPass123", "Other", "User", otherPermissions)

	// Create business owned by ownerUser
	business := ta.CreateTestBusiness(t, "Protected Business", schema.BTypeBakery, ownerUser.ID, "Should not be deleted")

	// Give otherUser permission but not owner role
	otherPermissions[business.ID] = []schema.UserRole{schema.URUser}
	otherUser.Permissions = otherPermissions
	ta.DB.Save(otherUser)

	token := ta.GenerateJWT(t, *otherUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/v1/businesses/%d", business.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for non-owner delete, got %d", resp.StatusCode)
	}

	// Verify business was NOT deleted
	var existingBusiness schema.Business
	result := ta.DB.First(&existingBusiness, business.ID)
	if result.Error != nil {
		t.Errorf("business should still exist: %v", result.Error)
	}
}

func TestDelete_NonExistentBusiness_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodDelete, "/v1/businesses/99999", nil, token)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent business, got %d", resp.StatusCode)
	}
}

// ==================== TYPES TESTS ====================

func TestTypes_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create user
	permissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	user := ta.CreateTestUser(t, 9123456789, "userPass123", "Test", "User", permissions)
	token := ta.GenerateJWT(t, *user)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/system/businesses/types", nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponseArray(t, resp)
	if len(result) == 0 {
		t.Error("expected business types to be returned")
	}

	// Verify types have Label and Value
	for _, typ := range result {
		if typ["Label"] == nil || typ["Value"] == nil {
			t.Errorf("expected Label and Value in type, got: %v", typ)
		}
	}
}

func TestTypes_WithoutAuth_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	resp := ta.MakeRequest(t, http.MethodGet, "/v1/system/businesses/types", nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

// ==================== MY BUSINESSES TESTS ====================

func TestMyBusinesses_WithPermissions_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create user
	userPermissions := schema.UserPermissions{}
	user := ta.CreateTestUser(t, 9123456789, "userPass123", "Test", "User", userPermissions)

	// Create businesses
	business1 := ta.CreateTestBusiness(t, "My Business 1", schema.BTypeBakery, user.ID, "First business")
	business2 := ta.CreateTestBusiness(t, "My Business 2", schema.BTypeGymManager, user.ID, "Second business")
	ta.CreateTestBusiness(t, "Other Business", schema.BTypeBakery, user.ID, "Not my business")

	// Update user permissions to include only business1 and business2
	userPermissions[business1.ID] = []schema.UserRole{schema.URBusinessOwner}
	userPermissions[business2.ID] = []schema.UserRole{schema.URBusinessOwner}
	user.Permissions = userPermissions
	ta.DB.Save(user)

	token := ta.GenerateJWT(t, *user)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/user/businesses", nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponseArray(t, resp)
	if len(result) != 2 {
		t.Errorf("expected 2 businesses (only those with permissions), got %d", len(result))
	}
}

func TestMyBusinesses_NoPermissions_EmptyArray(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create user with no permissions
	userPermissions := schema.UserPermissions{}
	user := ta.CreateTestUser(t, 9123456789, "userPass123", "Test", "User", userPermissions)
	token := ta.GenerateJWT(t, *user)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/user/businesses", nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponseArray(t, resp)
	if len(result) != 0 {
		t.Errorf("expected empty array for user with no permissions, got %d items", len(result))
	}
}

func TestMyBusinesses_WithoutAuth_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	resp := ta.MakeRequest(t, http.MethodGet, "/v1/user/businesses", nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

// ==================== OPERATOR SHOW TESTS ====================

func TestOperatorShow_AsBusinessOwner_WithMeta(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create owner user
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456789, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create a dummy business first to avoid business ID 1 (ROOT_BUSINESS_ID check in controller)
	ta.CreateTestBusiness(t, "Dummy Business", schema.BTypeBakery, ownerUser.ID, "Placeholder")

	// Create business with meta (will have ID > 1)
	business := &schema.Business{
		Title:       "Business with Meta",
		Type:        schema.BTypeBakery,
		OwnerID:     ownerUser.ID,
		Description: "Has meta info",
		Account:     schema.BusinessAccountDefault,
		Meta: schema.BusinessMeta{
			ShebaNumber:    "IR123456789012345678901234",
			BankCardNumber: "6037991234567890",
		},
	}
	ta.DB.Create(business)

	// Update owner permissions
	ownerPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ownerUser.Permissions = ownerPermissions
	ta.DB.Save(ownerUser)

	token := ta.GenerateJWT(t, *ownerUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, fmt.Sprintf("/v1/business-owner/businesses/%d", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Operator show should include Meta
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in operator response, got: %v", result)
		return
	}

	if meta["ShebaNumber"] != "IR123456789012345678901234" {
		t.Errorf("expected ShebaNumber in Meta, got: %v", meta)
	}
}

// ==================== INTEGRATION TESTS ====================

func TestBusinessCRUD_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	// 1. CREATE
	createReq := request.Business{
		Title:       "Integration Test Business",
		Type:        schema.BTypeBakery,
		Description: "Created for integration test",
		OwnerID:     adminUser.ID,
	}

	createResp := ta.MakeAuthenticatedRequest(t, http.MethodPost, "/v1/businesses", createReq, token)
	if createResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, createResp)
		t.Fatalf("create failed: %v", result)
	}

	// Find the created business
	var createdBusiness schema.Business
	ta.DB.Where("title = ?", "Integration Test Business").First(&createdBusiness)
	if createdBusiness.ID == 0 {
		t.Fatal("business was not created")
	}

	// 2. READ
	showResp := ta.MakeAuthenticatedRequest(t, http.MethodGet, fmt.Sprintf("/v1/businesses/%d", createdBusiness.ID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, showResp)
		t.Errorf("show failed: %v", result)
	}

	// 3. UPDATE
	updateReq := request.Business{
		Title:       "Updated Integration Test",
		Type:        schema.BTypeGymManager,
		Description: "Updated description",
		OwnerID:     adminUser.ID,
	}

	updateResp := ta.MakeAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/v1/businesses/%d", createdBusiness.ID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, updateResp)
		t.Errorf("update failed: %v", result)
	}

	// Verify update
	var updatedBusiness schema.Business
	ta.DB.First(&updatedBusiness, createdBusiness.ID)
	if updatedBusiness.Title != "Updated Integration Test" {
		t.Errorf("expected title 'Updated Integration Test', got '%s'", updatedBusiness.Title)
	}

	// 4. DELETE
	deleteResp := ta.MakeAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/v1/businesses/%d", createdBusiness.ID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, deleteResp)
		t.Errorf("delete failed: %v", result)
	}

	// Verify deletion (soft delete)
	var deletedBusiness schema.Business
	result := ta.DB.First(&deletedBusiness, createdBusiness.ID)
	if result.Error == nil {
		t.Error("expected business to be soft deleted")
	}
}

// ==================== EMPTY BODY TESTS ====================

func TestStore_EmptyBody(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)
	token := ta.GenerateJWT(t, *adminUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPost, "/v1/businesses", nil, token)

	// Empty body should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

func TestUpdate_EmptyBody(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)

	// Create business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeBakery, adminUser.ID, "Description")

	// Update permissions
	adminPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	adminUser.Permissions = adminPermissions
	ta.DB.Save(adminUser)

	token := ta.GenerateJWT(t, *adminUser)

	resp := ta.MakeAuthenticatedRequest(t, http.MethodPut, fmt.Sprintf("/v1/businesses/%d", business.ID), nil, token)

	// Empty body should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

// ==================== PERMISSION TESTS ====================

func TestIndex_WithBusinessOwnerPermission_Forbidden(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create user with only business owner permission (not admin)
	userPermissions := schema.UserPermissions{}
	user := ta.CreateTestUser(t, 9123456789, "userPass123", "Test", "User", userPermissions)

	// Create a business and give user business owner permission
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeBakery, user.ID, "Description")
	userPermissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	user.Permissions = userPermissions
	ta.DB.Save(user)

	token := ta.GenerateJWT(t, *user)

	// Try to access /v1/businesses which requires admin or specific business permission
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/businesses", nil, token)

	// Should be forbidden since BusinessPermission middleware requires businessID in path
	// and checks DBusiness + PReadAll permission
	if resp.StatusCode != http.StatusInternalServerError && resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 or 500 for non-admin accessing index without businessID, got %d", resp.StatusCode)
	}
}

func TestDelete_AsAdmin_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create regular user as owner
	ownerPermissions := schema.UserPermissions{}
	ownerUser := ta.CreateTestUser(t, 9123456780, "ownerPass123", "Owner", "User", ownerPermissions)

	// Create admin user
	adminPermissions := schema.UserPermissions{
		schema.ROOT_BUSINESS_ID: []schema.UserRole{schema.URAdmin},
	}
	adminUser := ta.CreateTestUser(t, 9123456789, "adminPass123", "Admin", "User", adminPermissions)

	// Create business owned by regular user
	business := ta.CreateTestBusiness(t, "Regular User Business", schema.BTypeBakery, ownerUser.ID, "Owned by regular user")

	token := ta.GenerateJWT(t, *adminUser)

	// Admin should be able to delete any business
	resp := ta.MakeAuthenticatedRequest(t, http.MethodDelete, fmt.Sprintf("/v1/businesses/%d", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200 for admin delete, got %d, response: %v", resp.StatusCode, result)
	}
}

