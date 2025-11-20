package models

import (
	"time"

	"github.com/google/uuid"
)

// ConversationSummary represents a conversation overview between the current user and another user
type ConversationSummary struct {
	User                PublicUser `json:"user"`
	LastMessage         string     `json:"last_message,omitempty"`
	LastMessageTime     *time.Time `json:"last_message_time,omitempty"`
	LastMessageSenderID uuid.UUID  `json:"last_message_sender_id"`
	LastMessageIsRead   bool       `json:"last_message_is_read"`
	UnreadCount         int64      `json:"unread_count"`
}
