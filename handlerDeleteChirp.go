package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samersawan/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access", err)
		return
	}
	chirp_id_string := r.PathValue("chirpID")
	chirp_id, err := uuid.Parse(chirp_id_string)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse chirp id into uuid", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found!", err)
	}
	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Unauthorized access", err)
		return
	}
	err = cfg.db.DeleteChirpByID(r.Context(), chirp_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
