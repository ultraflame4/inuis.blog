package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/teenyblog/utils"
	"golang.org/x/oauth2"
)

type OidcProvider struct {
	Verifier* oidc.IDTokenVerifier
	Config   oauth2.Config
	Provider* oidc.Provider
}

var oauthConfigs map[string]OidcProvider
func RegisterOIDCProviders(ctx context.Context) {
	redirect_origin := utils.GetServerAddress()

	oauthConfigs = make(map[string]OidcProvider)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		value := parts[1]
		key_parts := strings.SplitN(key, "_", 3)
		if key_parts[0] == "OIDC" && key_parts[2] == "URL" {
			oidc_name := key_parts[1]
			oidc_url := value
			client_id := os.Getenv("OIDC_" + oidc_name + "_CLIENT_ID")
			provider, err := oidc.NewProvider(ctx, oidc_url)
			if err != nil {

				log.Println("Failed to create provider for " + oidc_name + ":\n\t", err.Error())
				continue
			}
			oidcConfig := oidc.Config{
				ClientID: client_id,
			}
			verifier := provider.Verifier(&oidcConfig)
			config := oauth2.Config{
				ClientID:  client_id,
				ClientSecret: os.Getenv("OIDC_" + oidc_name + "_CLIENT_SECRET"),
				RedirectURL: fmt.Sprintf("http://%s/api/oauth2/%s/callback", redirect_origin, strings.ToLower(oidc_name)) ,
				Endpoint: provider.Endpoint(),
				Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
			}
			oauthConfigs[oidc_name] = OidcProvider{
				Verifier: verifier,
				Config:   config,
				Provider: provider,
			}

			log.Println("Registered provider for oidc name " + oidc_name)
			log.Println("Callback url for oidc '" + oidc_name + "' is " + config.RedirectURL)


		}
	}

}

func RegisterLoginHandlers(router *gin.RouterGroup){
	router.GET("/login/:oidc", loginFlow)
	router.GET("/oauth2/:oidc/callback", acceptOAuthCallback)
}


func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

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
func loginFlow(c *gin.Context) {
	oidc_name := strings.ToUpper(c.Param("oidc"))
	config, ok := oauthConfigs[oidc_name]
	if !ok{
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Unknown OIDC Provider",
		})
		return
	}

	state, err := randString(16)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	setCallbackCookie(c, "state", state)

	c.Redirect(http.StatusFound, config.Config.AuthCodeURL(state))

}

func acceptOAuthCallback(c *gin.Context){
	oidc_name := strings.ToUpper(c.Param("oidc"))
	conf, ok := oauthConfigs[oidc_name]
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

	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}
	c.JSON(http.StatusOK, resp);

}
