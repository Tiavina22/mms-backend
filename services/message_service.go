package services

import (
	"errors"

	"github.com/google/uuid"
	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/utils"
)

// MessageService handles message business logic
type MessageService struct {
	messageRepo      *repositories.MessageRepository
	userRepo         *repositories.UserRepository
	notificationRepo *repositories.NotificationRepository
	pushService      *PushService
}

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo *repositories.MessageRepository,
	userRepo *repositories.UserRepository,
	notificationRepo *repositories.NotificationRepository,
	pushService *PushService,
) *MessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		pushService:      pushService,
	}
}

// SendMessageRequest represents a message send request
type SendMessageRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Content    string    `json:"content" binding:"required"`
}

// SendMessage sends a message from one user to another
func (s *MessageService) SendMessage(senderID uuid.UUID, req SendMessageRequest) (*models.MessageResponse, error) {
	// Validate receiver exists
	receiver, err := s.userRepo.FindByID(req.ReceiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}

	// Get sender info
	sender, err := s.userRepo.FindByID(senderID)
	if err != nil {
		return nil, errors.New("sender not found")
	}

	// Encrypt message content
	encryptedContent, err := utils.Encrypt(req.Content)
	if err != nil {
		return nil, errors.New("failed to encrypt message")
	}

	// Create message
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    encryptedContent,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	// Create notification for receiver
	notificationContent := req.Content
	if len(notificationContent) > 50 {
		notificationContent = notificationContent[:50] + "..."
	}

	notification := &models.Notification{
		UserID:      req.ReceiverID,
		Type:        models.NotificationTypeMessage,
		Content:     utils.T(receiver.Language, "new_message_notification", sender.Username, notificationContent),
		ReferenceID: &message.ID,
	}
	_ = s.notificationRepo.Create(notification)

	// Send push notification
	if receiver.DeviceToken != "" {
		_ = s.pushService.SendMessageNotification(receiver, sender.Username, notificationContent)
	}

	// Return decrypted message response
	return &models.MessageResponse{
		ID:         message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    req.Content, // Original unencrypted content
		IsRead:     message.IsRead,
		CreatedAt:  message.CreatedAt,
		Sender:     sender.ToPublicUser(),
	}, nil
}

// GetConversation retrieves messages between two users
func (s *MessageService) GetConversation(userID1, userID2 uuid.UUID, limit, offset int) ([]models.MessageResponse, error) {
	messages, err := s.messageRepo.GetConversation(userID1, userID2, limit, offset)
	if err != nil {
		return nil, err
	}

	// Decrypt messages and convert to response format
	responses := make([]models.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		decryptedContent, err := utils.Decrypt(msg.Content)
		if err != nil {
			// If decryption fails, skip the message or use placeholder
			decryptedContent = "[Encrypted]"
		}

		responses = append(responses, models.MessageResponse{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    decryptedContent,
			IsRead:     msg.IsRead,
			ReadAt:     msg.ReadAt,
			CreatedAt:  msg.CreatedAt,
			Sender:     msg.Sender.ToPublicUser(),
		})
	}

	return responses, nil
}

// MarkAsRead marks a message or conversation as read
func (s *MessageService) MarkAsRead(receiverID, senderID uuid.UUID) error {
	return s.messageRepo.MarkConversationAsRead(receiverID, senderID)
}

// GetUnreadCount gets the count of unread messages for a user
func (s *MessageService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

// GetRecentConversations gets recent conversations for a user
func (s *MessageService) GetRecentConversations(userID uuid.UUID, limit int) ([]models.PublicUser, error) {
	users, err := s.messageRepo.GetRecentConversations(userID, limit)
	if err != nil {
		return nil, err
	}

	publicUsers := make([]models.PublicUser, 0, len(users))
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublicUser())
	}

	return publicUsers, nil
}

// DeleteMessage deletes a message
func (s *MessageService) DeleteMessage(messageID, userID uuid.UUID) error {
	// Verify the message belongs to the user (either sender or receiver)
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return err
	}

	if message.SenderID != userID && message.ReceiverID != userID {
		return errors.New("unauthorized to delete this message")
	}

	return s.messageRepo.Delete(messageID)
}

