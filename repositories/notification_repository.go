package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mms-backend/models"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create creates a new notification
func (r *NotificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// FindByID finds a notification by ID
func (r *NotificationRepository) FindByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

// GetUserNotifications retrieves all notifications for a user
func (r *NotificationRepository) GetUserNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetUnreadNotifications retrieves unread notifications for a user
func (r *NotificationRepository) GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND read_status = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"read_status": true,
			"read_at":     gorm.Expr("NOW()"),
		}).Error
}

// MarkAllAsRead marks all notifications for a user as read
func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read_status = ?", userID, false).
		Updates(map[string]interface{}{
			"read_status": true,
			"read_at":     gorm.Expr("NOW()"),
		}).Error
}

// Delete deletes a notification
func (r *NotificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, id).Error
}

// DeleteUserNotifications deletes all notifications for a user
func (r *NotificationRepository) DeleteUserNotifications(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Notification{}).Error
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read_status = ?", userID, false).
		Count(&count).Error
	return count, err
}

