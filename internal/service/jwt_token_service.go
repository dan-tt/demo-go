package service

import (
	"time"

	"demo-go/internal/config"
	"demo-go/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// jwtTokenService implements domain.TokenService using JWT
type jwtTokenService struct {
	secretKey      []byte
	expirationTime time.Duration
	issuer         string
}

// NewJWTTokenService creates a new JWT token service
func NewJWTTokenService(cfg *config.Config) domain.TokenService {
	return &jwtTokenService{
		secretKey:      []byte(cfg.JWT.SecretKey),
		expirationTime: cfg.JWT.Expiration,
		issuer:         "demo-go-api",
	}
}

// GenerateToken generates a JWT token for the given user
func (s *jwtTokenService) GenerateToken(user *domain.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(s.expirationTime)
	
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expirationTime.Unix(),
		"iat":     now.Unix(),
		"iss":     s.issuer,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *jwtTokenService) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure token's signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return s.secretKey, nil
	})
	
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	
	if !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	// Extract claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	email, ok := claims["email"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	role, ok := claims["role"].(string)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	
	return &domain.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Exp:    int64(exp),
		Iat:    int64(iat),
	}, nil
}

// ExtractUserIDFromToken extracts user ID from a JWT token
func (s *jwtTokenService) ExtractUserIDFromToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.UserID, nil
}
