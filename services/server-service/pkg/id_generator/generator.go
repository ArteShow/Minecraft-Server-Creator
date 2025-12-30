package idgenerator

import (
	"github.com/google/uuid"
)

func GenerateServerID() string {
	return uuid.NewString()
}
