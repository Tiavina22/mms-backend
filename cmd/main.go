package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"mms-backend/config"
	"mms-backend/controllers"
	"mms-backend/middleware"
	"mms-backend/models"
	"mms-backend/repositories"
	"mms-backend/routes"
	"mms-backend/services"
	"mms-backend/utils"
	"mms-backend/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")

	// Load translations
	i18n := utils.GetI18n()
	log.Printf("Loaded translations for languages: %v", i18n.SupportedLanguages())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	groupMessageRepo := repositories.NewGroupMessageRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)

	// Initialize WebSocket hub first (needed by services)
	hub := websocket.NewHub()
	go hub.Run()
	wsHandler := websocket.NewHandler(hub)

	// Initialize services
	pushService := services.NewPushService(cfg)
	authService := services.NewAuthService(userRepo)
	messageService := services.NewMessageService(messageRepo, userRepo, notificationRepo, pushService, hub)
	groupService := services.NewGroupService(groupRepo, groupMessageRepo, userRepo, notificationRepo, pushService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userRepo)
	messageController := controllers.NewMessageController(messageService)
	groupController := controllers.NewGroupController(groupService)

	log.Println("WebSocket hub started")

	// Set up Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())

	// Set up routes
	routes.SetupRoutes(router, authController, userController, messageController, groupController, wsHandler)

	// Start server
	port := cfg.Server.Port
	log.Printf("Starting MMS Backend server on port %s", port)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("WebSocket endpoint: ws://localhost:%s/api/v1/ws", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runMigrations runs database migrations
func runMigrations(db interface{}) error {
	type Migrator interface {
		AutoMigrate(dst ...interface{}) error
	}

	migrator, ok := db.(Migrator)
	if !ok {
		return fmt.Errorf("database does not support migrations")
	}

	log.Println("Running database migrations...")

	// Auto-migrate all models
	if err := migrator.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Group{},
		&models.GroupMember{},
		&models.GroupMessage{},
		&models.Notification{},
	); err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
