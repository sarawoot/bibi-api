package model

import "github.com/google/uuid"

type ProductCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
