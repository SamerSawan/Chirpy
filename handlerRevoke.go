package main

import (
	"net/http"

	"github.com/samersawan/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve token", err)
		return
	}
	err = cfg.db.RevokeToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
