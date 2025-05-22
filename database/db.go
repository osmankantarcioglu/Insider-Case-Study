package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// DBConfig contains the database configuration
type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	ConnectionString string
	UseDirectURL    bool
}

// GetConnectionString returns the connection string for the database
func (c *DBConfig) GetConnectionString() string {
	// If we have a direct connection string and UseDirectURL is true, use it
	if c.UseDirectURL && c.ConnectionString != "" {
		log.Println("Using direct database connection string")
		return c.ConnectionString
	}
	
	// Otherwise build the connection string from components
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	
	log.Printf("DB Connection string: host=%s port=%d dbname=%s (password hidden)", 
		c.Host, c.Port, c.DBName)
	
	return connStr
}

// NewDBConfig creates a new DBConfig from environment variables
func NewDBConfig() *DBConfig {
	host := getEnv("DB_HOST", "localhost")
	
	// Get port from environment or use default
	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid DB_PORT: %s, using default 5432", portStr)
		port = 5432
	}
	
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "footballsim")
	sslMode := getEnv("DB_SSLMODE", "disable")

	log.Printf("Database configuration: Host=%s, Port=%d, User=%s, DB=%s", 
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

// ConnectDB establishes a connection to the database
func ConnectDB(config *DBConfig) (*sql.DB, error) {
	if config.UseDirectURL {
		log.Println("Connecting using direct database URL...")
		// For direct URLs, we use sql.Open with the URL string directly
		db, err := sql.Open("postgres", config.ConnectionString)
		if err != nil {
			log.Printf("Error opening database connection with URL: %v", err)
			return nil, err
		}
		
		log.Println("Pinging database to verify connection...")
		err = db.Ping()
		if err != nil {
			log.Printf("Error pinging database: %v", err)
			return nil, err
		}
		
		log.Println("Successfully connected to database")
		return db, nil
	}
	
	// Original component-based connection logic
	log.Printf("Attempting to connect to database at %s:%d...", config.Host, config.Port)
	
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	log.Println("Pinging database to verify connection...")
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// getEnv gets an environment variable value with a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// InitDB initializes the database with the schema and sample data
func InitDB(db *sql.DB) error {
	log.Println("Initializing database schema...")
	
	// Read the schema file
	schemaBytes, err := os.ReadFile("database/sql_schema.sql")
	if err != nil {
		log.Printf("Error reading schema file: %v", err)
		return err
	}

	log.Println("Executing schema SQL...")
	// Execute the schema
	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		log.Printf("Error executing schema SQL: %v", err)
		return err
	}

	log.Println("Database initialized successfully")
	return nil
} 