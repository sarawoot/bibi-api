package server

import (
	"context"
	"fmt"
	"os"

	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/repo/datastorepgx"
	"api/pkg/aws/s3"
	"api/pkg/log"
)

func RunServer() error {
	conf := &config.App
	srv, err := initServer(conf)
	if err != nil {
		return err
	}

	//  Ctrl+C    = syscall.SIGINT
	//  kubectl kill = syscall.SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	log.Info(fmt.Sprintf("Server is running on %s", conf.ServerAddr))

	<-ctx.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutting down gracefully, press Ctrl+C again to force")
	}

	log.Info("Server is stop")

	return nil
}

func initServer(conf *config.AppConfig) (*http.Server, error) {
	pgxConn, err := initPgx(config.App.PostgresConnection)
	if err != nil {
		return nil, err
	}

	dataStoreRepo := datastorepgx.New(pgxConn)
	s3Client, err := s3.New(conf.AWSBucket)
	if err != nil {
		return nil, err
	}

	h, err := handler.NewHandler(conf, dataStoreRepo, s3Client)
	if err != nil {
		return nil, err
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Location"}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(
		gin.Recovery(),
		h.LoggerMiddleware(),
		cors.New(corsConfig),
	)

	// Router
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	engine.POST("/signup", h.UserSignup)
	engine.POST("/login", h.UserLogin)

	engine.GET("/skin_types", h.ListSkinType)
	engine.GET("/skin_problems", h.ListSkinProblem)
	engine.GET("/banners/:area_code", h.GetBannerByAreaCode)

	mobile := engine.Group("/mobile")
	{
		home := mobile.Group("/home")
		{
			home.GET("/products/new_arrival", h.MobileListProductNewArrival)
			home.GET("/products/recommend", h.MobileListProductRecommend)
		}
	}

	engine.POST("/admin/login", h.AdminLogin)
	admin := engine.Group("/admin")
	{
		admin.Use(h.AdminBearerAuth)

		productType := admin.Group("/product_types")
		{
			productType.GET("", h.AdminListProductType)
			productType.POST("", h.AdminCreateProductType)
			productType.GET("/:id", h.AdminGetProductTypeByID)
			productType.PATCH("/:id", h.AdminUpdateProductType)
			productType.DELETE("/:id", h.AdminDeleteProductTypeByID)
		}

		productCategory := admin.Group("/product_categories")
		{
			productCategory.GET("", h.AdminListProductCategory)
			productCategory.POST("", h.AdminCreateProductCategory)
			productCategory.GET("/:id", h.AdminGetProductCategoryByID)
			productCategory.PATCH("/:id", h.AdminUpdateProductCategory)
			productCategory.DELETE("/:id", h.AdminDeleteProductCategoryByID)
		}

		country := admin.Group("/countries")
		{
			country.GET("", h.AdminListCountry)
			country.POST("", h.AdminCreateCountry)
			country.GET("/:id", h.AdminGetCountryByID)
			country.PATCH("/:id", h.AdminUpdateCountry)
			country.DELETE("/:id", h.AdminDeleteCountryByID)
		}

		banner := admin.Group("/banners")
		{
			banner.GET("", h.AdminListBanner)
			banner.POST("", h.AdminCreateBanner)
			banner.GET("/:id", h.AdminGetBannerByID)
			banner.PATCH("/:id", h.AdminUpdateBanner)
			banner.DELETE("/:id", h.AdminDeleteBannerByID)

			bannerImage := banner.Group("/:id/banner_images")
			{
				bannerImage.POST("", h.AdminCreateBannerImage)
				bannerImage.DELETE("/:banner_image_id", h.AdminDeleteBannerImageByID)
			}
		}

		product := admin.Group("/products")
		{
			product.GET("", h.AdminListProduct)
			product.POST("", h.AdminCreateProduct)
			product.GET("/:id", h.AdminGetProductByID)
			product.PATCH("/:id", h.AdminUpdateProduct)
			product.DELETE("/:id", h.AdminDeleteProductByID)

			productImage := product.Group("/:id/product_images")
			{
				productImage.POST("", h.AdminCreateProductImage)
				productImage.DELETE("/:product_image_id", h.AdminDeleteProductImageByID)
			}
		}

		productRecommend := admin.Group("/product_recommends")
		{
			productRecommend.GET("", h.AdminListProductRecommend)
			productRecommend.POST("", h.AdminCreateProductRecommend)
			productRecommend.DELETE("/:id", h.AdminDeleteProductRecommendByID)
		}
	}

	srv := &http.Server{
		Addr:    conf.ServerAddr,
		Handler: engine,
	}

	return srv, nil
}

func initPgx(postgresConnection string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(postgresConnection)
	if err != nil {
		return nil, err
	}

	looger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	conf.ConnConfig.Logger = logrusadapter.NewLogger(looger)

	return pgxpool.ConnectConfig(context.Background(), conf)
}
