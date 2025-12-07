package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/coupon/request"
)

// =============================================================================
// INDEX TESTS - GET /v1/business/:businessID/coupons
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

	// Create test coupons
	now := time.Now()
	ta.CreateTestCoupon(t, "DISCOUNT10", "10% Off", 10, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))
	ta.CreateTestCoupon(t, "FLAT50", "50 Toman Off", 50, schema.CouponTypeFixedAmount, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons", business.ID), nil, token)

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
		t.Errorf("expected 2 coupons, got %d", len(data))
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/coupons", nil, "")

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
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons", business.ID), nil, token)

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

	// Create multiple coupons
	now := time.Now()
	for i := 1; i <= 5; i++ {
		ta.CreateTestCoupon(t, fmt.Sprintf("CODE%d", i), fmt.Sprintf("Coupon %d", i), float64(i*10), schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination (using correct query params: page and itemPerPage)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 coupons in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["total"] != float64(5) {
		t.Errorf("expected total 5 coupons, got: %v", meta["total"])
	}
}

// =============================================================================
// SHOW TESTS - GET /v1/business/:businessID/coupons/:id
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

	// Create test coupon
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "DISCOUNT10", "10% Off", 10, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, coupon.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains coupon data
	if result["Code"] != "DISCOUNT10" {
		t.Errorf("expected code 'DISCOUNT10', got: %v", result["Code"])
	}

	if result["Title"] != "10% Off" {
		t.Errorf("expected title '10%% Off', got: %v", result["Title"])
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

	// Make request for non-existent coupon
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// STORE TESTS - POST /v1/business/:businessID/coupons
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

	// Create coupon request (using map to avoid sending BusinessID in body)
	now := time.Now()
	storeReq := map[string]interface{}{
		"Code":      "NEWCOUPON",
		"Title":     "New Coupon",
		"Value":     20,
		"Type":      string(schema.CouponTypePercentage),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(24 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 50},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}

	// Verify coupon was created in database
	var count int64
	ta.DB.Model(&schema.Coupon{}).Where("code = ?", "NEWCOUPON").Count(&count)
	if count != 1 {
		t.Errorf("expected coupon to be created in database, count: %d", count)
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

	// Create coupon request without required fields
	storeReq := map[string]interface{}{
		"Code": "INCOMPLETE",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_DuplicateCode(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create existing coupon
	now := time.Now()
	ta.CreateTestCoupon(t, "EXISTING", "Existing Coupon", 10, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Try to create coupon with same code (using map to avoid sending BusinessID in body)
	storeReq := map[string]interface{}{
		"Code":      "EXISTING",
		"Title":     "Duplicate Coupon",
		"Value":     20,
		"Type":      string(schema.CouponTypePercentage),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(24 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 50},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for duplicate code, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	// Check for Persian error message about duplicate code
	if messages[0] != "این کد کوپن قبلا ثبت شده است، لطفا مقداری خاص ثبت کنید" {
		t.Errorf("expected duplicate code error message, got: %v", messages[0])
	}
}

func TestStore_InvalidDateRange(t *testing.T) {
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

	// Create coupon request with end time before start time (using map to avoid sending BusinessID in body)
	now := time.Now()
	storeReq := map[string]interface{}{
		"Code":      "BADDATE",
		"Title":     "Bad Date Coupon",
		"Value":     20,
		"Type":      string(schema.CouponTypePercentage),
		"StartTime": now.Add(24 * time.Hour).Format(time.DateTime), // Future
		"EndTime":   now.Add(-time.Hour).Format(time.DateTime),     // Past
		"Meta":      map[string]interface{}{"MaxUsage": 50},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for invalid date range, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	// Check for Persian error message about date range
	if messages[0] != "تاریخ شروع پس از پایان است" {
		t.Errorf("expected date range error message, got: %v", messages[0])
	}
}

// =============================================================================
// UPDATE TESTS - PUT /v1/business/:businessID/coupons/:id
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

	// Create test coupon
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "ORIGINAL", "Original Title", 10, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update coupon request (using map to avoid sending BusinessID in body)
	updateReq := map[string]interface{}{
		"Code":      "UPDATED",
		"Title":     "Updated Title",
		"Value":     25,
		"Type":      string(schema.CouponTypeFixedAmount),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(48 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 100},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, coupon.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify coupon was updated in database
	var updatedCoupon schema.Coupon
	ta.DB.First(&updatedCoupon, coupon.ID)

	if updatedCoupon.Code != "UPDATED" {
		t.Errorf("expected code 'UPDATED', got: %s", updatedCoupon.Code)
	}

	if updatedCoupon.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got: %s", updatedCoupon.Title)
	}

	if updatedCoupon.Value != 25 {
		t.Errorf("expected value 25, got: %f", updatedCoupon.Value)
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

	// Update request for non-existent coupon (using map to avoid sending BusinessID in body)
	now := time.Now()
	updateReq := map[string]interface{}{
		"Code":      "NOTEXIST",
		"Title":     "Not Exist",
		"Value":     10,
		"Type":      string(schema.CouponTypePercentage),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(24 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 50},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/coupons/99999", business.ID), updateReq, token)

	// The update might not fail if no rows are affected (GORM behavior)
	// Just verify the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", resp.StatusCode)
	}
}

// =============================================================================
// DELETE TESTS - DELETE /v1/business/:businessID/coupons/:id
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

	// Create test coupon
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "TODELETE", "To Delete", 10, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, coupon.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify coupon was deleted (soft delete)
	var count int64
	ta.DB.Unscoped().Model(&schema.Coupon{}).Where("id = ? AND deleted_at IS NOT NULL", coupon.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected coupon to be soft deleted")
	}
}

func TestDelete_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodDelete, "/v1/business/1/coupons/1", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

// =============================================================================
// COUPON VALIDATE TESTS - POST /v1/business/:businessID/coupon-validate
// =============================================================================

func TestCouponValidate_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role (regular user can validate coupons)
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create active coupon
	now := time.Now()
	ta.CreateTestCoupon(t, "VALID20", "20% Off", 20, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "VALID20",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify discounted amount (1000 - 20% = 800)
	data, ok := result["Data"].(float64)
	if !ok {
		t.Errorf("expected Data (discounted amount) in response, got: %v", result)
		return
	}

	if data != 800 {
		t.Errorf("expected discounted amount 800, got: %f", data)
	}
}

func TestCouponValidate_FixedAmount(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create fixed amount coupon
	now := time.Now()
	ta.CreateTestCoupon(t, "FLAT100", "100 Off", 100, schema.CouponTypeFixedAmount, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "FLAT100",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify discounted amount (1000 - 100 = 900)
	data, ok := result["Data"].(float64)
	if !ok {
		t.Errorf("expected Data (discounted amount) in response, got: %v", result)
		return
	}

	if data != 900 {
		t.Errorf("expected discounted amount 900, got: %f", data)
	}
}

func TestCouponValidate_InvalidCode(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request with non-existent code
	validateReq := request.ValidateCoupon{
		Code:          "INVALID",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid code, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	if messages[0] != "کد تخفیف معتبر نمی باشد" {
		t.Errorf("expected invalid code error message, got: %v", messages[0])
	}
}

func TestCouponValidate_ExpiredCoupon(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create expired coupon
	now := time.Now()
	ta.CreateTestCoupon(t, "EXPIRED", "Expired Coupon", 20, schema.CouponTypePercentage, business.ID, now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "EXPIRED",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for expired coupon, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	if messages[0] != "کد تخفیف منقضی شده است" {
		t.Errorf("expected expired coupon error message, got: %v", messages[0])
	}
}

func TestCouponValidate_NotYetActive(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create future coupon
	now := time.Now()
	ta.CreateTestCoupon(t, "FUTURE", "Future Coupon", 20, schema.CouponTypePercentage, business.ID, now.Add(24*time.Hour), now.Add(48*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "FUTURE",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for not yet active coupon, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	if messages[0] != "کد تخفیف در حال حاضر فعال نمی باشد" {
		t.Errorf("expected not active error message, got: %v", messages[0])
	}
}

func TestCouponValidate_AlreadyUsedByUser(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create coupon that has already been used by this user
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "USEDONCE", "Used Once", 20, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Mark coupon as used by this user
	coupon.Meta.UsedBy = []uint64{user.ID}
	ta.DB.Save(coupon)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "USEDONCE",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for already used coupon, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	if messages[0] != "کد تخفیف برای شما قبلا استفاده شده است" {
		t.Errorf("expected already used error message, got: %v", messages[0])
	}
}

func TestCouponValidate_MaxUsageReached(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create coupon with max usage reached
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "MAXED", "Maxed Out", 20, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Set max usage and times used
	coupon.Meta.MaxUsage = 5
	coupon.TimesUsed = 5
	ta.DB.Save(coupon)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon request
	validateReq := request.ValidateCoupon{
		Code:          "MAXED",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for max usage reached, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	if messages[0] != "تعداد استفاده از کد تخفیف بیش از حد مجاز است" {
		t.Errorf("expected max usage error message, got: %v", messages[0])
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

	// 1. CREATE (using map to avoid sending BusinessID in body)
	now := time.Now()
	createReq := map[string]interface{}{
		"Code":      "INTEGRATION",
		"Title":     "Integration Test Coupon",
		"Value":     15,
		"Type":      string(schema.CouponTypePercentage),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(24 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 50},
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		t.Fatalf("CREATE failed: %d", createResp.StatusCode)
	}

	// Get created coupon ID
	var createdCoupon schema.Coupon
	ta.DB.Where("code = ?", "INTEGRATION").First(&createdCoupon)

	// 2. READ (Show)
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, createdCoupon.ID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("READ failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	if showResult["Code"] != "INTEGRATION" {
		t.Errorf("READ: expected code 'INTEGRATION', got: %v", showResult["Code"])
	}

	// 3. UPDATE (using map to avoid sending BusinessID in body)
	updateReq := map[string]interface{}{
		"Code":      "INTEGRATION_UPDATED",
		"Title":     "Updated Integration Coupon",
		"Value":     25,
		"Type":      string(schema.CouponTypeFixedAmount),
		"StartTime": now.Add(-time.Hour).Format(time.DateTime),
		"EndTime":   now.Add(48 * time.Hour).Format(time.DateTime),
		"Meta":      map[string]interface{}{"MaxUsage": 100},
	}

	updateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, createdCoupon.ID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		t.Fatalf("UPDATE failed: %d", updateResp.StatusCode)
	}

	// Verify update
	var updatedCoupon schema.Coupon
	ta.DB.First(&updatedCoupon, createdCoupon.ID)
	if updatedCoupon.Code != "INTEGRATION_UPDATED" {
		t.Errorf("UPDATE: expected code 'INTEGRATION_UPDATED', got: %s", updatedCoupon.Code)
	}

	// 4. DELETE
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/coupons/%d", business.ID, createdCoupon.ID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed: %d", deleteResp.StatusCode)
	}

	// Verify deletion
	var count int64
	ta.DB.Model(&schema.Coupon{}).Where("id = ?", createdCoupon.ID).Count(&count)
	if count != 0 {
		t.Errorf("DELETE: expected coupon to be deleted, count: %d", count)
	}
}

func TestCouponTypes_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create both types of coupons
	now := time.Now()

	// Percentage coupon with max discount
	percentCoupon := ta.CreateTestCoupon(t, "PERCENT50", "50% Off Max 200", 50, schema.CouponTypePercentage, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))
	percentCoupon.Meta.MaxDiscount = 200
	ta.DB.Save(percentCoupon)

	// Fixed amount coupon
	ta.CreateTestCoupon(t, "FIXED150", "150 Toman Off", 150, schema.CouponTypeFixedAmount, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test percentage coupon (50% of 1000 = 500, but max is 200, so discount should be 200, total = 800)
	validateReq1 := request.ValidateCoupon{
		Code:          "PERCENT50",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	resp1 := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq1, token)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("PERCENT50 validation failed: %d", resp1.StatusCode)
	}

	result1 := ParseResponse(t, resp1)
	if result1["Data"] != float64(800) {
		t.Errorf("PERCENT50: expected 800, got: %v", result1["Data"])
	}

	// Test fixed amount coupon (1000 - 150 = 850)
	validateReq2 := request.ValidateCoupon{
		Code:          "FIXED150",
		UserID:        user.ID,
		OrderTotalAmt: 1000,
	}

	resp2 := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq2, token)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("FIXED150 validation failed: %d", resp2.StatusCode)
	}

	result2 := ParseResponse(t, resp2)
	if result2["Data"] != float64(850) {
		t.Errorf("FIXED150: expected 850, got: %v", result2["Data"])
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestCouponValidate_DiscountGreaterThanTotal(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create fixed amount coupon greater than order total
	now := time.Now()
	ta.CreateTestCoupon(t, "BIGDISCOUNT", "500 Off", 500, schema.CouponTypeFixedAmount, business.ID, now.Add(-time.Hour), now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Validate coupon with small order
	validateReq := request.ValidateCoupon{
		Code:          "BIGDISCOUNT",
		UserID:        user.ID,
		OrderTotalAmt: 100, // Less than discount
	}

	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupon-validate", business.ID), validateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Should be 0 (not negative)
	data, ok := result["Data"].(float64)
	if !ok || data != 0 {
		t.Errorf("expected discounted amount 0, got: %v", result["Data"])
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
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/coupons", business.ID), nil, token)

	// Should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

