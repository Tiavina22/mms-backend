package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupMessage represents a message in a group
type GroupMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	SenderID  uuid.UUID `gorm:"type:uuid;not null;index" json:"sender_id"`
	Content   string    `gorm:"type:text;not null" json:"content"` // Encrypted content
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Group  Group `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"group,omitempty"`
	Sender User  `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender,omitempty"`
}

// BeforeCreate hook to generate UUID before creating group message
func (gm *GroupMessage) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == uuid.Nil {
		gm.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for GroupMessage model
func (GroupMessage) TableName() string {
	return "group_messages"
}

// GroupMessageResponse is the structure returned to clients
type GroupMessageResponse struct {
	ID        uuid.UUID  `json:"id"`
	GroupID   uuid.UUID  `json:"group_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Content   string     `json:"content"` // Decrypted content
	CreatedAt time.Time  `json:"created_at"`
	Sender    PublicUser `json:"sender,omitempty"`
}

