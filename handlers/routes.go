package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, teamHandler *TeamHandler, matchHandler *MatchHandler, leagueHandler *LeagueHandler) {
	// API group
	api := app.Group("/api")

	// Teams routes
	teams := api.Group("/teams")
	teams.Get("/", teamHandler.GetAllTeams)
	teams.Get("/:id", teamHandler.GetTeamByID)
	teams.Post("/", teamHandler.CreateTeam)
	teams.Put("/:id", teamHandler.UpdateTeam)
	teams.Delete("/:id", teamHandler.DeleteTeam)

	// Matches routes
	matches := api.Group("/matches")
	matches.Get("/", matchHandler.GetAllMatches)
	matches.Get("/week/:week", matchHandler.GetMatchesByWeek)
	matches.Post("/week/:week/simulate", matchHandler.SimulateWeek)
	matches.Post("/simulate-all", matchHandler.SimulateAllRemainingMatches)
	matches.Put("/:id", matchHandler.UpdateMatchResult)

	// League routes
	league := api.Group("/league")
	league.Get("/", leagueHandler.GetCurrentLeague)
	league.Get("/table", leagueHandler.GetLeagueTable)
	league.Get("/prediction", leagueHandler.GetPrediction)
	league.Post("/", leagueHandler.CreateLeague)
	league.Post("/reset", leagueHandler.ResetLeague)
} 