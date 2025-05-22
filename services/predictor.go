package services

import (
	"sort"

	"github.com/user/footballsim/models"
)

// TablePredictor implements the Predictor interface
type TablePredictor struct {
	TeamRepo   TeamRepository
	MatchRepo  MatchRepository
	LeagueRepo LeagueRepository
	Simulator  Simulator
}

// NewTablePredictor creates a new table predictor
func NewTablePredictor(teamRepo TeamRepository, matchRepo MatchRepository, leagueRepo LeagueRepository, simulator Simulator) *TablePredictor {
	return &TablePredictor{
		TeamRepo:   teamRepo,
		MatchRepo:  matchRepo,
		LeagueRepo: leagueRepo,
		Simulator:  simulator,
	}
}

// PredictFinalTable predicts the final league table based on current standings and team strengths
func (p *TablePredictor) PredictFinalTable() ([]*models.TeamStats, error) {
	// Get all teams
	teams, err := p.TeamRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Create copies of teams for prediction
	teamCopies := make([]*models.Team, len(teams))
	for i, team := range teams {
		teamCopy := *team // Create a copy
		teamCopies[i] = &teamCopy
	}

	// Get matches that have not been played yet
	unplayedMatches, err := p.MatchRepo.GetUnplayed()
	if err != nil {
		return nil, err
	}

	// Create a map of team ID to team object for easy lookup
	teamMap := make(map[int]*models.Team)
	for _, team := range teamCopies {
		teamMap[team.ID] = team
	}

	// Simulate remaining matches
	for _, match := range unplayedMatches {
		homeTeam := teamMap[match.HomeTeamID]
		awayTeam := teamMap[match.AwayTeamID]

		// Simulate the match
		simulatedMatch, err := p.Simulator.SimulateMatch(homeTeam, awayTeam)
		if err != nil {
			return nil, err
		}

		// Update team stats based on simulated match result
		updateTeamStats(homeTeam, awayTeam, simulatedMatch)
	}

	// Convert team data to team stats
	teamStats := make([]*models.TeamStats, len(teamCopies))
	for i, team := range teamCopies {
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

	return teamStats, nil
} 