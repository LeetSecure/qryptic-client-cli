package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"regexp"
	"strings"
)

func IsValidEmailId(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,10}$`)
	return emailRegex.MatchString(email)
}

// base64URLEncode performs the base64 URL-safe encoding
// without padding (i.e., no '=' characters).
func Base64URLEncode(input []byte) string {
	encoded := base64.StdEncoding.EncodeToString(input)
	encoded = strings.TrimRight(encoded, "=")       // Remove padding
	encoded = strings.ReplaceAll(encoded, "+", "-") // Replace '+' with '-'
	encoded = strings.ReplaceAll(encoded, "/", "_") // Replace '/' with '_'
	return encoded
}

func RandomStringGenerator(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Verify that the code verifier matches the code challenge
func GetCodeChallenge(codeVerifier string) string {
	hash := sha256.New()
	hash.Write([]byte(codeVerifier))
	hashedVerifier := hash.Sum(nil)
	codeChallenge := Base64URLEncode(hashedVerifier)
	return codeChallenge
}
