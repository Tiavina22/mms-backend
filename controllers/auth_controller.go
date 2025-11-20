package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mms-backend/middleware"
	"mms-backend/services"
	"mms-backend/utils"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Signup handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.SignupRequest true "Signup Request"
// @Success 201 {object} services.AuthResponse
// @Router /auth/signup [post]
func (ctrl *AuthController) Signup(c *gin.Context) {
	var req services.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := ctrl.authService.Signup(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	lang := response.User.Language
	c.JSON(http.StatusCreated, gin.H{
		"message": utils.T(lang, "signup_success"),
		"data":    response,
	})
}

// Login handles user authentication
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login Request"
// @Success 200 {object} services.AuthResponse
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := ctrl.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	lang := response.User.Language
	c.JSON(http.StatusOK, gin.H{
		"message": utils.T(lang, "login_success"),
		"data":    response,
	})
}

// CheckUsername checks username availability
func (ctrl *AuthController) CheckUsername(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	available, err := ctrl.authService.CheckUsernameAvailability(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}

// CheckEmail checks email availability
func (ctrl *AuthController) CheckEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	available, err := ctrl.authService.CheckEmailAvailability(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}

// CheckPhone checks phone availability
func (ctrl *AuthController) CheckPhone(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	available, err := ctrl.authService.CheckPhoneAvailability(req.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}

// GetMe returns current user info
// @Summary Get current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.PublicUser
// @Router /auth/me [get]
func (ctrl *AuthController) GetMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// Get token from header and validate to get user info
	token := c.GetHeader("Authorization")
	token = token[len("Bearer "):]

	user, err := ctrl.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	// Verify userID matches
	if user.ID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token mismatch",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToPublicUser(),
	})
}
