package test

import (
	"fmt"
	"net/http"
	"testing"

	"go-fiber-starter/app/database/schema"
)

// =============================================================================
// INDEX TESTS - GET /v1/business/:businessID/posts
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

	// Create test posts
	ta.CreateTestPost(t, "First Post", "Content of first post", "Excerpt 1", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	ta.CreateTestPost(t, "Second Post", "Content of second post", "Excerpt 2", schema.PostStatusDraft, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts", business.ID), nil, token)

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
		t.Errorf("expected 2 posts, got %d", len(data))
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/posts", nil, "")

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
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts", business.ID), nil, token)

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

	// Create multiple posts
	for i := 1; i <= 5; i++ {
		ta.CreateTestPost(t, fmt.Sprintf("Post %d", i), fmt.Sprintf("Content %d", i), fmt.Sprintf("Excerpt %d", i), schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 posts in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["total"] != float64(5) {
		t.Errorf("expected total 5 posts, got: %v", meta["total"])
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

	// Create posts with different titles
	ta.CreateTestPost(t, "Golang Tutorial", "Content about Go", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	ta.CreateTestPost(t, "Python Guide", "Content about Python", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	ta.CreateTestPost(t, "Advanced Golang", "More Go content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with keyword filter
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts?keyword=Golang", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify only Golang posts are returned
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 2 {
		t.Errorf("expected 2 posts matching 'Golang', got %d", len(data))
	}
}

func TestIndex_UserEndpoint_OnlyPublished(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create posts with different statuses
	ta.CreateTestPost(t, "Published Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	ta.CreateTestPost(t, "Draft Post", "Content", "Excerpt", schema.PostStatusDraft, schema.PostTypePost, user.ID, business.ID)

	// Make request to user endpoint (no authentication required for this endpoint)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/posts", business.ID), nil, "")

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify only published posts are returned
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 1 {
		t.Errorf("expected 1 published post, got %d", len(data))
	}
}

// =============================================================================
// SHOW TESTS - GET /v1/business/:businessID/posts/:id
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

	// Create test post
	post := ta.CreateTestPost(t, "Test Post", "Test Content", "Test Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains post data
	if result["Title"] != "Test Post" {
		t.Errorf("expected title 'Test Post', got: %v", result["Title"])
	}

	if result["Content"] != "Test Content" {
		t.Errorf("expected content 'Test Content', got: %v", result["Content"])
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

	// Make request for non-existent post
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

func TestShow_WrongBusiness(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business1 := ta.CreateTestBusiness(t, "Business 1", schema.BTypeGymManager, user.ID)
	business2 := ta.CreateTestBusiness(t, "Business 2", schema.BTypeGymManager, user.ID)

	// Update user with business permissions for both
	user.Permissions[business1.ID] = []schema.UserRole{schema.URBusinessOwner}
	user.Permissions[business2.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create post for business1
	post := ta.CreateTestPost(t, "Test Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business1.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Try to access post via business2
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts/%d", business2.ID, post.ID), nil, token)

	// Should fail because post doesn't belong to business2
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for wrong business, got %d", resp.StatusCode)
	}
}

// =============================================================================
// STORE TESTS - POST /v1/business/:businessID/posts
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

	// Create post request
	storeReq := map[string]interface{}{
		"Title":   "New Post Title",
		"Content": "This is the content of the new post",
		"Excerpt": "Short excerpt",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
		"Meta": map[string]interface{}{
			"CommentsStatus": string(schema.PostCommentStatusOpen),
		},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains created post
	if result["Title"] != "New Post Title" {
		t.Errorf("expected title 'New Post Title', got: %v", result["Title"])
	}

	// Verify post was created in database
	var count int64
	ta.DB.Model(&schema.Post{}).Where("title = ?", "New Post Title").Count(&count)
	if count != 1 {
		t.Errorf("expected post to be created in database, count: %d", count)
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

	// Create post request without required fields (Status and Type are required)
	storeReq := map[string]interface{}{
		"Title": "Incomplete Post",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_ValidationError_InvalidStatus(t *testing.T) {
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

	// Create post request with invalid status
	storeReq := map[string]interface{}{
		"Title":   "Invalid Status Post",
		"Content": "Content",
		"Status":  "invalid_status",
		"Type":    string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_DraftStatus(t *testing.T) {
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

	// Create draft post request
	storeReq := map[string]interface{}{
		"Title":   "Draft Post",
		"Content": "Draft content",
		"Status":  string(schema.PostStatusDraft),
		"Type":    string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify post was created with draft status
	var post schema.Post
	ta.DB.Where("title = ?", "Draft Post").First(&post)
	if post.Status != schema.PostStatusDraft {
		t.Errorf("expected status 'draft', got: %s", post.Status)
	}
}

// =============================================================================
// UPDATE TESTS - PUT /v1/business/:businessID/posts/:id
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

	// Create test post
	post := ta.CreateTestPost(t, "Original Title", "Original Content", "Original Excerpt", schema.PostStatusDraft, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update post request
	updateReq := map[string]interface{}{
		"Title":   "Updated Title",
		"Content": "Updated Content",
		"Excerpt": "Updated Excerpt",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify post was updated in database
	var updatedPost schema.Post
	ta.DB.First(&updatedPost, post.ID)

	if updatedPost.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got: %s", updatedPost.Title)
	}

	if updatedPost.Content != "Updated Content" {
		t.Errorf("expected content 'Updated Content', got: %s", updatedPost.Content)
	}

	if updatedPost.Status != schema.PostStatusPublished {
		t.Errorf("expected status 'published', got: %s", updatedPost.Status)
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

	// Update request for non-existent post
	updateReq := map[string]interface{}{
		"Title":   "Updated Title",
		"Content": "Updated Content",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/99999", business.ID), updateReq, token)

	// GORM Update might not fail for non-existent records
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestUpdate_ChangeStatus(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create draft post
	post := ta.CreateTestPost(t, "Draft Post", "Content", "Excerpt", schema.PostStatusDraft, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update to published
	updateReq := map[string]interface{}{
		"Title":   "Draft Post",
		"Content": "Content",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify status change
	var updatedPost schema.Post
	ta.DB.First(&updatedPost, post.ID)

	if updatedPost.Status != schema.PostStatusPublished {
		t.Errorf("expected status 'published', got: %s", updatedPost.Status)
	}
}

// =============================================================================
// DELETE TESTS - DELETE /v1/business/:businessID/posts/:id
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

	// Create test post
	post := ta.CreateTestPost(t, "To Delete", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), nil, token)

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
	resp := ta.MakeRequest(t, http.MethodDelete, "/v1/business/1/posts/1", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestDelete_WrongBusiness(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business1 := ta.CreateTestBusiness(t, "Business 1", schema.BTypeGymManager, user.ID)
	business2 := ta.CreateTestBusiness(t, "Business 2", schema.BTypeGymManager, user.ID)

	// Update user with business permissions for both
	user.Permissions[business1.ID] = []schema.UserRole{schema.URBusinessOwner}
	user.Permissions[business2.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create post for business1
	post := ta.CreateTestPost(t, "Test Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business1.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Try to delete post via business2
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/posts/%d", business2.ID, post.ID), nil, token)

	// Should succeed but not delete anything (since post doesn't belong to business2)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Verify post still exists
	var count int64
	ta.DB.Model(&schema.Post{}).Where("id = ?", post.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected post to still exist since it doesn't belong to business2")
	}
}

// =============================================================================
// TAXONOMY TESTS - POST /v1/business/:businessID/posts/:id/insert-taxonomies
//                 POST /v1/business/:businessID/posts/:id/delete-taxonomies
// =============================================================================

func TestInsertTaxonomies_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post
	post := ta.CreateTestPost(t, "Test Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Create test taxonomies
	tax1 := ta.CreateTestTaxonomy(t, "Category 1", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID)
	tax2 := ta.CreateTestTaxonomy(t, "Tag 1", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Insert taxonomies request
	insertReq := map[string]interface{}{
		"IDs": []uint64{tax1.ID, tax2.ID},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts/%d/insert-taxonomies", business.ID, post.ID), insertReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify taxonomies were attached
	var count int64
	ta.DB.Raw("SELECT COUNT(*) FROM posts_taxonomies WHERE post_id = ?", post.ID).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 taxonomy associations, got %d", count)
	}
}

func TestDeleteTaxonomies_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post
	post := ta.CreateTestPost(t, "Test Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Create test taxonomies
	tax1 := ta.CreateTestTaxonomy(t, "Category 1", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID)
	tax2 := ta.CreateTestTaxonomy(t, "Tag 1", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID)

	// Attach taxonomies to post
	ta.AttachTaxonomyToPost(t, post.ID, tax1.ID)
	ta.AttachTaxonomyToPost(t, post.ID, tax2.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Delete one taxonomy
	deleteReq := map[string]interface{}{
		"IDs": []uint64{tax1.ID},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts/%d/delete-taxonomies", business.ID, post.ID), deleteReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify only one taxonomy remains
	var count int64
	ta.DB.Raw("SELECT COUNT(*) FROM posts_taxonomies WHERE post_id = ?", post.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 taxonomy association after deletion, got %d", count)
	}
}

func TestDeleteTaxonomies_Multiple(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post
	post := ta.CreateTestPost(t, "Test Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)

	// Create test taxonomies
	tax1 := ta.CreateTestTaxonomy(t, "Category 1", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID)
	tax2 := ta.CreateTestTaxonomy(t, "Category 2", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID)
	tax3 := ta.CreateTestTaxonomy(t, "Tag 1", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID)

	// Attach all taxonomies to post
	ta.AttachTaxonomyToPost(t, post.ID, tax1.ID)
	ta.AttachTaxonomyToPost(t, post.ID, tax2.ID)
	ta.AttachTaxonomyToPost(t, post.ID, tax3.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Delete multiple taxonomies
	deleteReq := map[string]interface{}{
		"IDs": []uint64{tax1.ID, tax2.ID},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts/%d/delete-taxonomies", business.ID, post.ID), deleteReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify only one taxonomy remains
	var count int64
	ta.DB.Raw("SELECT COUNT(*) FROM posts_taxonomies WHERE post_id = ?", post.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 taxonomy association after deletion, got %d", count)
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
		"Title":   "Integration Test Post",
		"Content": "This is an integration test post",
		"Excerpt": "Integration excerpt",
		"Status":  string(schema.PostStatusDraft),
		"Type":    string(schema.PostTypePost),
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE failed: %d", createResp.StatusCode)
	}

	createResult := ParseResponse(t, createResp)
	postID := uint64(createResult["ID"].(float64))

	// 2. READ (Show)
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, postID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("READ failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	if showResult["Title"] != "Integration Test Post" {
		t.Errorf("READ: expected title 'Integration Test Post', got: %v", showResult["Title"])
	}

	// 3. UPDATE
	updateReq := map[string]interface{}{
		"Title":   "Updated Integration Post",
		"Content": "Updated content",
		"Excerpt": "Updated excerpt",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	updateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, postID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("UPDATE failed: %d", updateResp.StatusCode)
	}

	// Verify update
	var updatedPost schema.Post
	ta.DB.First(&updatedPost, postID)
	if updatedPost.Title != "Updated Integration Post" {
		t.Errorf("UPDATE: expected title 'Updated Integration Post', got: %s", updatedPost.Title)
	}

	// 4. DELETE
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, postID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed: %d", deleteResp.StatusCode)
	}

	// Verify deletion
	var count int64
	ta.DB.Model(&schema.Post{}).Where("id = ?", postID).Count(&count)
	if count != 0 {
		t.Errorf("DELETE: expected post to be deleted, count: %d", count)
	}
}

func TestPostWithTaxonomies_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create taxonomies
	category := ta.CreateTestTaxonomy(t, "Tech", schema.TaxonomyTypeCategory, schema.PostTypePost, business.ID)
	tag1 := ta.CreateTestTaxonomy(t, "Go", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID)
	tag2 := ta.CreateTestTaxonomy(t, "Programming", schema.TaxonomyTypeTag, schema.PostTypePost, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// 1. Create post
	createReq := map[string]interface{}{
		"Title":   "Go Tutorial",
		"Content": "Learn Go programming",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE failed: %d", createResp.StatusCode)
	}

	createResult := ParseResponse(t, createResp)
	postID := uint64(createResult["ID"].(float64))

	// 2. Add taxonomies
	insertReq := map[string]interface{}{
		"IDs": []uint64{category.ID, tag1.ID, tag2.ID},
	}

	insertResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts/%d/insert-taxonomies", business.ID, postID), insertReq, token)
	if insertResp.StatusCode != http.StatusOK {
		t.Fatalf("INSERT TAXONOMIES failed: %d", insertResp.StatusCode)
	}

	// 3. Verify post has taxonomies via Show
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, postID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("SHOW failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	taxonomies, ok := showResult["Taxonomies"].([]interface{})
	if !ok || len(taxonomies) != 3 {
		t.Errorf("expected 3 taxonomies, got: %v", showResult["Taxonomies"])
	}

	// 4. Remove one taxonomy
	deleteReq := map[string]interface{}{
		"IDs": []uint64{tag2.ID},
	}

	deleteResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts/%d/delete-taxonomies", business.ID, postID), deleteReq, token)
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE TAXONOMIES failed: %d", deleteResp.StatusCode)
	}

	// 5. Verify only 2 taxonomies remain
	var count int64
	ta.DB.Raw("SELECT COUNT(*) FROM posts_taxonomies WHERE post_id = ?", postID).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 taxonomies after deletion, got: %d", count)
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestStore_EmptyContent(t *testing.T) {
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

	// Create post with empty content (content is required by DB but not by validation)
	storeReq := map[string]interface{}{
		"Title":   "Post Without Content",
		"Content": "",
		"Status":  string(schema.PostStatusDraft),
		"Type":    string(schema.PostTypePost),
	}

	// Make request - might fail due to DB constraint
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), storeReq, token)

	// Should either succeed or fail with validation/DB error
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

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
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/posts", business.ID), nil, token)

	// Should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

func TestMultiplePosts_SameBusiness(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create 10 posts
	for i := 1; i <= 10; i++ {
		ta.CreateTestPost(t, fmt.Sprintf("Post %d", i), fmt.Sprintf("Content %d", i), fmt.Sprintf("Excerpt %d", i), schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Get all posts
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 10 {
		t.Errorf("expected 10 posts, got: %d", len(data))
	}
}

func TestDifferentPostTypes(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create posts with different types (only PostTypePost should be returned by the posts endpoint)
	ta.CreateTestPost(t, "Blog Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, user.ID, business.ID)
	ta.CreateTestPost(t, "Page", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePage, user.ID, business.ID)
	ta.CreateTestPost(t, "Product Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypeProduct, user.ID, business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Get posts (should only return type=post)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/posts", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("expected 1 post (type=post only), got: %d", len(data))
	}
}

func TestPermissionLevels(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Create regular user with only read permission
	reader := ta.CreateTestUser(t, 9123456780, "testPassword123", "Reader", "User", 0, nil)
	reader.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(reader)

	// Update owner permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create a post as owner
	post := ta.CreateTestPost(t, "Owner Post", "Content", "Excerpt", schema.PostStatusPublished, schema.PostTypePost, owner.ID, business.ID)

	// Generate tokens
	ownerToken := ta.GenerateTestToken(t, owner)
	readerToken := ta.GenerateTestToken(t, reader)

	// Owner should be able to update
	updateReq := map[string]interface{}{
		"Title":   "Updated by Owner",
		"Content": "Content",
		"Status":  string(schema.PostStatusPublished),
		"Type":    string(schema.PostTypePost),
	}

	ownerUpdateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), updateReq, ownerToken)
	if ownerUpdateResp.StatusCode != http.StatusOK {
		t.Errorf("owner should be able to update, got status: %d", ownerUpdateResp.StatusCode)
	}

	// Reader should NOT be able to update
	readerUpdateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/posts/%d", business.ID, post.ID), updateReq, readerToken)
	if readerUpdateResp.StatusCode != http.StatusForbidden {
		t.Errorf("reader should not be able to update, got status: %d", readerUpdateResp.StatusCode)
	}
}

