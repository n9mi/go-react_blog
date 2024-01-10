package utils

import "golang.org/x/crypto/bcrypt"

// HashUserPassword returns hashed password and error while hashing the password
// by taking plain text password
func HashUserPassword(plainPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	return string(hashed), err
}

// IsUserPasswordValid check if user password from request (plainPassword) is matched with
// password that stored in database (hashedPassword)
// Returns true if valid
func IsUserPasswordValid(plainPassword string, hashedPassword string) bool {
	isValid := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil

	return isValid
}
