package handler

import (
	"api/internal/model"
	"mime/multipart"

	"github.com/google/uuid"
)

type ProductImageResponse struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type ProductResponse struct {
	ID                model.UUID             `json:"id"`
	Brand             string                 `json:"brand" binding:"required"`
	Name              string                 `json:"name" binding:"required"`
	ShortDescription  string                 `json:"short_description"`
	Description       string                 `json:"description"`
	Size              string                 `json:"size"`
	Price             model.Decimal          `json:"price"`
	ProductTypeID     model.UUID             `json:"product_type_id,omitempty"`
	ProductCategoryID model.UUID             `json:"product_category_id,omitempty"`
	SkinTypeID        model.UUID             `json:"skin_type_id,omitempty"`
	CountryID         model.UUID             `json:"country_id,omitempty"`
	Tags              []string               `json:"tags"`
	Images            []ProductImageResponse `json:"images"`
}
type AdminListProductResponse struct {
	TotalCount int               `json:"total_count"`
	Items      []ProductResponse `json:"items"`
}

type AdminCreateProductRequest struct {
	Brand             *string                 `form:"brand" binding:"required"`
	Name              *string                 `form:"name" binding:"required"`
	ShortDescription  *string                 `form:"short_description"`
	Description       *string                 `form:"description"`
	Size              *string                 `form:"size"`
	Price             *model.Decimal          `form:"price"`
	ProductTypeID     *string                 `form:"product_type_id"  binding:"omitempty,uuid"`
	ProductCategoryID *string                 `form:"product_category_id" binding:"omitempty,uuid"`
	SkinTypeID        *string                 `form:"skin_type_id" binding:"omitempty,uuid"`
	CountryID         *string                 `form:"country_id" binding:"omitempty,uuid"`
	Tags              []string                `form:"tags[]"`
	Images            []*multipart.FileHeader `form:"images[]"`
}

type AdminCreateProductResponse struct {
	ID model.UUID `json:"id"`
}

type AdminCreateProductImageRequest struct {
	Images []*multipart.FileHeader `form:"images[]" binding:"required"`
}

type AdminListProductRequest struct {
	Limit int `form:"limit,default=10" binding:"gt=0,lte=100"`
	Page  int `form:"page,default=1" binding:"gte=1"`
}

type AdminUpdateProductRequest struct {
	Brand             *string        `json:"brand"`
	Name              *string        `json:"name"`
	ShortDescription  *string        `json:"short_description"`
	Description       *string        `json:"description"`
	Size              *string        `json:"size"`
	Price             *model.Decimal `json:"price"`
	ProductTypeID     *string        `json:"product_type_id"  binding:"omitempty,uuid"`
	ProductCategoryID *string        `json:"product_category_id" binding:"omitempty,uuid"`
	SkinTypeID        *string        `json:"skin_type_id" binding:"omitempty,uuid"`
	CountryID         *string        `json:"country_id" binding:"omitempty,uuid"`
	Tags              []string       `json:"tags"`
}

type MobileProductNewArrivalRequest struct {
	Limit int `form:"limit,default=10" binding:"gt=0,lte=20"`
}
