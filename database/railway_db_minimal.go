package database

import (
	"log"
	"os"
)

// GetRailwayDBConfig is a simplified version for testing
// Use this if the full implementation causes build issues
func GetRailwayDBConfig() *DBConfig {
	// Check if we're on Railway
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("PGHOST") != "" {
		log.Println("Railway environment detected, using Railway database configuration")
		
		// Use Railway's PostgreSQL environment variables
		config := &DBConfig{
			Host:     os.Getenv("PGHOST"),
			User:     os.Getenv("PGUSER"),
			Password: os.Getenv("PGPASSWORD"),
			DBName:   os.Getenv("PGDATABASE"),
			SSLMode:  "require",
		}
		
		// Only try to parse port if it's provided
		if os.Getenv("PGPORT") != "" {
			config.Port = 5432 // Default
		}
		
		return config
	}
	
	// Fallback to standard config
	return NewDBConfig()
} 