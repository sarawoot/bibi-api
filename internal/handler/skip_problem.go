package handler

import (
	"api/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListSkinProblem(c *gin.Context) {
	ctx := c.Request.Context()
	skinProblems, err := h.datastoreRepo.ListSkinProblem(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(skinProblems))
}
