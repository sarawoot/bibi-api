package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AdminListProductRecommend(c *gin.Context) {
	var req AdminListProductRecommendRequest
	if ok := bindQuery(c, &req); !ok {
		return
	}

	offset := (req.Page - 1) * req.Limit
	ctx := c.Request.Context()
	productRecommends, err := h.datastoreRepo.AdminListProductRecommend(ctx, req.Limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	totalCount, err := h.datastoreRepo.AdminCountProductRecommend(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	items := make([]AdminProductRecommendResponse, 0, len(productRecommends))
	for _, product := range productRecommends {
		items = append(items, h.modelProductRecommendToResponse(product))
	}

	c.JSON(http.StatusOK, RenderJSON(AdminListProductRecommendResponse{
		TotalCount: totalCount,
		Items:      items,
	}))
}

func (h *Handler) AdminCreateProductRecommend(c *gin.Context) {
	var req AdminCreateProductRecommendRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	if _, err := h.datastoreRepo.GetProductByID(ctx, uuid.UUID(req.ProductID)); err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	id, err := h.datastoreRepo.CreateProductRecommend(ctx, uuid.UUID(req.ProductID))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateProductRecommendResponse{ID: model.UUID(id)}))
}

func (h *Handler) AdminDeleteProductRecommendByID(c *gin.Context) {
	productRecommendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteProductRecommendByID(ctx, productRecommendID); err != nil {
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

func (h *Handler) MobileListProductRecommend(c *gin.Context) {
	var req MobileProductRecommendRequest
	if ok := bindQuery(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	products, err := h.datastoreRepo.MobileListProductRecommend(ctx, req.Limit)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	res := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		res = append(res, h.modelProductToResponse(product.Product))
	}

	c.JSON(http.StatusOK, RenderJSON(res))
}
