package model

import "github.com/google/uuid"

type Banner struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name"`
	AreaCode string        `json:"area_code"`
	Images   []BannerImage `json:"banner_images"`
}

type BannerImage struct {
	ID       uuid.UUID `json:"id"`
	BannerID uuid.UUID `json:"banner_id"`
	Path     string    `json:"path"`
}
