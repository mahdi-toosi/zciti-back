package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/uniwash/request"
)

// =============================================================================
// SEND COMMAND TESTS - POST /v1/business/:businessID/uni-wash/send-command
// =============================================================================

func TestSendCommand_BusinessRoute_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post and product (washing machine)
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}
}

func TestSendCommand_UserRoute_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product (washing machine)
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation for the user - starting in a few minutes
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(5*time.Minute), now.Add(65*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/user/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}
}

func TestSendCommand_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	sendReq := request.SendCommand{
		ReservationID: 1,
		ProductID:     1,
		Command:       schema.UniWashCommandON,
	}

	resp := ta.MakeRequest(t, http.MethodPost, "/v1/business/1/uni-wash/send-command", sendReq, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestSendCommand_Forbidden_NoPermission(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user without business permissions
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)

	// Create test business
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Generate token (user has no permissions for this business)
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: 1,
		ProductID:     1,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 for forbidden, got %d", resp.StatusCode)
	}
}

func TestSendCommand_ValidationError_MissingRequired(t *testing.T) {
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

	// Send command request without required fields
	sendReq := map[string]interface{}{
		"Command": "ON",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for validation error, got %d", resp.StatusCode)
	}
}

func TestSendCommand_InvalidCommand(t *testing.T) {
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

	// Send command request with invalid command
	sendReq := map[string]interface{}{
		"ReservationID": 1,
		"ProductID":     1,
		"Command":       "INVALID_COMMAND",
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422 for invalid command, got %d", resp.StatusCode)
	}
}

func TestSendCommand_MachineOff(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post and product (washing machine OFF)
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusOFF)

	// Create reservation
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for machine off, got %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	// Check for Persian error message about machine being off
	expectedMsg := "در حال حاضر دستگاه دچار مشکل شده و در دسترس نیست، در صورتی که دستگاه را رزرو کرده اید، با پشتیبانی تماس بگیرید"
	if messages[0] != expectedMsg {
		t.Errorf("expected machine off error message, got: %v", messages[0])
	}
}

func TestSendCommand_AllCommands(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post and product (washing machine)
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	commands := []schema.UniWashCommand{
		schema.UniWashCommandON,
		schema.UniWashCommandOFF,
		schema.UniWashCommandRewash,
		schema.UniWashCommandEvacuation,
	}

	for i, cmd := range commands {
		// Create reservation
		now := time.Now()
		reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

		// Send command request
		sendReq := request.SendCommand{
			ReservationID: reservation.ID,
			ProductID:     product.ID,
			Command:       cmd,
		}

		// Make request
		resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-command", business.ID), sendReq, token)

		if resp.StatusCode != http.StatusOK {
			result := ParseResponse(t, resp)
			t.Errorf("command %d (%s): expected status 200, got %d, response: %v", i, cmd, resp.StatusCode, result)
		}
	}
}

// =============================================================================
// INDEX RESERVED MACHINES TESTS - GET /v1/user/business/:businessID/uni-wash/reserved-machines
// =============================================================================

func TestIndexReservedMachines_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create multiple reservations
	now := time.Now()
	ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(1*time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/uni-wash/reserved-machines", business.ID), nil, token)

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
		t.Errorf("expected 2 reservations, got %d", len(data))
	}
}

func TestIndexReservedMachines_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/user/business/1/uni-wash/reserved-machines", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

func TestIndexReservedMachines_WithDateFilter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservations on different dates
	loc, _ := time.LoadLocation("Asia/Tehran")
	today := time.Now().In(loc)
	tomorrow := today.Add(24 * time.Hour)

	ta.CreateTestReservation(t, user.ID, product.ID, business.ID,
		time.Date(today.Year(), today.Month(), today.Day(), 10, 0, 0, 0, loc),
		time.Date(today.Year(), today.Month(), today.Day(), 11, 0, 0, 0, loc),
		schema.ReservationStatusReserved)

	ta.CreateTestReservation(t, user.ID, product.ID, business.ID,
		time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, loc),
		time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 11, 0, 0, 0, loc),
		schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with date filter for today
	dateFilter := today.Format(time.DateOnly)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/uni-wash/reserved-machines?Date=%s", business.ID, dateFilter), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only today's reservation
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 1 {
		t.Errorf("expected 1 reservation for today, got %d", len(data))
	}
}

func TestIndexReservedMachines_WithPagination(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create multiple reservations
	now := time.Now()
	for i := 0; i < 5; i++ {
		ta.CreateTestReservation(t, user.ID, product.ID, business.ID,
			now.Add(time.Duration(i*2)*time.Hour),
			now.Add(time.Duration(i*2+1)*time.Hour),
			schema.ReservationStatusReserved)
	}

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with pagination
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/uni-wash/reserved-machines?page=1&itemPerPage=2", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains paginated data
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 2 {
		t.Errorf("expected 2 reservations in paginated response, got: %v", data)
	}

	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return
	}

	if meta["total"] != float64(5) {
		t.Errorf("expected total 5 reservations, got: %v", meta["total"])
	}
}

func TestIndexReservedMachines_WithProductIDFilter(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product1 := ta.CreateTestProduct(t, post1.ID, business.ID, 50000, "09123456781", schema.UniWashMachineStatusON)

	post2 := ta.CreateTestPost(t, "Washing Machine 2", business.ID, user.ID)
	product2 := ta.CreateTestProduct(t, post2.ID, business.ID, 50000, "09123456782", schema.UniWashMachineStatusON)

	// Create reservations for different products
	now := time.Now()
	ta.CreateTestReservation(t, user.ID, product1.ID, business.ID, now.Add(1*time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, user.ID, product2.ID, business.ID, now.Add(1*time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request with ProductID filter
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/uni-wash/reserved-machines?ProductID=%d", business.ID, product1.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains only product1's reservation
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != 1 {
		t.Errorf("expected 1 reservation for product1, got %d", len(data))
	}
}

// =============================================================================
// CHECK LAST COMMAND STATUS TESTS - GET /v1/business/:businessID/uni-wash/check-last-command-status/:reservationID
// =============================================================================

func TestCheckLastCommandStatus_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation with last command
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Update reservation with last command
	reservation.Meta.UniWashLastCommand = schema.UniWashCommandON
	reservation.Meta.UniWashLastCommandReferenceID = "test-ref-123"
	ta.DB.Save(reservation)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/check-last-command-status/%d", business.ID, reservation.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains status data
	data, ok := result["Data"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Data in response, got: %v", result)
		return
	}

	// The mock MessageWay should return some status
	if data["status"] == nil && data["Status"] == nil {
		t.Logf("Response data: %v", data)
	}
}

func TestCheckLastCommandStatus_NoCommandSent(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation without any command
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/check-last-command-status/%d", business.ID, reservation.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response indicates no command sent
	data, ok := result["Data"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Data in response, got: %v", result)
		return
	}

	if data["OTPStatus"] != "دستوری ارسال نشده" {
		t.Errorf("expected 'no command sent' status, got: %v", data)
	}
}

func TestCheckLastCommandStatus_ReservationNotFound(t *testing.T) {
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

	// Make request for non-existent reservation
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/check-last-command-status/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// GET RESERVATION OPTIONS TESTS - GET /v1/business/:businessID/uni-wash/device/reservation-options
// =============================================================================

func TestGetReservationOptions_Success(t *testing.T) {
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

	// Make request
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/device/reservation-options", business.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)

	// Verify response contains data for all days of the week
	data, ok := result["Data"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Data map in response, got: %v", result)
		return
	}

	// Should have 7 days (0-6, Sunday to Saturday)
	if len(data) != 7 {
		t.Errorf("expected 7 days in reservation options, got %d", len(data))
	}

	// Each day should have 24 hours
	for dayNum, hours := range data {
		hoursArray, ok := hours.([]interface{})
		if !ok {
			t.Errorf("expected hours array for day %s, got: %v", dayNum, hours)
			continue
		}

		if len(hoursArray) != 24 {
			t.Errorf("expected 24 hours for day %s, got %d", dayNum, len(hoursArray))
		}
	}
}

func TestGetReservationOptions_Unauthorized(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Make request without token
	resp := ta.MakeRequest(t, http.MethodGet, "/v1/business/1/uni-wash/device/reservation-options", nil, "")

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for unauthorized, got %d", resp.StatusCode)
	}
}

// =============================================================================
// SEND DEVICE IS OFF MSG TO USER TESTS - POST /v1/business/:businessID/uni-wash/send-device-is-off-msg-to-user/:reservationID
// =============================================================================

func TestSendDeviceIsOffMsgToUser_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create customer user
	customer := ta.CreateTestUser(t, 9123456780, "customerPass", "Customer", "User", 0, nil)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusOFF)

	// Create reservation for customer
	now := time.Now()
	reservation := ta.CreateTestReservation(t, customer.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token for business owner
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-device-is-off-msg-to-user/%d", business.ID, reservation.ID), nil, token)

	if resp.StatusCode != http.StatusOK {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 200, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	if result["result"] != "success" {
		t.Errorf("expected success response, got: %v", result)
	}
}

func TestSendDeviceIsOffMsgToUser_ReservationNotFound(t *testing.T) {
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

	// Make request for non-existent reservation
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-device-is-off-msg-to-user/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// SEND FULL COUPON TO USER TESTS - POST /v1/business/:businessID/uni-wash/send-full-coupon-to-user/:reservationID
// =============================================================================

func TestSendFullCouponToUser_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user with business owner role
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with business permissions
	user.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(user)

	// Create customer user
	customer := ta.CreateTestUser(t, 9123456780, "customerPass", "Customer", "User", 0, nil)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusOFF)

	// Create reservation for customer
	now := time.Now()
	reservation := ta.CreateTestReservation(t, customer.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token for business owner
	token := ta.GenerateTestToken(t, user)

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-full-coupon-to-user/%d", business.ID, reservation.ID), nil, token)

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
	ta.DB.Model(&schema.Coupon{}).Where("business_id = ?", business.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 coupon to be created, got %d", count)
	}
}

func TestSendFullCouponToUser_ReservationNotFound(t *testing.T) {
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

	// Make request for non-existent reservation
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-full-coupon-to-user/99999", business.ID), nil, token)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 for not found, got %d", resp.StatusCode)
	}
}

// =============================================================================
// USER ROUTE TIME RESTRICTIONS TESTS
// =============================================================================

func TestSendCommand_UserRoute_TooEarly(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation starting in 30 minutes (more than 10 minutes from now)
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(30*time.Minute), now.Add(90*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/user/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 400 for too early, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	expectedMsg := "در بازه زمانی که رزرو کرده اید، دوباره تلاش کنید"
	if messages[0] != expectedMsg {
		t.Errorf("expected time restriction error message, got: %v", messages[0])
	}
}

func TestSendCommand_UserRoute_TooLate(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation that started 15 minutes ago (more than 10 minutes)
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now.Add(-15*time.Minute), now.Add(45*time.Minute), schema.ReservationStatusReserved)

	// Generate token
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/user/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 400 for too late, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	expectedMsg := "شما تا نهایتاً 10 دقیقه پس از شروع زمان فرصت داشتید به دستگاه فرمان بدهید."
	if messages[0] != expectedMsg {
		t.Errorf("expected time restriction error message, got: %v", messages[0])
	}
}

func TestSendCommand_UserRoute_NotReservedByUser(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Update user with user role
	user.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(user)

	// Create another user who owns the reservation
	otherUser := ta.CreateTestUser(t, 9123456780, "otherPass", "Other", "User", 0, nil)

	// Create test post and product
	post := ta.CreateTestPost(t, "Washing Machine 1", business.ID, user.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 50000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation for other user
	now := time.Now()
	reservation := ta.CreateTestReservation(t, otherUser.ID, product.ID, business.ID, now.Add(5*time.Minute), now.Add(65*time.Minute), schema.ReservationStatusReserved)

	// Generate token for the first user (not the owner)
	token := ta.GenerateTestToken(t, user)

	// Send command request
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}

	// Make request
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/user/business/%d/uni-wash/send-command", business.ID), sendReq, token)

	if resp.StatusCode != http.StatusBadRequest {
		result := ParseResponse(t, resp)
		t.Errorf("expected status 400 for not reserved by user, got %d, response: %v", resp.StatusCode, result)
		return
	}

	result := ParseResponse(t, resp)
	messages, ok := result["Messages"].([]interface{})
	if !ok || len(messages) == 0 {
		t.Errorf("expected error message, got: %v", result)
		return
	}

	expectedMsg := "شما این دستگاه را رزرو نکرده اید"
	if messages[0] != expectedMsg {
		t.Errorf("expected not reserved error message, got: %v", messages[0])
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestUniWash_Integration_FullWorkflow(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "ownerPass", "Business", "Owner", 0, nil)
	business := ta.CreateTestBusiness(t, "UniWash Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create customer user
	customer := ta.CreateTestUser(t, 9123456780, "customerPass", "Customer", "User", 0, nil)
	customer.Permissions[business.ID] = []schema.UserRole{schema.URUser}
	ta.DB.Save(customer)

	// Create washing machine
	post := ta.CreateTestPost(t, "Washing Machine Premium", business.ID, owner.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 75000, "09123456789", schema.UniWashMachineStatusON)

	// Create reservation for customer
	now := time.Now()
	reservation := ta.CreateTestReservation(t, customer.ID, product.ID, business.ID, now.Add(5*time.Minute), now.Add(65*time.Minute), schema.ReservationStatusReserved)

	// Generate tokens
	ownerToken := ta.GenerateTestToken(t, owner)
	customerToken := ta.GenerateTestToken(t, customer)

	// 1. Customer sends ON command
	sendReq := request.SendCommand{
		ReservationID: reservation.ID,
		ProductID:     product.ID,
		Command:       schema.UniWashCommandON,
	}
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/user/business/%d/uni-wash/send-command", business.ID), sendReq, customerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("customer send command failed: %d", resp.StatusCode)
	}

	// 2. Business owner checks command status
	resp = ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/check-last-command-status/%d", business.ID, reservation.ID), nil, ownerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("check command status failed: %d", resp.StatusCode)
	}

	// 3. Customer views their reservations
	resp = ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/user/business/%d/uni-wash/reserved-machines", business.ID), nil, customerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("view reservations failed: %d", resp.StatusCode)
	}

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("expected 1 reservation in response, got: %v", data)
	}

	// 4. Business owner gets reservation options
	resp = ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/uni-wash/device/reservation-options", business.ID), nil, ownerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get reservation options failed: %d", resp.StatusCode)
	}

	t.Log("UniWash full workflow integration test passed!")
}

func TestUniWash_Integration_MachineFailure_Compensation(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create business owner
	owner := ta.CreateTestUser(t, 9123456789, "ownerPass", "Business", "Owner", 0, nil)
	business := ta.CreateTestBusiness(t, "UniWash Business", schema.BTypeGymManager, owner.ID)

	// Update owner with business permissions
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create customer user
	customer := ta.CreateTestUser(t, 9123456780, "customerPass", "Customer", "User", 0, nil)

	// Create washing machine that failed (OFF)
	post := ta.CreateTestPost(t, "Washing Machine Premium", business.ID, owner.ID)
	product := ta.CreateTestProduct(t, post.ID, business.ID, 75000, "09123456789", schema.UniWashMachineStatusOFF)

	// Create reservation for customer
	now := time.Now()
	reservation := ta.CreateTestReservation(t, customer.ID, product.ID, business.ID, now.Add(-5*time.Minute), now.Add(55*time.Minute), schema.ReservationStatusReserved)

	// Generate token for business owner
	ownerToken := ta.GenerateTestToken(t, owner)

	// 1. Business owner sends device off message to customer
	resp := ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-device-is-off-msg-to-user/%d", business.ID, reservation.ID), nil, ownerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("send device off message failed: %d", resp.StatusCode)
	}

	// 2. Business owner sends compensation coupon to customer
	resp = ta.MakeRequest(t, http.MethodPost, fmt.Sprintf("/v1/business/%d/uni-wash/send-full-coupon-to-user/%d", business.ID, reservation.ID), nil, ownerToken)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("send coupon failed: %d", resp.StatusCode)
	}

	// 3. Verify coupon was created with correct value (product price * 1.1 for tax)
	var coupon schema.Coupon
	if err := ta.DB.Where("business_id = ?", business.ID).First(&coupon).Error; err != nil {
		t.Fatalf("failed to find created coupon: %v", err)
	}

	expectedValue := product.Price * 1.1
	if coupon.Value != expectedValue {
		t.Errorf("expected coupon value %f, got %f", expectedValue, coupon.Value)
	}

	if coupon.Type != schema.CouponTypeFixedAmount {
		t.Errorf("expected coupon type FixedAmount, got %s", coupon.Type)
	}

	t.Log("UniWash machine failure compensation workflow test passed!")
}

