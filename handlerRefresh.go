package main

import (
	"net/http"
	"time"

	"github.com/samersawan/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Token string `json:"token"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve token", err)
		return
	}
	token, err := cfg.db.GetToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found or expired", err)
	}
	new_token, err := auth.MakeJWT(token.UserID, cfg.secret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new JWT token", err)
	}
	respondWithJSON(w, http.StatusOK, resp{Token: new_token})
}
