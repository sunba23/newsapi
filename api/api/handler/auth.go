package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/internal/database"
	"github.com/sunba23/news/internal/news"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	App          news.App
	Oauth        oauth2.Config
	SessionStore *sessions.CookieStore
}

func NewAuthHandler(app news.App) *AuthHandler {
	authHandler := AuthHandler{
		App:          app,
		Oauth:        *news.OauthConfigFromConfig(*app.Config()),
		SessionStore: sessions.NewCookieStore([]byte(app.Config().SessionSecret)),
	}
	return &authHandler
}

func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateStateToken()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	session, _ := h.SessionStore.Get(r, "news-session")
	session.Values["oauth_state"] = state
	if err := session.Save(r, w); err != nil {
		log.Error().Err(err).Msg("Session save failed")
	}

	url := h.Oauth.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var googleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, "news-session")
	storedState, ok := session.Values["oauth_state"].(string)

	if !ok || storedState != r.URL.Query().Get("state") {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.Oauth.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	client := h.Oauth.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&googleUser)

	log.Debug().Msg(fmt.Sprintf("Successfully authenticated as %v with id %v", googleUser.Email, googleUser.ID))

	// upsert user in DB
	user := &database.User{
		GoogleID: googleUser.ID,
		Email:    googleUser.Email,
	}
	repository := *h.App.Repository()
	if err := repository.UpsertUser(r.Context(), user); err != nil {
		log.Error().Err(err).Msg("Failed to save user in database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = googleUser.ID
	session.Values["email"] = googleUser.Email

	if err := session.Save(r, w); err != nil {
		log.Error().Err(err).Msg("Session save failed")
		http.Error(w, "Session save failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, "news-session")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
