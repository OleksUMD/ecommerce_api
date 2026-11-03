package services

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/events"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/notifications"
	"github.com/OleksUMD/ecommerce_api/internal/repositories"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo       repositories.IUserRepository
	cartRepo       repositories.ICartRepository
	config         *config.Config
	eventPublisher events.Publisher
}

func NewAuthService(
	userRepo repositories.IUserRepository,
	cartRepo repositories.ICartRepository,
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

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
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
	tmpTokenStr, err := s.generateTmpTokenStr()
	if err != nil {
		return nil, err
	}
	return s.generateUnverifiedUserResponse(tmpTokenStr, &user)
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
	err = s.userRepo.DeleteRefreshTokenByID(refreshToken.ID)
	if err != nil {
		return nil, errors.New("failed to remove refresh token")
	}
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
	err = s.userRepo.CreateRefreshToken(&refreshTokenModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create access tokens: %w", err)
	}
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

func (s *AuthService) generateUnverifiedUserResponse(tmpTokenStr string, user *models.User) (*dto.UserResponse, error) {
	baseAddr := net.JoinHostPort(s.config.Server.Host, s.config.Server.Port)
	tmpLink := fmt.Sprintf("http://%s/api/v1/auth/verify/%s", baseAddr, tmpTokenStr)
	payload := &models.UnverifiedUser{
		UserID:   user.ID,
		UserName: user.FirstName + " " + user.LastName,
		Link:     tmpLink,
		Email:    user.Email,
	}
	err := s.userRepo.SaveUnverifiedUser(tmpTokenStr, payload)
	if err != nil {
		return nil, err
	}
	err = s.eventPublisher.Publish(notifications.UserRegistered, payload, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate registration event message: %w", err)
	}
	regResponse := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return regResponse, nil
}

func (s *AuthService) VerifyEmail(token string) error {
	return s.userRepo.VerifyUserEmail(token)
}

func (s *AuthService) VerifyTOTP(token, code string) error {
	return nil
}

func (s *AuthService) generateTmpTokenStr() (string, error) {
	token, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (s *AuthService) ResubmitVerificationLink(userID uint) (*dto.UserResponse, error) {
	tmpTokenStr, err := s.generateTmpTokenStr()
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return s.generateUnverifiedUserResponse(tmpTokenStr, user)
}
