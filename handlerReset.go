package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Platform is ", cfg.platform)
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Unauthorized endpoint", nil)
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete users table", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
