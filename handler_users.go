package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NinjaCrusader/Chirpy/internal/auth"
	"github.com/NinjaCrusader/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userUpdate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updatedUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("There was an error decoding the request: %v\n", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("something went wrong while hashing the password: %v\n", err)
		return
	}

	createUserParam := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	}

	user, err := cfg.db.CreateUser(r.Context(), createUserParam)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			log.Printf("There was an error with postgres creating the user: %v\n", dbErr.Code)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			log.Printf("There was an error creating the user: %v\n", err)
			return
		}
	}

	userParams := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userParams)

}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error getting the request information: %v\n", err)
		return
	}

	user, err := cfg.db.FindUserToLogin(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("there was an error finding the user to login: %v\n", err)
		return
	}

	validCheck, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || validCheck == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("there was an error finding the user to login: %v\n", err)
		return
	}

	createToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error creating the JWT: %v\n", err)
		return
	}

	refresh := auth.MakeRefreshToken()
	refreshParams := database.InsertRefreshTokenParams{
		Token:     refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	refreshToken, err := cfg.db.InsertRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("something went wrong inserting the refresh token: %v\n", err)
		return
	}

	res := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        createToken,
		RefreshToken: refreshToken.Token,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		log.Printf("there was an authorization error: %v\n", err)
		return
	}

	userID, err := auth.ValidateJWT(headerToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		log.Printf("there was an error with validating the jwt: %v\n", err)
		return
	}

	requestBody := userUpdate{}
	decoder := json.NewDecoder(r.Body)
	decodeError := decoder.Decode(&requestBody)
	if decodeError != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		log.Printf("something went wrong while decoding the user update body: %v\n", decodeError)
		return
	}

	hashedPass, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error hashing the password: %v\n", err)
		return
	}

	updateParams := database.UpdateUserPassParams{
		Email:          requestBody.Email,
		HashedPassword: hashedPass,
		ID:             userID,
	}

	updated, err := cfg.db.UpdateUserPass(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error trying to do the user update: %v\n", err)
		return
	}

	res := updatedUser{
		ID:          updated.ID,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
		Email:       updated.Email,
		IsChirpyRed: updated.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, res)
}
