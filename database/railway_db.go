package database

import (
	"log"
	"os"
	"strconv"
)

// NewRailwayDBConfig creates a new DBConfig for Railway deployment
// using the environment variables provided by Railway
func NewRailwayDBConfig() *DBConfig {
	host := os.Getenv("PGHOST")
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
