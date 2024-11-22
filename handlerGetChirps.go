package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch chirps", err)
	}
	resp := make([]Chirp, 0, len(chirps))
	for _, v := range chirps {
		resp = append(resp, Chirp{ID: v.ID, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt, Body: v.Body, UserID: v.UserID})
	}
	respondWithJSON(w, http.StatusOK, resp)
}
