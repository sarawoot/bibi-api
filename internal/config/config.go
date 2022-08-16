package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var App = AppConfig{}

type AppConfig struct {
	ServerAddr string `long:"server-addr" env:"SERVER_ADDR"`
	JWKS       string `long:"jwks" env:"JWKS"`
	JWTTimeout int    `long:"jwt-timeout" env:"JWT_TIMEOUT"` // sec

	AdminJWKS       string `long:"admin-jwks" env:"ADMIN_JWKS"`
	AdminJWTTimeout int    `long:"admin-jwt-timeout" env:"ADMIN_JWT_TIMEOUT"` // sec

	PostgresConnection string `long:"postgres-connection" env:"POSTGRES_CONNECTION"`

	AssetsBaseURL string `long:"assets-base-url" env:"ASSETS_BASE_URL"`
	AWSBucket     string `long:"aws-bucket" env:"AWS_BUCKET"`
}

func init() {
	_, err := flags.Parse(&App)
	if err != nil {
		panic(err)
	}

	var port = os.Getenv("PORT")
	if port != "" {
		App.ServerAddr = ":" + port
	}
}
