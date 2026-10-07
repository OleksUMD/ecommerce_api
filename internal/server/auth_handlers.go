package server

import (
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

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

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Login failed", err)
		return
	}
	utils.SuccessResponse(c, "Login successful", response)
}

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

func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	response, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User  not found", err)
		return
	}
	utils.SuccessResponse(c, "Get profile success", response)
}

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
