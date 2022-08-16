package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
)

const userIDContextKey = "UserID"

func getBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	re := regexp.MustCompile(`(?i)^bearer\s(.*)`)
	match := re.FindStringSubmatch(authHeader)
	if len(match) == 0 {
		err := errors.New("not found token")
		return "", err
	}
	return match[1], nil
}

func setContextWithUserID(c *gin.Context, userID string) {
	c.Set(userIDContextKey, userID)
}

// func getContextUserID(c *gin.Context) string {
// 	return c.GetString(userIDContextKey)
// }

func hashPassword(password string) (string, error) {
	// example for making salt - https://play.golang.org/p/_Aw6WeWC42I
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func comparePasswords(storedPassword, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, err
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, err
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
