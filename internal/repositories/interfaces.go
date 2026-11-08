package repositories

import "github.com/OleksUMD/ecommerce_api/internal/models"

type IUserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	CreateRefreshToken(token *models.RefreshToken) error
	UnverifyUser(id uint) error
	GetValidRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByID(id uint) error
	DeleteRefreshTokenByUserID(userID uint) error
	SaveUnverifiedUser(token string, data *models.UnverifiedUser) error
	ConfirmUserEmail(user *models.User) error
	AddUserMFA(user *models.User, totpSecret, backupCode string) error
	GetMfaBackup(userID uint) (*models.UserMFABackup, error)
	UpdateMfaBackup(mfaBackup *models.UserMFABackup) error
	GetByTmpToken(token string) (*models.User, error)
}

type ICartRepository interface {
	GetByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(id uint) error
}
