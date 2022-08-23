package handler

import (
	"api/internal/model"
	"fmt"

	"github.com/google/uuid"
)

func (h *Handler) modelProductImageToResponse(image model.ProductImage) ProductImageResponse {
	return ProductImageResponse{
		ID:  image.ID,
		URL: fmt.Sprintf("%s/%s", h.appConfig.AssetsBaseURL, image.Path),
	}
}

func (h *Handler) modelProductToResponse(p model.Product) ProductResponse {
	resp := ProductResponse{
		ID:                model.UUID(p.ID),
		Brand:             *p.Brand,
		Name:              *p.Name,
		ShortDescription:  *p.ShortDescription,
		Description:       *p.Description,
		Size:              *p.Size,
		Price:             p.Price,
		Tags:              p.Tags,
		ProductTypeID:     model.UUID(uuid.Nil),
		ProductCategoryID: model.UUID(uuid.Nil),
		SkinTypeID:        model.UUID(uuid.Nil),
		CountryID:         model.UUID(uuid.Nil),
	}

	if p.ProductTypeID != nil {
		resp.ProductTypeID = model.UUID(*p.ProductTypeID)
	}

	if p.ProductCategoryID != nil {
		resp.ProductCategoryID = model.UUID(*p.ProductCategoryID)
	}

	if p.SkinTypeID != nil {
		resp.SkinTypeID = model.UUID(*p.SkinTypeID)
	}

	if p.CountryID != nil {
		resp.CountryID = model.UUID(*p.CountryID)
	}

	images := make([]ProductImageResponse, 0, len(p.Images))
	for _, image := range p.Images {
		images = append(images, h.modelProductImageToResponse(image))
	}
	resp.Images = images

	return resp
}

func (p AdminCreateProductRequest) toModel() model.Product {
	return model.Product{
		ID:                uuid.New(),
		Brand:             p.Brand,
		Name:              p.Name,
		ShortDescription:  p.ShortDescription,
		Description:       p.Description,
		Size:              p.Size,
		Price:             p.Price,
		ProductTypeID:     parseUUIDPtr(p.ProductTypeID),
		ProductCategoryID: parseUUIDPtr(p.ProductCategoryID),
		SkinTypeID:        parseUUIDPtr(p.SkinTypeID),
		CountryID:         parseUUIDPtr(p.CountryID),
		Tags:              p.Tags,
	}
}

func (p AdminUpdateProductRequest) toModel() model.Product {
	return model.Product{
		Brand:             p.Brand,
		Name:              p.Name,
		ShortDescription:  p.ShortDescription,
		Description:       p.Description,
		Size:              p.Size,
		Price:             p.Price,
		ProductTypeID:     parseUUIDPtr(p.ProductTypeID),
		ProductCategoryID: parseUUIDPtr(p.ProductCategoryID),
		SkinTypeID:        parseUUIDPtr(p.SkinTypeID),
		CountryID:         parseUUIDPtr(p.CountryID),
		Tags:              p.Tags,
	}
}
