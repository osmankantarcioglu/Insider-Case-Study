package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// SimpleDBConfig contains minimal database configuration
type SimpleDBConfig struct {
	ConnectionString string
	UseDirectURL     bool
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	SSLMode          string
}

// GetConnectionString returns the database connection string
func (c *SimpleDBConfig) GetConnectionString() string {
	if c.UseDirectURL && c.ConnectionString != "" {
		return c.ConnectionString
	}
	
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Main function for the debug app
func main() {
	log.Println("Starting minimal debug app...")
	
	// Get database config from environment
	var dbConfig *SimpleDBConfig
	
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		log.Println("Using DATABASE_URL")
		dbConfig = &SimpleDBConfig{
			ConnectionString: dbURL,
			UseDirectURL:     true,
		}
	} else {
		// Use individual variables or defaults
		host := getEnv("PGHOST", "localhost")
		user := getEnv("PGUSER", "postgres")
		password := getEnv("PGPASSWORD", "postgres")
		dbName := getEnv("PGDATABASE", "footballsim")
		
		// Determine SSL mode based on environment
		var sslMode string
		if os.Getenv("PGHOST") != "" {
			sslMode = "require"  // Use SSL for remote connections
		} else {
			sslMode = "disable"  // Disable for local development
		}
		
		dbConfig = &SimpleDBConfig{
			Host:     host,
			Port:     5432,
			User:     user,
			Password: password,
			DBName:   dbName,
			SSLMode:  sslMode,
		}
	}
	
	// Try to connect
	log.Println("Connecting to database...")
	connStr := dbConfig.GetConnectionString()
	
	// Mask the connection string for logging
	if dbConfig.UseDirectURL {
		log.Println("Using direct database URL (password masked)")
	} else {
		log.Printf("Connection: host=%s user=%s dbname=%s", 
			dbConfig.Host, dbConfig.User, dbConfig.DBName)
	}
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer db.Close()
	
	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	
	log.Println("Successfully connected to database!")
	
	// Keep app running for a bit
	for i := 0; i < 30; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	
	log.Println("Debug app completed successfully")
}

// Helper to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
} 