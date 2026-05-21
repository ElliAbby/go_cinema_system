package http

import (
	"context"
	"net/http"
	"strings"

	authjwt "github.com/ElliAbby/go_cinema_system/internal/platform/auth/jwt"

)

type contextKey string

const userIDContextKey contextKey = "user_id"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func JWTMiddleware(tokenManager *authjwt.Manager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, map[string]interface{}{
                "code":    "UNAUTHORIZED",
                "message": "Missing authorization header",
            }, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, map[string]interface{}{
						"code":    "UNAUTHORIZED",
						"message": "Invalid authorization header format",
					}, http.StatusUnauthorized)
					return
		}

		tokenString := parts[1]

		claims, err := tokenManager.ValidateToken(tokenString)
		if err != nil {
			respondError(w, map[string]interface{}{
						"code":    "UNAUTHORIZED",
						"message": "Invalid or expired token",
					}, http.StatusUnauthorized)
					return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDContextKey).(int)
	return userID, ok
}