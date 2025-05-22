// API Base URL - Use the port that's exposed to the host (8081)
const API_BASE_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' 
    ? `http://${window.location.hostname}:8081/api` 
    : '/api';

// DOM Elements
const tableBody = document.getElementById('table-body');
const matchWeekHeader = document.getElementById('match-week-header');
const matchResults = document.getElementById('match-results');
const predictionWeekHeader = document.getElementById('prediction-week-header');
const predictions = document.getElementById('predictions');
const playAllButton = document.getElementById('play-all');
const nextWeekButton = document.getElementById('next-week');
const resetLeagueButton = document.getElementById('reset-league');

// League data
let currentLeague = null;
let currentWeek = 0;
let totalWeeks = 6;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    console.log('App initialized, connecting to API at:', API_BASE_URL);
    
    // Load current league data
    loadLeagueData();
    
    // Add event listeners
    playAllButton.addEventListener('click', handlePlayAll);
    nextWeekButton.addEventListener('click', handleNextWeek);
    resetLeagueButton.addEventListener('click', handleResetLeague);
});

// Load league data
async function loadLeagueData() {
    try {
        console.log('Fetching league data...');
        
        // Get league table first since it's more reliable
        await loadLeagueTable();
        
        // Then try to get the current league data
        const leagueResponse = await fetch(`${API_BASE_URL}/league`);
        if (!leagueResponse.ok) {
            console.error('League response not OK:', leagueResponse.status, leagueResponse.statusText);
            throw new Error(`Failed to load league data: ${leagueResponse.statusText}`);
        }
        const league = await leagueResponse.json();
        console.log('League data received:', league);
        currentLeague = league;
        currentWeek = league.current_week;
        totalWeeks = league.total_weeks;
        
        // Update UI based on league status
        updateUIControlsBasedOnLeagueStatus();
        
        // On initial load, check if any matches have been played
        const weekMatchesResponse = await fetch(`${API_BASE_URL}/matches/week/${currentWeek}`);
        if (weekMatchesResponse.ok) {
            const weekMatchesData = await weekMatchesResponse.json();
            const hasPlayedMatches = weekMatchesData.matches && weekMatchesData.matches.some(m => m.played);
            
            if (hasPlayedMatches) {
                // If matches have been played (not initial state), show them
                await loadWeekMatches(currentWeek);
            } else if (currentWeek === 1) {
                // Initial state - no matches played yet
                matchWeekHeader.textContent = "No matches played yet";
                matchResults.innerHTML = '<div class="match-result">Press "Next Week" to simulate the first week\'s matches</div>';
            } else {
                // Matches for previous week should be shown
                await loadWeekMatches(currentWeek);
            }
        }
        
        // Get prediction if we're at week 4 or later
        if (currentWeek >= 4) {
            await loadPrediction();
        } else {
            // Clear prediction section if before week 4
            predictionWeekHeader.textContent = "Predictions available after week 4";
            predictions.innerHTML = '';
        }
    } catch (error) {
        console.error('Error loading league data:', error);
        displayErrorMessage('Failed to load league data. Please check if the server is running.');
    }
}

// Update UI controls based on league status
function updateUIControlsBasedOnLeagueStatus() {
    // Hide next week and play all buttons if league is completed
    const isLeagueCompleted = currentWeek >= totalWeeks || (currentLeague && currentLeague.is_completed);
    
    if (isLeagueCompleted) {
        nextWeekButton.style.display = 'none';
        playAllButton.style.display = 'none';
        matchWeekHeader.textContent = `League Completed`;
    } else {
        nextWeekButton.style.display = 'inline-block';
        playAllButton.style.display = 'inline-block';
    }
}

// Load and render the league table
async function loadLeagueTable() {
    try {
        console.log('Fetching league table...');
        const response = await fetch(`${API_BASE_URL}/league/table`);
        if (!response.ok) {
            console.error('League table response not OK:', response.status, response.statusText);
            throw new Error(`Failed to load league table: ${response.statusText}`);
        }
        const data = await response.json();
        console.log('Table data received:', data);
        
        // Clear existing table
        tableBody.innerHTML = '';
        
        // Create a set to track team IDs we've already added
        const addedTeamIds = new Set();
        
        // Add teams to table (with deduplication)
        if (data.teams && Array.isArray(data.teams)) {
            data.teams.forEach(team => {
                // Skip if we've already added this team
                if (addedTeamIds.has(team.team_id)) {
                    return;
                }
                
                addedTeamIds.add(team.team_id);
                
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${team.team_name}</td>
                    <td>${team.points}</td>
                    <td>${team.played}</td>
                    <td>${team.won}</td>
                    <td>${team.drawn}</td>
                    <td>${team.lost}</td>
                    <td>${team.goal_difference}</td>
                `;
                tableBody.appendChild(row);
            });
        } else {
            console.error('Invalid teams data:', data);
        }
        
        // If we got current week data from the table
        if (data.current_week) {
            currentWeek = data.current_week;
            totalWeeks = data.total_weeks || 6;
        }
        
    } catch (error) {
        console.error('Error loading league table:', error);
        displayErrorMessage('Failed to load league table. Please check if the server is running.');
    }
}

// Load and render matches for a specific week
async function loadWeekMatches(week) {
    try {
        console.log(`Fetching week ${week} matches...`);
        
        // Initialize the week header even if the fetch fails
        matchWeekHeader.textContent = `${week}${getOrdinalSuffix(week)} Week Match Result`;
        matchResults.innerHTML = '';
        
        const response = await fetch(`${API_BASE_URL}/matches/week/${week}`);
        if (!response.ok) {
            console.error(`Week ${week} matches response not OK:`, response.status, response.statusText);
            matchResults.innerHTML = '<div class="match-result">No matches available for this week.</div>';
            return; // Return without throwing so the app can still function
        }
        
        const data = await response.json();
        console.log(`Week ${week} matches received:`, data);
        
        // Check if there are any matches
        if (!data.matches || data.matches.length === 0) {
            matchResults.innerHTML = '<div class="match-result">No matches available for this week.</div>';
            return;
        }
        
        // Create container for matches
        const matchesContainer = document.createElement('div');
        matchesContainer.id = 'matches-container';
        
        // Add match results
        data.matches.forEach(match => {
            const matchElement = document.createElement('div');
            matchElement.className = 'match-result';
            matchElement.dataset.matchId = match.id;
            
            if (match.played) {
                // For all played matches, show score prominently
                matchElement.innerHTML = `
                    <div class="match-teams">${match.home_team_name} vs ${match.away_team_name}</div>
                    <div class="match-score">${match.home_team_name} <strong>${match.home_team_goals}</strong> - <strong>${match.away_team_goals}</strong> ${match.away_team_name}</div>
                    <div class="score-editor">
                        <input type="number" min="0" max="10" class="score-input home-score" value="${match.home_team_goals}" data-match-id="${match.id}">
                        -
                        <input type="number" min="0" max="10" class="score-input away-score" value="${match.away_team_goals}" data-match-id="${match.id}">
                    </div>
                `;
            } else {
                matchElement.innerHTML = `
                    <div class="match-teams">${match.home_team_name} vs ${match.away_team_name}</div>
                    <div class="match-status">(Not played yet)</div>
                    <div class="score-editor">
                        <input type="number" min="0" max="10" class="score-input home-score" value="0" data-match-id="${match.id}">
                        -
                        <input type="number" min="0" max="10" class="score-input away-score" value="0" data-match-id="${match.id}">
                    </div>
                `;
            }
            
            matchesContainer.appendChild(matchElement);
        });
        
        // Add matches to the DOM
        matchResults.appendChild(matchesContainer);
        
        // Add single save button for all matches
        const saveButtonContainer = document.createElement('div');
        saveButtonContainer.className = 'save-button-container';
        
        const saveAllButton = document.createElement('button');
        saveAllButton.id = 'save-all-scores';
        saveAllButton.className = 'button save-all-btn';
        saveAllButton.textContent = 'Save All Scores';
        saveAllButton.addEventListener('click', handleSaveAllScores);
        
        saveButtonContainer.appendChild(saveAllButton);
        matchResults.appendChild(saveButtonContainer);
        
    } catch (error) {
        console.error(`Error loading week ${week} matches:`, error);
        matchResults.innerHTML = '<div class="match-result">Failed to load matches.</div>';
    }
}

// Handle saving all match scores
async function handleSaveAllScores() {
    // Disable save button to prevent multiple clicks
    const saveButton = document.getElementById('save-all-scores');
    saveButton.disabled = true;
    saveButton.textContent = 'Saving...';
    
    try {
        const matchElements = document.querySelectorAll('.match-result');
        const updatePromises = [];
        
        // Process each match
        matchElements.forEach(matchElement => {
            const matchId = matchElement.dataset.matchId;
            const homeScore = parseInt(matchElement.querySelector('.home-score').value);
            const awayScore = parseInt(matchElement.querySelector('.away-score').value);
            
            if (isNaN(homeScore) || isNaN(awayScore) || homeScore < 0 || awayScore < 0) {
                throw new Error(`Invalid score for match ID ${matchId}`);
            }
            
            // Add update request to promises array
            updatePromises.push(
                fetch(`${API_BASE_URL}/matches/${matchId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        home_team_goals: homeScore,
                        away_team_goals: awayScore
                    })
                })
            );
        });
        
        // Wait for all updates to complete
        const results = await Promise.all(updatePromises);
        
        // Check if any requests failed
        const failedRequests = results.filter(response => !response.ok);
        if (failedRequests.length > 0) {
            throw new Error(`${failedRequests.length} match updates failed`);
        }
        
        // Reload data to update all views
        await loadLeagueData();
        displaySuccessMessage('All match scores updated successfully');
        
        // Update the match display - now matches should show as played
        const currentWeekMatchesResponse = await fetch(`${API_BASE_URL}/matches/week/${currentWeek}`);
        if (currentWeekMatchesResponse.ok) {
            const data = await currentWeekMatchesResponse.json();
            displayUpdatedMatches(data.matches);
        }
        
    } catch (error) {
        console.error('Error updating match scores:', error);
        displayErrorMessage(`Failed to update match scores: ${error.message}`);
    } finally {
        // Re-enable save button
        saveButton.disabled = false;
        saveButton.textContent = 'Save All Scores';
    }
}

// Display updated matches after saving
function displayUpdatedMatches(matches) {
    const matchesContainer = document.getElementById('matches-container');
    if (!matchesContainer) return;
    
    // Clear existing matches
    matchesContainer.innerHTML = '';
    
    // Add updated matches
    matches.forEach(match => {
        const matchElement = document.createElement('div');
        matchElement.className = 'match-result';
        matchElement.dataset.matchId = match.id;
        
        // All matches should now be played
        matchElement.innerHTML = `
            <div class="match-teams">${match.home_team_name} vs ${match.away_team_name}</div>
            <div class="match-score">${match.home_team_name} <strong>${match.home_team_goals}</strong> - <strong>${match.away_team_goals}</strong> ${match.away_team_name}</div>
            <div class="score-editor">
                <input type="number" min="0" max="10" class="score-input home-score" value="${match.home_team_goals}" data-match-id="${match.id}">
                -
                <input type="number" min="0" max="10" class="score-input away-score" value="${match.away_team_goals}" data-match-id="${match.id}">
            </div>
        `;
        
        matchesContainer.appendChild(matchElement);
    });
}

// Load and render prediction
async function loadPrediction() {
    try {
        console.log('Fetching predictions...');
        
        // Set default headers
        predictionWeekHeader.textContent = `${currentWeek}${getOrdinalSuffix(currentWeek)} Week Predictions of Championship`;
        predictions.innerHTML = '';
        
        const response = await fetch(`${API_BASE_URL}/league/prediction`);
        if (!response.ok) {
            console.error('Prediction response not OK:', response.status, response.statusText);
            predictions.innerHTML = '<div class="error-message">Predictions are only available after week 4.</div>';
            return;
        }
        
        const data = await response.json();
        console.log('Prediction data received:', data);
        
        // Check if the prediction data structure is as expected
        if (!data.prediction || !Array.isArray(data.prediction)) {
            console.error('Invalid prediction data structure:', data);
            predictions.innerHTML = '<div class="error-message">Invalid prediction data received.</div>';
            return;
        }
        
        // Calculate initial percentages
        let teamPredictions = data.prediction.map(team => {
            return {
                name: team.team_name,
                percentage: calculateWinningPercentage(team, data.prediction)
            };
        });
        
        // Normalize percentages to ensure they sum to exactly 100%
        const totalPercentage = teamPredictions.reduce((sum, team) => sum + team.percentage, 0);
        
        // If total is not 0, normalize each percentage
        if (totalPercentage > 0) {
            teamPredictions = teamPredictions.map(team => ({
                name: team.name,
                percentage: Math.round((team.percentage / totalPercentage) * 100)
            }));
            
            // Due to rounding, we might need to adjust the highest percentage to make sum exactly 100
            const newTotal = teamPredictions.reduce((sum, team) => sum + team.percentage, 0);
            if (newTotal !== 100) {
                // Find team with highest percentage and adjust
                const highestIndex = teamPredictions.reduce((maxIndex, team, index, arr) => 
                    team.percentage > arr[maxIndex].percentage ? index : maxIndex, 0);
                    
                teamPredictions[highestIndex].percentage += (100 - newTotal);
            }
        }
        
        // Sort by percentage (highest first)
        teamPredictions.sort((a, b) => b.percentage - a.percentage);
        
        // Add predictions
        teamPredictions.forEach(team => {
            const predictionRow = document.createElement('div');
            predictionRow.className = 'prediction-row';
            predictionRow.innerHTML = `
                <div class="team-name">${team.name}</div>
                <div class="percentage">%${team.percentage}</div>
            `;
            predictions.appendChild(predictionRow);
        });
    } catch (error) {
        console.error('Error loading prediction:', error);
        predictions.innerHTML = '<div class="error-message">Failed to load predictions.</div>';
    }
}

// Handle Play All button click
async function handlePlayAll() {
    try {
        console.log('Simulating all matches...');
        
        // Disable buttons to prevent multiple clicks
        playAllButton.disabled = true;
        nextWeekButton.disabled = true;
        
        const response = await fetch(`${API_BASE_URL}/matches/simulate-all`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error('Simulate all matches response not OK:', response.status, response.statusText);
            throw new Error(`Failed to simulate all matches: ${response.statusText}`);
        }
        
        // Reload data after simulation
        await loadLeagueData();
    } catch (error) {
        console.error('Error simulating all matches:', error);
        displayErrorMessage('Failed to simulate all matches. Please check if the server is running.');
    } finally {
        // Re-enable buttons
        playAllButton.disabled = false;
        nextWeekButton.disabled = false;
    }
}

// Handle Next Week button click
async function handleNextWeek() {
    if (currentWeek < totalWeeks) {
        try {
            console.log(`Simulating week ${currentWeek}...`);
            
            // Disable buttons to prevent multiple clicks
            playAllButton.disabled = true;
            nextWeekButton.disabled = true;
            
            const response = await fetch(`${API_BASE_URL}/matches/week/${currentWeek}/simulate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            if (!response.ok) {
                console.error(`Simulate week ${currentWeek} response not OK:`, response.status, response.statusText);
                throw new Error(`Failed to simulate week ${currentWeek}: ${response.statusText}`);
            }
            
            // Get and parse response to update current week
            const data = await response.json();
            console.log(`Week ${currentWeek} simulation result:`, data);
            
            // Need to reload all data to get the updated current week
            await loadLeagueData();
            
            // Now display results of the just-simulated matches
            if (data && data.matches) {
                displaySimulatedMatches(data.matches);
            }
            
        } catch (error) {
            console.error('Error simulating current week:', error);
            displayErrorMessage('Failed to simulate current week. Please check if the server is running.');
        } finally {
            // Re-enable buttons
            playAllButton.disabled = false;
            nextWeekButton.disabled = false;
        }
    } else {
        alert('The league is already finished!');
    }
}

// Display simulated match results immediately after simulation
function displaySimulatedMatches(matches) {
    // Update the match results display
    matchResults.innerHTML = '';
    
    // Create container for matches
    const matchesContainer = document.createElement('div');
    matchesContainer.id = 'matches-container';
    
    matches.forEach(match => {
        if (match.played) {
            const matchElement = document.createElement('div');
            matchElement.className = 'match-result';
            matchElement.dataset.matchId = match.id;
            
            // Show the match result prominently
            matchElement.innerHTML = `
                <div class="match-teams">${match.home_team_name} vs ${match.away_team_name}</div>
                <div class="match-score">${match.home_team_name} <strong>${match.home_team_goals}</strong> - <strong>${match.away_team_goals}</strong> ${match.away_team_name}</div>
                <div class="score-editor">
                    <input type="number" min="0" max="10" class="score-input home-score" value="${match.home_team_goals}" data-match-id="${match.id}">
                    -
                    <input type="number" min="0" max="10" class="score-input away-score" value="${match.away_team_goals}" data-match-id="${match.id}">
                </div>
            `;
            
            matchesContainer.appendChild(matchElement);
        }
    });
    
    // Add matches to DOM
    matchResults.appendChild(matchesContainer);
    
    // Add single save button
    const saveButtonContainer = document.createElement('div');
    saveButtonContainer.className = 'save-button-container';
    
    const saveAllButton = document.createElement('button');
    saveAllButton.id = 'save-all-scores';
    saveAllButton.className = 'button save-all-btn';
    saveAllButton.textContent = 'Save All Scores';
    saveAllButton.addEventListener('click', handleSaveAllScores);
    
    saveButtonContainer.appendChild(saveAllButton);
    matchResults.appendChild(saveButtonContainer);
}

// Handle Reset League button click
async function handleResetLeague() {
    try {
        if (!confirm('Are you sure you want to reset the league? This will clear all match results and statistics.')) {
            return;
        }
        
        console.log('Resetting league...');
        
        // Disable buttons to prevent multiple clicks
        playAllButton.disabled = true;
        nextWeekButton.disabled = true;
        resetLeagueButton.disabled = true;
        
        const response = await fetch(`${API_BASE_URL}/league/reset`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error('Reset league response not OK:', response.status, response.statusText);
            throw new Error(`Failed to reset league: ${response.statusText}`);
        }
        
        // Reset local currentWeek
        currentWeek = 1;
        
        // Reload data after reset
        await loadLeagueData();
        
        // Show confirmation message
        displaySuccessMessage('League has been reset successfully');
    } catch (error) {
        console.error('Error resetting league:', error);
        displayErrorMessage('Failed to reset the league. Please check if the server is running.');
    } finally {
        // Re-enable buttons
        playAllButton.disabled = false;
        nextWeekButton.disabled = false;
        resetLeagueButton.disabled = false;
    }
}

// Helper function to display error messages
function displayErrorMessage(message) {
    console.error('ERROR:', message);
    if (!document.getElementById('error-message')) {
        const errorDiv = document.createElement('div');
        errorDiv.id = 'error-message';
        errorDiv.style.cssText = 'position: fixed; top: 10px; left: 50%; transform: translateX(-50%); background-color: #f44336; color: white; padding: 15px; border-radius: 4px; z-index: 1000;';
        errorDiv.textContent = message;
        
        // Add close button
        const closeButton = document.createElement('span');
        closeButton.style.cssText = 'margin-left: 15px; color: white; font-weight: bold; float: right; font-size: 22px; line-height: 20px; cursor: pointer;';
        closeButton.textContent = '×';
        closeButton.onclick = function() {
            document.body.removeChild(errorDiv);
        };
        
        errorDiv.appendChild(closeButton);
        document.body.appendChild(errorDiv);
        
        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (document.body.contains(errorDiv)) {
                document.body.removeChild(errorDiv);
            }
        }, 5000);
    }
}

// Helper function to get ordinal suffix (1st, 2nd, 3rd, 4th, etc.)
function getOrdinalSuffix(num) {
    const j = num % 10;
    const k = num % 100;
    if (j === 1 && k !== 11) {
        return "st";
    }
    if (j === 2 && k !== 12) {
        return "nd";
    }
    if (j === 3 && k !== 13) {
        return "rd";
    }
    return "th";
}

// Calculate a winning percentage based on team statistics
function calculateWinningPercentage(team, allTeams) {
    // Sort teams by points
    const sortedTeams = [...allTeams].sort((a, b) => {
        if (a.points !== b.points) {
            return b.points - a.points;
        }
        if (a.goal_difference !== b.goal_difference) {
            return b.goal_difference - a.goal_difference;
        }
        return b.goals_for - a.goals_for;
    });
    
    // Find this team's position
    const teamIndex = sortedTeams.findIndex(t => t.team_id === team.team_id);
    
    // Calculate total points in the league
    const totalPoints = sortedTeams.reduce((sum, t) => sum + t.points, 0);
    
    // If no points yet (early in season), use position-based prediction
    if (totalPoints === 0) {
        if (teamIndex === 0) return 40;
        if (teamIndex === 1) return 30;
        if (teamIndex === 2) return 20;
        return 10;
    }
    
    // Calculate percentage based on proportion of points + goal difference factor
    const pointsPercentage = (team.points / totalPoints) * 100;
    
    // Add a bonus/penalty based on goal difference (up to ±10%)
    let goalDifferenceBonus = 0;
    if (team.goal_difference !== 0) {
        const maxGD = Math.max(...sortedTeams.map(t => Math.abs(t.goal_difference)));
        if (maxGD > 0) {
            goalDifferenceBonus = (team.goal_difference / maxGD) * 10;
        }
    }
    
    // Calculate weighted percentage based on position
    let positionFactor = 1;
    if (teamIndex === 0) {
        positionFactor = 1.2; // First place bonus
    } else if (teamIndex === sortedTeams.length - 1) {
        positionFactor = 0.8; // Last place penalty
    }
    
    // Combine factors and normalize to ensure reasonable percentages
    let finalPercentage = (pointsPercentage + goalDifferenceBonus) * positionFactor;
    
    // Make sure the percentage is between 5% and 70%
    finalPercentage = Math.max(5, Math.min(70, finalPercentage));
    
    // Round to integer
    return Math.round(finalPercentage);
}

// Helper function to display success messages
function displaySuccessMessage(message) {
    console.log('SUCCESS:', message);
    if (!document.getElementById('success-message')) {
        const successDiv = document.createElement('div');
        successDiv.id = 'success-message';
        successDiv.style.cssText = 'position: fixed; top: 10px; left: 50%; transform: translateX(-50%); background-color: #4CAF50; color: white; padding: 15px; border-radius: 4px; z-index: 1000;';
        successDiv.textContent = message;
        
        // Add close button
        const closeButton = document.createElement('span');
        closeButton.style.cssText = 'margin-left: 15px; color: white; font-weight: bold; float: right; font-size: 22px; line-height: 20px; cursor: pointer;';
        closeButton.textContent = '×';
        closeButton.onclick = function() {
            document.body.removeChild(successDiv);
        };
        
        successDiv.appendChild(closeButton);
        document.body.appendChild(successDiv);
        
        // Auto-remove after 3 seconds
        setTimeout(() => {
            if (document.body.contains(successDiv)) {
                document.body.removeChild(successDiv);
            }
        }, 3000);
    }
} 