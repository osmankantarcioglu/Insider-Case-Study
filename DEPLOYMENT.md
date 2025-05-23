# Free Deployment Guide

This guide will help you deploy the Football League Simulator application to Render.com for free.

## Option 1: Render.com (100% Free)

Render offers a completely free tier for both web services and PostgreSQL databases.

### Steps:

1. Create a [Render account](https://render.com/register)

2. Connect your GitHub repository to Render

3. Click "New" and select "Blueprint"
   - Choose your repository
   - Render will automatically detect the `render.yaml` file

4. Review the settings and click "Apply"
   - This will create both the web service and PostgreSQL database
   - The database connection will be automatically configured

5. Wait for deployment to complete (5-10 minutes for the first deploy)

6. Access your application via the provided URL (e.g., `https://football-sim.onrender.com`)

### Monitoring and Logs

- View logs: Go to your web service dashboard → Logs
- Check database: Go to your database dashboard → Connect → View Connection Details

### Limitations on Free Tier

- Web services sleep after 15 minutes of inactivity
  - They automatically wake up when receiving traffic
  - First request after sleeping may be slow
- PostgreSQL database has 1GB storage limit

## Option 2: Fly.io (Free Tier)

Fly.io also offers a generous free tier.

### Prerequisites

1. Install [flyctl](https://fly.io/docs/hands-on/install-flyctl/)
2. Create a [Fly.io account](https://fly.io/app/sign-up)

### Steps:

1. Login to Fly
   ```
   fly auth login
   ```

2. Launch the application
   ```
   fly launch
   ```
   - Follow the prompts
   - When asked to deploy, say Yes
   - It will detect the Dockerfile automatically

3. Create a PostgreSQL database
   ```
   fly postgres create
   ```

4. Connect the database to your application
   ```
   fly postgres attach --app your-app-name your-db-name
   ```

5. Deploy the application
   ```
   fly deploy
   ```

6. Open the application
   ```
   fly open
   ```

## Option 3: Free PaaS on Oracle Cloud (Always Free)

Oracle Cloud offers an "Always Free" tier that never expires.

1. Create an [Oracle Cloud account](https://www.oracle.com/cloud/free/)
2. Create a VM instance using the Always Free tier
3. SSH into the instance and set up Docker
4. Deploy your application using Docker Compose 