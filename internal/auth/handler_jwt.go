package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := []byte(tokenSecret)
	newJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
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
