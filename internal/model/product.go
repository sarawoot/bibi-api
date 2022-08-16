package model

import "github.com/google/uuid"

type Product struct {
	ID                uuid.UUID      `json:"id"`
	Brand             *string        `json:"brand"`
	Name              *string        `json:"name"`
	ShortDescription  *string        `json:"short_description"`
	Description       *string        `json:"description"`
	Size              *string        `json:"size"`
	Price             *float64       `json:"price"`
	ProductTypeID     *uuid.UUID     `json:"product_type_id"`
	ProductCategoryID *uuid.UUID     `json:"product_category_id"`
	SkinTypeID        *uuid.UUID     `json:"skin_type_id"`
	CountryID         *uuid.UUID     `json:"country_id"`
	Tags              []string       `json:"tags"`
	Images            []ProductImage `json:"product_images"`
}

type ProductImage struct {
	ID        uuid.UUID      `json:"id"`
	ProductID uuid.UUID      `json:"product_id"`
	Path      string         `json:"path"`
	Product   *Product       `json:"product"`
	Images    []ProductImage `json:"product_images"`
}
