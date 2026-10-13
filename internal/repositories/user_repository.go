package repositories

import (
	"errors"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	var err error
	err = r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
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
	var err error
	err = r.db.Where("email = ? AND is_active = ?", email, isActive).First(&user).Error
	return &user, err
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
	var err error
	return err
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

func (r *UserRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return r.db.Where("id = ?", id).Delete(&models.RefreshToken{}).Error
}
