package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// NewRailwayDBConfig creates a DBConfig for Railway deployment
func NewRailwayDBConfig() *DBConfig {
	// Railway provides PostgreSQL connection details as environment variables
	// Check for Railway-specific environment variables first
	if os.Getenv("PGHOST") != "" && os.Getenv("PGUSER") != "" && os.Getenv("PGDATABASE") != "" {
		// Using Railway's environment variables
		host := os.Getenv("PGHOST")
		
		// Get port from environment or use default
		portStr := os.Getenv("PGPORT")
		if portStr == "" {
			portStr = "5432"
		}
		
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Invalid PGPORT: %s, using default 5432", portStr)
			port = 5432
		}
		
		user := os.Getenv("PGUSER")
		password := os.Getenv("PGPASSWORD")
		dbName := os.Getenv("PGDATABASE")
		sslMode := "require" // Railway typically requires SSL
		
		log.Printf("Railway database configuration detected: Host=%s, Port=%d, User=%s, DB=%s", 
			host, port, user, dbName)
		
		return &DBConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbName,
			SSLMode:  sslMode,
		}
	}
	
	// Fallback to standard environment variables
	return NewDBConfig()
} 