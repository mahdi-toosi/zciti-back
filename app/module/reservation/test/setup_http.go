package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =============================================================================
// HTTP Request Helpers
// =============================================================================

// MakeRequest makes an HTTP request to the test server
func (ta *TestApp) MakeRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := ta.App.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	return resp
}

// MakeAuthenticatedRequest is a convenience method that uses the provided token
func (ta *TestApp) MakeAuthenticatedRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
	t.Helper()
	return ta.MakeRequest(t, method, path, body, token)
}

// MakeUnauthenticatedRequest makes a request without a token
func (ta *TestApp) MakeUnauthenticatedRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	t.Helper()
	return ta.MakeRequest(t, method, path, body, "")
}

// =============================================================================
// HTTP Response Helpers
// =============================================================================

// ParseResponse parses the response body into a map
func ParseResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Try parsing as string
		var strResult string
		if err := json.Unmarshal(body, &strResult); err != nil {
			t.Fatalf("failed to parse response: %v, body: %s", err, string(body))
		}
		return map[string]interface{}{"result": strResult}
	}

	return result
}

// ParseResponseTo parses the response body into the provided struct
func ParseResponseTo(t *testing.T, resp *http.Response, target interface{}) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to parse response: %v, body: %s", err, string(body))
	}
}

// GetResponseBody returns the raw response body as string
func GetResponseBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	defer resp.Body.Close()

	return string(body)
}

// =============================================================================
// Response Assertions
// =============================================================================

// AssertStatus asserts the response status code matches expected
func AssertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		result := ParseResponse(t, resp)
		t.Errorf("expected status %d, got %d, response: %v", expected, resp.StatusCode, result)
	}
}

// AssertOK asserts the response status is 200 OK
func AssertOK(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusOK)
}

// AssertUnauthorized asserts the response status is 401 Unauthorized
func AssertUnauthorized(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusUnauthorized)
}

// AssertForbidden asserts the response status is 403 Forbidden
func AssertForbidden(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusForbidden)
}

// AssertNotFound asserts the response status is 404 Not Found
func AssertNotFound(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusNotFound)
}

// AssertUnprocessableEntity asserts the response status is 422 Unprocessable Entity
func AssertUnprocessableEntity(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusUnprocessableEntity)
}

// AssertInternalServerError asserts the response status is 500 Internal Server Error
func AssertInternalServerError(t *testing.T, resp *http.Response) {
	t.Helper()
	AssertStatus(t, resp, http.StatusInternalServerError)
}

// =============================================================================
// Response Data Helpers
// =============================================================================

// GetDataFromResponse extracts the Data field from a standard response
func GetDataFromResponse(t *testing.T, resp *http.Response) []interface{} {
	t.Helper()

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return nil
	}

	return data
}

// GetMetaFromResponse extracts the Meta field from a standard response
func GetMetaFromResponse(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()

	result := ParseResponse(t, resp)
	meta, ok := result["Meta"].(map[string]interface{})
	if !ok {
		t.Errorf("expected Meta in response, got: %v", result)
		return nil
	}

	return meta
}

// AssertDataCount asserts the number of items in the Data array
func AssertDataCount(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	result := ParseResponse(t, resp)
	data, ok := result["Data"].([]interface{})
	if !ok {
		t.Errorf("expected Data array in response, got: %v", result)
		return
	}

	if len(data) != expected {
		t.Errorf("expected %d items in Data, got %d", expected, len(data))
	}
}
