package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"mms-backend/models"
)

// GroupMessageRepository handles database operations for group messages
type GroupMessageRepository struct {
	db *gorm.DB
}

// NewGroupMessageRepository creates a new group message repository
func NewGroupMessageRepository(db *gorm.DB) *GroupMessageRepository {
	return &GroupMessageRepository{db: db}
}

// Create creates a new group message
func (r *GroupMessageRepository) Create(message *models.GroupMessage) error {
	return r.db.Create(message).Error
}

// FindByID finds a group message by ID
func (r *GroupMessageRepository) FindByID(id uuid.UUID) (*models.GroupMessage, error) {
	var message models.GroupMessage
	err := r.db.Preload("Sender").Preload("Group").Where("id = ?", id).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group message not found")
		}
		return nil, err
	}
	return &message, nil
}

// GetGroupMessages retrieves messages for a specific group
func (r *GroupMessageRepository) GetGroupMessages(groupID uuid.UUID, limit, offset int) ([]models.GroupMessage, error) {
	var messages []models.GroupMessage
	err := r.db.Preload("Sender").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// Delete deletes a group message
func (r *GroupMessageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.GroupMessage{}, id).Error
}

// DeleteGroupMessages deletes all messages in a group
func (r *GroupMessageRepository) DeleteGroupMessages(groupID uuid.UUID) error {
	return r.db.Where("group_id = ?", groupID).Delete(&models.GroupMessage{}).Error
}

// GetMessageCount returns the count of messages in a group
func (r *GroupMessageRepository) GetMessageCount(groupID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.GroupMessage{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
}

