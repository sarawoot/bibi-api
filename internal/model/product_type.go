package model

import "github.com/google/uuid"

type ProductType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
