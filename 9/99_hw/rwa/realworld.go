package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"rwa/internal/db"
	internal "rwa/internal/handlers"
	"rwa/internal/handlers/article"
	"rwa/internal/handlers/fast_user"
	"rwa/internal/middleware"
	"rwa/internal/session"
	"sync"
)

func setAuthMiddleware(handler *fast_user.UserHandler, h http.HandlerFunc) http.Handler {
	return middleware.AuthMiddleware(handler, h)
}

func GetApp() http.Handler {
	r := mux.NewRouter()

	rwaHandler := internal.RwaHandler{
		SessionManager: session.Manager{
			Sessions: map[string]db.Session{},
			Mu:       &sync.Mutex{},
		},
	}
	userHandler := fast_user.UserHandler{
		Users:      map[string]*db.User{},
		UserMu:     sync.Mutex{},
		RwaHandler: rwaHandler,
	}
	articleHandler := article.ArticleHandler{
		RwaHandler: rwaHandler,
	}

	r.HandleFunc("/api/users", userHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/users/login", userHandler.Login).Methods(http.MethodPost)
	r.Handle(
		"/api/user/logout", setAuthMiddleware(&userHandler, userHandler.Logout)).Methods(http.MethodPost)
	r.Handle(
		"/api/user", setAuthMiddleware(&userHandler, userHandler.CurrentUser)).Methods(http.MethodGet)
	r.Handle(
		"/api/user", setAuthMiddleware(&userHandler, userHandler.UpdateUser)).Methods(http.MethodPut)

	r.Handle(
		"/api/articles", setAuthMiddleware(&userHandler, articleHandler.CreateArticle),
	).Methods(http.MethodPost)
	r.HandleFunc("/api/articles", articleHandler.GetArticles).Methods(http.MethodGet)
	return r
}
