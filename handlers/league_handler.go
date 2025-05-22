package handlers

import (
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/user/footballsim/models"
	"github.com/user/footballsim/services"
)

// LeagueHandler handles league related requests
type LeagueHandler struct {
	LeagueRepo services.LeagueRepository
	TeamRepo   services.TeamRepository
	MatchRepo  services.MatchRepository
	Predictor  services.Predictor
}

// NewLeagueHandler creates a new LeagueHandler
func NewLeagueHandler(leagueRepo services.LeagueRepository, teamRepo services.TeamRepository, matchRepo services.MatchRepository, predictor services.Predictor) *LeagueHandler {
	return &LeagueHandler{
		LeagueRepo: leagueRepo,
		TeamRepo:   teamRepo,
		MatchRepo:  matchRepo,
		Predictor:  predictor,
	}
}

// GetCurrentLeague returns the current league
func (h *LeagueHandler) GetCurrentLeague(c *fiber.Ctx) error {
	league, err := h.LeagueRepo.GetCurrent()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(league)
}

// CreateLeague creates a new league
func (h *LeagueHandler) CreateLeague(c *fiber.Ctx) error {
	league := new(models.League)
	if err := c.BodyParser(league); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.LeagueRepo.Create(league); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(league)
}

// ResetLeague resets the current league to the beginning
func (h *LeagueHandler) ResetLeague(c *fiber.Ctx) error {
	// Get current league
	league, err := h.LeagueRepo.GetCurrent()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Reset all teams
	teams, err := h.TeamRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, team := range teams {
		team.Played = 0
		team.Won = 0
		team.Drawn = 0
		team.Lost = 0
		team.GoalsFor = 0
		team.GoalsAgainst = 0
		team.GoalDifference = 0
		team.Points = 0

		if err := h.TeamRepo.Update(team); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Reset all matches
	matches, err := h.MatchRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, match := range matches {
		match.HomeTeamGoals = 0
		match.AwayTeamGoals = 0
		match.Played = false
		match.IsEdited = false

		if err := h.MatchRepo.Update(match); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Reset league to week 1
	league.CurrentWeek = 1
	league.IsCompleted = false
	if err := h.LeagueRepo.Update(league); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "League reset successfully",
	})
}

// GetLeagueTable returns the current league table
func (h *LeagueHandler) GetLeagueTable(c *fiber.Ctx) error {
	// Get current league
	league, err := h.LeagueRepo.GetCurrent()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get all teams
	teams, err := h.TeamRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert to TeamStats
	teamStats := make([]*models.TeamStats, len(teams))
	for i, team := range teams {
		teamStats[i] = &models.TeamStats{
			TeamID:         team.ID,
			TeamName:       team.Name,
			Played:         team.Played,
			Won:            team.Won,
			Drawn:          team.Drawn,
			Lost:           team.Lost,
			GoalsFor:       team.GoalsFor,
			GoalsAgainst:   team.GoalsAgainst,
			GoalDifference: team.GoalDifference,
			Points:         team.Points,
		}
	}

	// Sort the teams by points, then goal difference, then goals for
	sort.Slice(teamStats, func(i, j int) bool {
		if teamStats[i].Points != teamStats[j].Points {
			return teamStats[i].Points > teamStats[j].Points
		}
		if teamStats[i].GoalDifference != teamStats[j].GoalDifference {
			return teamStats[i].GoalDifference > teamStats[j].GoalDifference
		}
		return teamStats[i].GoalsFor > teamStats[j].GoalsFor
	})

	leagueTable := &models.LeagueTable{
		Teams:       teamStats,
		CurrentWeek: league.CurrentWeek,
		TotalWeeks:  league.TotalWeeks,
		IsCompleted: league.IsCompleted,
	}

	return c.JSON(leagueTable)
}

// GetPrediction returns the predicted final league table after week 4
func (h *LeagueHandler) GetPrediction(c *fiber.Ctx) error {
	// Get current league
	league, err := h.LeagueRepo.GetCurrent()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if we're at least at week 4
	if league.CurrentWeek < 4 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Predictions are only available after week 4",
		})
	}

	// Get prediction
	predictedTable, err := h.Predictor.PredictFinalTable()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sort the teams by points, then goal difference, then goals for (just to be sure)
	sort.Slice(predictedTable, func(i, j int) bool {
		if predictedTable[i].Points != predictedTable[j].Points {
			return predictedTable[i].Points > predictedTable[j].Points
		}
		if predictedTable[i].GoalDifference != predictedTable[j].GoalDifference {
			return predictedTable[i].GoalDifference > predictedTable[j].GoalDifference
		}
		return predictedTable[i].GoalsFor > predictedTable[j].GoalsFor
	})

	// Always return with the expected structure: {"prediction": [...]}
	return c.JSON(fiber.Map{
		"prediction": predictedTable,
	})
} 