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
	
	// Create a http server on port 8080 that responds to any request with "OK"
	// Start this early so Railway's healthcheck passes even if DB connection fails
	go func() {
		log.Println("Starting HTTP server for healthcheck on port 8080")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK - App is running"))
		})
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	// Check for DATABASE_URL environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		log.Printf("DATABASE_URL found with length: %d", len(dbURL))
	} else {
		log.Println("DATABASE_URL not found")
	}
	
	// Get database config from environment
	dbConfig := getDBConfig()
	
	// Keep trying to connect to the database
	var db *sql.DB
	var err error
	maxRetries := 10
	
	for i := 0; i < maxRetries; i++ {
		log.Printf("Database connection attempt %d/%d", i+1, maxRetries)
		
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
		
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error opening DB connection: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		
		// Test the connection
		log.Println("Pinging database...")
		err = db.Ping()
		if err != nil {
			log.Printf("Error pinging database: %v", err)
			db.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		
		log.Println("Successfully connected to database!")
		break
	}
	
	if err != nil {
		log.Printf("Failed to connect to database after %d attempts: %v", maxRetries, err)
		log.Println("Application will continue running without database connection")
	} else {
		defer db.Close()
	}
	
	// Keep app running
	log.Println("App is running and will continue indefinitely")
	for {
		time.Sleep(30 * time.Second)
		log.Println("App is still alive")
		
		// Try pinging the database periodically if we're connected
		if db != nil {
			if err := db.Ping(); err != nil {
				log.Printf("Periodic database ping failed: %v", err)
			} else {
				log.Println("Database connection is healthy")
			}
		}
	}
}

// getDBConfig gets the database configuration from environment or defaults
func getDBConfig() *SimpleDBConfig {
	// Check for DATABASE_URL environment variable with case insensitivity
	foundDBURL := false
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.ToUpper(parts[0]) == "DATABASE_URL" && parts[1] != "" {
			log.Printf("Found DATABASE_URL as '%s' with value length: %d", parts[0], len(parts[1]))
			foundDBURL = true
			return &SimpleDBConfig{
				ConnectionString: parts[1],
				UseDirectURL:     true,
			}
		}
	}
	
	if !foundDBURL {
		log.Println("No DATABASE_URL found with any case, checking for Railway variables...")
		
		// RAILWAY DEPLOYMENT: Use hardcoded values for Railway deployment
		// These values match Railway's PostgreSQL service
		if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
			log.Println("Using hardcoded Railway PostgreSQL connection info")
			
			// Get password from environment variable - try multiple possible names
			password := ""
			passwordEnvVars := []string{"DB_PASSWORD", "PGPASSWORD", "POSTGRES_PASSWORD"}
			
			for _, envVar := range passwordEnvVars {
				if p := os.Getenv(envVar); p != "" {
					password = p
					log.Printf("Found password in %s environment variable", envVar)
					break
				}
			}
			
			// Try using the public URL from environment if available
			pgURL := os.Getenv("RAILWAY_SERVICE_POSTGRES_URL")
			if pgURL != "" {
				log.Printf("Found PostgreSQL service URL: %s", pgURL)
				
				// Try directly constructing the connection string
				if password == "" {
					// Use hardcoded password as last resort
					password = "hSKIkxqgCLEmPPKkDqxdqRfJdGTTzZfI" // FALLBACK: Using password provided in previous messages
					log.Println("Using hardcoded password as fallback")
				}
				
				// Use the public URL instead of internal hostname
				dbURL := fmt.Sprintf("postgresql://postgres:%s@%s:5432/railway?sslmode=disable", 
					password, pgURL)
				
				log.Println("Using constructed DATABASE_URL with public PostgreSQL endpoint")
				return &SimpleDBConfig{
					ConnectionString: dbURL,
					UseDirectURL:     true,
				}
			}
			
			// Fallback to standard approach
			if password == "" {
				password = "postgres" // Try default password
				log.Println("WARNING: No password environment variable found! Trying default 'postgres'")
			}
			
			// Try different hostnames
			hostnames := []string{
				"postgres.railway.internal",
				"postgres",
				os.Getenv("RAILWAY_SERVICE_POSTGRES_URL"),
			}
			
			for _, hostname := range hostnames {
				if hostname != "" {
					log.Printf("Trying PostgreSQL hostname: %s", hostname)
					
					config := &SimpleDBConfig{
						Host:     hostname,
						Port:     5432,
						User:     "postgres",
						Password: password,
						DBName:   "railway",
						SSLMode:  "disable",
					}
					
					// Test this connection
					connStr := config.GetConnectionString()
					db, err := sql.Open("postgres", connStr)
					if err == nil {
						err = db.Ping()
						if err == nil {
							log.Printf("Successfully connected using hostname: %s", hostname)
							db.Close()
							return config
						} else {
							log.Printf("Failed to ping with hostname %s: %v", hostname, err)
							db.Close()
						}
					} else {
						log.Printf("Failed to open connection with hostname %s: %v", hostname, err)
					}
				}
			}
			
			// If we get here, none of the hostnames worked
			log.Println("All hostname attempts failed, using 'postgres' as last resort")
			return &SimpleDBConfig{
				Host:     "postgres",
				Port:     5432,
				User:     "postgres",
				Password: password,
				DBName:   "railway",
				SSLMode:  "disable",
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
			
			return &SimpleDBConfig{
				Host:     host,
				Port:     5432,
				User:     user,
				Password: password,
				DBName:   dbName,
				SSLMode:  sslMode,
			}
		}
	}
	
	return dbConfig
}

// Helper to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
} 