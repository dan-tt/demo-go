package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"demo-go/internal/domain"
	"demo-go/internal/logger"

	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService domain.UserService
	logger      *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger.GetGlobal().ForComponent("handler"),
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := h.logger.ForRequest(r.Method, r.URL.Path, h.getRequestID(r))

	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Invalid request body for registration", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	log.Info("User registration attempt", "email", req.Email)

	user, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		log.Error("User registration failed", "email", req.Email, "error", err)
		h.handleServiceError(w, err)
		return
	}

	log.Info("User registered successfully", "user_id", user.ID, "email", user.Email)
	h.writeSuccessResponse(w, http.StatusCreated, "User registered successfully", user)
}

// Login handles user authentication
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := h.logger.ForRequest(r.Method, r.URL.Path, h.getRequestID(r))

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Invalid request body for login", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	log.Info("User login attempt", "email", req.Email)

	token, user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		log.Error("User login failed", "email", req.Email, "error", err)
		h.handleServiceError(w, err)
		return
	}

	log.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	h.writeSuccessResponse(w, http.StatusOK, "Login successful", response)
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log := h.logger.ForRequest(r.Method, r.URL.Path, h.getRequestID(r))

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		log.Warn("Unauthorized profile access attempt")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	log.Debug("Getting user profile", "user_id", userID)

	user, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		log.Error("Failed to get user profile", "user_id", userID, "error", err)
		h.handleServiceError(w, err)
		return
	}

	log.Info("User profile retrieved successfully", "user_id", userID)
	h.writeSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log := h.logger.ForRequest(r.Method, r.URL.Path, h.getRequestID(r))

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		log.Warn("Unauthorized profile update attempt")
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Invalid request body for profile update", "user_id", userID, "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	log.Info("Profile update attempt", "user_id", userID)

	user, err := h.userService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		log.Error("Profile update failed", "user_id", userID, "error", err)
		h.handleServiceError(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, "Profile updated successfully", user)
}

// GetUsers handles getting all users (admin only)
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, total, err := h.userService.GetUsers(r.Context(), limit, offset)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	h.writeSuccessResponse(w, http.StatusOK, "Users retrieved successfully", response)
}

// GetUserByID handles getting a specific user by ID (admin only)
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Missing user ID", "User ID is required")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// DeleteUser handles deleting a user (admin only)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Missing user ID", "User ID is required")
		return
	}

	err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// RefreshToken handles token refresh
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}

	token, err := h.userService.RefreshToken(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := map[string]string{
		"token": token,
	}

	h.writeSuccessResponse(w, http.StatusOK, "Token refreshed successfully", response)
}

// Health check endpoint
func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "clean-architecture-api",
		"timestamp": "2025-09-18T00:00:00Z",
	}

	h.writeSuccessResponse(w, http.StatusOK, "Service is healthy", response)
}

// Helper methods

func (h *UserHandler) getUserIDFromContext(r *http.Request) string {
	if userID := r.Context().Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
	if domainErr, ok := err.(*domain.DomainError); ok {
		switch domainErr.Code {
		case "USER_NOT_FOUND":
			h.writeErrorResponse(w, http.StatusNotFound, domainErr.Message, domainErr.Code)
		case "USER_ALREADY_EXISTS":
			h.writeErrorResponse(w, http.StatusConflict, domainErr.Message, domainErr.Code)
		case "INVALID_CREDENTIALS":
			h.writeErrorResponse(w, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
		case "INVALID_TOKEN":
			h.writeErrorResponse(w, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
		case "UNAUTHORIZED":
			h.writeErrorResponse(w, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
		case "FORBIDDEN":
			h.writeErrorResponse(w, http.StatusForbidden, domainErr.Message, domainErr.Code)
		case "VALIDATION_FAILED":
			h.writeErrorResponse(w, http.StatusBadRequest, domainErr.Message, domainErr.Code)
		default:
			h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
		}
	} else {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
	}
}

func (h *UserHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we can't encode the response, there's not much we can do
		// The status code has already been set
		return
	}
}

func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	response := map[string]interface{}{
		"success": false,
		"message": message,
		"error": map[string]string{
			"code": code,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we can't encode the response, there's not much we can do
		// The status code has already been set
		return
	}
}

// Helper methods
func (h *UserHandler) getRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	if requestID := r.Context().Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return "unknown"
}
