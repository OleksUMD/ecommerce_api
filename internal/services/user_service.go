package services

import (
	"errors"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db     *gorm.DB
	config *config.Config
}

func NewUserService(db *gorm.DB, config *config.Config) *UserService {
	return &UserService{
		db:     db,
		config: config,
	}
}

func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return s.GetProfile(userID)
}
