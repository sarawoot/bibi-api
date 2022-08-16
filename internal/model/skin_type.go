package model

import "github.com/google/uuid"

type SkinType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
