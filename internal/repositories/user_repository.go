package repositories

import (
	"errors"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/cache"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	cache *cache.CacheProvider
}

func NewUserRepository(db *gorm.DB, rdb *cache.CacheProvider) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: rdb,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	var user models.User
	if err := r.db.
		Where("email = ? AND is_active = ?", email, isActive).
		First(&user).
		Error; err != nil {

		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(&user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(&user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(&token).Error
}

func (r *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where(
		"token = ? AND expires_at > ?", token, time.Now(),
	).Find(&refreshToken).Error; err != nil {
		return nil, errors.New("invalid token")
	}
	return &refreshToken, nil
}

func (r *UserRepository) GetRefreshTokenByUserID(userID uint) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where(
		"user_id = ? AND expires_at > ?", userID, time.Now(),
	).Find(&refreshToken).Error; err != nil {
		return nil, errors.New("invalid token")
	}
	return &refreshToken, nil
}

func (r *UserRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.RefreshToken{}).Error
}

func (r *UserRepository) DeleteRefreshTokenByUserID(id uint) error {
	return r.db.Where("user_id = ?", id).Delete(&models.RefreshToken{}).Error
}

func (r *UserRepository) SaveUnverifiedUser(token string, data *models.UnverifiedUser) error {
	return r.cache.SetTempToken(token, data)
}

func (r *UserRepository) ConfirmUserEmail(user *models.User) error {
	user.EmailConfirmed = true
	return r.db.Save(&user).Error
}

func (r *UserRepository) AddUserMFA(user *models.User, totpSecret, backupCode string) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	if err := tx.Delete(&models.UserMFABackup{}, map[string]any{"user_id": user.ID}).Error; err != nil {
		return err
	}
	user.TotpSecret = totpSecret
	if err := tx.Save(&user).Error; err != nil {
		return err
	}
	userBackup := models.UserMFABackup{
		UserID:     user.ID,
		BackupCode: backupCode,
	}
	if err := tx.Save(&userBackup).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (r *UserRepository) GetMfaBackup(userID uint) (*models.UserMFABackup, error) {
	var mfaBackup models.UserMFABackup
	err := r.db.Where("user_id = ? AND is_used = ?", userID, false).First(&mfaBackup).Error
	if err != nil {
		return nil, err
	}
	return &mfaBackup, nil
}

func (r *UserRepository) UpdateMfaBackup(mfaBackup *models.UserMFABackup) error {
	return r.db.Save(&mfaBackup).Error
}

func (r *UserRepository) GetByTmpToken(token string) (*models.User, error) {
	unverifiedUser, err := r.cache.GetUnverifiedUser(token)
	if err != nil {
		return nil, err
	}
	user, err := r.GetByID(unverifiedUser.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UnverifyUser(id uint) error {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}
	user.IsVerified = false
	return r.db.Save(&user).Error
}
