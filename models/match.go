package models

import "time"

type Match struct {
	ID            int       `json:"id" db:"id"`
	Week          int       `json:"week" db:"week"`
	HomeTeamID    int       `json:"home_team_id" db:"home_team_id"`
	AwayTeamID    int       `json:"away_team_id" db:"away_team_id"`
	HomeTeamName  string    `json:"home_team_name" db:"home_team_name"`
	AwayTeamName  string    `json:"away_team_name" db:"away_team_name"`
	HomeTeamGoals int       `json:"home_team_goals" db:"home_team_goals"`
	AwayTeamGoals int       `json:"away_team_goals" db:"away_team_goals"`
	Played        bool      `json:"played" db:"played"`
	PlayedAt      time.Time `json:"played_at,omitempty" db:"played_at"`
	IsEdited      bool      `json:"is_edited" db:"is_edited"`
}

// MatchResult represents the result of a match
type MatchResult struct {
	Match      Match `json:"match"`
	HomePoints int   `json:"home_points"`
	AwayPoints int   `json:"away_points"`
}

// GetResult returns the points gained by each team
func (m *Match) GetResult() (homePoints, awayPoints int) {
	if m.HomeTeamGoals > m.AwayTeamGoals {
		homePoints = 3
		awayPoints = 0
	} else if m.HomeTeamGoals < m.AwayTeamGoals {
		homePoints = 0
		awayPoints = 3
	} else {
		homePoints = 1
		awayPoints = 1
	}
	return
}

// IsHomeWin returns true if the home team won
func (m *Match) IsHomeWin() bool {
	return m.HomeTeamGoals > m.AwayTeamGoals
}

// IsAwayWin returns true if the away team won
func (m *Match) IsAwayWin() bool {
	return m.HomeTeamGoals < m.AwayTeamGoals
}

// IsDraw returns true if the match ended in a draw
func (m *Match) IsDraw() bool {
	return m.HomeTeamGoals == m.AwayTeamGoals
} 