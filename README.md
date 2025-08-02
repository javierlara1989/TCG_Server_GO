# JWT Authentication Server in Go

A modular Go server that handles JWT authentication with MariaDB database integration, email verification system, and game user management, organized in separate packages for better maintainability and production readiness.

## Features

- User registration with comprehensive validation
- **Email verification system** with validation codes
- **Game user management** with level, experience, and money tracking
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
- **Email verification workflow** with expiration and resend functionality
- **Game progression system** with automatic level up and rewards

## Game User Management

### UserInfo Model
Each user account has an associated `UserInfo` record that contains:
- **Level**: Current player level (starts at 1)
- **Experience**: Current experience points (starts at 0)
- **Money**: Current money balance (starts at 100)

### Game Features
- **Automatic level up**: When experience reaches level * 1000
- **Level up rewards**: Bonus money when leveling up
- **Money management**: Add and spend money with validation
- **Experience tracking**: Add experience points with automatic progression

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

### Validation Code
- **6-character alphanumeric code**
- **Case-insensitive** matching
- **Required for email verification**

### Game Values
- **Level**: Minimum 1
- **Experience**: Minimum 0
- **Money**: Minimum 0

## Project Structure

```
TCG_Server_GO/
├── main.go              # Main entry point
├── go.mod               # Dependency management
├── README.md            # Documentation
├── ENVIRONMENT.md       # Environment variables documentation
├── EMAIL_VERIFICATION.md # Email verification documentation
├── models/              # Data structures
│   ├── user.go
│   └── user_info.go     # Game user information model
├── database/            # Database operations
│   ├── database.go      # Database connection and configuration
│   ├── users.go         # User database operations
│   └── user_info.go     # UserInfo database operations
├── auth/                # Authentication logic
│   ├── jwt.go           # JWT token handling
│   └── users.go         # User management
├── handlers/            # HTTP handlers
│   ├── auth.go          # Login, register, and email verification handlers
│   ├── user_info.go     # Game user info handlers
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
    validation_code VARCHAR(255) NULL,
    validation_code_expires_at TIMESTAMP NULL,
    validated_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### User Info Table

```sql
CREATE TABLE user_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 1,
    experience INT NOT NULL DEFAULT 0,
    money INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
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

### Authentication Endpoints

#### POST /register
Registers a new user with validation. Generates a validation code that should be sent via email.

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

#### POST /verify-email
Verifies a user's email with the provided validation code.

**Request:**
```json
{
  "email": "juan@example.com",
  "validation_code": "A1B2C3"
}
```

**Response:**
```json
{
  "message": "Email verified successfully",
  "user_id": 1
}
```

#### POST /resend-code
Resends a new validation code to the user's email.

**Request:**
```json
{
  "email": "juan@example.com"
}
```

**Response:**
```json
{
  "message": "Validation code sent successfully"
}
```

#### POST /login
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

### Game User Info Endpoints (All require authentication)

#### GET /api/user-info
Retrieves the current user's game information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_info": {
    "id": 1,
    "user_id": 1,
    "level": 1,
    "experience": 0,
    "money": 100,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "User info retrieved successfully"
}
```

**Important:** User game information (level, experience, money) is **read-only** and can only be modified through server-side game logic during actual gameplay. This ensures complete game integrity and prevents any form of cheating or manipulation.

### Other Endpoints

#### GET /api/validate
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

#### GET /health
Checks server status.

**Response:**
```json
{
  "status": "OK"
}
```

## Database Operations

The application includes comprehensive database operations:

### User Operations
- **CreateUser**: Create new users with validation codes
- **GetUserByEmail**: Retrieve user by email
- **GetUserByID**: Retrieve user by ID
- **VerifyEmail**: Verify user email with validation code
- **ResendValidationCode**: Generate and send new validation code
- **UpdateUser**: Update user information
- **UpdatePassword**: Update user password
- **SoftDeleteUser**: Mark user as deleted (soft delete)
- **HardDeleteUser**: Permanently delete user
- **EmailExists**: Check if email already exists

### UserInfo Operations
- **CreateUserInfo**: Create new user info record
- **GetUserInfoByUserID**: Retrieve user info by user ID
- **GetUserInfoByID**: Retrieve user info by its own ID
- **UpdateUserInfo**: Update user info (internal use only)
- **UpdateUserInfoPartial**: Update specific fields (internal use only)
- **AddExperience**: Add experience with automatic level up (internal use only)
- **AddMoney**: Add money to user account (internal use only)
- **SpendMoney**: Spend money with validation (internal use only)
- **DeleteUserInfo**: Delete user info
- **UserInfoExists**: Check if user info exists
- **CreateDefaultUserInfo**: Create default user info for new users

**Note:** Game modification functions (AddExperience, AddMoney, SpendMoney) are only available internally for server-side game logic and cannot be accessed directly by clients.

## Usage Examples

### 1. Complete User Registration and Game Setup
```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "María García",
    "email": "maria@example.com",
    "password": "password123"
  }'

# Verify email (after receiving code via email)
curl -X POST http://localhost:8080/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "validation_code": "A1B2C3"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria@example.com",
    "password": "password123"
  }'
```

### 2. Game Operations (using token from login)
```bash
# Get user game info (read-only)
curl -X GET http://localhost:8080/api/user-info \
  -H "Authorization: Bearer <token_from_login>"
```

**Important:** User game information (level, experience, money) is **read-only** and controlled entirely by server-side game logic. Clients can only:
- View their current game information

All modifications to game stats (experience, money, level) happen automatically during gameplay through server-side logic to maintain complete game integrity and prevent any form of cheating. 