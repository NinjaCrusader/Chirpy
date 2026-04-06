package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("there was something wrong with creating the password hash: %v\n", err)
		return hash, err
	}

	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {

	valid, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("there was an error with comparing the pass and hash: %v\n", err)
		return valid, err
	}

	return valid, err
}
