package services

import (
	"github.com/google/uuid"
	"mms-backend/models"
	"mms-backend/repositories"
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	return s.notificationRepo.GetUserNotifications(userID, limit, offset)
}

// GetUnreadNotifications retrieves unread notifications for a user
func (s *NotificationService) GetUnreadNotifications(userID uuid.UUID) ([]models.Notification, error) {
	return s.notificationRepo.GetUnreadNotifications(userID)
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

// GetUnreadCount gets the count of unread notifications
func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(notificationID uuid.UUID) error {
	return s.notificationRepo.Delete(notificationID)
}

