package server

import (
	"errors"
	"fmt"

	"github.com/OleksUMD/ecommerce_api/internal/cache"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.UserResponse} "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request dtaa or user already exists"
// @Router /auth/register [post]
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}
	utils.CreatedResponse(
		c,
		fmt.Sprintf("Check the email(%s) to finish the registration", response.Email),
		response,
	)
}

// @Summary User Login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "User logged in successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Router /auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Login failed", err)

	}
	utils.SuccessResponse(c, "Login successful", response)
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Token refreshed successfully"
// @Failure 401 {object} utils.Response "Invalid refresh token"
// @Router /auth/refresh [post]
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Refresh token failed", err)
		return
	}
	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// @Summary User Logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=nil} "Logout successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Authorization required"
// @Router /users/logout [post]
func (s *Server) logout(c *gin.Context) {
	userID := c.GetUint(utils.UserIDKey)
	err := s.authService.Logout(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to logout", err)
		return
	}
	utils.SuccessResponse(c, "Logout successful", nil)
}

// @Summary Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint(utils.UserIDKey)
	response, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User  not found", err)
		return
	}
	utils.SuccessResponse(c, "Get profile success", response)
}

// @Summary Update user profile
// @Description Update current authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile's data to update"
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /users/profile [put]
func (s *Server) updateProfile(c *gin.Context) {

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	userID := c.GetUint(utils.UserIDKey)
	response, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", response)
}

// @Summary Verify an email during registration
// @Description Create a new user account with email and password
// @Tags Authentication
// @Param token path string true "Temporary verification token"
// @Success 200 {object} utils.Response{data=nil} "User's email verified successfully"
// @Failure 404 {object} utils.Response "Invalid or expired verification link"
// @Router /auth/verify-email/{token} [get]
func (s *Server) verifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		s.logger.Warn().Msg("invalid attempt to verify an email")
		utils.NotFoundResponse(c, "Not found", nil)
		return
	}
	err := s.authService.VerifyEmail(token)
	if err != nil {
		if !errors.Is(err, cache.ErrRedisNotFound) {
			s.logger.Err(err).Msg("invalid attempt to verify an email")
		}
		utils.NotFoundResponse(c, "Email verification failed due to expired link. Make request to resubmit another one.", nil)
		return
	}
	utils.SuccessResponse(c, "User's email verified successfully", nil)
}

// @Summary Generate new verification link and send it to the user's email
// @Description Send a temporary confirmation link to the user's email
// @Tags User
// @Produce json
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Verification link has been sent successfully"
// @Failure 404 {object} utils.Response "User not found"
// @Router /auth/verify-email/{token} [post]
func (s *Server) resendVerificationLink(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		s.logger.Warn().Msg("invalid attempt to verify an email")
		utils.NotFoundResponse(c, "Not found", nil)
		return
	}
	response, err := s.authService.ResubmitVerificationLink(token)
	if err != nil {
		utils.NotFoundResponse(c, "User not found", err)
		return
	}
	utils.SuccessResponse(
		c,
		fmt.Sprintf("Check the email(%s) to finish the registration", response.Email),
		response,
	)
}

// @Summary Get totp link to connect Authenticator
// @Description Get link and backup code for MFA device connection
// @Tags Authentication
// @Success 200 {object} utils.Response{data=dto.AddMFAResponse} "All parameters necessary to connect authenticator."
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Server side error"
// @Router /auth/add-mfa/{token} [get]
func (s *Server) addMFA(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		s.logger.Warn().Msg("invalid attempt to connect authenticator")
		utils.NotFoundResponse(c, "Not found", nil)
		return
	}
	user, err := s.authService.GetUnverifiedUser(token)
	if err != nil {
		s.logger.Warn().Msg("invalid attempt to connect authenticator")
		utils.NotFoundResponse(c, "Not found", nil)
		return
	}
	response, err := s.authService.AddMFA(user)
	if err != nil {
		s.logger.Err(err).Msg("failed to generate totp parameters")
		utils.InternalServerErrorResponse(c, "Failed to get TOTP parameters", nil)
		return
	}
	utils.SuccessResponse(c, "Multifactor authentication added successfully", response)
}

// @Summary Second step of authentication.
// @Description Verify user's identity with totp code and temporary token
// @Tags Authentication
// @Security BearerAuth
// @Param request body dto.TotpRequest true "Temporary AuthToken and TOTP code from Authenticator"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Access and refresh tokens"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Server side error"
// @Router /auth/verify-totp [post]
func (s *Server) verifyTOTP(c *gin.Context) {
	var req dto.TotpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.VerifyTOTP(&req)
	if err != nil {
		s.logger.Err(err).Msg("failed totp verification")
		utils.BadRequestResponse(c, "Invalid code", nil)
		return
	}
	utils.SuccessResponse(c, "Totp verification successful", response)
}

// @Summary Use backup code to recover user's MFA
// @Description Use backup code to connect new device as MFA
// @Tags Authentication
// @Param request body dto.BackupRequest true "User backup data"
// @Success 200 {object} utils.Response{data=dto.AddMFAResponse} "All parameters necessary to connect authenticator."
// @Failure 500 {object} utils.Response "Server side error"
// @Router /auth/backup-recovery [post]
func (s *Server) useBackup(c *gin.Context) {
	var req dto.BackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.UseBackup(&req)
	if err != nil {
		s.logger.Err(err).Msg("failed to use backup code")
		utils.InternalServerErrorResponse(c, "Failed to use backup code", nil)
		return
	}
	utils.SuccessResponse(c, "Account recovered successfully", response)
}

// @Summary Use backup code to recover user's MFA
// @Description Use backup code to connect new device as MFA
// @Tags Authentication
// @Param request body dto.BackupRequest true "User backup data"
// @Success 200 {object} utils.Response{data=dto.AddMFAResponse} "All parameters necessary to connect authenticator."
// @Failure 500 {object} utils.Response "Server side error"
// @Router /auth/reset-password [post]
func (s *Server) sendChangePasswordLink(c *gin.Context) {
	var req dto.ResetPasswordRequest
	err := s.authService.ResetPassword(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to reset the password", err)
		return
	}
	utils.SuccessResponse(
		c,
		fmt.Sprintf("Check the email(%s) to finish the password change flow", req.Email),
		nil,
	)
}

// @Summary Use backup code to recover user's MFA
// @Description Use backup code to connect new device as MFA
// @Tags Authentication
// @Param request body dto.BackupRequest true "User backup data"
// @Success 200 {object} utils.Response{data=dto.AddMFAResponse} "All parameters necessary to connect authenticator."
// @Failure 500 {object} utils.Response "Server side error"
// @Router /auth/change-password [post]
func (s *Server) changePassword(c *gin.Context) {
	var req dto.BackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.UseBackup(&req)
	if err != nil {
		s.logger.Err(err).Msg("failed to use backup code")
		utils.InternalServerErrorResponse(c, "Failed to use backup code", nil)
		return
	}
	utils.SuccessResponse(c, "Account recovered successfully", response)
}

// @Summary Use backup code to recover user's MFA
// @Description Use backup code to connect new device as MFA
// @Tags Authentication
// @Param request body dto.BackupRequest true "User backup data"
// @Success 200 {object} utils.Response{data=dto.AddMFAResponse} "All parameters necessary to connect authenticator."
// @Failure 404 {object} utils.Response "Not Found"
// @Router /auth/confirm-password-change/{token} [get]
func (s *Server) confirmPasswordChange(c *gin.Context) {
	var req dto.BackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.UseBackup(&req)
	if err != nil {
		s.logger.Err(err).Msg("failed to use backup code")
		utils.InternalServerErrorResponse(c, "Failed to use backup code", nil)
		return
	}
	utils.SuccessResponse(c, "Account recovered successfully", response)
}
