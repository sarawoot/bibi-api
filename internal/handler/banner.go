package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"errors"
	"fmt"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	prefixBannerImageKey = "banners"
)

func (h *Handler) AdminCreateBanner(c *gin.Context) {
	var req AdminCreateBannerRequest
	if ok := bindFormMultipart(c, &req); !ok {
		return
	}

	banner := model.Banner{
		ID:       uuid.New(),
		Name:     req.Name,
		AreaCode: req.AreaCode,
		Images:   make([]model.BannerImage, 0, len(req.Images)),
	}

	ctx := c.Request.Context()
	for _, image := range req.Images {
		file, err := image.Open()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}
		defer file.Close()

		key := fmt.Sprintf("%s/%v/%s", prefixBannerImageKey, banner.ID, generateFileName(image.Filename))
		if err := h.s3Client.UploadFile(ctx, key, file); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		banner.Images = append(banner.Images, model.BannerImage{
			ID:       uuid.New(),
			BannerID: banner.ID,
			Path:     key,
		})
	}

	if err := h.datastoreRepo.CreateBanner(ctx, &banner); err != nil {
		log.Error(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == "banners_area_code_unique" {
				c.JSON(http.StatusConflict, RenderJSON(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(AdminCreateBannerResponse{ID: model.UUID(banner.ID)}))
}

func (h *Handler) AdminCreateBannerImage(c *gin.Context) {
	var req AdminCreateBannerImageRequest
	if ok := bindFormMultipart(c, &req); !ok {
		return
	}

	bannerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if _, err := h.datastoreRepo.GetBannerByID(ctx, bannerID); err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	response := make([]AdminBannerImageResponse, 0, len(req.Images))
	for _, image := range req.Images {
		file, err := image.Open()
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}
		defer file.Close()

		key := fmt.Sprintf("%s/%v/%s", prefixBannerImageKey, bannerID, generateFileName(image.Filename))
		if err := h.s3Client.UploadFile(ctx, key, file); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		imageID, err := h.datastoreRepo.CreateBannerImage(ctx, bannerID, key)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, RenderJSON(err))
			return
		}

		response = append(response, h.modelBannerImageToResponse(model.BannerImage{
			ID:   imageID,
			Path: key,
		}))
	}

	c.JSON(http.StatusOK, RenderJSON(response))
}

func (h *Handler) AdminGetBannerByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	banner, err := h.datastoreRepo.GetBannerByID(ctx, id)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(h.modelBannerToResponse(*banner)))
}

func (h *Handler) GetBannerByAreaCode(c *gin.Context) {
	ctx := c.Request.Context()

	banner, err := h.datastoreRepo.GetBannerByAreaCode(ctx, c.Param("area_code"))
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusNotFound, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(h.modelBannerToResponse(*banner)))
}

func (h *Handler) AdminListBanner(c *gin.Context) {
	ctx := c.Request.Context()

	banners, err := h.datastoreRepo.ListBanner(ctx)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	res := make([]AdminBannerResponse, 0, len(banners))
	for _, banner := range banners {
		res = append(res, h.modelBannerToResponse(banner))
	}

	c.JSON(http.StatusOK, RenderJSON(res))
}

func (h *Handler) AdminUpdateBanner(c *gin.Context) {
	var req AdminUpdateBannerRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	banner := model.Banner{
		Name:     req.Name,
		AreaCode: req.AreaCode,
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.UpdateBanner(ctx, id, &banner); err != nil {
		log.Error(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == "banners_area_code_unique" {
				c.JSON(http.StatusConflict, RenderJSON(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON("Successful"))
}

func (h *Handler) AdminDeleteBannerByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteBannerByID(ctx, id); err != nil {
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

func (h *Handler) AdminDeleteBannerImageByID(c *gin.Context) {
	bannerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	bannerImageID, err := uuid.Parse(c.Param("banner_image_id"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, RenderJSON(err))
		return
	}

	ctx := c.Request.Context()
	if err := h.datastoreRepo.DeleteBannerImageByID(ctx, bannerID, bannerImageID); err != nil {
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
