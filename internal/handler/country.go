package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) AdminListCountry(c *gin.Context) {
	ctx := c.Request.Context()
	countries, err := h.datastoreRepo.ListCountry(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(countries))
}

func (h *Handler) AdminGetCountryByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	country, err := h.datastoreRepo.GetCountryByID(ctx, id)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminGetCountryByIDResponse{
		ID:   model.UUID(country.ID),
		Name: country.Name,
	}))
}

func (h *Handler) AdminCreateCountry(c *gin.Context) {
	var req AdminCreateCountryRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	country := model.Country{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	err := h.datastoreRepo.CreateCountry(ctx, &country)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateCountryResponse{ID: model.UUID(country.ID)}))
}

func (h *Handler) AdminUpdateCountry(c *gin.Context) {
	var req AdminUpdateCountryRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	country := model.Country{
		Name: req.Name,
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.UpdateCountry(ctx, id, &country); err != nil {
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

func (h *Handler) AdminDeleteCountryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteCountryByID(ctx, id); err != nil {
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
