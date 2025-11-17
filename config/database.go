package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := config.GetDSN()

	// Set up GORM config
	gormConfig := &gorm.Config{}
	
	// Enable detailed logging in development
	if AppConfig.Server.Environment == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully")
	DB = db
	return db, nil
}

