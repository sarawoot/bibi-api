package handler

import "time"

const (
	jwtIssuer         = "bibi-api"
	jwtAud            = "bibi-app"
	adminJWTAud       = "bibi-admin"
	jwtExpirationTime = time.Hour * 87600
)
