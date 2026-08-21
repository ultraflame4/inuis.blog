package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)


func setCallbackCookie(ctx *gin.Context, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   false, // false for now
		HttpOnly: true,
	}
	ctx.SetCookieData(c)
}
func (*ApiRoute) LoginFlow(c *gin.Context) {
	oidc_name := strings.ToUpper(c.Param("oidc"))
	config, ok := app.GetOidcProvider(oidc_name)
	if !ok{
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Unknown OIDC Provider",
		})
		return
	}

	state, err := app.RandString(16)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	setCallbackCookie(c, "state", state)

	c.Redirect(http.StatusFound, config.Config.AuthCodeURL(state))

}

func (a *ApiRoute) AcceptOAuthCallback(c *gin.Context){
	oidc_name := strings.ToUpper(c.Param("oidc"))
	conf, ok := app.GetOidcProvider(oidc_name)
	if !ok{
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Unknown OIDC Provider",
		})
		return
	}
	state, err := c.Cookie("state")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "state not found",
		})
		return
	}
	_ = state
	ctx := context.Background()
	oauth2Token, err := conf.Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange token: " + err.Error(),
		})
		return
	}
	userInfo, err := conf.Provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error" : "Failed to get userinfo: " + err.Error(),
		})
		return
	}

	// Find existing identities with this email.
	existing, err := gorm.G[models.UserIdentity](a.db).Where("email = ?", userInfo.Email).Find(c.Request.Context())
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
	}

	if len(existing) == 0 {
		log.Println("Creating new user for email: " + userInfo.Email)
		// Create new user
		email_parts := strings.SplitN(userInfo.Email, "@", 2)
		discriminator_hash := sha256.Sum256([]byte(email_parts[1]))
		discriminator := base64.RawURLEncoding.EncodeToString(discriminator_hash[:])[:4]
		// Create user
		user := models.User{
			Username: email_parts[0] + "_" + discriminator,
		}
		gorm.G[models.User](a.db).Create(c.Request.Context(), &user)
		// Attach identity
		gorm.G[models.UserIdentity](a.db).Create(c.Request.Context(), &models.UserIdentity{
			Subject: userInfo.Subject,
			Email:  &userInfo.Email,
			Issuer: conf.Issuer,
			UserID: user.ID,
		})

		// Create session

		resp := struct {
			OAuth2Token *oauth2.Token
			UserInfo    *oidc.UserInfo
		}{oauth2Token, userInfo}

		c.JSON(http.StatusOK, resp);
		return
	} else {
		// Fetch identity
		identity, err := gorm.G[models.UserIdentity](a.db).
			Preload("User", nil).
			Where("subject = ? AND issuer = ?", userInfo.Subject, conf.Issuer).
			First(c.Request.Context())
		if err != nil {
		   if errors.Is(err, gorm.ErrRecordNotFound){
				// User with this email exists. But does not have this identity registered.
				// Should reject oauth, ask user to login with correct identity first.
				// then connect provider when logged in
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"error" : "User with email " + userInfo.Email + " already exists, but does not have this identity registered!",
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error" : "Failed to fetch user identity: " + err.Error(),
				})
			}
			return
		}
		// Continue to get user object
		log.Printf("Found user identity: %s %s (sub: %s)\n", identity.Email, identity.Issuer, identity.Subject)
		user := identity.User
		log.Println("User logged in: " + user.Username)
	}




	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}

	c.JSON(http.StatusOK, resp);

}
