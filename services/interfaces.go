package services

import "github.com/user/footballsim/models"

// TeamRepository defines the methods that any team repository must implement
type TeamRepository interface {
	GetAll() ([]*models.Team, error)
	GetByID(id int) (*models.Team, error)
	Create(team *models.Team) error
	Update(team *models.Team) error
	Delete(id int) error
}

// MatchRepository defines the methods that any match repository must implement
type MatchRepository interface {
	GetAll() ([]*models.Match, error)
	GetByID(id int) (*models.Match, error)
	GetByWeek(week int) ([]*models.Match, error)
	GetUnplayed() ([]*models.Match, error)
	Create(match *models.Match) error
	Update(match *models.Match) error
	Delete(id int) error
}

// LeagueRepository defines the methods that any league repository must implement
type LeagueRepository interface {
	GetCurrent() (*models.League, error)
	Create(league *models.League) error
	Update(league *models.League) error
	GetCurrentWeek() (int, error)
	GetTotalWeeks() (int, error)
	UpdateWeek(week int) error
	MarkAsCompleted() error
}

// Simulator defines the methods that any match simulator must implement
type Simulator interface {
	SimulateMatch(homeTeam, awayTeam *models.Team) (*models.Match, error)
	SimulateWeek(week int) ([]*models.Match, error)
	SimulateRemaining() ([]*models.Match, error)
}

// Predictor defines the methods that any predictor must implement
type Predictor interface {
	PredictFinalTable() ([]*models.TeamStats, error)
} 