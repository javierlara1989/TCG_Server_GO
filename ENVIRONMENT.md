# Environment Variables

This application uses the following environment variables for configuration:

## Database Configuration

- `DB_HOST`: MariaDB host (default: localhost)
- `DB_PORT`: MariaDB port (default: 3306)
- `DB_USER`: Database username (default: root)
- `DB_PASSWORD`: Database password (required)
- `DB_NAME`: Database name (default: tcg_server)

## Server Configuration

- `PORT`: Server port (default: 8080)

## JWT Configuration

- `JWT_SECRET`: Secret key for JWT tokens (optional, will use default if not set)

## Example .env file

Create a `.env` file in the root directory with the following content:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=tcg_server
PORT=8080
JWT_SECRET=your_jwt_secret_here
```

## Database Setup

1. Create a MariaDB database named `tcg_server` (or whatever you set in `DB_NAME`)
2. The application will automatically create the required tables on startup
3. Make sure the database user has the necessary permissions to create tables and perform CRUD operations 