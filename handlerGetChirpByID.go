package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirp_id_string := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(chirp_id_string)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse chirp id into uuid", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{ID: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserID: chirp.UserID})
}
