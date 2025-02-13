package database

import (
	"fitnesshub/server/config"
	"fitnesshub/server/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global variable to hold the database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	// Form the connection string with actual configuration values
	var dsn string
	if config.AppConfig.DBHost != "" && config.AppConfig.DBUser == "" {
		// Если используется DATABASE_URL (Render)
		dsn = config.AppConfig.DBHost
	} else {
		// Если используется локальное подключение
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.AppConfig.DBHost,
			config.AppConfig.DBUser,
			config.AppConfig.DBPassword,
			config.AppConfig.DBName,
			config.AppConfig.DBPort)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Migrate the User table (without dropping it)
	MigrateTables()

	// Additional setup for DB (if needed)
	log.Println("Successfully connected to the database")
}

// MigrateTables performs migrations for necessary tables
func MigrateTables() {
	// Migrate the schema to create the User table (if it doesn't exist or needs an update)
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error creating User table: %v", err)
	}
	log.Println("User table has been migrated")

}

// GetDB returns the current database connection
func GetDB() *gorm.DB {
	return DB
}
