package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"tcg-server-go/models"
)

var jwtSecret = []byte("mi_clave_secreta_muy_segura")

func GenerateToken(username string) (string, error) {
	claims := models.Claims{
		Username: username,
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