# Docker Setup Guide

This guide explains how to run the TCG Server Go application in Docker containers, completely isolated from your local environment.

## Prerequisites

- Docker Desktop installed (or Docker Engine + Docker Compose)
- At least 2GB of free disk space

## Quick Start

1. **Create environment file** (optional, defaults are provided):
   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings
   ```

2. **Build and start containers**:
   ```bash
   docker-compose up -d
   ```

3. **Check server status**:
   ```bash
   curl http://localhost:8080/health
   ```

4. **View logs**:
   ```bash
   # All services
   docker-compose logs -f
   
   # Just the server
   docker-compose logs -f tcg-server
   
   # Just the database
   docker-compose logs -f mariadb
   ```

## Environment Variables

Create a `.env` file in the project root with the following variables (or use the defaults):

```env
# Database Configuration
DB_HOST=mariadb
DB_PORT=3306
DB_USER=tcg_user
DB_PASSWORD=rootpassword
DB_NAME=tcg_server

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
```

**Important:** Change the `DB_PASSWORD` and `JWT_SECRET` in production!

## Docker Commands

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### Stop and Remove Volumes (⚠️ This deletes database data)
```bash
docker-compose down -v
```

### Rebuild Containers
```bash
docker-compose build --no-cache
docker-compose up -d
```

### View Running Containers
```bash
docker-compose ps
```

### Access Database
```bash
# Connect to MariaDB container
docker-compose exec mariadb mysql -u tcg_user -p tcg_server

# Or using root
docker-compose exec mariadb mysql -u root -p
```

### Access Server Container
```bash
docker-compose exec tcg-server sh
```

## Container Details

### MariaDB Container
- **Image:** mariadb:10.11
- **Port:** 3306 (mapped to host)
- **Data Persistence:** Stored in Docker volume `mariadb_data`
- **Health Check:** Automatically checks database connectivity

### TCG Server Container
- **Port:** 8080 (mapped to host)
- **Depends on:** MariaDB (waits for healthy database)
- **Health Check:** Checks `/health` endpoint
- **Auto-restart:** Restarts automatically on failure

## Network

Both containers run on a private Docker network (`tcg-network`), so they can communicate using service names:
- Server connects to database using hostname: `mariadb`
- Database is accessible from host on: `localhost:3306`
- Server is accessible from host on: `localhost:8080`

## Troubleshooting

### Server can't connect to database
```bash
# Check if MariaDB is healthy
docker-compose ps

# Check MariaDB logs
docker-compose logs mariadb

# Restart services
docker-compose restart
```

### Port already in use
If port 8080 or 3306 is already in use, change them in `.env`:
```env
PORT=8081
DB_PORT=3307
```

### Database data persistence
Database data is stored in a Docker volume. To completely reset:
```bash
docker-compose down -v
docker-compose up -d
```

### View real-time logs
```bash
docker-compose logs -f --tail=100
```

## Production Considerations

1. **Change default passwords** in `.env`
2. **Use strong JWT_SECRET** (at least 32 characters)
3. **Configure proper firewall rules**
4. **Set up database backups**
5. **Use Docker secrets** for sensitive data
6. **Configure resource limits** in docker-compose.yml
7. **Use reverse proxy** (nginx/traefik) for HTTPS

## Development Workflow

1. Make code changes
2. Rebuild container:
   ```bash
   docker-compose build tcg-server
   docker-compose up -d tcg-server
   ```
3. Or rebuild everything:
   ```bash
   docker-compose up -d --build
   ```

## Cleanup

To remove everything (containers, volumes, networks):
```bash
docker-compose down -v --remove-orphans
```
