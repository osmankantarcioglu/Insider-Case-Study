package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	
	// Print all environment variables for debugging
	log.Println("--- ENVIRONMENT VARIABLES ---")
	for _, env := range os.Environ() {
		// Don't print secrets directly to logs
		if strings.Contains(strings.ToLower(env), "password") || 
		   strings.Contains(strings.ToLower(env), "secret") ||
		   strings.Contains(strings.ToLower(env), "database_url") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				log.Printf("%s=<masked>", parts[0])
			}
		} else {
			log.Println(env)
		}
	}
	log.Println("--- END ENVIRONMENT VARIABLES ---")
	
	// Check for DATABASE_URL environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		log.Printf("DATABASE_URL found with length: %d", len(dbURL))
	} else {
		log.Println("DATABASE_URL not found")
	}
	
	// Get database config from environment
	var dbConfig *SimpleDBConfig
	
	// Check for DATABASE_URL environment variable with case insensitivity
	foundDBURL := false
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.ToUpper(parts[0]) == "DATABASE_URL" && parts[1] != "" {
			log.Printf("Found DATABASE_URL as '%s' with value length: %d", parts[0], len(parts[1]))
			foundDBURL = true
			dbConfig = &SimpleDBConfig{
				ConnectionString: parts[1],
				UseDirectURL:     true,
			}
			break
		}
	}
	
	if !foundDBURL {
		log.Println("No DATABASE_URL found with any case, checking for Railway variables...")
		
		// RAILWAY DEPLOYMENT: Use hardcoded values for Railway deployment
		// These values match Railway's PostgreSQL service
		if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
			log.Println("Using hardcoded Railway PostgreSQL connection info")
			
			// Get password from environment variable
			password := os.Getenv("DB_PASSWORD")
			if password == "" {
				password = "password" // Fallback, but this won't work
				log.Println("WARNING: DB_PASSWORD environment variable not set!")
			} else {
				log.Println("Found DB_PASSWORD in environment variables")
			}
			
			dbConfig = &SimpleDBConfig{
				Host:     "postgres", // Simplified hostname
				Port:     5432,
				User:     "postgres",
				Password: password,
				DBName:   "railway",
				SSLMode:  "prefer", // Try different SSL modes: prefer, require, disable
			}
		} else {
			// Local development fallback
			host := getEnv("PGHOST", "localhost")
			user := getEnv("PGUSER", "postgres")
			password := getEnv("PGPASSWORD", "postgres")
			dbName := getEnv("PGDATABASE", "footballsim")
			
			// Determine SSL mode based on environment
			var sslMode string
			if os.Getenv("PGHOST") != "" && os.Getenv("PGHOST") != "localhost" {
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
	}
	
	// Try to connect
	log.Println("Connecting to database...")
	connStr := dbConfig.GetConnectionString()
	
	// Mask the connection string for logging
	if dbConfig.UseDirectURL {
		log.Println("Using direct database URL (password masked)")
	} else {
		log.Printf("Connection: host=%s user=%s dbname=%s sslmode=%s", 
			dbConfig.Host, dbConfig.User, dbConfig.DBName, dbConfig.SSLMode)
	}
	
	// Create a http server on port 8080 that responds to any request with "OK"
	// to allow Railway to healthcheck the application
	go func() {
		log.Println("Starting HTTP server for healthcheck on port 8080")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	defer db.Close()
	
	// Test the connection
	log.Println("Pinging database...")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	
	log.Println("Successfully connected to database!")
	
	// Keep app running for a bit
	log.Println("App is running and will continue indefinitely")
	for {
		time.Sleep(10 * time.Second)
		log.Println("App is still alive")
	}
}

// Helper to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
} 