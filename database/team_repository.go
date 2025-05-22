package database

import (
	"database/sql"

	"github.com/user/footballsim/models"
)

// SQLTeamRepository implements the TeamRepository interface
type SQLTeamRepository struct {
	DB *sql.DB
}

// NewSQLTeamRepository creates a new SQLTeamRepository
func NewSQLTeamRepository(db *sql.DB) *SQLTeamRepository {
	return &SQLTeamRepository{
		DB: db,
	}
}

// GetAll returns all teams
func (r *SQLTeamRepository) GetAll() ([]*models.Team, error) {
	query := `
		SELECT DISTINCT id, name, played, won, drawn, lost, goals_for, goals_against, goal_difference, points, strength
		FROM teams
		ORDER BY points DESC, goal_difference DESC, goals_for DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]*models.Team, 0)
	for rows.Next() {
		team := &models.Team{}
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Played,
			&team.Won,
			&team.Drawn,
			&team.Lost,
			&team.GoalsFor,
			&team.GoalsAgainst,
			&team.GoalDifference,
			&team.Points,
			&team.Strength,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// GetByID returns a team by ID
func (r *SQLTeamRepository) GetByID(id int) (*models.Team, error) {
	query := `
		SELECT id, name, played, won, drawn, lost, goals_for, goals_against, goal_difference, points, strength
		FROM teams
		WHERE id = $1`

	team := &models.Team{}
	err := r.DB.QueryRow(query, id).Scan(
		&team.ID,
		&team.Name,
		&team.Played,
		&team.Won,
		&team.Drawn,
		&team.Lost,
		&team.GoalsFor,
		&team.GoalsAgainst,
		&team.GoalDifference,
		&team.Points,
		&team.Strength,
	)
	if err != nil {
		return nil, err
	}

	return team, nil
}

// Create creates a new team
func (r *SQLTeamRepository) Create(team *models.Team) error {
	query := `
		INSERT INTO teams (name, played, won, drawn, lost, goals_for, goals_against, goal_difference, points, strength)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		team.Name,
		team.Played,
		team.Won,
		team.Drawn,
		team.Lost,
		team.GoalsFor,
		team.GoalsAgainst,
		team.GoalDifference,
		team.Points,
		team.Strength,
	).Scan(&team.ID)

	return err
}

// Update updates an existing team
func (r *SQLTeamRepository) Update(team *models.Team) error {
	query := `
		UPDATE teams
		SET name = $1,
			played = $2,
			won = $3,
			drawn = $4,
			lost = $5,
			goals_for = $6,
			goals_against = $7,
			goal_difference = $8,
			points = $9,
			strength = $10
		WHERE id = $11`

	_, err := r.DB.Exec(
		query,
		team.Name,
		team.Played,
		team.Won,
		team.Drawn,
		team.Lost,
		team.GoalsFor,
		team.GoalsAgainst,
		team.GoalDifference,
		team.Points,
		team.Strength,
		team.ID,
	)

	return err
}

// Delete deletes a team
func (r *SQLTeamRepository) Delete(id int) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
} 