package handler

import (
	"api/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListSkinType(c *gin.Context) {
	ctx := c.Request.Context()
	skinTypes, err := h.datastoreRepo.ListSkinType(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(skinTypes))
}
