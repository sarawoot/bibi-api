package handler

import "api/internal/model"

type AdminCreateProductRecommendResponse struct {
	ID model.UUID `json:"id"`
}

type AdminCreateProductRecommendRequest struct {
	ProductID model.UUID `json:"product_id" binding:"required"`
}

type AdminListProductRecommendRequest struct {
	Limit int `form:"limit,default=10" binding:"gt=0,lte=100"`
	Page  int `form:"page,default=1" binding:"gte=1"`
}

type AdminProductRecommendResponse struct {
	ID      model.UUID      `json:"id"`
	Product ProductResponse `json:"product"`
}

type AdminListProductRecommendResponse struct {
	TotalCount int                             `json:"total_count"`
	Items      []AdminProductRecommendResponse `json:"items"`
}

type MobileProductRecommendRequest struct {
	Limit int `form:"limit,default=10" binding:"gt=0,lte=20"`
}
