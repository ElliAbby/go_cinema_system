package http

import (
	"context"
	"net/http"
	"strings"

	authjwt "github.com/ElliAbby/go_cinema_system/internal/platform/auth/jwt"

)

type contextKey string

const userIDContextKey contextKey = "user_id"

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