package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex:idx_username_lower;not null" json:"username"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"type:varchar(20);index" json:"phone"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // Never expose password in JSON
	Avatar    string    `gorm:"type:varchar(500)" json:"avatar"`
	Language  string    `gorm:"type:varchar(10);default:'en'" json:"language"` // e.g., 'en', 'fr', 'es'
	DeviceToken string  `gorm:"type:varchar(500)" json:"-"` // For push notifications
	Platform    string  `gorm:"type:varchar(20)" json:"-"` // 'ios', 'android'
	IsOnline    bool    `gorm:"default:false" json:"is_online"`
	LastSeen    *time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook to generate UUID before creating user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// PublicUser returns a user object safe for public viewing (without sensitive data)
type PublicUser struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	Language  string     `json:"language"`
	IsOnline  bool       `json:"is_online"`
	LastSeen  *time.Time `json:"last_seen"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToPublicUser converts User to PublicUser
func (u *User) ToPublicUser() PublicUser {
	return PublicUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Language:  u.Language,
		IsOnline:  u.IsOnline,
		LastSeen:  u.LastSeen,
		CreatedAt: u.CreatedAt,
	}
}
