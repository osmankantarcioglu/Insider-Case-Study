package services

import (
	"math/rand"
	"time"
	"log"

	"github.com/user/footballsim/models"
)

// MatchSimulator implements the Simulator interface
type MatchSimulator struct {
	TeamRepo   TeamRepository
	MatchRepo  MatchRepository
	LeagueRepo LeagueRepository
}

// NewMatchSimulator creates a new match simulator
func NewMatchSimulator(teamRepo TeamRepository, matchRepo MatchRepository, leagueRepo LeagueRepository) *MatchSimulator {
	rand.Seed(time.Now().UnixNano())
	return &MatchSimulator{
		TeamRepo:   teamRepo,
		MatchRepo:  matchRepo,
		LeagueRepo: leagueRepo,
	}
}

// SimulateMatch simulates a match between two teams
func (s *MatchSimulator) SimulateMatch(homeTeam, awayTeam *models.Team) (*models.Match, error) {
	// Simulate based on team strength
	homeTeamGoals := simulateGoals(homeTeam.Strength, true)
	awayTeamGoals := simulateGoals(awayTeam.Strength, false)

	match := &models.Match{
		HomeTeamID:    homeTeam.ID,
		AwayTeamID:    awayTeam.ID,
		HomeTeamName:  homeTeam.Name,
		AwayTeamName:  awayTeam.Name,
		HomeTeamGoals: homeTeamGoals,
		AwayTeamGoals: awayTeamGoals,
		Played:        true,
		PlayedAt:      time.Now(),
	}

	return match, nil
}

// SimulateWeek simulates all matches for a specific week
func (s *MatchSimulator) SimulateWeek(week int) ([]*models.Match, error) {
	log.Printf("SimulateWeek called for week %d", week)
	
	matches, err := s.MatchRepo.GetByWeek(week)
	if err != nil {
		log.Printf("Error getting matches for week %d: %v", week, err)
		return nil, err
	}

	log.Printf("Found %d matches for week %d", len(matches), week)
	
	playedMatches := make([]*models.Match, 0, len(matches))

	for _, match := range matches {
		if match.Played && !match.IsEdited {
			log.Printf("Skipping already played match: %s vs %s", match.HomeTeamName, match.AwayTeamName)
			continue
		}

		log.Printf("Simulating match: %s vs %s", match.HomeTeamName, match.AwayTeamName)
		
		homeTeam, err := s.TeamRepo.GetByID(match.HomeTeamID)
		if err != nil {
			log.Printf("Error getting home team (ID: %d): %v", match.HomeTeamID, err)
			return nil, err
		}

		awayTeam, err := s.TeamRepo.GetByID(match.AwayTeamID)
		if err != nil {
			log.Printf("Error getting away team (ID: %d): %v", match.AwayTeamID, err)
			return nil, err
		}

		simulatedMatch, err := s.SimulateMatch(homeTeam, awayTeam)
		if err != nil {
			log.Printf("Error simulating match: %v", err)
			return nil, err
		}

		match.HomeTeamGoals = simulatedMatch.HomeTeamGoals
		match.AwayTeamGoals = simulatedMatch.AwayTeamGoals
		match.Played = true
		match.PlayedAt = time.Now()

		// Update match in database
		log.Printf("Updating match in database: %s %d-%d %s", 
			match.HomeTeamName, match.HomeTeamGoals, match.AwayTeamGoals, match.AwayTeamName)
		
		err = s.MatchRepo.Update(match)
		if err != nil {
			log.Printf("Error updating match: %v", err)
			return nil, err
		}

		// Update team stats
		log.Printf("Updating team stats for %s and %s", homeTeam.Name, awayTeam.Name)
		updateTeamStats(homeTeam, awayTeam, match)

		err = s.TeamRepo.Update(homeTeam)
		if err != nil {
			log.Printf("Error updating home team: %v", err)
			return nil, err
		}

		err = s.TeamRepo.Update(awayTeam)
		if err != nil {
			log.Printf("Error updating away team: %v", err)
			return nil, err
		}

		playedMatches = append(playedMatches, match)
	}

	// Update current week
	currentWeek, err := s.LeagueRepo.GetCurrentWeek()
	if err != nil {
		log.Printf("Error getting current week: %v", err)
		return nil, err
	}

	log.Printf("Current league week: %d, simulated week: %d", currentWeek, week)
	
	if currentWeek == week {
		totalWeeks, err := s.LeagueRepo.GetTotalWeeks()
		if err != nil {
			log.Printf("Error getting total weeks: %v", err)
			return nil, err
		}

		log.Printf("Total league weeks: %d", totalWeeks)
		
		if currentWeek < totalWeeks {
			log.Printf("Advancing to week %d", currentWeek + 1)
			err = s.LeagueRepo.UpdateWeek(currentWeek + 1)
			if err != nil {
				log.Printf("Error updating league week: %v", err)
				return nil, err
			}
		} else {
			log.Printf("Marking league as completed")
			err = s.LeagueRepo.MarkAsCompleted()
			if err != nil {
				log.Printf("Error marking league as completed: %v", err)
				return nil, err
			}
		}
	}

	log.Printf("Successfully simulated %d matches for week %d", len(playedMatches), week)
	return playedMatches, nil
}

// SimulateRemaining simulates all remaining matches in the league
func (s *MatchSimulator) SimulateRemaining() ([]*models.Match, error) {
	unplayedMatches, err := s.MatchRepo.GetUnplayed()
	if err != nil {
		return nil, err
	}

	// Group matches by week
	matchesByWeek := make(map[int][]*models.Match)
	for _, match := range unplayedMatches {
		matchesByWeek[match.Week] = append(matchesByWeek[match.Week], match)
	}

	// Get current week
	currentWeek, err := s.LeagueRepo.GetCurrentWeek()
	if err != nil {
		return nil, err
	}

	// Get total weeks
	totalWeeks, err := s.LeagueRepo.GetTotalWeeks()
	if err != nil {
		return nil, err
	}

	// Simulate each week in order
	allPlayedMatches := make([]*models.Match, 0)
	for week := currentWeek; week <= totalWeeks; week++ {
		if _, ok := matchesByWeek[week]; ok {
			playedMatches, err := s.SimulateWeek(week)
			if err != nil {
				return nil, err
			}
			allPlayedMatches = append(allPlayedMatches, playedMatches...)
		}
	}

	return allPlayedMatches, nil
}

// Helper functions
func simulateGoals(teamStrength int, isHome bool) int {
	// Home advantage factor
	homeFactor := 0
	if isHome {
		homeFactor = 1
	}

	// Base goal probability adjusted by team strength
	baseProb := float64(teamStrength) / 10.0

	// Generate a random number of goals with more weight to stronger teams
	goals := 0
	for i := 0; i < 5; i++ { // Max 5 goals
		if rand.Float64() < (baseProb + float64(homeFactor)*0.1) {
			goals++
		}
	}

	return goals
}

// updateTeamStats updates the statistics for both teams after a match
func updateTeamStats(homeTeam, awayTeam *models.Team, match *models.Match) {
	// Update home team stats
	homeTeam.Played++
	homeTeam.GoalsFor += match.HomeTeamGoals
	homeTeam.GoalsAgainst += match.AwayTeamGoals

	// Update away team stats
	awayTeam.Played++
	awayTeam.GoalsFor += match.AwayTeamGoals
	awayTeam.GoalsAgainst += match.HomeTeamGoals

	// Update wins, draws, losses
	if match.HomeTeamGoals > match.AwayTeamGoals {
		homeTeam.Won++
		awayTeam.Lost++
	} else if match.HomeTeamGoals < match.AwayTeamGoals {
		homeTeam.Lost++
		awayTeam.Won++
	} else {
		homeTeam.Drawn++
		awayTeam.Drawn++
	}

	// Update points and goal difference
	homeTeam.UpdateStats()
	awayTeam.UpdateStats()
} 