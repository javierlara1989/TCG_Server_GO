# Email Verification System

## Overview

The email verification system ensures that users can only access their accounts after confirming ownership of their email address. This prevents fake registrations and improves security.

## How It Works

### 1. User Registration Flow
```
User submits registration → System creates account → Validation code generated → Email sent → User receives code
```

### 2. Email Verification Flow
```
User enters code → System validates code → Account marked as verified → User can now login
```

### 3. Code Resend Flow
```
User requests new code → System generates new code → Email sent → User receives new code
```

## Database Schema

### New Fields Added to Users Table

| Field | Type | Description |
|-------|------|-------------|
| `validation_code` | VARCHAR(255) | 6-character alphanumeric code |
| `validation_code_expires_at` | TIMESTAMP | When the code expires (24 hours) |
| `validated_at` | TIMESTAMP | When the email was verified |

### Example User Record
```sql
INSERT INTO users (
    name, 
    email, 
    password, 
    validation_code, 
    validation_code_expires_at, 
    validated_at
) VALUES (
    'John Doe',
    'john@example.com',
    '$2a$10$...', -- hashed password
    'A1B2C3',
    '2024-01-15 10:30:00',
    NULL -- will be set when verified
);
```

## API Endpoints

### POST /register
Creates a new user account and generates a validation code.

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

**What happens:**
1. Validates input data
2. Checks if email already exists
3. Hashes password with bcrypt
4. Generates 6-character validation code
5. Sets expiration time (24 hours from now)
6. Creates user record in database
7. Returns JWT token

### POST /verify-email
Verifies user's email with the provided validation code.

**Request:**
```json
{
  "email": "john@example.com",
  "validation_code": "A1B2C3"
}
```

**Success Response:**
```json
{
  "message": "Email verified successfully",
  "user_id": 1
}
```

**Error Responses:**
```json
{
  "error": "user not found"
}
```
```json
{
  "error": "invalid validation code"
}
```
```json
{
  "error": "validation code has expired"
}
```
```json
{
  "error": "email already verified"
}
```

### POST /resend-code
Generates and sends a new validation code.

**Request:**
```json
{
  "email": "john@example.com"
}
```

**Success Response:**
```json
{
  "message": "Validation code sent successfully"
}
```

**Error Responses:**
```json
{
  "error": "user not found"
}
```
```json
{
  "error": "email already verified"
}
```

## Validation Code Generation

### Algorithm
- **Length**: 6 characters
- **Characters**: Alphanumeric (A-Z, 0-9)
- **Case**: Uppercase
- **Randomness**: Cryptographically secure using `crypto/rand`

### Example Code Generation
```go
func generateValidationCode() string {
    bytes := make([]byte, 3)
    rand.Read(bytes)
    return strings.ToUpper(hex.EncodeToString(bytes)[:6])
}
```

**Example codes**: `A1B2C3`, `F9E8D7`, `123ABC`

## Security Features

### 1. Code Expiration
- **Duration**: 24 hours from generation
- **Automatic cleanup**: Expired codes are ignored
- **Resend capability**: Users can request new codes

### 2. One-Time Use
- **Single verification**: Each code can only be used once
- **Automatic deletion**: Code is removed after successful verification
- **Prevents replay attacks**: Old codes cannot be reused

### 3. Case Insensitive
- **Flexible input**: Users can enter codes in any case
- **Normalized storage**: Codes are stored in uppercase
- **User-friendly**: Reduces input errors

### 4. Duplicate Prevention
- **Already verified check**: Prevents multiple verifications
- **Status tracking**: `validated_at` field tracks verification status
- **Clear error messages**: Users know if already verified

## Email Integration (To Be Implemented)

### Current State
The system generates and stores validation codes but **does not send emails**. You need to implement email sending functionality.

### Recommended Email Services
1. **SendGrid** - Popular, reliable, good free tier
2. **AWS SES** - Cost-effective, scalable
3. **SMTP** - Direct email server integration
4. **Mailgun** - Developer-friendly API

### Implementation Steps

#### 1. Add Email Service Dependencies
```bash
go get github.com/sendgrid/sendgrid-go
# or
go get github.com/aws/aws-sdk-go-v2/service/ses
```

#### 2. Create Email Service
```go
// services/email.go
package services

type EmailService interface {
    SendValidationCode(email, code string) error
}

type SendGridService struct {
    apiKey string
}

func (s *SendGridService) SendValidationCode(email, code string) error {
    // Implementation here
}
```

#### 3. Update Registration Handler
```go
// In handlers/auth.go
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Send validation code via email
    if err := emailService.SendValidationCode(user.Email, *user.ValidationCode); err != nil {
        log.Printf("Failed to send validation code: %v", err)
        // Don't fail registration, just log the error
    }
    
    // ... rest of code ...
}
```

#### 4. Add Environment Variables
```bash
# For SendGrid
SENDGRID_API_KEY=your_api_key_here
SENDGRID_FROM_EMAIL=noreply@yourdomain.com

# For AWS SES
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
```

### Email Template Example
```html
<!DOCTYPE html>
<html>
<head>
    <title>Verify Your Email</title>
</head>
<body>
    <h1>Welcome to Our Service!</h1>
    <p>Please verify your email address by entering this code:</p>
    <h2 style="font-size: 24px; color: #007bff;">{{.Code}}</h2>
    <p>This code will expire in 24 hours.</p>
    <p>If you didn't create an account, please ignore this email.</p>
</body>
</html>
```

## Testing the System

### Without Email Integration
1. **Register a user** and note the validation code from database
2. **Use the code** to verify the email
3. **Test expiration** by waiting or manually updating expiration time
4. **Test resend** functionality

### Database Queries for Testing
```sql
-- Check user registration
SELECT id, email, validation_code, validation_code_expires_at, validated_at 
FROM users WHERE email = 'test@example.com';

-- Manually expire a code (for testing)
UPDATE users 
SET validation_code_expires_at = DATE_SUB(NOW(), INTERVAL 1 HOUR)
WHERE email = 'test@example.com';

-- Check verification status
SELECT email, validated_at IS NOT NULL as is_verified 
FROM users WHERE email = 'test@example.com';
```

## Best Practices

### 1. Rate Limiting
- **Limit resend requests**: Prevent abuse
- **Cooldown period**: Wait between resend attempts
- **IP-based limits**: Prevent spam

### 2. Error Handling
- **Generic error messages**: Don't reveal if email exists
- **Graceful degradation**: System works even if email fails
- **Logging**: Track email delivery success rates

### 3. User Experience
- **Clear instructions**: Tell users to check their email
- **Resend option**: Allow users to request new codes
- **Expiration warnings**: Notify users about code expiration

### 4. Security
- **Secure code generation**: Use cryptographically secure random
- **Code expiration**: Prevent long-term code validity
- **Input validation**: Validate all user inputs
- **SQL injection prevention**: Use parameterized queries

## Future Enhancements

### 1. Advanced Features
- **Email change verification**: Verify new email addresses
- **Phone verification**: Add SMS verification option
- **Two-factor authentication**: Use email as 2FA method

### 2. Monitoring
- **Delivery tracking**: Monitor email delivery rates
- **Verification analytics**: Track verification completion rates
- **Error monitoring**: Alert on email service failures

### 3. Performance
- **Code caching**: Cache frequently used codes
- **Batch processing**: Process multiple verifications
- **Async email sending**: Don't block user registration 