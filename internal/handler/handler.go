package handler

import (
	"api/internal/config"
	"api/internal/repo/datastorepgx"
	"api/pkg/aws/s3"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Handler struct {
	appConfig  *config.AppConfig
	jwks       jwk.Set
	jwksPublic jwk.Set

	adminJWKS       jwk.Set
	adminJWKSPublic jwk.Set

	datastoreRepo *datastorepgx.DataStoreRepo

	s3Client *s3.Client
}

func NewHandler(appConfig *config.AppConfig, datastoreRepo *datastorepgx.DataStoreRepo, s3Client *s3.Client) (*Handler, error) {
	jwks, err := jwk.Parse([]byte(appConfig.JWKS))
	if err != nil {
		return nil, err
	}

	jwksPublic, err := jwk.PublicSetOf(jwks)
	if err != nil {
		return nil, err
	}

	adminJWKS, err := jwk.Parse([]byte(appConfig.AdminJWKS))
	if err != nil {
		return nil, err
	}

	adminJWKSPublic, err := jwk.PublicSetOf(adminJWKS)
	if err != nil {
		return nil, err
	}

	return &Handler{
		appConfig:  appConfig,
		jwks:       jwks,
		jwksPublic: jwksPublic,

		adminJWKS:       adminJWKS,
		adminJWKSPublic: adminJWKSPublic,

		datastoreRepo: datastoreRepo,

		s3Client: s3Client,
	}, nil
}
