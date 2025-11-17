package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mms-backend/middleware"
	"mms-backend/services"
)

// GroupController handles group endpoints
type GroupController struct {
	groupService *services.GroupService
}

// NewGroupController creates a new group controller
func NewGroupController(groupService *services.GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

// CreateGroup creates a new group
// @Summary Create a group
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateGroupRequest true "Group Request"
// @Success 201 {object} models.Group
// @Router /groups [post]
func (ctrl *GroupController) CreateGroup(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req services.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	group, err := ctrl.groupService.CreateGroup(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "group created successfully",
		"data":    group,
	})
}

// GetGroup gets a group by ID
// @Summary Get a group
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param group_id path string true "Group ID"
// @Success 200 {object} models.Group
// @Router /groups/{group_id} [get]
func (ctrl *GroupController) GetGroup(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid group id",
		})
		return
	}

	group, err := ctrl.groupService.GetGroup(groupID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": group,
	})
}

// SendGroupMessage sends a message to a group
// @Summary Send a group message
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.SendGroupMessageRequest true "Message Request"
// @Success 201 {object} models.GroupMessageResponse
// @Router /groups/messages [post]
func (ctrl *GroupController) SendGroupMessage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req services.SendGroupMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message, err := ctrl.groupService.SendGroupMessage(userID, req)
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

// GetGroupMessages gets messages for a group
// @Summary Get group messages
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param group_id path string true "Group ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.GroupMessageResponse
// @Router /groups/{group_id}/messages [get]
func (ctrl *GroupController) GetGroupMessages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid group id",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := ctrl.groupService.GetGroupMessages(groupID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}

// GetUserGroups gets all groups for the current user
// @Summary Get user groups
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Group
// @Router /groups/my [get]
func (ctrl *GroupController) GetUserGroups(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	groups, err := ctrl.groupService.GetUserGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groups,
	})
}

// DeleteGroup deletes a group
// @Summary Delete a group
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param group_id path string true "Group ID"
// @Success 200 {object} map[string]string
// @Router /groups/{group_id} [delete]
func (ctrl *GroupController) DeleteGroup(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	groupIDStr := c.Param("group_id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid group id",
		})
		return
	}

	if err := ctrl.groupService.DeleteGroup(groupID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "group deleted successfully",
	})
}

