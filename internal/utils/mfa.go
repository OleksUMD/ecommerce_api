package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func GenerateMFAKey(issuer, userEmail string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userEmail,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}

func GenerateBackupCode() (string, string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	hexStr := hex.EncodeToString(bytes)
	hashedCode, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return hexStr, string(hashedCode), nil
}

func CheckBackupCode(hashedCode, inputCode string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(inputCode))
}

func EncryptKey(key, secret string) (string, error) {
	// NewCipher accepts up to 32 bytes
	block, err := aes.NewCipher([]byte(secret)[:32])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(key), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptKey(key, secret string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(secret)[:32])
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("invalid key")
	}

	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainKey, err := aesGCM.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainKey), nil
}
