package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/events"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/notifications"
	"github.com/OleksUMD/ecommerce_api/internal/repositories"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
)

type AuthService struct {
	userRepo       repositories.UserRepositoryInterface
	cartRepo       repositories.CartRepositoryInterface
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	cartRepo repositories.CartRepositoryInterface,
	config *config.Config,
	eventPublisher events.Publisher,
) *AuthService {
	return &AuthService{
		cartRepo:       cartRepo,
		userRepo:       userRepo,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	exists, _ := s.userRepo.GetByEmail(req.Email)
	if exists != nil {
		return nil, errors.New("this email is already registered")
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		return nil, err
	}
	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	refreshToken, err := s.userRepo.GetValidRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	s.userRepo.DeleteRefreshTokenByID(refreshToken.ID)
	return s.generateAuthResponse(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshtokenExpires),
	}
	s.userRepo.CreateRefreshToken(&refreshTokenModel)
	err = s.eventPublisher.Publish(notifications.UserLoggedIn, user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth event message: %w", err)
	}
	authResponse := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return authResponse, nil
}
