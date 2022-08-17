package handler

import (
	"api/internal/model"
	"api/pkg/log"
	"errors"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func (h *Handler) UserSignup(c *gin.Context) {
	var req UserSignupRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	password, err := hashPassword(req.Password)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	user := model.User{
		Email:         req.Email,
		PasswordHash:  password,
		Gender:        string(req.Gender),
		Birthdate:     time.Time(req.Birthdate),
		SkinTypeID:    uuid.UUID(req.SkinTypeID),
		SkinProblemID: uuid.UUID(req.SkinProblemID),
	}

	if err := h.datastoreRepo.CreateUser(ctx, &user); err != nil {
		log.Error(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == "users_email_unique" {
				c.JSON(http.StatusConflict, RenderJSON(err))
				return
			}

			if pgErr.ConstraintName == "users_skin_type_id_fkey" || pgErr.ConstraintName == "users_skin_problem_id_fkey" {
				c.JSON(http.StatusBadRequest, RenderJSON(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	token, err := h.signJWT(user.ID, model.UserPasswordAuth)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(UserSignupResponse{AccessToken: token}))
}

func (h *Handler) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if ok := bindJSON(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	user, err := h.datastoreRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error(err)
		if pgxscan.NotFound(err) {
			c.JSON(http.StatusBadRequest, RenderJSON(err))
			return
		}

		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	ok, err := comparePasswords(user.PasswordHash, req.Password)
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

	token, err := h.signJWT(user.ID, model.UserPasswordAuth)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, RenderJSON(err))
		return
	}

	c.JSON(http.StatusOK, RenderJSON(UserLoginResponse{AccessToken: token}))
}
