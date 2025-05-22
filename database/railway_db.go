package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// NewRailwayDBConfig creates a new DBConfig for Railway deployment
// using the environment variables provided by Railway
func NewRailwayDBConfig() *DBConfig {
	// Print all environment variables for debugging
	log.Println("Railway environment variables:")
	log.Printf("RAILWAY_ENVIRONMENT=%s", os.Getenv("RAILWAY_ENVIRONMENT"))
	log.Printf("PGHOST=%s", os.Getenv("PGHOST"))
	log.Printf("PGPORT=%s", os.Getenv("PGPORT"))
	log.Printf("PGUSER=%s", os.Getenv("PGUSER"))
	log.Printf("PGDATABASE=%s", os.Getenv("PGDATABASE"))
	log.Printf("DATABASE_URL=%s", maskPassword(os.Getenv("DATABASE_URL")))
	
	// Check if we have a DATABASE_URL and try to use it first
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Println("Found DATABASE_URL, attempting to use it directly")
		return &DBConfig{
			ConnectionString: dbURL,
			UseDirectURL:     true,
		}
	}
	
	// Otherwise use individual environment variables
	host := os.Getenv("PGHOST")
	if host == "" {
		log.Println("WARNING: PGHOST is empty, falling back to default database configuration")
		return NewDBConfig() // Fall back to default configuration
	}
	
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")
	
	// Get port from environment or use default
	portStr := os.Getenv("PGPORT")
	port := 5432 // Default PostgreSQL port
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err == nil {
			port = p
		} else {
			log.Printf("Invalid PGPORT: %s, using default 5432", portStr)
		}
	}
	
	// Use SSL by default on Railway
	sslMode := "require"
	
	log.Printf("Railway database configuration: Host=%s, Port=%d, User=%s, DB=%s", 
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

// Function to mask password in connection string for logging
func maskPassword(connString string) string {
	if connString == "" {
		return ""
	}
	// Simple masking, doesn't handle all cases but good enough for logs
	return fmt.Sprintf("%s...masked...", connString[:10])
}
