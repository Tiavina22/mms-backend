package services

import (
	"errors"
	"time"

	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/utils"
	"mms-backend/websocket"

	"github.com/google/uuid"
)

// MessageService handles message business logic
type MessageService struct {
	messageRepo      *repositories.MessageRepository
	userRepo         *repositories.UserRepository
	notificationRepo *repositories.NotificationRepository
	pushService      *PushService
	wsHub            *websocket.Hub
}

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo *repositories.MessageRepository,
	userRepo *repositories.UserRepository,
	notificationRepo *repositories.NotificationRepository,
	pushService *PushService,
	wsHub *websocket.Hub,
) *MessageService {
	return &MessageService{
		messageRepo:      messageRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		pushService:      pushService,
		wsHub:            wsHub,
	}
}

// SendMessageRequest represents a message send request
type SendMessageRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Content    string    `json:"content" binding:"required"`
}

// EditMessageRequest represents a message edit request
type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
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
		ID:              message.ID,
		SenderID:        message.SenderID,
		ReceiverID:      message.ReceiverID,
		Content:         req.Content, // Original unencrypted content
		IsRead:          message.IsRead,
		CreatedAt:       message.CreatedAt,
		IsDeleted:       message.IsDeleted,
		Edited:          message.Edited,
		PreviousContent: "",
		Sender:          sender.ToPublicUser(),
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

		previousContent := ""
		if msg.PreviousContent != "" {
			if prev, err := utils.Decrypt(msg.PreviousContent); err == nil {
				previousContent = prev
			}
		}

		displayContent := decryptedContent
		if msg.IsDeleted {
			displayContent = "[message deleted]"
		}

		responses = append(responses, models.MessageResponse{
			ID:              msg.ID,
			SenderID:        msg.SenderID,
			ReceiverID:      msg.ReceiverID,
			Content:         displayContent,
			IsRead:          msg.IsRead,
			ReadAt:          msg.ReadAt,
			IsDeleted:       msg.IsDeleted,
			DeletedAt:       msg.DeletedAt,
			DeletedBy:       msg.DeletedBy,
			Edited:          msg.Edited,
			EditedAt:        msg.EditedAt,
			PreviousContent: previousContent,
			CreatedAt:       msg.CreatedAt,
			Sender:          msg.Sender.ToPublicUser(),
		})
	}

	return responses, nil
}

// MarkAsRead marks a message or conversation as read
func (s *MessageService) MarkAsRead(receiverID, senderID uuid.UUID) error {
	err := s.messageRepo.MarkConversationAsRead(receiverID, senderID)
	if err != nil {
		return err
	}

	// Notify sender via WebSocket that their messages have been read
	if s.wsHub != nil {
		readReceipt := websocket.Message{
			Type:       "message_read",
			SenderID:   receiverID, // The one who read the messages
			ReceiverID: senderID,   // The one who sent the messages (to notify)
			Timestamp:  time.Now(),
		}
		s.wsHub.SendToUser(senderID, &readReceipt)
	}

	return nil
}

// GetUnreadCount gets the count of unread messages for a user
func (s *MessageService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}

// GetRecentConversations gets recent conversations for a user
func (s *MessageService) GetRecentConversations(userID uuid.UUID, limit int) ([]models.ConversationSummary, error) {
	partners, err := s.messageRepo.GetRecentConversations(userID, limit)
	if err != nil {
		return nil, err
	}

	summaries := make([]models.ConversationSummary, 0, len(partners))

	for _, partner := range partners {
		user, err := s.userRepo.FindByID(partner.UserID)
		if err != nil {
			return nil, err
		}

		summary := models.ConversationSummary{
			User: user.ToPublicUser(),
		}

		lastMessage, err := s.messageRepo.GetLastMessageBetween(userID, partner.UserID)
		if err != nil {
			return nil, err
		}

		if lastMessage != nil {
			decryptedContent, err := utils.Decrypt(lastMessage.Content)
			if err != nil {
				decryptedContent = "[Encrypted]"
			}

			displayContent := decryptedContent
			if lastMessage.IsDeleted {
				displayContent = "[message deleted]"
			}

			summary.LastMessage = displayContent
			summary.LastMessageTime = &lastMessage.CreatedAt
			summary.LastMessageSenderID = lastMessage.SenderID
			summary.LastMessageIsRead = lastMessage.IsRead
		}

		unreadCount, err := s.messageRepo.GetUnreadCountForConversation(userID, partner.UserID)
		if err != nil {
			return nil, err
		}
		summary.UnreadCount = unreadCount

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// EditMessage updates the content of a message
func (s *MessageService) EditMessage(messageID, userID uuid.UUID, req EditMessageRequest) (*models.MessageResponse, error) {
	if req.Content == "" {
		return nil, errors.New("content cannot be empty")
	}

	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, err
	}

	if message.SenderID != userID {
		return nil, errors.New("unauthorized to edit this message")
	}

	if message.IsDeleted {
		return nil, errors.New("cannot edit a deleted message")
	}

	previousEncrypted := message.Content
	previousDecrypted, err := utils.Decrypt(previousEncrypted)
	if err != nil {
		previousDecrypted = "[Encrypted]"
	}

	newEncrypted, err := utils.Encrypt(req.Content)
	if err != nil {
		return nil, errors.New("failed to encrypt message")
	}

	if err := s.messageRepo.UpdateContent(messageID, newEncrypted, previousEncrypted); err != nil {
		return nil, err
	}

	now := time.Now()
	message.Content = newEncrypted
	message.PreviousContent = previousEncrypted
	message.Edited = true
	message.EditedAt = &now

	sender := message.Sender
	if sender.ID == uuid.Nil {
		if user, err := s.userRepo.FindByID(message.SenderID); err == nil {
			sender = *user
		}
	}

	return &models.MessageResponse{
		ID:              message.ID,
		SenderID:        message.SenderID,
		ReceiverID:      message.ReceiverID,
		Content:         req.Content,
		IsRead:          message.IsRead,
		ReadAt:          message.ReadAt,
		IsDeleted:       message.IsDeleted,
		DeletedAt:       message.DeletedAt,
		DeletedBy:       message.DeletedBy,
		Edited:          true,
		EditedAt:        &now,
		PreviousContent: previousDecrypted,
		CreatedAt:       message.CreatedAt,
		Sender:          sender.ToPublicUser(),
	}, nil
}

// DeleteMessage marks a message as deleted
func (s *MessageService) DeleteMessage(messageID, userID uuid.UUID) (*models.MessageResponse, error) {
	message, err := s.messageRepo.FindByID(messageID)
	if err != nil {
		return nil, err
	}

	if message.SenderID != userID {
		return nil, errors.New("unauthorized to delete this message")
	}

	if message.IsDeleted {
		return nil, errors.New("message already deleted")
	}

	if err := s.messageRepo.SoftDelete(messageID, userID); err != nil {
		return nil, err
	}

	now := time.Now()
	message.IsDeleted = true
	message.DeletedAt = &now
	message.DeletedBy = &userID

	previousDecrypted := ""
	sender := message.Sender
	if sender.ID == uuid.Nil {
		if user, err := s.userRepo.FindByID(message.SenderID); err == nil {
			sender = *user
		}
	}

	return &models.MessageResponse{
		ID:              message.ID,
		SenderID:        message.SenderID,
		ReceiverID:      message.ReceiverID,
		Content:         "[message deleted]",
		IsRead:          message.IsRead,
		ReadAt:          message.ReadAt,
		IsDeleted:       true,
		DeletedAt:       message.DeletedAt,
		DeletedBy:       message.DeletedBy,
		Edited:          message.Edited,
		EditedAt:        message.EditedAt,
		PreviousContent: previousDecrypted,
		CreatedAt:       message.CreatedAt,
		Sender:          sender.ToPublicUser(),
	}, nil
}
