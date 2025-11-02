package services

import (
	"ciphersafe/models"

	"gorm.io/gorm"
)

// UserService handles user-related database operations
type UserService struct {
	DB *gorm.DB
}

// NewUserService creates a new UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(email, passwordHash string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Password: passwordHash,
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// FindUserByEmail finds a user by their email address
func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindUserByID finds a user by their ID
func (s *UserService) FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
