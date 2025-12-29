package id

import "github.com/google/uuid"

func GenerateID() string {
	return uuid.NewString()
}
