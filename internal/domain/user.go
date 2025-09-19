package domain

import (
	"context"
	"time"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` // Hidden from JSON
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Role  *string `json:"role,omitempty"`
}

// LoginRequest represents user login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents user data returned to clients (without sensitive data)
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User entity to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id string, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (string, *UserResponse, error) // returns token and user
	GetProfile(ctx context.Context, userID string) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateUserRequest) (*UserResponse, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*UserResponse, int64, error)
	GetUserByID(ctx context.Context, id string) (*UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	RefreshToken(ctx context.Context, userID string) (string, error)
}

// TokenService defines the interface for JWT token operations
type TokenService interface {
	GenerateToken(user *User) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	ExtractUserIDFromToken(tokenString string) (string, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// Common errors
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound       = &DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUserAlreadyExists  = &DomainError{Code: "USER_ALREADY_EXISTS", Message: "User with this email already exists"}
	ErrInvalidCredentials = &DomainError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrInvalidToken       = &DomainError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrUnauthorized       = &DomainError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrForbidden          = &DomainError{Code: "FORBIDDEN", Message: "Access forbidden"}
	ErrValidationFailed   = &DomainError{Code: "VALIDATION_FAILED", Message: "Validation failed"}
)
