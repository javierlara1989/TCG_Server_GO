# JWT Authentication Server in Go

A modular Go server that handles JWT authentication with MariaDB database integration, email verification system, and game user management, organized in separate packages for better maintainability and production readiness.

## Features

- User registration with comprehensive validation
- **Email verification system** with validation codes
- **Game user management** with level, experience, and money tracking
- **Card management system** with types, elements, and legends
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

### Name
- **Minimum 6 characters**
- **Only letters and spaces allowed**
- **Supports international characters** (á, é, í, ó, ú, ñ)

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

## API Endpoints

### Authentication Endpoints

#### POST /register
Registers a new user with validation. Generates a validation code that should be sent via email.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
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
  "email": "john@example.com",
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
  "email": "john@example.com"
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
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Card Endpoints (Public Access - Read Only)

**Important:** Cards are managed via server-side seeds and cannot be modified through the API. All card endpoints are read-only to ensure data integrity and prevent unauthorized modifications.

#### GET /cards
Retrieves all cards from the database.

**Response:**
```json
{
  "cards": [
    {
      "id": 1,
      "name": "Dragon Warrior",
      "type": "Monster",
      "legend": "A powerful dragon warrior with fire abilities",
      "element": "Fire",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Cards retrieved successfully"
}
```

#### GET /cards/{id}
Retrieves a specific card by ID.

**Response:**
```json
{
  "card": {
    "id": 1,
    "name": "Dragon Warrior",
    "type": "Monster",
    "legend": "A powerful dragon warrior with fire abilities",
    "element": "Fire",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "message": "Card retrieved successfully"
}
```

#### GET /cards/type/{type}
Retrieves all cards of a specific type (Monster, Spell, or Energy).

**Response:**
```json
{
  "cards": [
    {
      "id": 1,
      "name": "Dragon Warrior",
      "type": "Monster",
      "legend": "A powerful dragon warrior with fire abilities",
      "element": "Fire",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Cards retrieved successfully"
}
```

#### GET /cards/element/{element}
Retrieves all cards of a specific element.

**Response:**
```json
{
  "cards": [
    {
      "id": 1,
      "name": "Dragon Warrior",
      "type": "Monster",
      "legend": "A powerful dragon warrior with fire abilities",
      "element": "Fire",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Cards retrieved successfully"
}
```

#### GET /cards/search?q={search_term}
Searches for cards by name (partial match).

**Response:**
```json
{
  "cards": [
    {
      "id": 1,
      "name": "Dragon Warrior",
      "type": "Monster",
      "legend": "A powerful dragon warrior with fire abilities",
      "element": "Fire",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "message": "Cards found successfully"
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

### User Cards Endpoints (All require authentication)

#### GET /api/user-cards
Retrieves all cards in the authenticated user's inventory.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user_cards": [
    {
      "id": 1,
      "user_id": 1,
      "card_id": 1,
      "amount": 3,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "card": {
        "id": 1,
        "name": "Dragon Warrior",
        "type": "Monster",
        "legend": "A powerful dragon warrior with fire abilities",
        "element": "Fire",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    }
  ],
  "message": "User cards retrieved successfully"
}
```

#### GET /api/user-cards/{id}
Retrieves a specific card from the user's inventory by card ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "user_card": {
    "id": 1,
    "user_id": 1,
    "card_id": 1,
    "amount": 3,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "card": {
      "id": 1,
      "name": "Dragon Warrior",
      "type": "Monster",
      "legend": "A powerful dragon warrior with fire abilities",
      "element": "Fire",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  },
  "message": "User card retrieved successfully"
}
```

**Important:** User card management (adding, removing, updating amounts) is handled internally by the server during gameplay. All card modifications are controlled by server-side logic to ensure game integrity and prevent any form of cheating or manipulation.

### Deck Endpoints (All require authentication)

#### GET /api/decks
Retrieves all decks for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "decks": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Fire Dragon Deck",
      "valid": true
    }
  ],
  "message": "Decks retrieved successfully"
}
```

#### GET /api/decks/limit
Retrieves deck limit information for the authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "current_decks": 2,
  "deck_limit": 3,
  "user_level": 1,
  "can_create": true,
  "message": "Deck limit information retrieved successfully"
}
```

#### POST /api/decks
Creates a new deck. The deck will only be created if the user has all the required cards in their inventory.

**Restrictions:**
- **Minimum 40 cards**: Each deck must contain at least 40 cards
- **Card ownership**: User must own all cards in the deck with sufficient quantities
- **Deck limit**: Users can have 3 decks + (level / 25) additional decks (rounded down)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Fire Dragon Deck",
  "card_ids": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
  "card_count": [4, 3, 2, 4, 3, 2, 4, 3, 2, 13]
}
```

**Required Fields:**
- `name`: Deck name (1-100 characters)
- `card_ids`: Array of card IDs to include in the deck
- `card_count`: Array of quantities for each card (must match card_ids length)

**Response (201 Created):**
```json
{
  "deck": {
    "id": 1,
    "user_id": 1,
    "name": "Fire Dragon Deck",
    "valid": true
  },
  "message": "Deck created successfully"
}
```

**Error Responses (400 Bad Request):**
```json
{
  "error": "Cannot create deck: you do not have all the required cards"
}
```
```json
{
  "error": "Cannot create deck: deck must have at least 40 cards"
}
```
```json
{
  "error": "Cannot create deck: deck limit reached: you can only have 3 decks"
}
```

#### GET /api/decks/{id}
Retrieves a specific deck by ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "deck": {
    "id": 1,
    "user_id": 1,
    "name": "Fire Dragon Deck",
    "valid": true
  },
  "message": "Deck retrieved successfully"
}
```

#### GET /api/decks/{id}/cards
Retrieves a deck with all its cards.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "deck_with_cards": {
    "deck": {
      "id": 1,
      "user_id": 1,
      "name": "Fire Dragon Deck",
      "valid": true
    },
    "cards": [
      {
        "deck_id": 1,
        "card_id": 1,
        "number": 2,
        "card": {
          "id": 1,
          "name": "Dragon Warrior",
          "type": "Monster",
          "legend": "A powerful dragon warrior with fire abilities",
          "element": "Fire",
          "created_at": "2024-01-15T10:30:00Z",
          "updated_at": "2024-01-15T10:30:00Z"
        }
      }
    ]
  },
  "message": "Deck with cards retrieved successfully"
}
```

#### DELETE /api/decks/{id}
Deletes a deck and all its cards.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Deck deleted successfully"
}
```

**Important:** Deck creation is validated to ensure the user has all required cards in their inventory. The `valid` field indicates whether the deck can be used in gameplay.

**Deck System Rules:**
- **Minimum cards**: Each deck must contain at least 40 cards
- **Card ownership**: Users can only include cards they own in sufficient quantities
- **Deck limit**: Base limit of 3 decks + 1 additional deck per 25 levels
  - Level 1-24: 3 decks
  - Level 25-49: 4 decks
  - Level 50-74: 5 decks
  - Level 75-99: 6 decks
  - And so on...

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

### Cards Table

```sql
CREATE TABLE cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type ENUM('Monster', 'Spell', 'Energy') NOT NULL,
    legend TEXT NOT NULL,
    element ENUM('Fire', 'Water', 'Wind', 'Earth', 'Neutral', 'Holy', 'Dark') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_element (element),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### User Cards Table

```sql
CREATE TABLE user_cards (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    card_id INT NOT NULL,
    amount INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_card_id (card_id),
    UNIQUE KEY unique_user_card (user_id, card_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Decks Table

```sql
CREATE TABLE decks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_valid (valid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Deck Cards Table

```sql
CREATE TABLE deck_cards (
    deck_id INT NOT NULL,
    card_id INT NOT NULL,
    number INT NOT NULL DEFAULT 1,
    PRIMARY KEY (deck_id, card_id),
    FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    INDEX idx_deck_id (deck_id),
    INDEX idx_card_id (card_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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

### Card Operations
- **GetCardByID**: Retrieve card by ID
- **GetCardByName**: Retrieve card by name
- **GetAllCards**: Retrieve all cards
- **GetCardsByType**: Retrieve cards by type (Monster/Spell/Energy)
- **GetCardsByElement**: Retrieve cards by element
- **SearchCards**: Search cards by name (partial match)
- **CreateCard**: Create new card record (internal use only - for seeds)
- **UpdateCard**: Update card information (internal use only)
- **UpdateCardPartial**: Update specific card fields (internal use only)
- **DeleteCard**: Delete card by ID (internal use only)
- **CardExists**: Check if card exists by ID
- **CardNameExists**: Check if card exists by name

**Note:** Card modification functions (CreateCard, UpdateCard, DeleteCard) are only available internally for server-side seeds and cannot be accessed directly by clients.

## Usage Examples

### 1. Complete User Registration and Game Setup
```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mary Garcia",
    "email": "mary@example.com",
    "password": "password123"
  }'

# Verify email (after receiving code via email)
curl -X POST http://localhost:8080/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mary@example.com",
    "validation_code": "A1B2C3"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mary@example.com",
    "password": "password123"
  }'
```

### 2. Game Operations (using token from login)
```bash
# Get user game info (read-only)
curl -X GET http://localhost:8080/api/user-info \
  -H "Authorization: Bearer <token_from_login>"
```

### 3. Card Operations (Read-only)
```bash
# Get all cards
curl -X GET http://localhost:8080/cards

# Get cards by type
curl -X GET http://localhost:8080/cards/type/Monster

# Get cards by element
curl -X GET http://localhost:8080/cards/element/Fire

# Search cards
curl -X GET "http://localhost:8080/cards/search?q=dragon"

# Get specific card by ID
curl -X GET http://localhost:8080/cards/1
```

**Important:** Cards are managed via server-side seeds and cannot be modified through the API. All card operations are read-only to ensure data integrity.