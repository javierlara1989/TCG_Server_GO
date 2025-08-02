# JWT Authentication Server in Go

A modular Go server that handles JWT authentication with MariaDB database integration, organized in separate packages for better maintainability and production readiness.

## Features

- User registration with comprehensive validation
- Login endpoint that receives email and password
- JWT token generation with expiration (24 hours)
- Protected endpoint for token validation
- Authentication middleware
- Password encryption with bcrypt
- **MariaDB database integration** with proper user management
- **Modular architecture** with separation of concerns
- **Soft delete** functionality for users
- **Automatic table creation** on startup
- **Advanced input validation** with custom rules

## Input Validation Rules

### Nombre (Name)
- **Minimum 6 characters**
- **Only letters and spaces allowed**
- **Supports Spanish characters** (á, é, í, ó, ú, ñ)

### Password
- **Minimum 6 characters**
- **Must contain both letters and numbers**
- **Automatically encrypted** with bcrypt

### Email
- **Valid email format required**
- **Must be unique** in the database

## Project Structure

```
TCG_Server_GO/
├── main.go              # Main entry point
├── go.mod               # Dependency management
├── README.md            # Documentation
├── ENVIRONMENT.md       # Environment variables documentation
├── models/              # Data structures
│   └── user.go
├── database/            # Database operations
│   ├── database.go      # Database connection and configuration
│   └── users.go         # User database operations
├── auth/                # Authentication logic
│   ├── jwt.go           # JWT token handling
│   └── users.go         # User management
├── handlers/            # HTTP handlers
│   ├── auth.go          # Login and register handlers
│   ├── validate.go      # Validation handler and custom validation
│   ├── health.go        # Health check handler
│   └── routes.go        # Route configuration
└── middleware/          # Middlewares
    └── auth.go          # Authentication middleware
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Installation

1. Ensure you have Go installed (version 1.21 or higher)

2. Install MariaDB and create a database

3. Clone or download this project

4. Install dependencies:
```bash
go mod tidy
```

5. Configure environment variables (see `ENVIRONMENT.md` for details)

6. Run the server:
```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port using the `PORT` environment variable.

## API Endpoints

### POST /register
Registers a new user with validation.

**Request:**
```json
{
  "nombre": "Juan Pérez",
  "email": "juan@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Validation Errors:**
```json
{
  "errors": [
    {
      "field": "nombre",
      "message": "Nombre must be at least 6 characters and contain only letters"
    },
    {
      "field": "password",
      "message": "Password must be at least 6 characters and contain both letters and numbers"
    }
  ]
}
```

### POST /login
Authenticates user with email and password.

**Request:**
```json
{
  "email": "juan@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### GET /api/validate
Validates a JWT token (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "OK"
}
```

### GET /health
Checks server status.

**Response:**
```json
{
  "status": "OK"
}
```

## Database Operations

The application includes comprehensive database operations:

- **CreateUser**: Create new users
- **GetUserByEmail**: Retrieve user by email
- **GetUserByID**: Retrieve user by ID
- **UpdateUser**: Update user information
- **UpdatePassword**: Update user password
- **SoftDeleteUser**: Mark user as deleted (soft delete)
- **HardDeleteUser**: Permanently delete user
- **EmailExists**: Check if email already exists

## Usage Examples

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "María García",
    "email": "maria@example.com",
    "password": "password123"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "password123"
  }'
```

3. **Validate token:**
```bash
curl -X GET http://localhost:8080/api/validate \
  -H "Authorization: Bearer <token_from_login>"
```

## Modular Architecture

### Packages

- **`models`**: Defines data structures (User, LoginRequest, etc.)
- **`database`**: Handles database connection and operations
- **`auth`**: Handles authentication logic and user management
- **`handlers`**: Contains HTTP handlers for each endpoint
- **`middleware`**: Middlewares for validation and request processing

### Benefits of Modular Structure

- **Separation of concerns**: Each package has a specific function
- **Maintainability**: Easy to maintain and extend
- **Testability**: Each component can be tested independently
- **Reusability**: Packages can be reused in other projects
- **Scalability**: Easy to add new features

## Security Configuration

- **JWT Secret**: Configure via `JWT_SECRET` environment variable
- **Expiration**: Tokens expire in 24 hours by default
- **Passwords**: Stored hashed with bcrypt
- **Database**: Uses parameterized queries to prevent SQL injection
- **Soft Delete**: Users are marked as deleted rather than permanently removed
- **Input Validation**: Comprehensive validation for all user inputs

## Production Considerations

- Configure HTTPS
- Use environment variables for all secrets
- Implement rate limiting
- Add appropriate logging
- Set up proper monitoring and health checks
- Implement proper error handling and logging
- Use secure session management
- Regular security audits and updates
- Consider using connection pooling for database
- Implement database migrations for schema changes 