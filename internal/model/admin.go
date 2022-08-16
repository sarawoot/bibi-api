package model

import (
	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
