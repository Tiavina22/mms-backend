package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mms-backend/repositories"
)

// UserController handles user endpoints
type UserController struct {
	userRepo *repositories.UserRepository
}

// NewUserController creates a new user controller
func NewUserController(userRepo *repositories.UserRepository) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

// GetUser gets a user by ID
// @Summary Get a user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} models.PublicUser
// @Router /users/{user_id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToPublicUser(),
	})
}

// SearchUsers searches for users by username or email
// @Summary Search users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.PublicUser
// @Router /users/search [get]
func (ctrl *UserController) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "search query required",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := ctrl.userRepo.Search(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert to public users
	publicUsers := make([]interface{}, 0, len(users))
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublicUser())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": publicUsers,
	})
}

// ListUsers lists all users with pagination
// @Summary List users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.PublicUser
// @Router /users [get]
func (ctrl *UserController) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := ctrl.userRepo.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert to public users
	publicUsers := make([]interface{}, 0, len(users))
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublicUser())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": publicUsers,
	})
}

