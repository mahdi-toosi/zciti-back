package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
)

// =============================================================================
// INDEX TESTS - GET /v1/wallets/:walletID/transactions
// =============================================================================

func TestIndex_Success_UserWallet(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create test transactions
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Test transaction 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Test transaction 2")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 50, schema.TransactionStatusPending, schema.OrderPaymentMethodOnline, "Test transaction 3")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

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

	if len(data) != 3 {
		t.Errorf("expected 3 transactions, got %d", len(data))
	}
}

func TestIndex_Success_BusinessWallet(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create wallet for the business
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 5000)

	// Create test transactions
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 500, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Business transaction 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 1000, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Business transaction 2")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

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
		t.Errorf("expected 2 transactions, got %d", len(data))
	}
}

func TestIndex_Success_BusinessWallet_AsObserver(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner user
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Create observer user with business observer role
	observer := ta.CreateTestUser(t, 9987654321, "testPassword123", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver})

	// Create wallet for the business
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 5000)

	// Create test transactions
	ta.CreateTestTransaction(t, wallet.ID, owner.ID, 500, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Business transaction 1")
	ta.CreateTestTransaction(t, wallet.ID, owner.ID, 1000, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Business transaction 2")

	// Generate token for observer
	token := ta.GenerateTestToken(t, observer)

	// Make request - business observer should be able to view business wallet transactions
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200 for business observer, got %d, response: %v", resp.StatusCode, result)
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
		t.Errorf("expected 2 transactions, got %d", len(data))
	}
}

func TestIndex_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/wallets/1/transactions", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestIndex_Forbidden_NotOwner(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create first user (wallet owner)
	ownerUser := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)

	// Create wallet for the owner
	wallet := ta.CreateTestWallet(t, &ownerUser.ID, nil, 1000)

	// Create second user (trying to access wallet)
	otherUser := ta.CreateTestUser(t, 9987654321, "testPassword456", "Other", "User", 0, nil)

	// Generate token for other user
	token := ta.GenerateTestToken(t, otherUser)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for forbidden, got %d", resp.StatusCode)
	}
}

func TestIndex_Forbidden_NotBusinessOwner(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	ownerUser := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, ownerUser.ID)

	// Update owner user with business permissions
	ownerUser.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(ownerUser)

	// Create wallet for the business
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 5000)

	// Create another user (not business owner)
	otherUser := ta.CreateTestUser(t, 9987654321, "testPassword456", "Other", "User", 0, nil)

	// Generate token for other user
	token := ta.GenerateTestToken(t, otherUser)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for forbidden, got %d", resp.StatusCode)
	}
}

func TestIndex_WithPagination(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create multiple transactions
	for i := 1; i <= 5; i++ {
		ta.CreateTestTransaction(t, wallet.ID, user.ID, float64(i*100), schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, fmt.Sprintf("Transaction %d", i))
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?page=1&itemPerPage=2", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data and meta
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 transactions in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	// Check TotalAmount is present
	if _, ok := meta["TotalAmount"]; !ok {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
	}

	// Check pagination info
	metaInner, ok := meta["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta.Meta in response, got: %v", meta)
		return
	}

	if metaInner["total"] != float64(5) {
		t.Errorf("expected total 5 transactions, got: %v", metaInner["total"])
	}
}

func TestIndex_TotalAmount_OnlySuccessful(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create transactions with different statuses
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Successful 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Successful 2")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 50, schema.TransactionStatusPending, schema.OrderPaymentMethodOnline, "Pending")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 75, schema.TransactionStatusFailed, schema.OrderPaymentMethodOnline, "Failed")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify total amount only includes successful transactions (100 + 200 = 300)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
		return
	}

	if totalAmount != 300 {
		t.Errorf("expected total amount 300 (only successful), got: %f", totalAmount)
	}
}

func TestIndex_WalletNotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request for non-existent wallet
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/wallets/99999/transactions", nil, token)

	// Should return error (wallet not found)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for wallet not found, got %d", resp.StatusCode)
	}
}

func TestIndex_EmptyTransactions(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user (no transactions)
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 0)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains empty data or null
	data := result["Data"]
	if data != nil {
		dataArr, ok := data.([]interface{})
		if ok && len(dataArr) != 0 {
			t.Errorf("expected empty transactions, got: %v", dataArr)
		}
	}
}

// =============================================================================
// DATE FILTER TESTS
// =============================================================================

func TestIndex_WithDateFilters(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create transactions at different times
	now := time.Now()

	// Transaction 1: yesterday
	tx1 := ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Yesterday")
	ta.DB.Model(&tx1).Update("created_at", now.Add(-24*time.Hour))

	// Transaction 2: today
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Today")

	// Transaction 3: tomorrow (future, but for testing date filter)
	tx3 := ta.CreateTestTransaction(t, wallet.ID, user.ID, 300, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Tomorrow")
	ta.DB.Model(&tx3).Update("created_at", now.Add(24*time.Hour))

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with date filter (only today and yesterday)
	startTime := now.Add(-48 * time.Hour).Format("2006-01-02")
	endTime := now.Format("2006-01-02")

	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?StartTime=%s&EndTime=%s", wallet.ID, startTime, endTime), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains filtered data
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	// Should have 2 transactions (yesterday and today, not tomorrow)
	if len(data) != 2 {
		t.Errorf("expected 2 transactions with date filter, got %d", len(data))
	}
}

// =============================================================================
// TRANSACTION STATUS TESTS
// =============================================================================

func TestIndex_TransactionStatuses(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create transactions with all different statuses
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusPending, schema.OrderPaymentMethodOnline, "Pending transaction")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Success transaction")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 300, schema.TransactionStatusFailed, schema.OrderPaymentMethodOnline, "Failed transaction")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 400, schema.TransactionStatusRefunded, schema.OrderPaymentMethodOnline, "Refunded transaction")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 500, schema.TransactionStatusCancelled, schema.OrderPaymentMethodOnline, "Cancelled transaction")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains all transactions
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 5 {
		t.Errorf("expected 5 transactions (all statuses), got %d", len(data))
	}

	// Verify total amount only includes successful (200)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
		return
	}

	if totalAmount != 200 {
		t.Errorf("expected total amount 200 (only successful), got: %f", totalAmount)
	}
}

// =============================================================================
// PAYMENT METHOD TESTS
// =============================================================================

func TestIndex_PaymentMethods(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create transactions with different payment methods
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Online payment")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodCash, "Wallet payment")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 300, schema.TransactionStatusSuccess, schema.OrderPaymentMethodCashOnDelivery, "Offline payment")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains all transactions
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 3 {
		t.Errorf("expected 3 transactions (all payment methods), got %d", len(data))
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestIntegration_UserWalletTransactions(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 0)

	// Simulate multiple transactions over time
	amounts := []float64{100, 250, 75, 500, 125}
	for i, amount := range amounts {
		status := schema.TransactionStatusSuccess
		if i == 2 { // Make one failed
			status = schema.TransactionStatusFailed
		}
		ta.CreateTestTransaction(t, wallet.ID, user.ID, amount, status, schema.OrderPaymentMethodOnline, fmt.Sprintf("Transaction %d", i+1))
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// 1. Get all transactions
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get transactions: %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)

	// Verify all transactions are returned
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 5 {
		t.Errorf("expected 5 transactions, got: %v", len(data))
	}

	// Verify total amount (100 + 250 + 500 + 125 = 975, excluding failed 75)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Meta in response")
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok || totalAmount != 975 {
		t.Errorf("expected total amount 975, got: %f", totalAmount)
	}

	// 2. Test pagination
	paginatedResp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?page=1&itemPerPage=3", wallet.ID), nil, token)

	if paginatedResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get paginated transactions: %d", paginatedResp.StatusCode)
	}

	paginatedResult := ParseResponse(t, paginatedResp)
	paginatedData, ok := paginatedResult["Data"].([]interface{})
	if !ok || len(paginatedData) != 3 {
		t.Errorf("expected 3 paginated transactions, got: %v", len(paginatedData))
	}
}

func TestIntegration_BusinessWalletTransactions(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)

	// Create business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create transactions
	ta.CreateTestTransaction(t, wallet.ID, owner.ID, 1000, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Business sale 1")
	ta.CreateTestTransaction(t, wallet.ID, owner.ID, 2500, schema.TransactionStatusSuccess, schema.OrderPaymentMethodCashOnDelivery, "Business sale 2")
	ta.CreateTestTransaction(t, wallet.ID, owner.ID, 500, schema.TransactionStatusPending, schema.OrderPaymentMethodOnline, "Pending sale")

	// Generate token for owner
	token := ta.GenerateTestToken(t, owner)

	// Get transactions
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify transactions
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 3 {
		t.Errorf("expected 3 transactions, got: %v", len(data))
	}

	// Verify total amount (only successful: 1000 + 2500 = 3500)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Meta in response")
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok || totalAmount != 3500 {
		t.Errorf("expected total amount 3500, got: %f", totalAmount)
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestIndex_InvalidWalletID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with invalid wallet ID
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/wallets/invalid/transactions", nil, token)

	// Should return bad request or not found
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 400, 404, or 500 for invalid wallet ID, got %d", resp.StatusCode)
	}
}

func TestIndex_LargeAmountTransactions(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 0)

	// Create transactions with large amounts
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 1000000000, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Large transaction 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 2000000000, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Large transaction 2")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify total amount handles large numbers
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Meta in response")
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok {
		t.Errorf("expected TotalAmount in Meta, got: %v", meta)
		return
	}

	if totalAmount != 3000000000 {
		t.Errorf("expected total amount 3000000000, got: %f", totalAmount)
	}
}

func TestIndex_DecimalAmountTransactions(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 0)

	// Create transactions with decimal amounts
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100.50, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Decimal transaction 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200.75, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Decimal transaction 2")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 transactions, got: %v", result["Data"])
	}
}

func TestIndex_ZeroAmountTransaction(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 0)

	// Create transaction with zero amount
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 0, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Zero amount transaction")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("expected 1 transaction, got: %v", result["Data"])
	}
}

// =============================================================================
// CONCURRENT ACCESS TESTS
// =============================================================================

func TestIndex_MultipleUsersAccessOwnWallets(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create multiple users with their own wallets
	user1 := ta.CreateTestUser(t, 9123456781, "testPassword123", "User", "One", 0, nil)
	user2 := ta.CreateTestUser(t, 9123456782, "testPassword123", "User", "Two", 0, nil)

	wallet1 := ta.CreateTestWallet(t, &user1.ID, nil, 1000)
	wallet2 := ta.CreateTestWallet(t, &user2.ID, nil, 2000)

	ta.CreateTestTransaction(t, wallet1.ID, user1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "User1 transaction")
	ta.CreateTestTransaction(t, wallet2.ID, user2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "User2 transaction")

	// User1 accesses their wallet
	token1 := ta.GenerateTestToken(t, user1)
	resp1 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet1.ID), nil, token1)

	if resp1.StatusCode != http.StatusOK {
		t.Errorf("User1 should access their own wallet, got status: %d", resp1.StatusCode)
	}

	// User2 accesses their wallet
	token2 := ta.GenerateTestToken(t, user2)
	resp2 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet2.ID), nil, token2)

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("User2 should access their own wallet, got status: %d", resp2.StatusCode)
	}

	// User1 tries to access User2's wallet (should fail)
	resp3 := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet2.ID), nil, token1)

	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("User1 should NOT access User2's wallet, got status: %d", resp3.StatusCode)
	}
}

// =============================================================================
// RESPONSE FORMAT TESTS
// =============================================================================

func TestIndex_ResponseFormat(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create a transaction
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 500, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Test transaction")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)

	// Verify response structure
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Fatalf("expected Data array in response")
	}

	if len(data) != 1 {
		t.Fatalf("expected 1 transaction")
	}

	// Verify transaction fields
	tx, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected transaction object")
	}

	// Check required fields exist
	requiredFields := []string{"ID", "Amount", "Status", "UpdatedAt", "Description", "OrderPaymentMethod"}
	for _, field := range requiredFields {
		if _, exists := tx[field]; !exists {
			t.Errorf("expected field %s in transaction response", field)
		}
	}

	// Verify Meta structure
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Meta in response")
	}

	if _, exists := meta["TotalAmount"]; !exists {
		t.Errorf("expected TotalAmount in Meta")
	}

	if _, exists := meta["Meta"]; !exists {
		t.Errorf("expected Meta.Meta (pagination) in Meta")
	}
}

// =============================================================================
// TAXONOMY FILTER TESTS
// =============================================================================

func TestIndex_FilterByTaxonomies(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create taxonomies (city -> workspace -> dormitory)
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 200, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 200)

	// Create transactions linked to orders
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")

	// Generate token
	token := ta.GenerateTestToken(t, owner)

	// Filter by CityID - should only get transactions for city1
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?CityID=%d", wallet.ID, city1.ID), nil, token)

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
		t.Errorf("expected 1 transaction for city1 filter, got %d", len(data))
	}

	// Verify total amount only includes filtered transactions
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response")
		return
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok || totalAmount != 100 {
		t.Errorf("expected total amount 100 for city1 filter, got: %v", totalAmount)
	}
}

func TestIndex_FilterByDormitoryID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create hierarchical taxonomies
	city := ta.CreateTestTaxonomy(t, "City", schema.TaxonomyTypeCategory, business.ID, nil)
	workspace := ta.CreateTestTaxonomy(t, "Workspace", schema.TaxonomyTypeCategory, business.ID, &city.ID)
	dormitory := ta.CreateTestTaxonomy(t, "Dormitory", schema.TaxonomyTypeCategory, business.ID, &workspace.ID)

	// Create posts
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies - post1 to dormitory, post2 to workspace
	ta.AttachTaxonomyToPost(t, post1.ID, dormitory.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, workspace.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 200, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 200)

	// Create transactions linked to orders
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")

	// Generate token
	token := ta.GenerateTestToken(t, owner)

	// Filter by DormitoryID - should only get dormitory transactions
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?DormitoryID=%d", wallet.ID, dormitory.ID), nil, token)

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
		t.Errorf("expected 1 transaction for dormitory filter, got %d", len(data))
	}
}

// =============================================================================
// BUSINESS OBSERVER TESTS
// =============================================================================

func TestIndex_BusinessObserver_WithMeta(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create posts
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 200, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 200)

	// Create transactions linked to orders
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	// Generate token for observer
	token := ta.GenerateTestToken(t, observer)

	// Make request - observer should only see transactions for city1
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

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

	// Observer should only see transactions related to city1
	if len(data) != 1 {
		t.Errorf("expected 1 transaction for observer with city1 access, got %d", len(data))
	}

	// Verify total amount only includes observer's accessible transactions
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response")
		return
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok || totalAmount != 100 {
		t.Errorf("expected total amount 100 for observer, got: %v", totalAmount)
	}
}

func TestIndex_BusinessObserver_WithoutMeta(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create observer without meta (should be forbidden)
	observer := ta.CreateTestUser(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver})

	// Generate token for observer
	token := ta.GenerateTestToken(t, observer)

	// Make request - observer without meta should get forbidden
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions", wallet.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for observer without meta, got %d", resp.StatusCode)
	}
}

func TestIndex_BusinessObserver_FilterWithTaxonomyAccess(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create posts
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 200, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 200)

	// Create transactions linked to orders
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	// Generate token for observer
	token := ta.GenerateTestToken(t, observer)

	// Observer requests with CityID filter for city1 (which they have access to)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?CityID=%d", wallet.ID, city1.ID), nil, token)

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

	// Should only see city1 transaction
	if len(data) != 1 {
		t.Errorf("expected 1 transaction for observer with city1 filter, got %d", len(data))
	}
}

func TestIndex_BusinessObserver_FilterWithNoTaxonomyAccess(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create posts
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 200, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 200)

	// Create transactions linked to orders
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	// Generate token for observer
	token := ta.GenerateTestToken(t, observer)

	// Observer requests with CityID filter for city2 (which they DON'T have access to)
	// Should return empty results because city2 is not in observer's access list
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?CityID=%d", wallet.ID, city2.ID), nil, token)

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

	// Observer requested city2 but only has access to city1, so should get only city1 results
	// Based on the controller logic, if the requested taxonomy is not in observer's list,
	// it falls back to showing all taxonomies the observer has access to
	if len(data) != 1 {
		t.Errorf("expected 1 transaction (fallback to observer's accessible taxonomies), got %d", len(data))
	}
}

func TestIndex_TotalAmount_WithTaxonomyFilter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create business wallet
	wallet := ta.CreateTestWallet(t, nil, &business.ID, 10000)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create posts
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)
	post3 := ta.CreateTestPost(t, "Product 3", schema.PostTypeProduct, business.ID, owner.ID)

	// Attach taxonomies
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city1.ID) // Also city1
	ta.AttachTaxonomyToPost(t, post3.ID, city2.ID)

	// Create orders
	order1 := ta.CreateTestOrder(t, owner.ID, business.ID, 100, schema.OrderStatusCompleted)
	order2 := ta.CreateTestOrder(t, owner.ID, business.ID, 250, schema.OrderStatusCompleted)
	order3 := ta.CreateTestOrder(t, owner.ID, business.ID, 500, schema.OrderStatusCompleted)

	// Create order items linked to posts
	ta.CreateTestOrderItem(t, order1.ID, post1.ID, 1, 100)
	ta.CreateTestOrderItem(t, order2.ID, post2.ID, 1, 250)
	ta.CreateTestOrderItem(t, order3.ID, post3.ID, 1, 500)

	// Create transactions - 2 for city1, 1 for city2
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order1.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 1")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order2.ID, 250, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 2")
	ta.CreateTestTransactionWithOrder(t, wallet.ID, owner.ID, order3.ID, 500, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Transaction 3")

	// Generate token
	token := ta.GenerateTestToken(t, owner)

	// Filter by city1 - should get total amount of 350 (100 + 250)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?CityID=%d", wallet.ID, city1.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify total amount for city1 filter
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response")
		return
	}

	totalAmount, ok := meta["TotalAmount"].(float64)
	if !ok || totalAmount != 350 {
		t.Errorf("expected total amount 350 for city1 filter, got: %v", totalAmount)
	}

	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 transactions for city1 filter, got %d", len(data))
	}
}

func TestIndex_StatusFilter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create wallet for the user
	wallet := ta.CreateTestWallet(t, &user.ID, nil, 1000)

	// Create transactions with different statuses
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 100, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Success 1")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 200, schema.TransactionStatusSuccess, schema.OrderPaymentMethodOnline, "Success 2")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 50, schema.TransactionStatusPending, schema.OrderPaymentMethodOnline, "Pending")
	ta.CreateTestTransaction(t, wallet.ID, user.ID, 75, schema.TransactionStatusFailed, schema.OrderPaymentMethodOnline, "Failed")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Filter by status=success
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/wallets/%d/transactions?Status=success", wallet.ID), nil, token)

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
		t.Errorf("expected 2 successful transactions, got %d", len(data))
	}
}

