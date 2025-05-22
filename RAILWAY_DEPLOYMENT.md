# Railway Deployment Guide for Football League Simulator

## Fix for Database Connection Issues

If you're encountering database connection errors when deploying to Railway, follow these steps:

## Step 1: Deploy Both PostgreSQL and Your Application

1. In Railway dashboard, create a new project
2. Add a PostgreSQL database service
3. Add another service for your application (GitHub or Docker deployment)

## Step 2: Link the Database to Your Application

1. Go to your application service in Railway
2. Click on "Variables" tab
3. Click "Add from service" and select your PostgreSQL database
4. This will add all the necessary PostgreSQL environment variables:
   - `PGHOST`
   - `PGPORT`
   - `PGUSER`
   - `PGPASSWORD`
   - `PGDATABASE`

## Step 3: Add Required Code Files

Make sure you've added these files to your project:

1. `database/railway_db.go` - Contains Railway configuration logic
2. `railway.toml` - Railway deployment configuration

## Step 4: Verify Database Connection in Logs

1. Once deployed, check your logs in Railway dashboard
2. Look for these messages to confirm proper connection:
   - "Railway environment detected, using Railway database configuration"
   - "Railway database configuration detected: Host=..."
   - "Successfully connected to database"

## Troubleshooting

If you still experience connection issues:

1. **Check SSL Mode**: Railway may require SSL for database connections
   - In `database/railway_db.go`, ensure `sslMode` is set to `"require"`

2. **Check Port Configuration**: Make sure your app is listening on the correct port
   - Verify that your app uses the `PORT` environment variable provided by Railway

3. **Database Initialization**: If your database schema isn't being applied
   - Check that the schema file path is correct in `database/db.go`
   - The path should be `"database/sql_schema.sql"` relative to your app's working directory

4. **Manual Schema Application**: You can use Railway's PostgreSQL service dashboard to run:
   ```sql
   -- Connect to your database in Railway dashboard
   -- Paste your schema.sql contents here
   ```

## Additional Railway Tips

1. **Custom Domains**: To use a custom domain for your app
   - Go to your application service in Railway
   - Click "Settings" tab
   - Under "Domains", click "Generate Domain" or "Add Custom Domain"

2. **Environment Variables**: Add any additional environment variables needed
   - Go to application service → Variables tab
   - Click "New Variable" to add custom configurations

3. **Monitoring**: Monitor your application's health
   - Go to application service → Metrics tab
   - Check CPU, Memory, and Disk usage

4. **Logs**: Continuously monitor logs for issues
   - Go to application service → Logs tab
   - Filter logs to identify any recurring errors 