package database

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// toFileURL converts a file path to a proper file:// URL that works cross-platform
func toFileURL(path string) string {
	// Convert to absolute path for file:// URLs
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Warning: Could not resolve absolute path for %s: %v", path, err)
		absPath = path
	}

	// Convert backslashes to forward slashes for URL format
	absPath = filepath.ToSlash(absPath)

	// For Windows: file://C:/path/to/migrations
	// For Unix: file:///path/to/migrations
	if runtime.GOOS == "windows" {
		return "file://" + absPath
	}

	// Unix paths already start with /, so we only need file://
	return "file://" + absPath
}

// RunMigrations runs all pending database migrations
func RunMigrations(migrationsPath string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	sourceURL := toFileURL(migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(migrationsPath string) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		toFileURL(migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Rollback one step
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(migrationsPath string) (uint, bool, error) {
	if db == nil {
		return 0, false, fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying database: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		toFileURL(migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// MigrateTo migrates to a specific version
func MigrateTo(migrationsPath string, version uint) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		toFileURL(migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Printf("Migrated to version %d successfully", version)
	return nil
}
