# JWT Authentication Server in Go

A modular Go server that handles JWT authentication, organized in separate packages for better maintainability and production readiness.

## Features

- Login endpoint that receives username and password
- JWT token generation with expiration (24 hours)
- Protected endpoint for token validation
- Authentication middleware
- Password encryption with bcrypt
- **Modular architecture** with separation of concerns

## Project Structure

```
TCG_Server_GO/
├── main.go              # Main entry point
├── go.mod               # Dependency management
├── README.md            # Documentation
├── models/              # Data structures
│   └── user.go
├── auth/                # Authentication logic
│   ├── jwt.go           # JWT token handling
│   └── users.go         # User management
├── handlers/            # HTTP handlers
│   ├── auth.go          # Login handler
│   ├── validate.go      # Validation handler
│   ├── health.go        # Health check handler
│   └── routes.go        # Route configuration
└── middleware/          # Middlewares
    └── auth.go          # Authentication middleware
```

## Installation

1. Ensure you have Go installed (version 1.21 or higher)

2. Clone or download this project

3. Install dependencies:
```bash
go mod tidy
```

4. Run the server:
```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port using the `PORT` environment variable.

## API Endpoints

### POST /login
Authenticates user with username and password.

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
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

## Test Users

The server includes two sample users:

- **admin** / **admin123**
- **user** / **user123**

## Usage Examples

1. **Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

2. **Validate token:**
```bash
curl -X GET http://localhost:8080/api/validate \
  -H "Authorization: Bearer <token_from_login>"
```

## Modular Architecture

### Packages

- **`models`**: Defines data structures (LoginRequest, LoginResponse, etc.)
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

- **JWT Secret**: Change the `jwtSecret` variable in `auth/jwt.go` for production
- **Expiration**: Tokens expire in 24 hours by default
- **Passwords**: Stored hashed with bcrypt

## Production Considerations

- Use a real database instead of the user map
- Configure HTTPS
- Use environment variables for secrets
- Implement rate limiting
- Add appropriate logging
- Consider using an ORM like GORM for database operations
- Set up proper monitoring and health checks
- Implement proper error handling and logging
- Use secure session management
- Regular security audits and updates 