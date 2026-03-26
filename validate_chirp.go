package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type cleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		log.Printf("There was an error decoding the request: %v\n", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedParams := removeBadWords(params.Body)

	res := cleaned{
		CleanedBody: cleanedParams,
	}

	respondWithJSON(w, 200, res)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnVals struct {
		Error string `json:"error"`
	}

	respBody := returnVals{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling response JSON: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func removeBadWords(s string) string {

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	split := strings.Split(s, " ")

	for i := 0; i < len(split); i++ {
		copy := split[i]
		if slices.Contains(bannedWords, strings.ToLower(copy)) {
			split[i] = "****"
			continue
		}
	}

	final := strings.Join(split, " ")

	return final
}
