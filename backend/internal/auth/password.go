package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword applies Bcrypt with an optimal work factor calculation
func HashPassword(password string) (string, error) {
	// Cost 12 balances modern CPU execution times with strong brute force deterrence
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// VerifyPassword performs a constant-time comparison against timing-attack vectors
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}