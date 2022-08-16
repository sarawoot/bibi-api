package handler

import "api/internal/model"

type AdminCreateCountryRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminCreateCountryResponse struct {
	ID model.UUID `json:"id"`
}

type AdminUpdateCountryRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminGetCountryByIDResponse struct {
	ID   model.UUID `json:"id"`
	Name string     `json:"name"`
}
