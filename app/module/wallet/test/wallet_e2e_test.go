package test

import (
	"net/http"
	"strconv"
	"testing"

	"go-fiber-starter/app/module/wallet/response"
)

func TestShowWallet_UserWallet_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	testMobile := uint64(9123456789)
	testPassword := "testPassword123"
	user := ta.CreateTestUser(t, testMobile, testPassword, "Test", "User")

	// Create wallet for user
	amount := float64(1000)
	ta.CreateTestWallet(t, &user.ID, nil, amount)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID="+uintToString(user.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	var wallet response.Wallet
	ParseResponseTo(t, resp, &wallet)

	if wallet.Amount != amount {
		t.Errorf("expected wallet amount %f, got %f", amount, wallet.Amount)
	}

	if wallet.UserID == nil || *wallet.UserID != user.ID {
		t.Errorf("expected wallet user ID %d, got %v", user.ID, wallet.UserID)
	}
}

func TestShowWallet_CreateWalletIfNotExists(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user without wallet
	testMobile := uint64(9123456789)
	testPassword := "testPassword123"
	user := ta.CreateTestUser(t, testMobile, testPassword, "Test", "User")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request - should create wallet
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID="+uintToString(user.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	var wallet response.Wallet
	ParseResponseTo(t, resp, &wallet)

	// New wallet should have 0 amount
	if wallet.Amount != 0 {
		t.Errorf("expected new wallet amount 0, got %f", wallet.Amount)
	}

	if wallet.UserID == nil || *wallet.UserID != user.ID {
		t.Errorf("expected wallet user ID %d, got %v", user.ID, wallet.UserID)
	}
}

func TestShowWallet_BusinessWallet_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user who will own the business
	testMobile := uint64(9123456789)
	testPassword := "testPassword123"
	owner := ta.CreateTestUser(t, testMobile, testPassword, "Business", "Owner")

	// Create business
	business := ta.CreateTestBusiness(t, "Test Business", owner.ID)

	// Create user with business owner permissions
	businessOwner := ta.CreateTestUserWithBusinessOwner(t, 9123456780, testPassword, "Owner", "User", business.ID)

	// Create wallet for business
	amount := float64(5000)
	ta.CreateTestWallet(t, nil, &business.ID, amount)

	// Generate token for business owner
	token := ta.GenerateTestToken(t, businessOwner)

	// Make request
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?BusinessID="+uintToString(business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	var wallet response.Wallet
	ParseResponseTo(t, resp, &wallet)

	if wallet.Amount != amount {
		t.Errorf("expected wallet amount %f, got %f", amount, wallet.Amount)
	}

	if wallet.BusinessID == nil || *wallet.BusinessID != business.ID {
		t.Errorf("expected wallet business ID %d, got %v", business.ID, wallet.BusinessID)
	}
}

func TestShowWallet_Unauthorized_NoToken(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID=1", nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestShowWallet_Unauthorized_InvalidToken(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request with invalid token
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID=1", nil, "invalid-token")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestShowWallet_Forbidden_OtherUserWallet(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create two users
	user1 := ta.CreateTestUser(t, 9123456789, "password123", "User", "One")
	user2 := ta.CreateTestUser(t, 9123456780, "password123", "User", "Two")

	// Create wallet for user2
	ta.CreateTestWallet(t, &user2.ID, nil, 1000)

	// Generate token for user1
	token := ta.GenerateTestToken(t, user1)

	// User1 tries to access user2's wallet
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID="+uintToString(user2.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestShowWallet_Forbidden_NotBusinessOwner(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user who will own the business
	ownerMobile := uint64(9123456789)
	owner := ta.CreateTestUser(t, ownerMobile, "password123", "Business", "Owner")

	// Create business
	business := ta.CreateTestBusiness(t, "Test Business", owner.ID)

	// Create regular user (not business owner)
	regularUser := ta.CreateTestUser(t, 9123456780, "password123", "Regular", "User")

	// Create wallet for business
	ta.CreateTestWallet(t, nil, &business.ID, 5000)

	// Generate token for regular user
	token := ta.GenerateTestToken(t, regularUser)

	// Regular user tries to access business wallet
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?BusinessID="+uintToString(business.ID), nil, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestShowWallet_BadRequest_NoIDProvided(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "password123", "Test", "User")

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request without UserID or BusinessID
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet", nil, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	// Check for Persian error message
	if messages[0] != "id ارسال نشده است" {
		t.Errorf("expected 'id ارسال نشده است' error message, got: %v", messages[0])
	}
}

func TestShowWallet_MultipleWallets_UserAndBusiness(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user who owns a business
	user := ta.CreateTestUser(t, 9123456789, "password123", "Test", "User")

	// Create business owned by user
	business := ta.CreateTestBusiness(t, "Test Business", user.ID)

	// Update user with business owner permissions
	userWithPerm := ta.CreateTestUserWithBusinessOwner(t, 9123456780, "password123", "Owner", "User", business.ID)

	// Create user wallet
	userWalletAmount := float64(1000)
	ta.CreateTestWallet(t, &userWithPerm.ID, nil, userWalletAmount)

	// Create business wallet
	businessWalletAmount := float64(5000)
	ta.CreateTestWallet(t, nil, &business.ID, businessWalletAmount)

	// Generate token
	token := ta.GenerateTestToken(t, userWithPerm)

	// Get user wallet
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?UserID="+uintToString(userWithPerm.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200 for user wallet, got %d, response: %v", resp.StatusCode, result)
		return
	}

	var userWallet response.Wallet
	ParseResponseTo(t, resp, &userWallet)

	if userWallet.Amount != userWalletAmount {
		t.Errorf("expected user wallet amount %f, got %f", userWalletAmount, userWallet.Amount)
	}

	// Get business wallet
	resp2 := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?BusinessID="+uintToString(business.ID), nil, token)

	if resp2.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp2)
		t.Errorf("expected status 200 for business wallet, got %d, response: %v", resp2.StatusCode, result)
		return
	}

	var businessWallet response.Wallet
	ParseResponseTo(t, resp2, &businessWallet)

	if businessWallet.Amount != businessWalletAmount {
		t.Errorf("expected business wallet amount %f, got %f", businessWalletAmount, businessWallet.Amount)
	}
}

func TestShowWallet_CreateBusinessWalletIfNotExists(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user who owns a business
	user := ta.CreateTestUser(t, 9123456789, "password123", "Test", "User")

	// Create business
	business := ta.CreateTestBusiness(t, "Test Business", user.ID)

	// Create user with business owner permissions
	businessOwner := ta.CreateTestUserWithBusinessOwner(t, 9123456780, "password123", "Owner", "User", business.ID)

	// Generate token
	token := ta.GenerateTestToken(t, businessOwner)

	// Make request - should create business wallet
	resp := ta.MakeAuthenticatedRequest(t, http.MethodGet, "/v1/wallets/wallet?BusinessID="+uintToString(business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	var wallet response.Wallet
	ParseResponseTo(t, resp, &wallet)

	// New wallet should have 0 amount
	if wallet.Amount != 0 {
		t.Errorf("expected new wallet amount 0, got %f", wallet.Amount)
	}

	if wallet.BusinessID == nil || *wallet.BusinessID != business.ID {
		t.Errorf("expected wallet business ID %d, got %v", business.ID, wallet.BusinessID)
	}
}

// Helper function to convert uint64 to string
func uintToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}

