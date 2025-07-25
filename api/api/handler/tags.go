package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/sunba23/news/internal/news"
)

type TagsHandler struct {
	App news.App
}

func (h *TagsHandler) HandleGetAllTags(w http.ResponseWriter, r *http.Request) {
	repository := *h.App.Repository()
	tags, err := repository.GetAllTags(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("getting all tags has failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Error().Err(err).Msg("encoding tags to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *TagsHandler) HandleGetNewsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid news id", http.StatusBadRequest)
		return
	}
	repository := *h.App.Repository()
	news, err := repository.GetNewsByTag(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting news by tag id %v has failed", id))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(news); err != nil {
		log.Error().Err(err).Msg("encoding news to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
