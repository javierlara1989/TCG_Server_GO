package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"tcg-server-go/auth"
	"tcg-server-go/database"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get user from database to check validation status
		user, err := database.GetUserByID(claims.UserID)
		if err != nil {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Check if user's email has been validated
		if user.ValidatedAt == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email not verified. Please verify your email before accessing this resource.",
			})
			return
		}

		// Set user information in headers for downstream handlers
		r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
		r.Header.Set("X-User-Email", claims.Email)

		next.ServeHTTP(w, r)
	})
}
