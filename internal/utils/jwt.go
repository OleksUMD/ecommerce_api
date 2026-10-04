package utils

import (
	"errors"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims a custom implementation of JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokenPair creates jwt access/refresh token pair
func GenerateTokenPair(cfg *config.JWTConfig, userID uint, email, role string) (accessToken, refreshToken string, err error) {
	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := at.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshtokenExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := rt.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

// ValidateToken checks if token claims are valid according to the custom one
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(_ *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
