package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"demo-go/internal/domain"
	"demo-go/internal/handler"
)

// mockUserService implements domain.UserService for testing
type mockUserService struct {
	registerFunc      func(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error)
	loginFunc         func(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error)
	getProfileFunc    func(ctx context.Context, userID string) (*domain.UserResponse, error)
	updateProfileFunc func(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error)
	getUsersFunc      func(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error)
	getUserByIDFunc   func(ctx context.Context, id string) (*domain.UserResponse, error)
	deleteUserFunc    func(ctx context.Context, id string) error
	refreshTokenFunc  func(ctx context.Context, userID string) (string, error)
}

func (m *mockUserService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserService) Login(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return "", nil, fmt.Errorf("not implemented")
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	if m.updateProfileFunc != nil {
		return m.updateProfileFunc(ctx, userID, req)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserService) GetUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
	if m.getUsersFunc != nil {
		return m.getUsersFunc(ctx, limit, offset)
	}
	return nil, 0, fmt.Errorf("not implemented")
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(ctx, id)
	}
	return fmt.Errorf("not implemented")
}

func (m *mockUserService) RefreshToken(ctx context.Context, userID string) (string, error) {
	if m.refreshTokenFunc != nil {
		return m.refreshTokenFunc(ctx, userID)
	}
	return "", fmt.Errorf("not implemented")
}

func (m *mockUserService) GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error) {
	if m.getProfileFunc != nil {
		return m.getProfileFunc(ctx, userID)
	}
	return nil, fmt.Errorf("not implemented")
}

// Test data
var (
	testUser = &domain.UserResponse{
		ID:        "test-user-1",
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testAdmin = &domain.UserResponse{
		ID:        "admin-user-1",
		Name:      "Admin User",
		Email:     "admin@example.com",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*mockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful user registration",
			requestBody: domain.CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.registerFunc = func(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
					return testUser, nil
				}
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				if body["message"].(string) != "User registered successfully" {
					t.Error("Unexpected message")
				}
				data := body["data"].(map[string]interface{})
				if data["email"].(string) != testUser.Email {
					t.Error("Unexpected email in response")
				}
			},
		},
		{
			name: "duplicate email registration",
			requestBody: domain.CreateUserRequest{
				Name:     "Jane Doe",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.registerFunc = func(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
					return nil, domain.ErrUserAlreadyExists
				}
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != domain.ErrUserAlreadyExists.Message {
					t.Error("Unexpected error message")
				}
			},
		},
		{
			name: "validation error - invalid email",
			requestBody: domain.CreateUserRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.registerFunc = func(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
					return nil, &domain.Error{Code: "VALIDATION_FAILED", Message: "Invalid email format"}
				}
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
			},
		},
		{
			name:           "invalid JSON request body",
			requestBody:    `{"invalid": json}`,
			mockSetup:      func(m *mockUserService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if !strings.Contains(body["message"].(string), "Invalid request body") {
					t.Error("Expected invalid request body message")
				}
			},
		},
		{
			name: "internal server error",
			requestBody: domain.CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.registerFunc = func(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
					return nil, fmt.Errorf("database connection failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
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
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.Register(rr, req)

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

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    domain.LoginRequest
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful login",
			requestBody: domain.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.loginFunc = func(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
					return "jwt-token-123", testUser, nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("Expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["token"].(string) != "jwt-token-123" {
					t.Error("Expected JWT token in response")
				}
				user := data["user"].(map[string]interface{})
				if user["email"].(string) != testUser.Email {
					t.Error("Expected user data in response")
				}
			},
		},
		{
			name: "invalid credentials",
			requestBody: domain.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mockUserService) {
				m.loginFunc = func(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
					return "", nil, domain.ErrInvalidCredentials
				}
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != domain.ErrInvalidCredentials.Message {
					t.Error("Expected invalid credentials message")
				}
			},
		},
		{
			name: "user not found",
			requestBody: domain.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockUserService) {
				m.loginFunc = func(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
					return "", nil, domain.ErrInvalidCredentials
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

			// Create request
			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.Login(rr, req)

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

func TestUserHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mockUserService)
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "successful get profile",
			userID: "test-user-1",
			mockSetup: func(m *mockUserService) {
				m.getProfileFunc = func(ctx context.Context, userID string) (*domain.UserResponse, error) {
					if userID == "test-user-1" {
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
				if data["email"].(string) != testUser.Email {
					t.Error("Expected user email in response")
				}
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent-user",
			mockSetup: func(m *mockUserService) {
				m.getProfileFunc = func(ctx context.Context, userID string) (*domain.UserResponse, error) {
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
			name:   "missing user ID in context",
			userID: "",
			mockSetup: func(m *mockUserService) {
				// No mock setup needed as handler should return early
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("Expected success to be false")
				}
				if body["message"].(string) != "Unauthorized" {
					t.Error("Expected unauthorized message")
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
			req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", http.NoBody)

			// Add user ID to context if provided
			if tt.userID != "" {
				ctx := context.WithValue(req.Context(), "user_id", tt.userID)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			userHandler.GetProfile(rr, req)

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
