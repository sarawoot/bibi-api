package model

import "github.com/google/uuid"

type ProductRecommend struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Product   Product   `json:"product"`
}
