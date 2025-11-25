package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/web-dashboard-made-by-renz/backend/internal/models"
)

type AuthService interface {
	Login(username, password string) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type authService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

// Hardcoded admin credentials
const (
	AdminUsername = "admin"
	AdminPassword = "admin123"
)

func (s *authService) Login(username, password string) (*models.LoginResponse, error) {
	// Validate credentials
	if username != AdminUsername || password != AdminPassword {
		return nil, errors.New("invalid username or password")
	}

	// Create user object
	user := models.User{
		Username: AdminUsername,
		Role:     "admin",
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		"iat":      time.Now().Unix(),
	})

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
