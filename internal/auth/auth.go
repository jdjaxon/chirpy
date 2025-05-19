package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeader = errors.New("no authorization header")
var ErrMalformedAuthHeader = errors.New("malformed authorization header")

// HashPassword -
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPasswordHash -
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) < 2 || splitAuthHeader[0] != "Bearer" {
		return "", ErrMalformedAuthHeader
	}

	return splitAuthHeader[1], nil
}

// GetApiKey -
func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}
	splitAuthHeader := strings.Split(authHeader, " ")
	if len(splitAuthHeader) < 2 || splitAuthHeader[0] != "ApiKey" {
		return "", ErrMalformedAuthHeader
	}

	return splitAuthHeader[1], nil
}

// MakeRefreshToken -
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	encodedToken := hex.EncodeToString(token)

	return encodedToken, nil
}
