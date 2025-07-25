package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sunba23/news/api/handler"
	"github.com/sunba23/news/api/middleware"
	"github.com/sunba23/news/internal/news"
)

func NewHttpHandler(app news.App) http.Handler {
	router := mux.NewRouter()

	authHandler := handler.NewAuthHandler(app)
	newsHandler := handler.NewsHandler{App: app}
	tagsHandler := handler.TagsHandler{App: app}
	userHandler := handler.UserHandler{App: app}

	authenticationMiddleware := middleware.NewAuthenticationMiddleware()
	userContextMiddleware := middleware.NewUserContextMiddleware(authHandler.SessionStore, app)

	router.Use(middleware.LoggingMiddleware, userContextMiddleware)
	router.HandleFunc("/", handler.HandleRoot)

	authSubRouter := router.PathPrefix("/auth/google").Subrouter()
	authSubRouter.HandleFunc("/login", authHandler.HandleGoogleLogin)
	authSubRouter.HandleFunc("/callback", authHandler.HandleGoogleCallback)
	authSubRouter.HandleFunc("/logout", authHandler.HandleLogout)

	newsSubRouter := router.PathPrefix("/news").Subrouter()
	newsSubRouter.HandleFunc("", newsHandler.HandleGetAllNews).Methods(http.MethodGet)
	newsSubRouter.HandleFunc("/{id:[0-9]+}", newsHandler.HandleGetNewsById).Methods(http.MethodGet)
	newsSubRouter.HandleFunc("/{id:[0-9]+}/tags", newsHandler.HandleGetTagsForNews).Methods(http.MethodGet)
	newsSubRouter.Use(authenticationMiddleware)

	tagsSubRouter := router.PathPrefix("/tags").Subrouter()
	tagsSubRouter.HandleFunc("", tagsHandler.HandleGetAllTags).Methods(http.MethodGet)
	tagsSubRouter.HandleFunc("/{id:[0-9]+}/news", tagsHandler.HandleGetNewsByTag).Methods(http.MethodGet)
	tagsSubRouter.Use(authenticationMiddleware)

	userSubRouter := router.PathPrefix("/user").Subrouter()
	userSubRouter.HandleFunc("/tags", userHandler.HandleGetFavoriteTags).Methods(http.MethodGet)
	userSubRouter.HandleFunc("/tags/{id:[0-9]+}", userHandler.HandleAddFavoriteTag).Methods(http.MethodPost)
	userSubRouter.HandleFunc("/tags/{id:[0-9]+}", userHandler.HandleDeleteFavoriteTag).Methods(http.MethodDelete)
	userSubRouter.HandleFunc("/news", userHandler.HandleGetFavoriteNews).Methods(http.MethodGet)
	userSubRouter.Use(authenticationMiddleware)

	return router
}
