package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message represents a direct message between two users
type Message struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SenderID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"sender_id"`
	ReceiverID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"receiver_id"`
	Content         string     `gorm:"type:text;not null" json:"content"` // Encrypted content
	IsRead          bool       `gorm:"default:false" json:"is_read"`
	ReadAt          *time.Time `json:"read_at"`
	IsDeleted       bool       `gorm:"default:false" json:"is_deleted"`
	DeletedAt       *time.Time `json:"deleted_at"`
	DeletedBy       *uuid.UUID `gorm:"type:uuid" json:"deleted_by"`
	Edited          bool       `gorm:"default:false" json:"edited"`
	EditedAt        *time.Time `json:"edited_at"`
	PreviousContent string     `gorm:"type:text" json:"previous_content"` // Encrypted previous content
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationships
	Sender   User `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender,omitempty"`
	Receiver User `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"receiver,omitempty"`
}

// BeforeCreate hook to generate UUID before creating message
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Message model
func (Message) TableName() string {
	return "messages"
}

// MessageResponse is the structure returned to clients (with decrypted content)
type MessageResponse struct {
	ID              uuid.UUID  `json:"id"`
	SenderID        uuid.UUID  `json:"sender_id"`
	ReceiverID      uuid.UUID  `json:"receiver_id"`
	Content         string     `json:"content"` // Decrypted content
	IsRead          bool       `json:"is_read"`
	ReadAt          *time.Time `json:"read_at"`
	IsDeleted       bool       `json:"is_deleted"`
	DeletedAt       *time.Time `json:"deleted_at"`
	DeletedBy       *uuid.UUID `json:"deleted_by"`
	Edited          bool       `json:"edited"`
	EditedAt        *time.Time `json:"edited_at"`
	PreviousContent string     `json:"previous_content"`
	CreatedAt       time.Time  `json:"created_at"`
	Sender          PublicUser `json:"sender,omitempty"`
}
