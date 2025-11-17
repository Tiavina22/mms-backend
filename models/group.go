package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupType defines the type of group
type GroupType string

const (
	GroupTypePublic  GroupType = "public"
	GroupTypePrivate GroupType = "private"
)

// Group represents a chat group
type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        GroupType `gorm:"type:varchar(20);not null;default:'private'" json:"type"`
	Avatar      string    `gorm:"type:varchar(500)" json:"avatar"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Creator User          `gorm:"foreignKey:CreatedBy;constraint:OnDelete:CASCADE" json:"creator,omitempty"`
	Members []GroupMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
}

// BeforeCreate hook to generate UUID before creating group
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Group model
func (Group) TableName() string {
	return "groups"
}

