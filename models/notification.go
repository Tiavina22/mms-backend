package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationTypeMessage      NotificationType = "message"
	NotificationTypeGroupMessage NotificationType = "group_message"
	NotificationTypeGroupInvite  NotificationType = "group_invite"
	NotificationTypeSystem       NotificationType = "system"
)

// Notification represents a user notification
type Notification struct {
	ID         uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type       NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Content    string           `gorm:"type:text;not null" json:"content"` // Preview/summary of the notification
	ReferenceID *uuid.UUID      `gorm:"type:uuid" json:"reference_id"` // ID of related message/group
	ReadStatus bool             `gorm:"default:false" json:"read_status"`
	ReadAt     *time.Time       `json:"read_at"`
	CreatedAt  time.Time        `json:"created_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// BeforeCreate hook to generate UUID before creating notification
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

