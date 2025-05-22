# Football League Simulator

A football league simulation system built with Go that simulates match results and provides league table predictions.

## Features

- 4-team football league simulation
- Match simulation based on team strengths
- League table with Premier League rules (points, goal difference, etc.)
- Week-by-week match results
- Prediction of final league standings after week 4
- Edit match results and recalculate standings
- Play all remaining matches in one go

## Technology Stack

- **Backend**: Go with Fiber framework
- **Database**: PostgreSQL
- **Architecture**: Interface-based design with clean code principles

## API Endpoints

### Teams

- `GET /api/teams` - Get all teams
- `GET /api/teams/:id` - Get team by ID
- `POST /api/teams` - Create a new team
- `PUT /api/teams/:id` - Update a team
- `DELETE /api/teams/:id` - Delete a team

### Matches

- `GET /api/matches` - Get all matches
- `GET /api/matches/week/:week` - Get matches for a specific week
- `POST /api/matches/week/:week/simulate` - Simulate matches for a specific week
- `POST /api/matches/simulate-all` - Simulate all remaining matches
- `PUT /api/matches/:id` - Update match result

### League

- `GET /api/league` - Get current league information
- `GET /api/league/table` - Get current league table
- `GET /api/league/prediction` - Get final league table prediction (after week 4)
- `POST /api/league` - Create a new league
- `POST /api/league/reset` - Reset the league to the beginning

## Setup and Installation

### Prerequisites

- Go (version 1.20 or higher)
- PostgreSQL

### Environment Variables

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=footballsim
DB_SSLMODE=disable
```

### Installation Steps

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/footballsim.git
   cd footballsim
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Create the database:
   ```
   createdb footballsim
   ```

4. Run the application:
   ```
   go run cmd/main.go
   ```

5. The API will be available at:
   ```
   http://localhost:8080
   ```

## Database Schema

The database schema is defined in `database/sql_schema.sql`. It contains the following tables:

- `teams` - Team information
- `leagues` - League information
- `matches` - Match information
- `predictions` - Prediction information

## Usage Examples

### Simulate a Week

```
curl -X POST http://localhost:8080/api/matches/week/1/simulate
```

### Get Current League Table

```
curl http://localhost:8080/api/league/table
```

### Get Prediction

```
curl http://localhost:8080/api/league/prediction
```

### Edit Match Result

```
curl -X PUT http://localhost:8080/api/matches/1 -H "Content-Type: application/json" -d '{"home_team_goals": 3, "away_team_goals": 1}'
```

## Docker Deployment

Build the Docker image:

```
docker build -t footballsim .
```

Run the container:

```
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=footballsim \
  -e DB_SSLMODE=disable \
  footballsim
```

## Code Structure

- `cmd/` - Application entry point
- `models/` - Data models
- `services/` - Business logic implementation
- `handlers/` - HTTP handlers
- `database/` - Database access layer
- `utils/` - Utility functions

## License

MIT License

## Author

Football League Simulator Team 