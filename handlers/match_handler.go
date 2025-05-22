package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/user/footballsim/services"
	"log"
)

// MatchHandler handles match related requests
type MatchHandler struct {
	MatchRepo services.MatchRepository
	TeamRepo  services.TeamRepository
	Simulator services.Simulator
}

// NewMatchHandler creates a new MatchHandler
func NewMatchHandler(matchRepo services.MatchRepository, teamRepo services.TeamRepository, simulator services.Simulator) *MatchHandler {
	return &MatchHandler{
		MatchRepo: matchRepo,
		TeamRepo:  teamRepo,
		Simulator: simulator,
	}
}

// GetAllMatches returns all matches
func (h *MatchHandler) GetAllMatches(c *fiber.Ctx) error {
	matches, err := h.MatchRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(matches)
}

// GetMatchesByWeek returns matches for a specific week
func (h *MatchHandler) GetMatchesByWeek(c *fiber.Ctx) error {
	week, err := strconv.Atoi(c.Params("week"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}

	// Debug logging
	log.Printf("Getting matches for week %d", week)
	
	matches, err := h.MatchRepo.GetByWeek(week)
	if err != nil {
		log.Printf("Error getting matches for week %d: %v", week, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("Found %d matches for week %d", len(matches), week)
	
	return c.JSON(fiber.Map{
		"week":    week,
		"matches": matches,
	})
}

// SimulateWeek simulates all matches for a specific week
func (h *MatchHandler) SimulateWeek(c *fiber.Ctx) error {
	week, err := strconv.Atoi(c.Params("week"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid week number",
		})
	}

	// Debug logging
	log.Printf("Simulating matches for week %d", week)
	
	playedMatches, err := h.Simulator.SimulateWeek(week)
	if err != nil {
		log.Printf("Error simulating week %d: %v", week, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("Successfully simulated %d matches for week %d", len(playedMatches), week)
	
	return c.JSON(fiber.Map{
		"week":    week,
		"matches": playedMatches,
	})
}

// SimulateAllRemainingMatches simulates all remaining matches in the league
func (h *MatchHandler) SimulateAllRemainingMatches(c *fiber.Ctx) error {
	// Debug logging
	log.Printf("Simulating all remaining matches")
	
	playedMatches, err := h.Simulator.SimulateRemaining()
	if err != nil {
		log.Printf("Error simulating all remaining matches: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("Successfully simulated %d remaining matches", len(playedMatches))
	
	return c.JSON(fiber.Map{
		"matches": playedMatches,
	})
}

// UpdateMatchResult updates the result of a match
func (h *MatchHandler) UpdateMatchResult(c *fiber.Ctx) error {
	matchID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid match ID",
		})
	}

	var updateData struct {
		HomeTeamGoals int `json:"home_team_goals"`
		AwayTeamGoals int `json:"away_team_goals"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get match
	match, err := h.MatchRepo.GetByID(matchID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Match not found",
		})
	}

	// Get teams
	homeTeam, err := h.TeamRepo.GetByID(match.HomeTeamID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Home team not found",
		})
	}

	awayTeam, err := h.TeamRepo.GetByID(match.AwayTeamID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Away team not found",
		})
	}

	// If match was already played, revert the old result
	if match.Played {
		// Revert home team stats
		homeTeam.Played--
		homeTeam.GoalsFor -= match.HomeTeamGoals
		homeTeam.GoalsAgainst -= match.AwayTeamGoals

		// Revert away team stats
		awayTeam.Played--
		awayTeam.GoalsFor -= match.AwayTeamGoals
		awayTeam.GoalsAgainst -= match.HomeTeamGoals

		// Revert wins, draws, losses
		if match.HomeTeamGoals > match.AwayTeamGoals {
			homeTeam.Won--
			awayTeam.Lost--
		} else if match.HomeTeamGoals < match.AwayTeamGoals {
			homeTeam.Lost--
			awayTeam.Won--
		} else {
			homeTeam.Drawn--
			awayTeam.Drawn--
		}
	}

	// Update match result
	match.HomeTeamGoals = updateData.HomeTeamGoals
	match.AwayTeamGoals = updateData.AwayTeamGoals
	match.Played = true
	match.IsEdited = true

	// Update teams with new result
	homeTeam.Played++
	homeTeam.GoalsFor += match.HomeTeamGoals
	homeTeam.GoalsAgainst += match.AwayTeamGoals

	awayTeam.Played++
	awayTeam.GoalsFor += match.AwayTeamGoals
	awayTeam.GoalsAgainst += match.HomeTeamGoals

	// Update wins, draws, losses
	if match.HomeTeamGoals > match.AwayTeamGoals {
		homeTeam.Won++
		awayTeam.Lost++
	} else if match.HomeTeamGoals < match.AwayTeamGoals {
		homeTeam.Lost++
		awayTeam.Won++
	} else {
		homeTeam.Drawn++
		awayTeam.Drawn++
	}

	// Update team stats
	homeTeam.UpdateStats()
	awayTeam.UpdateStats()

	// Save changes
	if err := h.MatchRepo.Update(match); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.TeamRepo.Update(homeTeam); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.TeamRepo.Update(awayTeam); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(match)
} 