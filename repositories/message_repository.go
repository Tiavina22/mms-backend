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

// ConversationPartner represents another user in a conversation and the last message timestamp
type ConversationPartner struct {
	UserID      uuid.UUID
	LastMessage time.Time
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

// GetRecentConversations gets list of conversation partners ordered by last message time
func (r *MessageRepository) GetRecentConversations(userID uuid.UUID, limit int) ([]ConversationPartner, error) {
	var conversations []ConversationPartner

	err := r.db.Model(&models.Message{}).
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as user_id, MAX(created_at) as last_message", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("user_id").
		Order("last_message DESC").
		Limit(limit).
		Scan(&conversations).Error

	return conversations, err
}

// GetLastMessageBetween returns the most recent message between two users
func (r *MessageRepository) GetLastMessageBetween(userID1, userID2 uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Sender").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		Limit(1).
		First(&message).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &message, nil
}

// GetUnreadCountForConversation returns unread messages count between receiver and sender
func (r *MessageRepository) GetUnreadCountForConversation(receiverID, senderID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND is_read = ?", receiverID, senderID, false).
		Count(&count).Error
	return count, err
}

// UpdateContent updates message content and tracks previous content
func (r *MessageRepository) UpdateContent(messageID uuid.UUID, newEncryptedContent string, previousEncryptedContent string) error {
	return r.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"content":          newEncryptedContent,
			"previous_content": previousEncryptedContent,
			"edited":           true,
			"edited_at":        time.Now(),
		}).Error
}

// SoftDelete marks a message as deleted without removing it
func (r *MessageRepository) SoftDelete(messageID, userID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
			"deleted_by": userID,
		}).Error
}
