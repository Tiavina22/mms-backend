package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mms-backend/models"
)

// MessageRepository handles database operations for messages
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// FindByID finds a message by ID
func (r *MessageRepository) FindByID(id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Sender").Preload("Receiver").Where("id = ?", id).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

// GetConversation retrieves messages between two users
func (r *MessageRepository) GetConversation(userID1, userID2 uuid.UUID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// GetUserMessages retrieves all messages for a user
func (r *MessageRepository) GetUserMessages(userID uuid.UUID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("Sender").Preload("Receiver").
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// MarkAsRead marks a message as read
func (r *MessageRepository) MarkAsRead(messageID uuid.UUID) error {
	return r.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

// MarkConversationAsRead marks all messages in a conversation as read
func (r *MessageRepository) MarkConversationAsRead(receiverID, senderID uuid.UUID) error {
	return r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = ?", receiverID, senderID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

// GetUnreadCount returns the count of unread messages for a user
func (r *MessageRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// Delete deletes a message
func (r *MessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Message{}, id).Error
}

// GetRecentConversations gets list of users with recent conversations
func (r *MessageRepository) GetRecentConversations(userID uuid.UUID, limit int) ([]models.User, error) {
	type ConversationUser struct {
		UserID      uuid.UUID
		LastMessage time.Time
	}
	
	var conversations []ConversationUser
	
	// Get user IDs with last message time
	err := r.db.Model(&models.Message{}).
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as user_id, MAX(created_at) as last_message", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("user_id").
		Order("last_message DESC").
		Limit(limit).
		Scan(&conversations).Error
	
	if err != nil {
		return nil, err
	}
	
	// Extract user IDs
	var userIDs []uuid.UUID
	for _, conv := range conversations {
		userIDs = append(userIDs, conv.UserID)
	}
	
	// Get user details
	var users []models.User
	if len(userIDs) > 0 {
		err = r.db.Where("id IN ?", userIDs).Find(&users).Error
	}
	
	return users, err
}

