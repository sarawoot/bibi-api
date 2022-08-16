package handler

import (
	"api/internal/model"
	"mime/multipart"

	"github.com/google/uuid"
)

type AdminBannerImageResponse struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

type AdminBannerResponse struct {
	ID       model.UUID                 `json:"id"`
	Name     string                     `json:"name"`
	AreaCode string                     `json:"area_code"`
	Images   []AdminBannerImageResponse `json:"images"`
}

type AdminCreateBannerRequest struct {
	Name     string                  `form:"name" binding:"required"`
	AreaCode string                  `form:"area_code" binding:"required"`
	Images   []*multipart.FileHeader `form:"images[]"`
}

type AdminCreateBannerResponse struct {
	ID model.UUID `json:"id"`
}

type AdminCreateBannerImageRequest struct {
	Images []*multipart.FileHeader `form:"images[]" binding:"required"`
}

type AdminUpdateBannerRequest struct {
	Name     string `json:"name" binding:"required"`
	AreaCode string `json:"area_code" binding:"required"`
}
