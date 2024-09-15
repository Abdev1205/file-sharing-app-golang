package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Abdev1205/21BCE11045_Backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extracting the Authorization header
		log.Println("Auth Header:", r.Header.Get("Authorization"))
		authHeader := r.Header.Get("Authorization")

		// Checking if the Authorization header is present and starts with "Bearer"
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// Extracting the token part, removing "Bearer " prefix
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parseing and validatinge the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTSecret()), nil
		})

		// Check if the token is valid
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extracting claims and then passing them to context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
