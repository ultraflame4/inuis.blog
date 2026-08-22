package app

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/gin-gonic/gin"
)

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HasCookie(c *gin.Context, key string) bool {
	_, err := c.Cookie(key)
	return err == nil
}