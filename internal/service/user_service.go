// Package service provides business logic implementations for the demo-go application.
// It includes user management services, JWT token services, and caching functionality.
package service

import (
	"context"
	"fmt"
	"strings"

	"demo-go/internal/domain"
	"demo-go/internal/logger"

	"golang.org/x/crypto/bcrypt"
)

// userService implements domain.UserService
type userService struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
	logger       *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository, tokenService domain.TokenService) domain.UserService {
	return &userService{
		userRepo:     userRepo,
		tokenService: tokenService,
		logger:       logger.GetGlobal().ForComponent("user-service"),
	}
}

// Register creates a new user account
func (s *userService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	log := s.logger.ForService("user", "register").WithField("email", req.Email)

	log.Debug("Starting user registration")

	// Validate request
	if err := s.validateCreateUserRequest(req); err != nil {
		log.Warn("User registration validation failed", "error", err)
		return nil, err
	}

	// Check if user already exists
	log.Debug("Checking if user already exists")
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != domain.ErrUserNotFound {
		log.Error("Error checking existing user", "error", err)
		return nil, err
	}
	if existingUser != nil {
		log.Warn("User already exists")
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	log.Debug("Hashing password")
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		log.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	log.Debug("Creating user entity", "role", role)

	// Create user entity
	user := &domain.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
		Role:     role,
	}

	// Save user
	log.Debug("Saving user to repository")
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error("Failed to create user in repository", "error", err)
		return nil, err
	}

	log.Info("User registered successfully", "user_id", user.ID)
	return user.ToResponse(), nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (string, *domain.UserResponse, error) {
	log := s.logger.ForService("user", "login").WithField("email", req.Email)

	log.Debug("Starting user login")

	// Validate request
	if err := s.validateLoginRequest(req); err != nil {
		log.Warn("Login validation failed", "error", err)
		return "", nil, err
	}

	// Get user by email
	log.Debug("Looking up user by email")
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if err == domain.ErrUserNotFound {
			log.Warn("Login attempt with non-existent email")
			return "", nil, domain.ErrInvalidCredentials
		}
		log.Error("Error retrieving user", "error", err)
		return "", nil, err
	}

	// Verify password
	if err := s.verifyPassword(user.Password, req.Password); err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user.ToResponse(), nil
}

// GetProfile retrieves user profile by user ID
func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// UpdateProfile updates user profile
func (s *userService) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate update request
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Prepare updated user
	updatedUser := *existingUser

	// Update fields if provided
	if req.Name != nil {
		updatedUser.Name = strings.TrimSpace(*req.Name)
	}

	if req.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*req.Email))
		if newEmail != existingUser.Email {
			// Check if new email already exists
			_, err := s.userRepo.GetByEmail(ctx, newEmail)
			if err != nil && err != domain.ErrUserNotFound {
				return nil, err
			}
			if err == nil {
				return nil, domain.ErrUserAlreadyExists
			}
		}
		updatedUser.Email = newEmail
	}

	if req.Role != nil {
		updatedUser.Role = *req.Role
	}

	// Update user
	if err := s.userRepo.Update(ctx, userID, &updatedUser); err != nil {
		return nil, err
	}

	return updatedUser.ToResponse(), nil
}

// GetUsers retrieves all users with pagination
func (s *userService) GetUsers(ctx context.Context, limit, offset int) ([]*domain.UserResponse, int64, error) {
	// Set default and max limits
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response format
	var userResponses []*domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, count, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// RefreshToken generates a new token for the user
func (s *userService) RefreshToken(ctx context.Context, userID string) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	token, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// Helper methods

func (s *userService) validateCreateUserRequest(req *domain.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Name is required"}
	}

	if len(strings.TrimSpace(req.Name)) < 2 {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Name must be at least 2 characters long"}
	}

	if strings.TrimSpace(req.Email) == "" {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Email is required"}
	}

	if !s.isValidEmail(req.Email) {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Invalid email format"}
	}

	if len(req.Password) < 6 {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Password must be at least 6 characters long"}
	}

	return nil
}

func (s *userService) validateLoginRequest(req *domain.LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Email is required"}
	}

	if req.Password == "" {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Password is required"}
	}

	return nil
}

func (s *userService) validateUpdateUserRequest(req *domain.UpdateUserRequest) error {
	if req.Name != nil && len(strings.TrimSpace(*req.Name)) < 2 {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Name must be at least 2 characters long"}
	}

	if req.Email != nil && !s.isValidEmail(*req.Email) {
		return &domain.DomainError{Code: "VALIDATION_FAILED", Message: "Invalid email format"}
	}

	return nil
}

func (s *userService) isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (s *userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *userService) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
