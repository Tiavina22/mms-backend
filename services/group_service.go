package services

import (
	"errors"

	"github.com/google/uuid"
	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/utils"
)

// GroupService handles group business logic
type GroupService struct {
	groupRepo        *repositories.GroupRepository
	groupMessageRepo *repositories.GroupMessageRepository
	userRepo         *repositories.UserRepository
	notificationRepo *repositories.NotificationRepository
	pushService      *PushService
}

// NewGroupService creates a new group service
func NewGroupService(
	groupRepo *repositories.GroupRepository,
	groupMessageRepo *repositories.GroupMessageRepository,
	userRepo *repositories.UserRepository,
	notificationRepo *repositories.NotificationRepository,
	pushService *PushService,
) *GroupService {
	return &GroupService{
		groupRepo:        groupRepo,
		groupMessageRepo: groupMessageRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		pushService:      pushService,
	}
}

// CreateGroupRequest represents a group creation request
type CreateGroupRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Type        models.GroupType  `json:"type"`
	MemberIDs   []uuid.UUID       `json:"member_ids"`
}

// SendGroupMessageRequest represents a group message send request
type SendGroupMessageRequest struct {
	GroupID uuid.UUID `json:"group_id" binding:"required"`
	Content string    `json:"content" binding:"required"`
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(creatorID uuid.UUID, req CreateGroupRequest) (*models.Group, error) {
	// Validate type
	if req.Type != models.GroupTypePublic && req.Type != models.GroupTypePrivate {
		req.Type = models.GroupTypePrivate
	}

	// Create group
	group := &models.Group{
		Name:        utils.SanitizeString(req.Name),
		Description: utils.SanitizeString(req.Description),
		Type:        req.Type,
		CreatedBy:   creatorID,
	}

	if err := s.groupRepo.Create(group); err != nil {
		return nil, err
	}

	// Add creator as admin
	creatorMember := &models.GroupMember{
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    models.MemberRoleAdmin,
	}
	if err := s.groupRepo.AddMember(creatorMember); err != nil {
		return nil, err
	}

	// Add other members
	for _, memberID := range req.MemberIDs {
		if memberID == creatorID {
			continue // Skip creator, already added
		}
		member := &models.GroupMember{
			GroupID: group.ID,
			UserID:  memberID,
			Role:    models.MemberRoleMember,
		}
		_ = s.groupRepo.AddMember(member)

		// Send notification to invited members
		user, err := s.userRepo.FindByID(memberID)
		if err == nil {
			notification := &models.Notification{
				UserID:      memberID,
				Type:        models.NotificationTypeGroupInvite,
				Content:     utils.T(user.Language, "group_invite_notification", group.Name),
				ReferenceID: &group.ID,
			}
			_ = s.notificationRepo.Create(notification)
		}
	}

	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(groupID, userID uuid.UUID) (*models.Group, error) {
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return nil, err
	}

	// Check if user is a member (for private groups)
	if group.Type == models.GroupTypePrivate {
		isMember, err := s.groupRepo.IsMember(groupID, userID)
		if err != nil || !isMember {
			return nil, errors.New("not a member of this group")
		}
	}

	return group, nil
}

// SendGroupMessage sends a message to a group
func (s *GroupService) SendGroupMessage(senderID uuid.UUID, req SendGroupMessageRequest) (*models.GroupMessageResponse, error) {
	// Check if user is a member
	isMember, err := s.groupRepo.IsMember(req.GroupID, senderID)
	if err != nil || !isMember {
		return nil, errors.New("not a member of this group")
	}

	// Get group info
	group, err := s.groupRepo.FindByID(req.GroupID)
	if err != nil {
		return nil, errors.New("group not found")
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

	// Create group message
	message := &models.GroupMessage{
		GroupID:  req.GroupID,
		SenderID: senderID,
		Content:  encryptedContent,
	}

	if err := s.groupMessageRepo.Create(message); err != nil {
		return nil, err
	}

	// Notify group members
	members, err := s.groupRepo.GetGroupMembers(req.GroupID)
	if err == nil {
		notificationContent := req.Content
		if len(notificationContent) > 50 {
			notificationContent = notificationContent[:50] + "..."
		}

		for _, member := range members {
			if member.UserID == senderID {
				continue // Don't notify sender
			}

			user, err := s.userRepo.FindByID(member.UserID)
			if err != nil {
				continue
			}

			// Create notification
			notification := &models.Notification{
				UserID:      member.UserID,
				Type:        models.NotificationTypeGroupMessage,
				Content:     utils.T(user.Language, "new_group_message_notification", group.Name, sender.Username, notificationContent),
				ReferenceID: &message.ID,
			}
			_ = s.notificationRepo.Create(notification)

			// Send push notification
			if user.DeviceToken != "" {
				_ = s.pushService.SendGroupMessageNotification(user, group.Name, sender.Username, notificationContent)
			}
		}
	}

	return &models.GroupMessageResponse{
		ID:        message.ID,
		GroupID:   message.GroupID,
		SenderID:  message.SenderID,
		Content:   req.Content,
		CreatedAt: message.CreatedAt,
		Sender:    sender.ToPublicUser(),
	}, nil
}

// GetGroupMessages retrieves messages for a group
func (s *GroupService) GetGroupMessages(groupID, userID uuid.UUID, limit, offset int) ([]models.GroupMessageResponse, error) {
	// Check if user is a member
	isMember, err := s.groupRepo.IsMember(groupID, userID)
	if err != nil || !isMember {
		return nil, errors.New("not a member of this group")
	}

	messages, err := s.groupMessageRepo.GetGroupMessages(groupID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Decrypt messages
	responses := make([]models.GroupMessageResponse, 0, len(messages))
	for _, msg := range messages {
		decryptedContent, err := utils.Decrypt(msg.Content)
		if err != nil {
			decryptedContent = "[Encrypted]"
		}

		responses = append(responses, models.GroupMessageResponse{
			ID:        msg.ID,
			GroupID:   msg.GroupID,
			SenderID:  msg.SenderID,
			Content:   decryptedContent,
			CreatedAt: msg.CreatedAt,
			Sender:    msg.Sender.ToPublicUser(),
		})
	}

	return responses, nil
}

// AddMember adds a member to a group
func (s *GroupService) AddMember(groupID, userID, newMemberID uuid.UUID) error {
	// Check if requester is an admin
	isAdmin, err := s.groupRepo.IsAdmin(groupID, userID)
	if err != nil || !isAdmin {
		return errors.New("only admins can add members")
	}

	// Add member
	member := &models.GroupMember{
		GroupID: groupID,
		UserID:  newMemberID,
		Role:    models.MemberRoleMember,
	}

	return s.groupRepo.AddMember(member)
}

// RemoveMember removes a member from a group
func (s *GroupService) RemoveMember(groupID, userID, memberID uuid.UUID) error {
	// Check if requester is an admin
	isAdmin, err := s.groupRepo.IsAdmin(groupID, userID)
	if err != nil || !isAdmin {
		return errors.New("only admins can remove members")
	}

	return s.groupRepo.RemoveMember(groupID, memberID)
}

// GetUserGroups gets all groups a user belongs to
func (s *GroupService) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	return s.groupRepo.GetUserGroups(userID)
}

// DeleteGroup deletes a group
func (s *GroupService) DeleteGroup(groupID, userID uuid.UUID) error {
	group, err := s.groupRepo.FindByID(groupID)
	if err != nil {
		return err
	}

	// Only creator can delete group
	if group.CreatedBy != userID {
		return errors.New("only the creator can delete this group")
	}

	return s.groupRepo.Delete(groupID)
}

