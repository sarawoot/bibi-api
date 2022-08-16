package handler

import (
	"api/pkg/log"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

func bindJSON(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		log.Error(msg)
		c.JSON(http.StatusUnsupportedMediaType, RenderJSON(msg))
		return false
	}

	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return false
	}

	return true
}

func bindFormMultipart(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "multipart/form-data" {
		msg := fmt.Sprintf("%s only accepts Content-Type multipart/form-data", c.FullPath())
		log.Error(msg)
		c.JSON(http.StatusUnsupportedMediaType, RenderJSON(msg))
		return false
	}

	if err := c.ShouldBindWith(req, binding.FormMultipart); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return false
	}

	return true
}

func bindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindWith(req, binding.Query); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return false
	}

	return true
}

func generateFileName(path string) string {
	ext := filepath.Ext(path)
	name := shortuuid.New()

	return name + ext
}

func parseUUIDPtr(uid *string) *uuid.UUID {
	if uid != nil {
		rs, err := uuid.Parse(*uid)
		if err != nil {
			rs = uuid.Nil
		}
		return &rs
	}

	return nil
}

// func parseUUID(uid string) uuid.UUID {
// 	rs, err := uuid.Parse(uid)
// 	if err != nil {
// 		return uuid.Nil
// 	}

// 	return rs
// }
