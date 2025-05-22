package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/user/footballsim/database"
	"github.com/user/footballsim/handlers"
	"github.com/user/footballsim/services"
)

func main() {
	// Initialize database connection - use Railway config if available
	var dbConfig *database.DBConfig
	
	// Check if running on Railway
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("PGHOST") != "" {
		log.Println("Railway environment detected, using Railway database configuration")
		dbConfig = database.NewRailwayDBConfig()
	} else {
		dbConfig = database.NewDBConfig()
	}
	
	db, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema and sample data
	err = database.InitDB(db)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialize repositories
	teamRepo := database.NewSQLTeamRepository(db)
	matchRepo := database.NewSQLMatchRepository(db)
	leagueRepo := database.NewSQLLeagueRepository(db)

	// Initialize services
	simulator := services.NewMatchSimulator(teamRepo, matchRepo, leagueRepo)
	predictor := services.NewTablePredictor(teamRepo, matchRepo, leagueRepo, simulator)

	// Initialize handlers
	teamHandler := handlers.NewTeamHandler(teamRepo)
	matchHandler := handlers.NewMatchHandler(matchRepo, teamRepo, simulator)
	leagueHandler := handlers.NewLeagueHandler(leagueRepo, teamRepo, matchRepo, predictor)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Football League Simulator",
	})

	// Add middleware
	app.Use(logger.New())
	
	// Configure CORS to allow requests from any origin
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	handlers.SetupRoutes(app, teamHandler, matchHandler, leagueHandler)

	// Serve static files
	app.Static("/", "./utils/static")

	// Default route
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Football League Simulator API - Use /api endpoints")
	})

	// Get port from environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Printf("Visit http://localhost:%s to view the application", port)
	log.Fatal(app.Listen(":" + port))
} 