package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mms-backend/middleware"
	"mms-backend/services"
)

// MessageController handles message endpoints
type MessageController struct {
	messageService *services.MessageService
}

// NewMessageController creates a new message controller
func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// SendMessage sends a message to another user
// @Summary Send a message
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.SendMessageRequest true "Message Request"
// @Success 201 {object} models.MessageResponse
// @Router /messages [post]
func (ctrl *MessageController) SendMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req services.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message, err := ctrl.messageService.SendMessage(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "message sent successfully",
		"data":    message,
	})
}

// GetConversation retrieves messages between current user and another user
// @Summary Get conversation
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.MessageResponse
// @Router /messages/conversation/{user_id} [get]
func (ctrl *MessageController) GetConversation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	otherUserIDStr := c.Param("user_id")
	otherUserID, err := uuid.Parse(otherUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := ctrl.messageService.GetConversation(userID, otherUserID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}

// MarkAsRead marks messages in a conversation as read
// @Summary Mark messages as read
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "Sender User ID"
// @Success 200 {object} map[string]string
// @Router /messages/read/{user_id} [put]
func (ctrl *MessageController) MarkAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	senderIDStr := c.Param("user_id")
	senderID, err := uuid.Parse(senderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if err := ctrl.messageService.MarkAsRead(userID, senderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "marked as read",
	})
}

// GetUnreadCount gets unread message count
// @Summary Get unread message count
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int64
// @Router /messages/unread/count [get]
func (ctrl *MessageController) GetUnreadCount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	count, err := ctrl.messageService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// GetRecentConversations gets recent conversations
// @Summary Get recent conversations
// @Tags messages
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} models.PublicUser
// @Router /messages/conversations [get]
func (ctrl *MessageController) GetRecentConversations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, err := ctrl.messageService.GetRecentConversations(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// EditMessage updates a message content
func (ctrl *MessageController) EditMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	messageID, err := uuid.Parse(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}

	var req services.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message, err := ctrl.messageService.EditMessage(messageID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "message updated",
		"data":    message,
	})
}

// DeleteMessage marks a message as deleted
func (ctrl *MessageController) DeleteMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	messageID, err := uuid.Parse(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid message id",
		})
		return
	}

	message, err := ctrl.messageService.DeleteMessage(messageID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "message deleted",
		"data":    message,
	})
}
