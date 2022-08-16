package handler

import "api/internal/model"

type AdminCreateProductCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminCreateProductCategoryResponse struct {
	ID model.UUID `json:"id"`
}

type AdminUpdateProductCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type AdminGetProductCategoryByIDResponse struct {
	ID   model.UUID `json:"id"`
	Name string     `json:"name"`
}
