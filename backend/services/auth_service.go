package services

import (
	"ciphersafe/config"
	"ciphersafe/models"
	"ciphersafe/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService handles registration, login, and JWT generation
type AuthService struct {
	UserService *UserService
}

// NewAuthService creates a new AuthService
func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{UserService: userService}
}

// Register creates a new user, hashes their password, and saves them
func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	_, err := s.UserService.FindUserByEmail(email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash the password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create the user
	return s.UserService.CreateUser(email, passwordHash)
}

// Login validates user credentials and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.UserService.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	// Check the password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	// Generate and return JWT token
	return s.generateJWT(user.ID)
}

// generateJWT creates a new JWT token for a given user ID
func (s *AuthService) generateJWT(userID uint) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"sub": userID,                                // 'sub' (subject) is the standard claim for user ID
		"iat": time.Now().Unix(),                     // 'iat' (issued at)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 'exp' (expiration time) - 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with our secret
	tokenString, err := token.SignedString(config.AppConfig.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
