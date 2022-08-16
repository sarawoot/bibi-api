package handler

import (
	"api/pkg/log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func (h *Handler) BearerAuth(c *gin.Context) {
	token, err := getBearerToken(c)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, RenderJSON(err))
		return
	}

	claims, err := jwt.ParseString(
		string(token),
		jwt.WithKeySet(h.jwksPublic),
	)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, RenderJSON(err))
		return
	}

	setContextWithUserID(c, claims.Subject())

	c.Next()
}

func (h *Handler) AdminBearerAuth(c *gin.Context) {
	token, err := getBearerToken(c)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, RenderJSON(err))
		return
	}

	claims, err := jwt.ParseString(
		string(token),
		jwt.WithKeySet(h.adminJWKSPublic),
	)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, RenderJSON(err))
		return
	}

	setContextWithUserID(c, claims.Subject())

	c.Next()
}

func (h *Handler) LoggerMiddleware(notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		logger := log.WithFields(map[string]interface{}{
			"hostname":    hostname,
			"referer":     referer,
			"data-lenght": dataLength,
			"status":      statusCode,
			"method":      c.Request.Method,
			"path":        path,
			"query":       c.Request.URL.RawQuery,
			"ip":          clientIP,
			"user-agent":  clientUserAgent,
			"elapsed":     stop.String(),
		})
		if len(c.Errors) > 0 {
			logger.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode >= http.StatusInternalServerError {
				logger.Error(path)
			} else if statusCode >= http.StatusBadRequest {
				logger.Warn(path)
			} else {
				logger.Info(path)
			}
		}
	}
}

func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://art-sarawoot.stoplight.io/")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,Accept,Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
