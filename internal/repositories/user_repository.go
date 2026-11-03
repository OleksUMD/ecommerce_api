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

func (r *UserRepository) SaveUnverifiedUser(token string, data *models.UnverifiedUser) error {
	return r.cache.SetTempToken(token, data)
}

func (r *UserRepository) VerifyUserEmail(token string) error {
	unverifiedUser, err := r.cache.GetUnverifiedUser(token)
	if err != nil {
		return err
	}
	user, err := r.GetByID(unverifiedUser.UserID)
	if err != nil {
		return err
	}
	user.IsVerified = true
	return r.db.Save(&user).Error
}
