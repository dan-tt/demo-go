package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demo-go/internal/domain"
	"demo-go/internal/handler"

	"github.com/gorilla/mux"
)

const testUserID = "test-user-1"

func TestUserHandler_UpdateProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    domain.UpdateUserRequest
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful profile update",
			userID: testUserID,
			requestBody: domain.UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			},
			mockSetup: func(m *mockUserService) {
				m.updateProfileFunc = func(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
					updatedUser := *testUser
					if req.Name != nil {
						updatedUser.Name = *req.Name
					}
					return &updatedUser, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["name"].(string) != "Updated Name" {
					t.Error("Expected updated name in response")
				}
			},
		},
		{
			name:   "update with duplicate email",
			userID: testUserID,
			requestBody: domain.UpdateUserRequest{
				Email: stringPtr("existing@example.com"),
			},
			mockSetup: func(m *mockUserService) {
				m.updateProfileFunc = func(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
					return nil, domain.ErrUserAlreadyExists
				}
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != domain.ErrUserAlreadyExists.Message {
					t.Error("Expected user already exists message")
				}
			},
		},
		{
			name:        "missing user ID in context",
			userID:      "",
			requestBody: domain.UpdateUserRequest{},
			mockSetup: func(m *mockUserService) {
				// No setup needed for this test case
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockUserService{}
			tt.mockSetup(mockService)

			// Create handler
			userHandler := handler.NewUserHandler(mockService)

			// Create request
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Add user ID to context if provided
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "user_id", tt.userID)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.UpdateProfile(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:        "successful get users with default pagination",
			queryParams: map[string]string{},
			mockSetup: func(m *mockUserService) {
				m.getUsersFunc = func(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
					if limit == 10 && offset == 0 {
						return []*domain.UserResponse{testUser, testAdmin}, 2, nil
					}
					return []*domain.UserResponse{}, 0, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				users := data["users"].([]interface{})
				if len(users) != 2 {
					t.Errorf("Expected 2 users, got %d", len(users))
				}
				if data["total"].(float64) != 2 {
					t.Error("Expected total count of 2")
				}
			},
		},
		{
			name: "successful get users with custom pagination",
			queryParams: map[string]string{
				"limit":  "5",
				"offset": "10",
			},
			mockSetup: func(m *mockUserService) {
				m.getUsersFunc = func(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
					if limit == 5 && offset == 10 {
						return []*domain.UserResponse{testUser}, 1, nil
					}
					return []*domain.UserResponse{}, 0, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["limit"].(float64) != 5 {
					t.Error("Expected limit to be 5")
				}
				if data["offset"].(float64) != 10 {
					t.Error("Expected offset to be 10")
				}
			},
		},
		{
			name: "invalid query parameters (non-numeric)",
			queryParams: map[string]string{
				"limit":  "invalid",
				"offset": "invalid",
			},
			mockSetup: func(m *mockUserService) {
				m.getUsersFunc = func(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
					// Should use defaults (10, 0) when invalid params provided
					if limit == 10 && offset == 0 {
						return []*domain.UserResponse{}, 0, nil
					}
					return []*domain.UserResponse{}, 0, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				// Should use default values when invalid params provided
				data := body["data"].(map[string]interface{})
				if data["limit"].(float64) != 10 {
					t.Error("Expected default limit of 10")
				}
				if data["offset"].(float64) != 0 {
					t.Error("Expected default offset of 0")
				}
			},
		},
		{
			name:        "service error",
			queryParams: map[string]string{},
			mockSetup: func(m *mockUserService) {
				m.getUsersFunc = func(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
					return nil, 0, domain.ErrUnauthorized
				}
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockUserService{}
			tt.mockSetup(mockService)

			// Create handler
			userHandler := handler.NewUserHandler(mockService)

			// Build URL with query parameters
			url := "/api/v1/admin/users"
			if len(tt.queryParams) > 0 {
				url += "?"
				for key, value := range tt.queryParams {
					url += key + "=" + value + "&"
				}
				url = url[:len(url)-1] // Remove trailing &
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.GetUsers(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful get user by ID",
			userID: testUserID,
			mockSetup: func(m *mockUserService) {
				m.getUserByIDFunc = func(ctx context.Context, id string) (*domain.UserResponse, error) {
					if id == testUserID {
						return testUser, nil
					}
					return nil, domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["id"].(string) != testUserID {
					t.Error("Expected correct user ID in response")
				}
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			mockSetup: func(m *mockUserService) {
				m.getUserByIDFunc = func(ctx context.Context, id string) (*domain.UserResponse, error) {
					return nil, domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != domain.ErrUserNotFound.Message {
					t.Error("Expected user not found message")
				}
			},
		},
		{
			name:           "missing user ID parameter",
			userID:         "",
			mockSetup:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != "Missing user ID" {
					t.Error("Expected missing user ID message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockUserService{}
			tt.mockSetup(mockService)

			// Create handler
			userHandler := handler.NewUserHandler(mockService)

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/"+tt.userID, http.NoBody)

			// Setup mux vars
			vars := map[string]string{}
			if tt.userID != "" {
				vars["id"] = tt.userID
			}
			req = mux.SetURLVars(req, vars)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.GetUserByID(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful user deletion",
			userID: testUserID,
			mockSetup: func(m *mockUserService) {
				m.deleteUserFunc = func(ctx context.Context, id string) error {
					if id == testUserID {
						return nil
					}
					return domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				if body["message"].(string) != "User deleted successfully" {
					t.Error("Expected success message")
				}
			},
		},
		{
			name:   "user not found for deletion",
			userID: "nonexistent-user",
			mockSetup: func(m *mockUserService) {
				m.deleteUserFunc = func(ctx context.Context, id string) error {
					return domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != domain.ErrUserNotFound.Message {
					t.Error("Expected user not found message")
				}
			},
		},
		{
			name:           "missing user ID parameter",
			userID:         "",
			mockSetup:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != "Missing user ID" {
					t.Error("Expected missing user ID message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockUserService{}
			tt.mockSetup(mockService)

			// Create handler
			userHandler := handler.NewUserHandler(mockService)

			// Create request with mux vars
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/users/"+tt.userID, http.NoBody)

			// Setup mux vars
			vars := map[string]string{}
			if tt.userID != "" {
				vars["id"] = tt.userID
			}
			req = mux.SetURLVars(req, vars)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.DeleteUser(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

func TestUserHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful token refresh",
			userID: testUserID,
			mockSetup: func(m *mockUserService) {
				m.refreshTokenFunc = func(ctx context.Context, userID string) (string, error) {
					if userID == testUserID {
						return "new-jwt-token-456", nil
					}
					return "", domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["token"].(string) != "new-jwt-token-456" {
					t.Error("Expected new JWT token in response")
				}
			},
		},
		{
			name:   "user not found for token refresh",
			userID: "nonexistent-user",
			mockSetup: func(m *mockUserService) {
				m.refreshTokenFunc = func(ctx context.Context, userID string) (string, error) {
					return "", domain.ErrUserNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
			},
		},
		{
			name:           "missing user ID in context",
			userID:         "",
			mockSetup:      func(m *mockUserService) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := &mockUserService{}
			tt.mockSetup(mockService)

			// Create handler
			userHandler := handler.NewUserHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", http.NoBody)

			// Add user ID to context if provided
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "user_id", tt.userID)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.RefreshToken(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

func TestUserHandler_Health(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:           "health check successful",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				if body["message"].(string) != "Service is healthy" {
					t.Error("Expected health message")
				}
				data := body["data"].(map[string]interface{})
				if data["status"].(string) != "healthy" {
					t.Error("Expected healthy status")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with nil service (health check doesn't use it)
			userHandler := handler.NewUserHandler(nil)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.Health(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Parse response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// Run custom checks
			if tt.checkResponse != nil {
				tt.checkResponse(t, responseBody)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
