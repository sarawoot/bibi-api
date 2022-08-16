package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AdminListProductType(c *gin.Context) {
	ctx := c.Request.Context()
	productTypes, err := h.datastoreRepo.ListProductType(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(productTypes))
}

func (h *Handler) AdminGetProductTypeByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	productType, err := h.datastoreRepo.GetProductTypeByID(ctx, id)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminGetProductTypeByIDResponse{
		ID:   model.UUID(productType.ID),
		Name: productType.Name,
	}))
}

func (h *Handler) AdminCreateProductType(c *gin.Context) {
	var req AdminCreateProductTypeRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	productType := model.ProductType{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	err := h.datastoreRepo.CreateProductType(ctx, &productType)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateProductTypeResponse{ID: model.UUID(productType.ID)}))
}

func (h *Handler) AdminUpdateProductType(c *gin.Context) {
	var req AdminUpdateProductTypeRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	productType := model.ProductType{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.UpdateProductType(ctx, id, &productType); err != nil {
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

func (h *Handler) AdminDeleteProductTypeByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteProductTypeByID(ctx, id); err != nil {
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
