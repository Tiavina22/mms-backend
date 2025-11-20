package routes

import (
	"github.com/gin-gonic/gin"
	"mms-backend/controllers"
	"mms-backend/middleware"
	"mms-backend/websocket"
)

// SetupRoutes sets up all application routes
func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	messageController *controllers.MessageController,
	groupController *controllers.GroupController,
	wsHandler *websocket.Handler,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "MMS Backend is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authController.Signup)
			auth.POST("/login", authController.Login)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes
			protected.GET("/auth/me", authController.GetMe)

			// User routes
			users := protected.Group("/users")
			{
				users.GET("", userController.ListUsers)
				users.GET("/search", userController.SearchUsers)
				users.GET("/:user_id", userController.GetUser)
			}

			// Message routes
			messages := protected.Group("/messages")
			{
				messages.POST("", messageController.SendMessage)
				messages.GET("/conversations", messageController.GetRecentConversations)
				messages.GET("/conversation/:user_id", messageController.GetConversation)
				messages.PUT("/read/:user_id", messageController.MarkAsRead)
				messages.GET("/unread/count", messageController.GetUnreadCount)
			}

			// Group routes
			groups := protected.Group("/groups")
			{
				groups.POST("", groupController.CreateGroup)
				groups.GET("/my", groupController.GetUserGroups)
				groups.GET("/:group_id", groupController.GetGroup)
				groups.DELETE("/:group_id", groupController.DeleteGroup)
				groups.GET("/:group_id/messages", groupController.GetGroupMessages)
				groups.POST("/messages", groupController.SendGroupMessage)
				groups.GET("/:group_id/members", groupController.GetGroupMembers)
				groups.POST("/:group_id/members", groupController.AddGroupMember)
				groups.DELETE("/:group_id/members/:user_id", groupController.RemoveGroupMember)
			}

			// WebSocket route (protected)
			protected.GET("/ws", wsHandler.HandleWebSocket)
		}
	}
}
