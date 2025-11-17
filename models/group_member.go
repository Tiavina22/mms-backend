package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberRole defines the role of a group member
type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

// GroupMember represents a member in a group
type GroupMember struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	GroupID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"group_id"`
	UserID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Role     MemberRole `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	JoinedAt time.Time  `json:"joined_at"`
	
	// Relationships
	Group Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"group,omitempty"`
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

// BeforeCreate hook to generate UUID before creating group member
func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == uuid.Nil {
		gm.ID = uuid.New()
	}
	if gm.JoinedAt.IsZero() {
		gm.JoinedAt = time.Now()
	}
	return nil
}

// TableName specifies the table name for GroupMember model
func (GroupMember) TableName() string {
	return "group_members"
}

