package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sunba23/news/constants"
)

func NewAuthenticationMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(constants.UserIdContextKey)
			if userId == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
