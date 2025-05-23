-- Schema for the football simulation database

-- Teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    played INTEGER NOT NULL DEFAULT 0,
    won INTEGER NOT NULL DEFAULT 0,
    drawn INTEGER NOT NULL DEFAULT 0,
    lost INTEGER NOT NULL DEFAULT 0,
    goals_for INTEGER NOT NULL DEFAULT 0,
    goals_against INTEGER NOT NULL DEFAULT 0,
    goal_difference INTEGER NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0,
    strength INTEGER NOT NULL DEFAULT 5
);

-- League table
CREATE TABLE IF NOT EXISTS leagues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    season VARCHAR(20) NOT NULL,
    current_week INTEGER NOT NULL DEFAULT 1,
    total_weeks INTEGER NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    week INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    home_team_name VARCHAR(100) NOT NULL,
    away_team_name VARCHAR(100) NOT NULL,
    home_team_goals INTEGER,
    away_team_goals INTEGER,
    played BOOLEAN NOT NULL DEFAULT FALSE,
    played_at TIMESTAMP,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT different_teams CHECK (home_team_id != away_team_id)
);

-- Predictions table
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    league_id INTEGER NOT NULL REFERENCES leagues(id),
    team_id INTEGER NOT NULL REFERENCES teams(id),
    predicted_position INTEGER NOT NULL,
    predicted_points INTEGER NOT NULL,
    predicted_goal_difference INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Delete any existing data to prevent duplicates (in correct dependency order)
TRUNCATE TABLE predictions CASCADE;
TRUNCATE TABLE matches CASCADE;
TRUNCATE TABLE leagues CASCADE;
TRUNCATE TABLE teams RESTART IDENTITY CASCADE;

-- Insert sample teams
INSERT INTO teams (name, strength) VALUES 
('Manchester City', 9),
('Liverpool', 8),
('Arsenal', 7),
('Chelsea', 7);

-- Create a new league
INSERT INTO leagues (name, season, total_weeks) VALUES 
('Premier League', '2023-2024', 18);

-- First round: each team plays against each other team once
-- Week 1
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(1, 1, 2, 'Manchester City', 'Liverpool'),
(1, 3, 4, 'Arsenal', 'Chelsea');

-- Week 2
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(2, 1, 3, 'Manchester City', 'Arsenal'),
(2, 2, 4, 'Liverpool', 'Chelsea');

-- Week 3
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(3, 1, 4, 'Manchester City', 'Chelsea'),
(3, 2, 3, 'Liverpool', 'Arsenal');

-- Second round: each team plays against each other team again (reversed venues)
-- Week 4
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(4, 2, 1, 'Liverpool', 'Manchester City'),
(4, 4, 3, 'Chelsea', 'Arsenal');

-- Week 5
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(5, 3, 1, 'Arsenal', 'Manchester City'),
(5, 4, 2, 'Chelsea', 'Liverpool');

-- Week 6
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(6, 4, 1, 'Chelsea', 'Manchester City'),
(6, 3, 2, 'Arsenal', 'Liverpool');

-- Third round: each team plays against each other team a third time
-- Week 7
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(7, 1, 2, 'Manchester City', 'Liverpool'),
(7, 3, 4, 'Arsenal', 'Chelsea');

-- Week 8
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(8, 1, 3, 'Manchester City', 'Arsenal'),
(8, 2, 4, 'Liverpool', 'Chelsea');

-- Week 9
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(9, 1, 4, 'Manchester City', 'Chelsea'),
(9, 2, 3, 'Liverpool', 'Arsenal');

-- Week 10
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(10, 2, 1, 'Liverpool', 'Manchester City'),
(10, 4, 3, 'Chelsea', 'Arsenal');

-- Week 11
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(11, 3, 1, 'Arsenal', 'Manchester City'),
(11, 4, 2, 'Chelsea', 'Liverpool');

-- Week 12
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(12, 4, 1, 'Chelsea', 'Manchester City'),
(12, 3, 2, 'Arsenal', 'Liverpool');

-- Fourth round (making sure each team plays with others exactly 3 times)
-- Week 13
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(13, 1, 2, 'Manchester City', 'Liverpool'),
(13, 4, 3, 'Chelsea', 'Arsenal');

-- Week 14
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(14, 3, 1, 'Arsenal', 'Manchester City'),
(14, 2, 4, 'Liverpool', 'Chelsea');

-- Week 15
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(15, 1, 4, 'Manchester City', 'Chelsea'),
(15, 3, 2, 'Arsenal', 'Liverpool');

-- Week 16
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(16, 2, 1, 'Liverpool', 'Manchester City'),
(16, 3, 4, 'Arsenal', 'Chelsea');

-- Week 17
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(17, 1, 3, 'Manchester City', 'Arsenal'),
(17, 4, 2, 'Chelsea', 'Liverpool');

-- Week 18
INSERT INTO matches (week, home_team_id, away_team_id, home_team_name, away_team_name) VALUES 
(18, 4, 1, 'Chelsea', 'Manchester City'),
(18, 2, 3, 'Liverpool', 'Arsenal'); 