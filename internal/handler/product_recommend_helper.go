package handler

import "api/internal/model"

func (h *Handler) modelProductRecommendToResponse(p model.ProductRecommend) AdminProductRecommendResponse {
	return AdminProductRecommendResponse{
		ID:      model.UUID(p.ID),
		Product: h.modelProductToResponse(p.Product),
	}
}
