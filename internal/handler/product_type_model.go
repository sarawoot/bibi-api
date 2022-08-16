package handler

import "api/internal/model"

type AdminCreateProductTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminCreateProductTypeResponse struct {
	ID model.UUID `json:"id"`
}

type AdminUpdateProductTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminGetProductTypeByIDResponse struct {
	ID   model.UUID `json:"id"`
	Name string     `json:"name"`
}
