package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/constants"
	"github.com/sunba23/news/internal/news"
)

func NewUserContextMiddleware(store *sessions.CookieStore, app news.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "news-session")
			if err != nil {
				log.Warn().Err(err).Msg("Session error")
				next.ServeHTTP(w, r)
				return
			}

			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				next.ServeHTTP(w, r)
				return
			}

			if exp, ok := session.Values["expires_at"].(int64); ok && time.Now().Unix() > exp {
				next.ServeHTTP(w, r)
				return
			}

			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				next.ServeHTTP(w, r)
				return
			}

			repository := *app.Repository()
			user, err := repository.GetUserByGoogleID(r.Context(), userID)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("Failed to load user with id %v", userID))
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserIdContextKey, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
