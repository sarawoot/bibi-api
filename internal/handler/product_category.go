package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AdminListProductCategory(c *gin.Context) {
	ctx := c.Request.Context()
	productCategorys, err := h.datastoreRepo.ListProductCategory(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(productCategorys))
}

func (h *Handler) AdminGetProductCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	productCategory, err := h.datastoreRepo.GetProductCategoryByID(ctx, id)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminGetProductCategoryByIDResponse{
		ID:   model.UUID(productCategory.ID),
		Name: productCategory.Name,
	}))
}

func (h *Handler) AdminCreateProductCategory(c *gin.Context) {
	var req AdminCreateProductCategoryRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	productCategory := model.ProductCategory{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	err := h.datastoreRepo.CreateProductCategory(ctx, &productCategory)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateProductCategoryResponse{ID: model.UUID(productCategory.ID)}))
}

func (h *Handler) AdminUpdateProductCategory(c *gin.Context) {
	var req AdminUpdateProductCategoryRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	productCategory := model.ProductCategory{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.UpdateProductCategory(ctx, id, &productCategory); err != nil {
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

func (h *Handler) AdminDeleteProductCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteProductCategoryByID(ctx, id); err != nil {
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
