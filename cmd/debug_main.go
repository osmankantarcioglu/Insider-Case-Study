package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/user/footballsim/database"
)

// SimpleDebugMain is a simplified version of main.go for testing
// Railway build and database connectivity
func SimpleDebugMain() {
	log.Println("Starting debug app...")
	
	// Initialize database connection
	var dbConfig *database.DBConfig
	
	// First explicitly check for DATABASE_URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Println("Found DATABASE_URL environment variable")
		dbConfig = &database.DBConfig{
			ConnectionString: dbURL,
			UseDirectURL:     true,
		}
	} else if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("PGHOST") != "" {
		log.Println("Using Railway individual environment variables")
		// Fall back to minimal implementation for simplicity
		dbConfig = database.GetRailwayDBConfig()
	} else {
		log.Println("Using local database configuration")
		dbConfig = database.NewDBConfig()
	}
	
	// Try to connect to database
	log.Println("Attempting database connection...")
	db, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return
	}
	defer db.Close()
	
	log.Println("Database connection successful!")
	
	// Just stay alive for a bit
	log.Println("Debug app running...")
	for i := 0; i < 10; i++ {
		fmt.Printf(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	
	log.Println("Debug app completed successfully")
} 