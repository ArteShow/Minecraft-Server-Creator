package hashing

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
}

func ComparePasswords(password string, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(
		hashedPassword,
		[]byte(password),
	)
	return err == nil
}

