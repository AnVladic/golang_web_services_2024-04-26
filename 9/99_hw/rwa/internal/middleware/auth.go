package middleware

import (
	"context"
	"net/http"
	"rwa/internal/handlers/fast_user"
	"strings"
)

func AuthMiddleware(userHandler *fast_user.UserHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Token ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Token ")
		session, ok := userHandler.SessionManager.Sessions[tokenString]
		if !ok {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), "session", &session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
