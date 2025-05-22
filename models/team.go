package models

type Team struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Played        int    `json:"played" db:"played"`
	Won           int    `json:"won" db:"won"`
	Drawn         int    `json:"drawn" db:"drawn"`
	Lost          int    `json:"lost" db:"lost"`
	GoalsFor      int    `json:"goals_for" db:"goals_for"`
	GoalsAgainst  int    `json:"goals_against" db:"goals_against"`
	GoalDifference int    `json:"goal_difference" db:"goal_difference"`
	Points        int    `json:"points" db:"points"`
	Strength      int    `json:"strength" db:"strength"` // 1-10 scale to determine team's strength
}

// Calculate points based on Premier League rules
func (t *Team) CalculatePoints() {
	t.Points = (t.Won * 3) + t.Drawn
}

// Calculate goal difference
func (t *Team) CalculateGoalDifference() {
	t.GoalDifference = t.GoalsFor - t.GoalsAgainst
}

// Update team stats after a match
func (t *Team) UpdateStats() {
	t.CalculatePoints()
	t.CalculateGoalDifference()
}

// TeamStats represents a summary of team statistics
type TeamStats struct {
	TeamID        int    `json:"team_id" db:"team_id"`
	TeamName      string `json:"team_name" db:"team_name"`
	Played        int    `json:"played" db:"played"`
	Won           int    `json:"won" db:"won"`
	Drawn         int    `json:"drawn" db:"drawn"`
	Lost          int    `json:"lost" db:"lost"`
	GoalsFor      int    `json:"goals_for" db:"goals_for"`
	GoalsAgainst  int    `json:"goals_against" db:"goals_against"`
	GoalDifference int    `json:"goal_difference" db:"goal_difference"`
	Points        int    `json:"points" db:"points"`
} 