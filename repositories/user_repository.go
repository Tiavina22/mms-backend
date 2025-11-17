package repositories

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mms-backend/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username (case-insensitive)
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("LOWER(username) = LOWER(?)", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List returns a paginated list of users
func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Search searches users by username or email
func (r *UserRepository) Search(query string, limit int) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + strings.ToLower(query) + "%"
	err := r.db.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern).
		Limit(limit).
		Find(&users).Error
	return users, err
}

// UpdateOnlineStatus updates user's online status
func (r *UserRepository) UpdateOnlineStatus(userID uuid.UUID, isOnline bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_online": isOnline,
			"last_seen": gorm.Expr("NOW()"),
		}).Error
}

// UpdateDeviceToken updates user's device token for push notifications
func (r *UserRepository) UpdateDeviceToken(userID uuid.UUID, deviceToken, platform string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"device_token": deviceToken,
			"platform":     platform,
		}).Error
}

