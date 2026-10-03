package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the current role for registered users
type UserRole string

// Logged in user is either a customer or an admin
const (
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
)

// User is a database table for registered users
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone" gorm:"uniqueIndex;not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Role      UserRole       `json:"role" gorm:"default:customer"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	RefreshTokens []RefreshToken `json:"-"`
	Orders        []Order        `json:"-"`
	Cart          Cart           `json:"-"`
}

// RefreshToken is a database table for authenticated users's JWT tokens
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"-"`
}
