package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

var (
	UserPasswordAuth AuthProvider = "password"
	FacebookAuth     AuthProvider = "facebook"
	GmailAuth        AuthProvider = "gmail"
)

type User struct {
	ID            uuid.UUID
	Email         string    `json:"email"`
	PasswordHash  string    `json:"password_hash"`
	Gender        string    `json:"gender"`
	Birthdate     time.Time `json:"birthdate"`
	SkinTypeID    uuid.UUID `json:"skin_type_id"`
	SkinProblemID uuid.UUID `json:"skin_problem_id"`
}
