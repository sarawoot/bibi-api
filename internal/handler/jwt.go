package handler

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func (h *Handler) signJWT(userID uuid.UUID) (string, error) {
	token := jwt.New()
	currentTime := time.Now()
	tokenSetMap := map[string]interface{}{
		jwt.IssuerKey:     jwtIssuer,
		jwt.SubjectKey:    userID.String(),
		jwt.AudienceKey:   []string{jwtAud},
		jwt.ExpirationKey: currentTime.Add(jwtExpirationTime),
		jwt.IssuedAtKey:   currentTime,
	}

	for k, v := range tokenSetMap {
		if err := token.Set(k, v); err != nil {
			return "", err
		}
	}

	key, ok := h.jwks.Key(0)
	if !ok {
		return "", errors.New("not jwk key 0")
	}

	tokenJWT, err := jwt.Sign(token, jwt.WithKey(key.Algorithm(), key))
	if err != nil {
		return "", err
	}

	return string(tokenJWT), nil
}

func (h *Handler) signAdminJWT(adminID uuid.UUID) (string, error) {
	token := jwt.New()
	currentTime := time.Now()
	tokenSetMap := map[string]interface{}{
		jwt.IssuerKey:     jwtIssuer,
		jwt.SubjectKey:    adminID.String(),
		jwt.AudienceKey:   []string{adminJWTAud},
		jwt.ExpirationKey: currentTime.Add(jwtExpirationTime),
		jwt.IssuedAtKey:   currentTime,
	}

	for k, v := range tokenSetMap {
		if err := token.Set(k, v); err != nil {
			return "", err
		}
	}

	key, ok := h.adminJWKS.Key(0)
	if !ok {
		return "", errors.New("not jwk key 0")
	}

	tokenJWT, err := jwt.Sign(token, jwt.WithKey(key.Algorithm(), key))
	if err != nil {
		return "", err
	}

	return string(tokenJWT), nil
}
