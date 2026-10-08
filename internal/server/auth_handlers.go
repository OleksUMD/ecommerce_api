package server

import (
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summery Register a new user
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=dto.AuthResponse} "User registered successfully"
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
	utils.CreatedResponse(c, "User registered successfully", response)
}

// @Summery User Login
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

// @Summery Refresh access token
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

// @Summery User Logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} utils.Response{data=dto.AuthResponse} "Logout successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Router /auth/logout [post]
func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	err := s.authService.Logout(req.RefreshToken)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to logout", err)
		return
	}
	utils.SuccessResponse(c, "Logout successful", nil)
}

// @Summery Get user profile
// @Description Get current authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.UserResponse} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Router /users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	response, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User  not found", err)
		return
	}
	utils.SuccessResponse(c, "Get profile success", response)
}

// @Summery Update user profile
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
	userID := c.GetUint("user_id")
	response, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile updated successfully", response)
}
