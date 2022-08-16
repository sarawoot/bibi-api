package model

import "github.com/google/uuid"

type SkinProblem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
