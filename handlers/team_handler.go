package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/user/footballsim/models"
	"github.com/user/footballsim/services"
)

// TeamHandler handles team related requests
type TeamHandler struct {
	TeamRepo services.TeamRepository
}

// NewTeamHandler creates a new TeamHandler
func NewTeamHandler(teamRepo services.TeamRepository) *TeamHandler {
	return &TeamHandler{
		TeamRepo: teamRepo,
	}
}

// GetAllTeams returns all teams
func (h *TeamHandler) GetAllTeams(c *fiber.Ctx) error {
	teams, err := h.TeamRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(teams)
}

// GetTeamByID returns a team by ID
func (h *TeamHandler) GetTeamByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	team, err := h.TeamRepo.GetByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return c.JSON(team)
}

// CreateTeam creates a new team
func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	team := new(models.Team)
	if err := c.BodyParser(team); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.TeamRepo.Create(team); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(team)
}

// UpdateTeam updates an existing team
func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	team := new(models.Team)
	if err := c.BodyParser(team); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	team.ID = id
	if err := h.TeamRepo.Update(team); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(team)
}

// DeleteTeam deletes a team
func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	if err := h.TeamRepo.Delete(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
} 