package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samersawan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "asc"
	}
	var chirps []database.Chirp
	var err error
	if authorID == "" {
		if sort == "asc" {
			chirps, err = cfg.db.GetAllChirpsAsc(r.Context())
		} else {
			chirps, err = cfg.db.GetAllChirpsDesc(r.Context())
		}
	} else {
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to parse uuid string", err)
			return
		}
		if sort == "asc" {
			chirps, err = cfg.db.GetChirpsByAuthorAsc(r.Context(), userID)
		} else {
			chirps, err = cfg.db.GetChirpsByAuthorDesc(r.Context(), userID)
		}

	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch chirps", err)
		return
	}
	resp := make([]Chirp, 0, len(chirps))
	for _, v := range chirps {
		resp = append(resp, Chirp{ID: v.ID, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, Body: v.Body, UserID: v.UserID})
	}
	respondWithJSON(w, http.StatusOK, resp)
}
