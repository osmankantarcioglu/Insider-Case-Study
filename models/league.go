package models

type League struct {
	ID      int     `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Season  string  `json:"season" db:"season"`
	Teams   []*Team `json:"teams,omitempty"`
	Matches []*Match `json:"matches,omitempty"`
	CurrentWeek int `json:"current_week" db:"current_week"`
	TotalWeeks  int `json:"total_weeks" db:"total_weeks"`
	IsCompleted bool `json:"is_completed" db:"is_completed"`
}

// LeagueTable represents the current league standings
type LeagueTable struct {
	Teams        []*TeamStats `json:"teams"`
	CurrentWeek  int          `json:"current_week"`
	TotalWeeks   int          `json:"total_weeks"`
	IsCompleted  bool         `json:"is_completed"`
}

// WeeklyMatches represents all matches for a specific week
type WeeklyMatches struct {
	Week    int      `json:"week"`
	Matches []*Match `json:"matches"`
}

// PredictionTable represents the predicted final league standings
type PredictionTable struct {
	Teams []*TeamStats `json:"teams"`
} 