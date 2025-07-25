package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/constants"
	"github.com/sunba23/news/internal/news"
)

type UserHandler struct {
	App news.App
}

func (h *UserHandler) HandleAddFavoriteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid news id", http.StatusBadRequest)
		return
	}

	uid := r.Context().Value(constants.UserIdContextKey).(string)

	repository := *h.App.Repository()
	err = repository.AddFavoriteTag(r.Context(), uid, id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("adding tag for user %v failed", uid))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) HandleDeleteFavoriteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid news id", http.StatusBadRequest)
		return
	}

	uid := r.Context().Value(constants.UserIdContextKey).(string)

	repository := *h.App.Repository()
	err = repository.RemoveFavoriteTag(r.Context(), uid, id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("deleting tag for user %v failed", uid))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) HandleGetFavoriteTags(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(constants.UserIdContextKey).(string)

	repository := *h.App.Repository()
	tags, err := repository.GetFavoriteTags(r.Context(), uid)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting favorite tags for user %v failed", uid))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Error().Err(err).Msg("encoding tags to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) HandleGetFavoriteNews(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(constants.UserIdContextKey).(string)

	repository := *h.App.Repository()
	tags, err := repository.GetFavoriteNews(r.Context(), uid)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting favorite news for user %v failed", uid))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Error().Err(err).Msg("encoding news to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
