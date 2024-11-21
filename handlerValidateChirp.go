package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Clean string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	split_body := strings.Split(params.Body, " ")
	for i, v := range split_body {
		lower_string := strings.ToLower(v)
		if lower_string == "kerfuffle" || lower_string == "sharbert" || lower_string == "fornax" {
			split_body[i] = "****"
		}
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Clean: strings.Join(split_body, " "),
	})
}
