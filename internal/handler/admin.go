package handler

import (
	"api/pkg/log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	admin, err := h.datastoreRepo.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusBadRequest, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	ok, err := comparePasswords(admin.PasswordHash, req.Password)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	if !ok {
		msg := "email or password is incorrect"
		log.Error(msg)
		c.JSON(http.StatusBadRequest, RenderJSON(msg))
		return
	}

	token, err := h.signAdminJWT(admin.ID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(UserLoginResponse{AccessToken: token}))

}
