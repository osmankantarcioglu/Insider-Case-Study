package database

import (
	"database/sql"

	"github.com/user/footballsim/models"
)

// SQLMatchRepository implements the MatchRepository interface
type SQLMatchRepository struct {
	DB *sql.DB
}

// NewSQLMatchRepository creates a new SQLMatchRepository
func NewSQLMatchRepository(db *sql.DB) *SQLMatchRepository {
	return &SQLMatchRepository{
		DB: db,
	}
}

// GetAll returns all matches
func (r *SQLMatchRepository) GetAll() ([]*models.Match, error) {
	query := `
		SELECT id, week, home_team_id, away_team_id, home_team_name, away_team_name, 
			   home_team_goals, away_team_goals, played, played_at, is_edited
		FROM matches
		ORDER BY week ASC, id ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		match := &models.Match{}
		var playedAt sql.NullTime
		var homeTeamGoals, awayTeamGoals sql.NullInt32

		err := rows.Scan(
			&match.ID,
			&match.Week,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeTeamName,
			&match.AwayTeamName,
			&homeTeamGoals,
			&awayTeamGoals,
			&match.Played,
			&playedAt,
			&match.IsEdited,
		)
		if err != nil {
			return nil, err
		}

		if homeTeamGoals.Valid {
			match.HomeTeamGoals = int(homeTeamGoals.Int32)
		}

		if awayTeamGoals.Valid {
			match.AwayTeamGoals = int(awayTeamGoals.Int32)
		}

		if playedAt.Valid {
			match.PlayedAt = playedAt.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetByID returns a match by ID
func (r *SQLMatchRepository) GetByID(id int) (*models.Match, error) {
	query := `
		SELECT id, week, home_team_id, away_team_id, home_team_name, away_team_name, 
		       home_team_goals, away_team_goals, played, played_at, is_edited
		FROM matches
		WHERE id = $1`

	match := &models.Match{}
	var playedAt sql.NullTime
	var homeTeamGoals, awayTeamGoals sql.NullInt32

	err := r.DB.QueryRow(query, id).Scan(
		&match.ID,
		&match.Week,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.HomeTeamName,
		&match.AwayTeamName,
		&homeTeamGoals,
		&awayTeamGoals,
		&match.Played,
		&playedAt,
		&match.IsEdited,
	)
	if err != nil {
		return nil, err
	}

	if homeTeamGoals.Valid {
		match.HomeTeamGoals = int(homeTeamGoals.Int32)
	}

	if awayTeamGoals.Valid {
		match.AwayTeamGoals = int(awayTeamGoals.Int32)
	}

	if playedAt.Valid {
		match.PlayedAt = playedAt.Time
	}

	return match, nil
}

// GetByWeek returns all matches for a specific week
func (r *SQLMatchRepository) GetByWeek(week int) ([]*models.Match, error) {
	query := `
		SELECT id, week, home_team_id, away_team_id, home_team_name, away_team_name, 
		       home_team_goals, away_team_goals, played, played_at, is_edited
		FROM matches
		WHERE week = $1
		ORDER BY id ASC`

	rows, err := r.DB.Query(query, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		match := &models.Match{}
		var playedAt sql.NullTime
		var homeTeamGoals, awayTeamGoals sql.NullInt32

		err := rows.Scan(
			&match.ID,
			&match.Week,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeTeamName,
			&match.AwayTeamName,
			&homeTeamGoals,
			&awayTeamGoals,
			&match.Played,
			&playedAt,
			&match.IsEdited,
		)
		if err != nil {
			return nil, err
		}

		if homeTeamGoals.Valid {
			match.HomeTeamGoals = int(homeTeamGoals.Int32)
		}

		if awayTeamGoals.Valid {
			match.AwayTeamGoals = int(awayTeamGoals.Int32)
		}

		if playedAt.Valid {
			match.PlayedAt = playedAt.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// GetUnplayed returns all unplayed matches
func (r *SQLMatchRepository) GetUnplayed() ([]*models.Match, error) {
	query := `
		SELECT id, week, home_team_id, away_team_id, home_team_name, away_team_name, 
		       home_team_goals, away_team_goals, played, played_at, is_edited
		FROM matches
		WHERE played = false
		ORDER BY week ASC, id ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]*models.Match, 0)
	for rows.Next() {
		match := &models.Match{}
		var playedAt sql.NullTime
		var homeTeamGoals, awayTeamGoals sql.NullInt32

		err := rows.Scan(
			&match.ID,
			&match.Week,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.HomeTeamName,
			&match.AwayTeamName,
			&homeTeamGoals,
			&awayTeamGoals,
			&match.Played,
			&playedAt,
			&match.IsEdited,
		)
		if err != nil {
			return nil, err
		}

		if homeTeamGoals.Valid {
			match.HomeTeamGoals = int(homeTeamGoals.Int32)
		}

		if awayTeamGoals.Valid {
			match.AwayTeamGoals = int(awayTeamGoals.Int32)
		}

		if playedAt.Valid {
			match.PlayedAt = playedAt.Time
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// Create creates a new match
func (r *SQLMatchRepository) Create(match *models.Match) error {
	query := `
		INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name, 
		                    home_team_goals, away_team_goals, played, played_at, is_edited)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	var playedAt sql.NullTime
	if !match.PlayedAt.IsZero() {
		playedAt = sql.NullTime{Time: match.PlayedAt, Valid: true}
	}

	err := r.DB.QueryRow(
		query,
		match.Week,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeTeamName,
		match.AwayTeamName,
		match.HomeTeamGoals,
		match.AwayTeamGoals,
		match.Played,
		playedAt,
		match.IsEdited,
	).Scan(&match.ID)

	return err
}

// Update updates an existing match
func (r *SQLMatchRepository) Update(match *models.Match) error {
	query := `
		UPDATE matches
		SET week = $1,
			home_team_id = $2,
			away_team_id = $3,
			home_team_name = $4,
			away_team_name = $5,
			home_team_goals = $6,
			away_team_goals = $7,
			played = $8,
			played_at = $9,
			is_edited = $10
		WHERE id = $11`

	var playedAt sql.NullTime
	if !match.PlayedAt.IsZero() {
		playedAt = sql.NullTime{Time: match.PlayedAt, Valid: true}
	}

	_, err := r.DB.Exec(
		query,
		match.Week,
		match.HomeTeamID,
		match.AwayTeamID,
		match.HomeTeamName,
		match.AwayTeamName,
		match.HomeTeamGoals,
		match.AwayTeamGoals,
		match.Played,
		playedAt,
		match.IsEdited,
		match.ID,
	)

	return err
}

// Delete deletes a match
func (r *SQLMatchRepository) Delete(id int) error {
	query := `DELETE FROM matches WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
} 