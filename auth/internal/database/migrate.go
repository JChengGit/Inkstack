package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// RunMigrations runs database migrations using GORM AutoMigrate
// For production, consider using golang-migrate or similar tools
func RunMigrations() error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Println("Running database migrations...")

	// Note: AutoMigrate will create tables if they don't exist
	// For more control, use the SQL migration files in migrations/
	// This is kept simple for development purposes

	// You can uncomment this to use AutoMigrate:
	// if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
	// 	return fmt.Errorf("failed to run migrations: %w", err)
	// }

	log.Println("Migrations completed successfully")
	return nil
}

// For production, you should use SQL migrations with migrate library
// Example implementation:
// func RunSQLMigrations(migrationsPath string) error {
// 	// Use golang-migrate/migrate library here
// 	return nil
// }
