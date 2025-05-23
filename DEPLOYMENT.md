# Deployment Guide

This guide will help you deploy the Football League Simulator application to Heroku using Docker.

## Prerequisites

1. [Heroku account](https://signup.heroku.com/) (free)
2. [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) installed
3. [Git](https://git-scm.com/downloads) installed

## Deployment Steps

### 1. Login to Heroku

```
heroku login
```

### 2. Create a new Heroku application

```
heroku create football-league-sim
```

Replace `football-league-sim` with your preferred application name.

### 3. Add a PostgreSQL database

```
heroku addons:create heroku-postgresql:mini
```

This will create a PostgreSQL database and automatically set the DATABASE_URL environment variable.

### 4. Configure the application for Docker deployment

```
heroku stack:set container
```

### 5. Deploy the application

```
git push heroku main
```

Replace `main` with your branch name if different.

### 6. Open the application

```
heroku open
```

## Troubleshooting

### Check logs

```
heroku logs --tail
```

### Connect to the PostgreSQL database

```
heroku pg:psql
```

### Restart the application

```
heroku restart
```

## Manual Database Setup

If you need to manually set up the database:

1. Get the database connection details:
   ```
   heroku pg:credentials:url
   ```

2. Use these credentials to connect with a PostgreSQL client like pgAdmin or DBeaver.

## Updating the Application

To update your application after making changes:

```
git add .
git commit -m "Your commit message"
git push heroku main
```

## Scaling the Application

To scale the application:

```
heroku ps:scale web=1
```

Replace `1` with the number of dynos you want to run. 