package auth

import (
	"fmt"
	"os"
	"time"

	"tcg-server-go/database"
	"tcg-server-go/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("mi_clave_secreta_muy_segura")

func init() {
	// Try to get JWT secret from environment variable
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		jwtSecret = []byte(secret)
	}
}

func GenerateToken(email string) (string, error) {
	// Get user from database to get user ID
	var userID int
	if database.DB != nil {
		user, err := database.GetUserByEmail(email)
		if err != nil || user == nil {
			return "", fmt.Errorf("user not found")
		}
		userID = user.ID
	} else {
		// Fallback for testing
		userID = 1 // Default user ID for testing
	}

	claims := models.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// SetJWTSecret allows changing the secret key (useful for testing)
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}
