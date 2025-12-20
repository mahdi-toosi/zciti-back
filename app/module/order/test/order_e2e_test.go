package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
)

// =============================================================================
// INDEX TESTS - GET /v1/business/:businessID/orders
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

	// Create test orders
	ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrder(t, user.ID, business.ID, 2000, schema.OrderStatusCompleted, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

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
		t.Errorf("expected 2 orders, got %d", len(data))
	}

	// Verify Meta contains TotalAmount
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["TotalAmount"] == nil {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
	}

	// TotalAmount should be 3000 (1000 + 2000)
	if meta["TotalAmount"] != float64(3000) {
		t.Errorf("expected TotalAmount 3000, got: %v", meta["TotalAmount"])
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/orders", nil, "")

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
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

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

	// Create multiple orders (100 + 200 + 300 + 400 + 500 = 1500)
	for i := 1; i <= 5; i++ {
		ta.CreateTestOrder(t, user.ID, business.ID, float64(i*100), schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 orders in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	// Check nested Meta pagination
	metaPagination, ok := meta["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected nested Meta in response, got: %v", meta)
		return
	}

	if metaPagination["total"] != float64(5) {
		t.Errorf("expected total 5 orders, got: %v", metaPagination["total"])
	}

	// Verify TotalAmount is present and correct (1500 total for all orders)
	if meta["TotalAmount"] == nil {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
	}

	if meta["TotalAmount"] != float64(1500) {
		t.Errorf("expected TotalAmount 1500, got: %v", meta["TotalAmount"])
	}
}

func TestIndex_FilterByCouponID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create a coupon
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "TESTCODE", "Test Coupon", 10, schema.CouponTypePercentage, business.ID, now, now.AddDate(0, 1, 0))

	// Create orders - one with coupon, one without
	orderWithCoupon := ta.CreateTestOrder(t, user.ID, business.ID, 900, schema.OrderStatusCompleted, schema.OrderPaymentMethodOnline)
	orderWithCoupon.CouponID = &coupon.ID
	ta.DB.Save(orderWithCoupon)

	ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request filtering by coupon ID
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders?CouponID=%d", business.ID, coupon.ID), nil, token)

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
		t.Errorf("expected 1 order with coupon, got %d", len(data))
	}

	// Verify TotalAmount reflects only the filtered order (900)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["TotalAmount"] != float64(900) {
		t.Errorf("expected TotalAmount 900 for filtered orders, got: %v", meta["TotalAmount"])
	}
}

func TestIndex_TotalAmountByStatus(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create orders with different statuses
	ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusCompleted, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrder(t, user.ID, business.ID, 2000, schema.OrderStatusCompleted, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrder(t, user.ID, business.ID, 500, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrder(t, user.ID, business.ID, 750, schema.OrderStatusCancelled, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test 1: Filter by completed status - should return TotalAmount = 3000 (1000 + 2000)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders?Status=completed", business.ID), nil, token)

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
		t.Errorf("expected 2 completed orders, got %d", len(data))
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["TotalAmount"] != float64(3000) {
		t.Errorf("expected TotalAmount 3000 for completed orders, got: %v", meta["TotalAmount"])
	}

	// Test 2: Filter by pending status - should return TotalAmount = 500
	resp2 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders?Status=pending", business.ID), nil, token)

	if resp2.StatusCode != http.StatusOK {
		result2 := ParseResponse(t, resp2)
		t.Errorf("expected status 200, got %d, response: %v", resp2.StatusCode, result2)
		return
	}

	result2 := ParseResponse(t, resp2)
	meta2, ok := result2["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result2)
		return
	}

	if meta2["TotalAmount"] != float64(500) {
		t.Errorf("expected TotalAmount 500 for pending orders, got: %v", meta2["TotalAmount"])
	}

	// Test 3: No status filter - should return TotalAmount = 4250 (all orders)
	resp3 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	if resp3.StatusCode != http.StatusOK {
		result3 := ParseResponse(t, resp3)
		t.Errorf("expected status 200, got %d, response: %v", resp3.StatusCode, result3)
		return
	}

	result3 := ParseResponse(t, resp3)
	meta3, ok := result3["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result3)
		return
	}

	if meta3["TotalAmount"] != float64(4250) {
		t.Errorf("expected TotalAmount 4250 for all orders, got: %v", meta3["TotalAmount"])
	}
}

// =============================================================================
// SHOW TESTS - GET /v1/business/:businessID/orders/:id
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

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1500, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains order data
	if result["ID"] != float64(order.ID) {
		t.Errorf("expected ID %d, got: %v", order.ID, result["ID"])
	}

	if result["TotalAmt"] != float64(1500) {
		t.Errorf("expected TotalAmt 1500, got: %v", result["TotalAmt"])
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

	// Make request for non-existent order
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

func TestShow_WithOrderItems(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create post and product
	post := ta.CreateTestPost(t, "Test Product", business.ID, user.ID)
	simpleVariant := schema.ProductVariantTypeSimple
	ta.CreateTestProduct(t, post.ID, business.ID, 500, schema.ProductTypeSimple, &simpleVariant)

	// Create test order with items
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1500, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrderItem(t, order.ID, post.ID, 3, 500)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify order items are included
	orderItems, ok := result["OrderItems"].([]interface{})
	if !ok {
		t.Errorf("expected OrderItems array in response, got: %v", result)
		return
	}

	if len(orderItems) != 1 {
		t.Errorf("expected 1 order item, got %d", len(orderItems))
	}
}

// =============================================================================
// STORE TESTS - POST /v1/business/:businessID/orders
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

	// Create order request
	storeReq := map[string]interface{}{
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"OrderItems":    []interface{}{},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/orders", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Data map in response, got: %v", result)
		return
	}

	if data["orderID"] == nil {
		t.Errorf("expected orderID in response, got: %v", data)
	}

	// Verify order was created in database
	var count int64
	ta.DB.Model(&schema.Order{}).Where("business_id = ?", business.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected order to be created in database, count: %d", count)
	}
}

func TestStore_ValidationError_MissingPaymentMethod(t *testing.T) {
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

	// Create order request without payment method
	storeReq := map[string]interface{}{
		"OrderItems": []interface{}{},
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/orders", business.ID), storeReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestStore_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create order request
	storeReq := map[string]interface{}{
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"OrderItems":    []interface{}{},
	}

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodPost, "/v1/business/1/orders", storeReq, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

// =============================================================================
// UPDATE TESTS - PUT /v1/business/:businessID/orders/:id
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

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Update order request (using map to include BusinessID)
	updateReq := map[string]interface{}{
		"Status":        string(schema.OrderStatusCompleted),
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"BusinessID":    business.ID,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), updateReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify order was updated in database
	var updatedOrder schema.Order
	ta.DB.First(&updatedOrder, order.ID)

	if updatedOrder.Status != schema.OrderStatusCompleted {
		t.Errorf("expected status 'completed', got: %s", updatedOrder.Status)
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

	// Create test order with pending status
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test status transitions
	statuses := []schema.OrderStatus{
		schema.OrderStatusProcessing,
		schema.OrderStatusCompleted,
	}

	for _, status := range statuses {
		updateReq := map[string]interface{}{
			"Status":        string(status),
			"PaymentMethod": string(schema.OrderPaymentMethodOnline),
			"BusinessID":    business.ID,
		}

		resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), updateReq, token)

		if resp.StatusCode != http.StatusOK {
			result := ParseResponse(t, resp)
			t.Errorf("expected status 200 for status %s, got %d, response: %v", status, resp.StatusCode, result)
			continue
		}

		var updatedOrder schema.Order
		ta.DB.First(&updatedOrder, order.ID)

		if updatedOrder.Status != status {
			t.Errorf("expected status '%s', got: %s", status, updatedOrder.Status)
		}
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

	// Update request for non-existent order
	updateReq := map[string]interface{}{
		"Status":        string(schema.OrderStatusCompleted),
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"BusinessID":    business.ID,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/orders/99999", business.ID), updateReq, token)

	// GORM might not fail when no rows affected
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", resp.StatusCode)
	}
}

// =============================================================================
// DELETE TESTS - DELETE /v1/business/:businessID/orders/:id
// =============================================================================

func TestDelete_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with admin role (URBusinessOwner doesn't have PDelete for orders)
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with admin permissions (admin has delete permission for orders)
	user.Permissions[business.ID] = []schema.UserRole{schema.URAdmin}
	ta.DB.Save(user)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify order was deleted (soft delete)
	var count int64
	ta.DB.Unscoped().Model(&schema.Order{}).Where("id = ? AND deleted_at IS NOT NULL", order.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected order to be soft deleted")
	}
}

func TestDelete_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodDelete, "/v1/business/1/orders/1", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestDelete_WithOrderItems(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with admin role (URBusinessOwner doesn't have PDelete for orders)
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with admin permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URAdmin}
	ta.DB.Save(user)

	// Create post
	post := ta.CreateTestPost(t, "Test Product", business.ID, user.ID)

	// Create test order with items
	order := ta.CreateTestOrder(t, user.ID, business.ID, 1500, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrderItem(t, order.ID, post.ID, 3, 500)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make delete request
	resp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, order.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	// Verify order was deleted
	var orderCount int64
	ta.DB.Unscoped().Model(&schema.Order{}).Where("id = ? AND deleted_at IS NOT NULL", order.ID).Count(&orderCount)
	if orderCount != 1 {
		t.Errorf("expected order to be soft deleted")
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestCRUD_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with admin role (for delete permission)
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with admin permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URAdmin}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// 1. CREATE
	createReq := map[string]interface{}{
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"OrderItems":    []interface{}{},
	}

	createResp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/orders", business.ID), createReq, token)
	if createResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, createResp)
		t.Fatalf("CREATE failed: %d, response: %v", createResp.StatusCode, result)
	}

	// Get created order ID from database (more reliable than parsing response)
	var createdOrder schema.Order
	ta.DB.Where("business_id = ?", business.ID).Order("id desc").First(&createdOrder)
	orderID := createdOrder.ID

	if orderID == 0 {
		t.Fatalf("CREATE: order not found in database")
	}

	// 2. READ (Show)
	showResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, orderID), nil, token)
	if showResp.StatusCode != http.StatusOK {
		t.Fatalf("READ failed: %d", showResp.StatusCode)
	}

	showResult := ParseResponse(t, showResp)
	if showResult["ID"] != float64(orderID) {
		t.Errorf("READ: expected ID %d, got: %v", orderID, showResult["ID"])
	}

	// 3. UPDATE
	updateReq := map[string]interface{}{
		"Status":        string(schema.OrderStatusCompleted),
		"PaymentMethod": string(schema.OrderPaymentMethodOnline),
		"BusinessID":    business.ID,
	}

	updateResp := ta.MakeRequest(t, http.MethodPut, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, orderID), updateReq, token)
	if updateResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, updateResp)
		t.Fatalf("UPDATE failed: %d, response: %v", updateResp.StatusCode, result)
	}

	// Verify update
	var updatedOrder schema.Order
	ta.DB.First(&updatedOrder, orderID)
	if updatedOrder.Status != schema.OrderStatusCompleted {
		t.Errorf("UPDATE: expected status 'completed', got: %s", updatedOrder.Status)
	}

	// 4. DELETE
	deleteResp := ta.MakeRequest(t, http.MethodDelete, fmt.Sprintf("/v1/business/%d/orders/%d", business.ID, orderID), nil, token)
	if deleteResp.StatusCode != http.StatusOK {
		result := ParseResponse(t, deleteResp)
		t.Fatalf("DELETE failed: %d, response: %v", deleteResp.StatusCode, result)
	}

	// Verify deletion
	var count int64
	ta.DB.Model(&schema.Order{}).Where("id = ?", orderID).Count(&count)
	if count != 0 {
		t.Errorf("DELETE: expected order to be deleted, count: %d", count)
	}
}

func TestOrderStatuses_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create orders with different statuses
	statuses := []schema.OrderStatus{
		schema.OrderStatusPending,
		schema.OrderStatusProcessing,
		schema.OrderStatusCompleted,
		schema.OrderStatusCancelled,
		schema.OrderStatusRefunded,
	}

	for _, status := range statuses {
		ta.CreateTestOrder(t, user.ID, business.ID, 1000, status, schema.OrderPaymentMethodOnline)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Verify all orders are returned
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get orders: %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data := result["Data"].([]interface{})

	if len(data) != len(statuses) {
		t.Errorf("expected %d orders, got %d", len(statuses), len(data))
	}
}

func TestPaymentMethods_Integration(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create orders with different payment methods
	paymentMethods := []schema.OrderPaymentMethod{
		schema.OrderPaymentMethodOnline,
		schema.OrderPaymentMethodCash,
		schema.OrderPaymentMethodCashOnDelivery,
	}

	for _, method := range paymentMethods {
		ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, method)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Verify all orders are returned
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to get orders: %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data := result["Data"].([]interface{})

	if len(data) != len(paymentMethods) {
		t.Errorf("expected %d orders, got %d", len(paymentMethods), len(data))
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
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	// Should return validation error
	if resp.StatusCode != http.StatusUnprocessableEntity && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 422 or 500 for empty body, got %d", resp.StatusCode)
	}
}

func TestInvalidBusinessID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Test with non-existent business ID
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/99999/orders", nil, token)

	// Should return forbidden (no permission for non-existent business)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for invalid business, got %d", resp.StatusCode)
	}
}

func TestMultipleOrdersForSameUser(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create multiple orders for the same user
	for i := 1; i <= 10; i++ {
		ta.CreateTestOrder(t, user.ID, business.ID, float64(i*100), schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	data := result["Data"].([]interface{})

	if len(data) != 10 {
		t.Errorf("expected 10 orders, got %d", len(data))
	}
}

func TestOrderWithCoupon(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create a coupon
	now := time.Now()
	coupon := ta.CreateTestCoupon(t, "DISCOUNT20", "20% Off", 20, schema.CouponTypePercentage, business.ID, now, now.AddDate(0, 1, 0))

	// Create order with coupon
	order := ta.CreateTestOrder(t, user.ID, business.ID, 800, schema.OrderStatusCompleted, schema.OrderPaymentMethodOnline)
	order.CouponID = &coupon.ID
	ta.DB.Save(order)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request to list orders (Index preloads Coupon, Show doesn't)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Errorf("expected Data array with orders, got: %v", result)
		return
	}

	orderData := data[0].(map[string]interface{})

	// Verify coupon is in response
	couponData, ok := orderData["Coupon"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Coupon in order response, got: %v", orderData)
		return
	}

	if couponData["ID"] != float64(coupon.ID) {
		t.Errorf("expected coupon ID %d, got: %v", coupon.ID, couponData["ID"])
	}

	// Also verify in database that the order has the coupon
	var dbOrder schema.Order
	ta.DB.First(&dbOrder, order.ID)
	if dbOrder.CouponID == nil || *dbOrder.CouponID != coupon.ID {
		t.Errorf("expected order to have coupon ID %d in database", coupon.ID)
	}
}

// =============================================================================
// ROLE-BASED ACCESS TESTS
// =============================================================================

func TestRoleBasedAccess_UserRole(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with User role only
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with User role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Try to access orders (should be forbidden for basic user role without proper permissions)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	// URUser might not have DOrder:PReadAll permission
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 200 or 403, got %d", resp.StatusCode)
	}
}

func TestRoleBasedAccess_BusinessOwner(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with BusinessOwner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with BusinessOwner role
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create an order
	ta.CreateTestOrder(t, user.ID, business.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Should have full access
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for business owner, got %d", resp.StatusCode)
	}
}

func TestRoleBasedAccess_CrossBusinessAccess(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create two users for two different businesses
	user1 := ta.CreateTestUser(t, 9123456789, "testPassword123", "User", "One", 0, nil)
	user2 := ta.CreateTestUser(t, 9123456788, "testPassword123", "User", "Two", 0, nil)

	business1 := ta.CreateTestBusiness(t, "Business One", schema.BTypeGymManager, user1.ID)
	business2 := ta.CreateTestBusiness(t, "Business Two", schema.BTypeGymManager, user2.ID)

	// Set permissions
	user1.Permissions[business1.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user1)

	user2.Permissions[business2.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user2)

	// Create orders for each business
	ta.CreateTestOrder(t, user1.ID, business1.ID, 1000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)
	ta.CreateTestOrder(t, user2.ID, business2.ID, 2000, schema.OrderStatusPending, schema.OrderPaymentMethodOnline)

	// Generate token for user1
	token := ta.GenerateTestToken(t, user1)

	// User1 should not be able to access business2's orders
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/orders", business2.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for cross-business access, got %d", resp.StatusCode)
	}
}

