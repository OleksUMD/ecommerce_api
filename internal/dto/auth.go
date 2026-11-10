// Package dto contains data transfer objects through API
package dto

import "time"

// RegisterRequest send by user when the one is registering in the system
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}

// LoginRequest send by user when the one is logging into the system
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest send by the user when the one is want to refresh the access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse is received by the user when the one has successfully logged in
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

// UserResponse is received by the user when the one ask for its own info from the system
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProfileRequest request to change user's personal data
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
}

type RegisterResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

// AuthResponse is received by the user when the one has successfully logged in
type AddMFAResponse struct {
	TOTPSecret string `json:"totp_secret"`
	TOTPLink   string `json:"totp_link"`
	BackupCode string `json:"backup_code"`
}

type BackupRequest struct {
	Token      string `json:"token" binding:"required"`
	BackupCode string `json:"backup_code" binding:"required"`
}

type TotpRequest struct {
	Token    string `json:"token" binding:"required"`
	TotpCode string `json:"totp_code" binding:"required"`
}

type TmpAuthResponse struct {
	User     UserResponse `json:"user"`
	TmpToken string       `json:"tmp_token"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}
