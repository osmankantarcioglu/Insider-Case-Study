# Football League Simulator

A full-stack football league simulation web application with Go backend and web frontend that allows users to simulate a football league season with realistic results.

## Features

- **League Simulation**: Simulate an entire season of football matches
- **Team Management**: View team statistics and performance
- **Match Results**: Dynamic match simulation with realistic scores
- **Championship Predictions**: Real-time win probability calculations
- **Editable Scores**: Manual override of simulated match scores

## Tech Stack

- **Backend**: Go with Fiber web framework
- **Database**: PostgreSQL 
- **Frontend**: HTML, CSS, JavaScript
- **Deployment**: Render.com (free tier)

## Key Components

- **Models**: Team, Match, League data structures
- **Services**: Match simulation, table prediction
- **Repositories**: Database interaction layer
- **Handlers**: HTTP request processing
- **Database**: PostgreSQL with SQL schema

## Local Development

### Prerequisites

- Go 1.20+
- PostgreSQL
- Docker and Docker Compose (optional)

### Run with Docker

```bash
# Start the application and database
docker-compose up
```

### Run Locally

1. Start PostgreSQL database
2. Set environment variables:
   ```
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=footballsim
   ```
3. Run the application:
   ```bash
   go run ./cmd
   ```

4. Access the application at http://localhost:8080

## Deployment

The application is deployed on Render.com using the free tier:

- Web service running Go application in Docker
- PostgreSQL database for data storage

See [DEPLOYMENT.md](DEPLOYMENT.md) for deployment instructions.

## Project Structure

```
├── cmd/            # Application entry point
├── database/       # Database interaction code
│   └── sql_schema.sql
├── handlers/       # HTTP request handlers
├── models/         # Data models
├── services/       # Business logic
└── utils/
    └── static/     # Frontend files
```

## Key Challenges Solved

1. **Week Progression Issue**: Fixed the "Next Week" button to correctly advance the league simulation
2. **UI Improvements**: Added dynamic match results display, score editing capabilities, and championship probability calculation
3. **Database Integration**: Created proper repository pattern for database operations
4. **Deployment**: Successfully deployed to Render.com with proper database configuration

## Next Steps

- Add user authentication
- Implement team creation/editing
- Add player management
- Provide more detailed statistics

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
