package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"fmt"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	prefixProductImageKey = "products"
)

func (h *Handler) AdminCreateProduct(c *gin.Context) {
	var req AdminCreateProductRequest
	if ok := bindFormMultipart(c, &req); !ok {
		return
	}

	product := req.toModel()
	ctx := c.Request.Context()
	for _, image := range req.Images {
		file, err := image.Open()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}
		defer file.Close()

		key := fmt.Sprintf("%s/%v/%s", prefixProductImageKey, product.ID, generateFileName(image.Filename))
		if err := h.s3Client.UploadFile(ctx, key, file); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		product.Images = append(product.Images, model.ProductImage{
			ID:        uuid.New(),
			ProductID: product.ID,
			Path:      key,
		})
	}

	if err := h.datastoreRepo.CreateProduct(ctx, &product); err != nil {
		log.Error(err)

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateProductResponse{ID: model.UUID(product.ID)}))
}

func (h *Handler) AdminCreateProductImage(c *gin.Context) {
	var req AdminCreateProductImageRequest
	if ok := bindFormMultipart(c, &req); !ok {
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if _, err := h.datastoreRepo.GetProductByID(ctx, productID); err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	response := make([]ProductImageResponse, 0, len(req.Images))
	for _, image := range req.Images {
		file, err := image.Open()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}
		defer file.Close()

		key := fmt.Sprintf("%s/%v/%s", prefixProductImageKey, productID, generateFileName(image.Filename))
		if err := h.s3Client.UploadFile(ctx, key, file); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		imageID, err := h.datastoreRepo.CreateProductImage(ctx, productID, key)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		response = append(response, h.modelProductImageToResponse(model.ProductImage{
			ID:   imageID,
			Path: key,
		}))
	}

	c.JSON(http.StatusOK, RenderJSON(response))
}

func (h *Handler) AdminGetProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	product, err := h.datastoreRepo.GetProductByID(ctx, id)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(h.modelProductToResponse(*product)))
}

func (h *Handler) AdminListProduct(c *gin.Context) {
	var req AdminListProductRequest
	if ok := bindQuery(c, &req); !ok {
		return
	}

	offset := (req.Page - 1) * req.Limit
	ctx := c.Request.Context()
	products, err := h.datastoreRepo.AdminListProduct(ctx, req.Limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	totalCount, err := h.datastoreRepo.AdminCountProduct(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	items := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		items = append(items, h.modelProductToResponse(product))
	}

	c.JSON(http.StatusOK, RenderJSON(AdminListProductResponse{
		TotalCount: totalCount,
		Items:      items,
	}))
}

func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	var req AdminUpdateProductRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	product := req.toModel()
	if err := h.datastoreRepo.UpdateProduct(ctx, id, &product); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON("Successful"))
}

func (h *Handler) AdminDeleteProductImageByID(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	productImageID, err := uuid.Parse(c.Param("product_image_id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteProductImageByID(ctx, productID, productImageID); err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON("Successful"))
}

func (h *Handler) AdminDeleteProductByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteProductByID(ctx, id); err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON("Successful"))
}

func (h *Handler) MobileListProductNewArrival(c *gin.Context) {
	var req MobileProductNewArrivalRequest
	if ok := bindQuery(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	products, err := h.datastoreRepo.ListProductNewArrival(ctx, req.Limit)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	res := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		res = append(res, h.modelProductToResponse(product))
	}

	c.JSON(http.StatusOK, RenderJSON(res))
}
