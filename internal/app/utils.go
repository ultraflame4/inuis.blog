package app

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"regexp"
	"strings"
	"github.com/google/uuid"
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
func Slugify(s string) string {
    s = strings.ToLower(s)
    s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
    return strings.Trim(s, "-")
}

func GenerateHID() string {
	uuid_ := uuid.New()
	return base64.RawStdEncoding.EncodeToString(uuid_[:])
}
