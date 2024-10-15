package user

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(userHandler Handler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Token ") {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Token ")
		session, ok := userHandler.SessionManager.Sessions[tokenString]
		if ok {
			ctx := context.WithValue(r.Context(), "session", session)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
