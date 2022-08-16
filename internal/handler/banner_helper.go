package handler

import (
	"api/internal/model"
	"fmt"
)

func (h *Handler) modelBannerImageToResponse(image model.BannerImage) AdminBannerImageResponse {
	return AdminBannerImageResponse{
		ID:  image.ID,
		URL: fmt.Sprintf("%s/%s", h.appConfig.AssetsBaseURL, image.Path),
	}
}

func (h *Handler) modelBannerToResponse(banner model.Banner) AdminBannerResponse {
	res := AdminBannerResponse{
		ID:       model.UUID(banner.ID),
		Name:     banner.Name,
		AreaCode: banner.AreaCode,
	}

	images := make([]AdminBannerImageResponse, 0, len(banner.Images))
	for _, image := range banner.Images {
		images = append(images, h.modelBannerImageToResponse(image))
	}
	res.Images = images

	return res
}
