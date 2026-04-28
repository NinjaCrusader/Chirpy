package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("there was something wrong with creating the password hash: %v\n", err)
		return "", err
	}

	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {

	valid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("there was an error with comparing the pass and hash: %v\n", err)
		return false, err
	}

	return valid, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {

	token := []byte(tokenSecret)
	newJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   userID.String(),
	})

	return newJWT.SignedString(token)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}

	userIDstring, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != "chirpy-access" {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDstring)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w\n", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	headerData := headers.Get("Authorization")
	headerSplice := strings.Fields(headerData)
	if (len(headerSplice) > 2) || (len(headerSplice) <= 1) {
		return "", errors.New("invalid token")
	}

	if headerSplice[0] != "Bearer" {
		return "", errors.New("invalid token scheme")
	}

	token := headerSplice[1]

	return token, nil
}

func MakeRefreshToken() string {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("there was a problem creating the token: %v\n", err)
		return ""
	}

	convert := hex.EncodeToString(b)

	return convert
}

func GetAPIKey(headers http.Header) (string, error) {

	headerData := headers.Get("Authorization")
	headerSplice := strings.Fields(headerData)
	if (len(headerSplice) > 2) || (len(headerData) <= 1) {
		return "", errors.New("invalid polka API key")
	}

	key := headerSplice[1]

	return key, nil
}
