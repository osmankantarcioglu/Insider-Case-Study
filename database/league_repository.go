package database

import (
	"database/sql"

	"github.com/user/footballsim/models"
)

// SQLLeagueRepository implements the LeagueRepository interface
type SQLLeagueRepository struct {
	DB *sql.DB
}

// NewSQLLeagueRepository creates a new SQLLeagueRepository
func NewSQLLeagueRepository(db *sql.DB) *SQLLeagueRepository {
	return &SQLLeagueRepository{
		DB: db,
	}
}

// GetCurrent returns the current league
func (r *SQLLeagueRepository) GetCurrent() (*models.League, error) {
	query := `
		SELECT id, name, season, current_week, total_weeks, is_completed
		FROM leagues
		ORDER BY id DESC
		LIMIT 1`

	league := &models.League{}
	err := r.DB.QueryRow(query).Scan(
		&league.ID,
		&league.Name,
		&league.Season,
		&league.CurrentWeek,
		&league.TotalWeeks,
		&league.IsCompleted,
	)
	if err != nil {
		return nil, err
	}

	return league, nil
}

// Create creates a new league
func (r *SQLLeagueRepository) Create(league *models.League) error {
	query := `
		INSERT INTO leagues (name, season, current_week, total_weeks, is_completed)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.DB.QueryRow(
		query,
		league.Name,
		league.Season,
		league.CurrentWeek,
		league.TotalWeeks,
		league.IsCompleted,
	).Scan(&league.ID)

	return err
}

// Update updates an existing league
func (r *SQLLeagueRepository) Update(league *models.League) error {
	query := `
		UPDATE leagues
		SET name = $1,
			season = $2,
			current_week = $3,
			total_weeks = $4,
			is_completed = $5
		WHERE id = $6`

	_, err := r.DB.Exec(
		query,
		league.Name,
		league.Season,
		league.CurrentWeek,
		league.TotalWeeks,
		league.IsCompleted,
		league.ID,
	)

	return err
}

// GetCurrentWeek returns the current week of the league
func (r *SQLLeagueRepository) GetCurrentWeek() (int, error) {
	query := `
		SELECT current_week
		FROM leagues
		ORDER BY id DESC
		LIMIT 1`

	var currentWeek int
	err := r.DB.QueryRow(query).Scan(&currentWeek)
	if err != nil {
		return 0, err
	}

	return currentWeek, nil
}

// GetTotalWeeks returns the total number of weeks in the league
func (r *SQLLeagueRepository) GetTotalWeeks() (int, error) {
	query := `
		SELECT total_weeks
		FROM leagues
		ORDER BY id DESC
		LIMIT 1`

	var totalWeeks int
	err := r.DB.QueryRow(query).Scan(&totalWeeks)
	if err != nil {
		return 0, err
	}

	return totalWeeks, nil
}

// UpdateWeek updates the current week of the league
func (r *SQLLeagueRepository) UpdateWeek(week int) error {
	query := `
		UPDATE leagues
		SET current_week = $1
		WHERE id = (
			SELECT id FROM leagues ORDER BY id DESC LIMIT 1
		)`

	_, err := r.DB.Exec(query, week)
	return err
}

// MarkAsCompleted marks the league as completed
func (r *SQLLeagueRepository) MarkAsCompleted() error {
	query := `
		UPDATE leagues
		SET is_completed = true
		WHERE id = (
			SELECT id FROM leagues ORDER BY id DESC LIMIT 1
		)`

	_, err := r.DB.Exec(query)
	return err
} 