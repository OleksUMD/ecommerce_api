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

	tmpTokenStr, err := s.generateTmpTokenStr()
	if err != nil {
		return nil, err
	}
	return s.generateUnverifiedUserResponse(tmpTokenStr, &user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.TmpAuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return s.generateTmpAuthResponse(user)
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

func (s *AuthService) Logout(userID uint) error {
	if err := s.userRepo.UnverifyUser(userID); err != nil {
		return err
	}
	return s.userRepo.DeleteRefreshTokenByUserID(userID)
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
	tmpLink := fmt.Sprintf("http://%s/api/v1/auth/verify-email/%s", baseAddr, tmpTokenStr)
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
	user, err := s.userRepo.GetByTmpToken(token)
	if err != nil {
		return err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		return err
	}
	return s.userRepo.ConfirmUserEmail(user)
}

func (s *AuthService) VerifyTOTP(req *dto.TotpRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByTmpToken(req.Token)
	if err != nil {
		return nil, err
	}
	plainSecret, err := utils.DecryptKey(user.TotpSecret, s.config.JWT.MfaEncKey)
	if err != nil {
		return nil, err
	}
	if !utils.ValidateTOTP(req.TotpCode, plainSecret) {
		return nil, errors.New("totp validation failed")
	}
	if !user.IsMfaEnabled {
		user.IsMfaEnabled = true
	}
	user.IsVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return s.generateAuthResponse(user)
}

func (s *AuthService) UseBackup(req *dto.BackupRequest) (*dto.AddMFAResponse, error) {
	user, err := s.userRepo.GetByTmpToken(req.Token)
	if err != nil {
		return nil, err
	}
	mfaBackup, err := s.userRepo.GetMfaBackup(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.checkBackupCode(mfaBackup.BackupCode, req.BackupCode); err != nil {
		return nil, err
	}
	mfaBackup.IsUsed = true
	err = s.userRepo.UpdateMfaBackup(mfaBackup)
	if err != nil {
		return nil, err
	}
	return s.generateAddMFAUserResponse(user)
}

func (s *AuthService) generateTmpTokenStr() (string, error) {
	token, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (s *AuthService) ResubmitVerificationLink(token string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByTmpToken(token)
	if err != nil {
		return nil, err
	}
	newToken, err := s.generateTmpTokenStr()
	if err != nil {
		return nil, err
	}
	return s.generateUnverifiedUserResponse(newToken, user)
}

func (s *AuthService) AddMFA(user *models.User) (*dto.AddMFAResponse, error) {
	return s.generateAddMFAUserResponse(user)
}

func (s *AuthService) generateAddMFAUserResponse(user *models.User) (*dto.AddMFAResponse, error) {
	otpKey, err := utils.GenerateMFAKey(s.config.JWT.MfaIssuer, user.Email)
	if err != nil {
		return nil, err
	}
	backupCode, hashedCode, err := utils.GenerateBackupCode()
	if err != nil {
		return nil, err
	}
	encTOTPSecret, err := utils.EncryptKey(otpKey.Secret(), s.config.JWT.MfaEncKey)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.AddUserMFA(user, encTOTPSecret, hashedCode); err != nil {
		return nil, err
	}
	return &dto.AddMFAResponse{
		TOTPSecret: otpKey.Secret(),
		TOTPLink:   otpKey.URL(),
		BackupCode: backupCode,
	}, nil
}

func (s *AuthService) checkBackupCode(backupCode, inputCode string) error {
	return utils.CheckBackupCode(backupCode, inputCode)
}

func (s *AuthService) generateTmpAuthResponse(user *models.User) (*dto.TmpAuthResponse, error) {
	tmpTokenStr, err := s.generateTmpTokenStr()
	if err != nil {
		return nil, err
	}
	unverifiedUser := &models.UnverifiedUser{
		UserID: user.ID,
		Email:  user.Email,
	}
	err = s.userRepo.SaveUnverifiedUser(tmpTokenStr, unverifiedUser)
	if err != nil {
		return nil, err
	}

	authResponse := &dto.TmpAuthResponse{
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
		TmpToken: tmpTokenStr,
	}
	return authResponse, nil
}

func (s *AuthService) GetUnverifiedUser(token string) (*models.User, error) {
	return s.userRepo.GetByTmpToken(token)
}
