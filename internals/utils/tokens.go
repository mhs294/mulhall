package utils

import (
	"math/rand"
)

// CreateAlphaNumToken creates a randomized token consisting only of alpha-numeric characters.
//
// length is the the desired amount of characters to use when creating the token.
func CreateAlphaNumToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return createToken(charset, length)
}

// CreateAlphaNumToken creates a randomized token consisting of standard ASCII characters.
//
// length is the the desired amount of characters to use when creating the token.
func CreateToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_+"
	return createToken(charset, length)
}

// createToken creates a randomized token consisting only of characters from the specified character set.
//
// charset is the set of characters that are allowed to be used when creating the token.
//
// length is the the desired amount of characters to use when creating the token.
func createToken(charset string, length int) string {
	token := make([]byte, length)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}

	return string(token)
}
