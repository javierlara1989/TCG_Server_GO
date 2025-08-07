# Tables API Documentation

## General Description

The table system allows users to create game tables for TCG matches. Each table can be public or private, and can have different categories and prizes.

## Data Models

### Table
- `id`: Unique identifier for the table
- `category`: Table category (S, A, B, C, D)
- `privacy`: Table privacy (private, public)
- `password`: Optional numeric password (maximum 10 digits)
- `prize`: Prize type (money, card, aura)
- `amount`: Bet amount (money or cards) - optional integer
- `winner`: Indicates if there's a winner (TRUE, FALSE, NULL)
- `created_at`: Creation date
- `updated_at`: Last update date
- `finished_at`: Completion date (optional)

### UserTable
- `id`: Unique identifier for the association
- `user_id`: ID of the table owner user
- `rival_id`: ID of the rival (NULL if waiting for rival)
- `table_id`: ID of the associated table
- `time`: Time elapsed in the match in seconds (integer)

## Endpoints

### 1. Create Table
**POST** `/api/tables`

Creates a new table and associates it with the authenticated user. The table is created in waiting state (without rival).

#### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "category": "A",
  "privacy": "public",
  "password": "1234",
  "prize": "money",
  "amount": 1000
}
```

#### Required Fields
- `category`: Must be S, A, B, C, or D
- `privacy`: Must be "private" or "public"
- `prize`: Must be "money", "card", or "aura"

#### Optional Fields
- `password`: Numeric password (maximum 10 digits)
- `amount`: Bet amount (positive integer)

#### Response (201 Created)
```json
{
  "message": "Table created successfully",
  "table_id": 1
}
```

### 2. Get User Tables
**GET** `/api/tables`

Gets all tables associated with the authenticated user (as owner or rival).

#### Headers
```
Authorization: Bearer <token>
```

#### Response (200 OK)
```json
{
  "tables": [
    {
      "id": 1,
      "user_id": 1,
      "rival_id": null,
      "table_id": 1,
      "time": 0,
      "table": {
         "id": 1,
         "category": "A",
         "privacy": "public",
         "password": null,
         "prize": "money",
         "amount": 1000,
         "winner": null,
         "created_at": "2024-01-01T12:00:00Z",
         "updated_at": "2024-01-01T12:00:00Z",
         "finished_at": null
       },
      "user_name": "User1",
      "user_email": "user1@example.com",
      "rival_name": null,
      "rival_email": null
    }
  ]
}
```

### 3. Update Table
**PUT** `/api/tables/{id}`

Updates table parameters. Can only be updated if:
- The user is the table owner
- The table is waiting for rival (rival_id is NULL)

#### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "category": "B",
  "privacy": "private",
  "password": "5678",
  "prize": "card",
  "amount": 500
}
```

#### Modifiable Fields
- `category`: New category (S, A, B, C, D)
- `privacy`: New privacy (private, public)
- `password`: New password (numeric, maximum 10 digits)
- `prize`: New prize (money, card, aura)
- `amount`: New bet amount (positive integer)

#### Response (200 OK)
```json
{
  "message": "Table updated successfully",
  "table_id": 1
}
```

### 4. Update Table Time
**PUT** `/api/tables/{id}/time`

Updates the time elapsed in the match for a specific table. Can only be updated if the user is the table owner.

#### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "time": 120
}
```

#### Required Fields
- `time`: Time elapsed in seconds (non-negative integer)

#### Response (200 OK)
```json
{
  "message": "Table time updated successfully",
  "time": 120
}
```

## Validations

### Valid Categories
- S, A, B, C, D

### Valid Privacy
- private, public

### Valid Prizes
- money, card, aura

### Password
- Maximum 10 characters
- Only numeric digits (0-9)
- Optional

### Amount
- Positive integer
- Represents the amount of money or cards bet
- Optional

## Table States

1. **Waiting for Rival**: `rival_id` is NULL
   - Parameters can be modified
   - The table is available for other users to join

2. **With Rival**: `rival_id` has a value
   - Parameters cannot be modified
   - The match can begin

3. **Finished**: `finished_at` has a value
   - `winner` indicates the result
   - The table is closed

## Error Codes

- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: You don't have permission to perform the action
- `500 Internal Server Error`: Internal server error

## Usage Examples

### Create a public table
```bash
curl -X POST http://localhost:8080/api/tables \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "A",
    "privacy": "public",
    "prize": "money",
    "amount": 1000
  }'
```

### Create a private table with password
```bash
curl -X POST http://localhost:8080/api/tables \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "S",
    "privacy": "private",
    "password": "123456",
    "prize": "card",
    "amount": 500
  }'
```

### Update table parameters
```bash
curl -X PUT http://localhost:8080/api/tables/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "B",
    "privacy": "private",
    "password": "9999",
    "amount": 2000
  }'
```

### Update table time
```bash
curl -X PUT http://localhost:8080/api/tables/1/time \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "time": 180
  }'
``` 