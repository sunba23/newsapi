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

type NewsHandler struct {
	App news.App
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome 2 newsapp"))
}

func (h *NewsHandler) HandleGetAllNews(w http.ResponseWriter, r *http.Request) {
	repository := *h.App.Repository()
	news, err := repository.GetAllNews(r.Context())
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting all news has failed"))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(news); err != nil {
		log.Error().Err(err).Msg("encoding all news to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *NewsHandler) HandleGetNewsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid news id", http.StatusBadRequest)
		return
	}

	repository := *h.App.Repository()
	news, err := repository.GetNewsByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting news with id %v has failed", id))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(news); err != nil {
		log.Error().Err(err).Msg("encoding news to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *NewsHandler) HandleGetTagsForNews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid news id", http.StatusBadRequest)
		return
	}

	repository := *h.App.Repository()
	tags, err := repository.GetTagsForNews(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("getting tags for news with  id %v has failed", id))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tags); err != nil {
		log.Error().Err(err).Msg("encoding tags to JSON failed")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
