package database

import (
	"database/sql"
	"fmt"
)

// HealthStatus represents the health status of the database
type HealthStatus struct {
	Status          string `json:"status"`
	Error           string `json:"error,omitempty"`
	OpenConnections int    `json:"open_connections"`
	IdleConnections int    `json:"idle_connections"`
	InUse           int    `json:"in_use"`
	WaitCount       int64  `json:"wait_count"`
}

// CheckHealth performs a health check on the database connection
func CheckHealth() (*HealthStatus, error) {
	status := &HealthStatus{
		Status: "disconnected",
	}

	if db == nil {
		status.Error = "database not initialized"
		return status, fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		status.Error = err.Error()
		return status, fmt.Errorf("failed to get underlying database: %w", err)
	}

	// Ping the database
	if err := sqlDB.Ping(); err != nil {
		status.Error = err.Error()
		return status, fmt.Errorf("failed to ping database: %w", err)
	}

	// Get connection pool stats
	stats := sqlDB.Stats()
	status.Status = "connected"
	status.OpenConnections = stats.OpenConnections
	status.IdleConnections = stats.Idle
	status.InUse = stats.InUse
	status.WaitCount = stats.WaitCount

	return status, nil
}

// IsHealthy returns true if the database connection is healthy
func IsHealthy() bool {
	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		return false
	}

	// Check if there are any open connections available
	stats := sqlDB.Stats()
	if stats.OpenConnections == 0 {
		return false
	}

	return true
}

// GetStats returns the current database connection pool statistics
func GetStats() sql.DBStats {
	if db == nil {
		return sql.DBStats{}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}
	}

	return sqlDB.Stats()
}