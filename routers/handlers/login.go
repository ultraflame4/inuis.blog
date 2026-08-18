package handlers

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OidcProvider struct {
	Verifier* oidc.IDTokenVerifier
	Config   oauth2.Config
}

var oauthConfigs map[string]OidcProvider
func RegisterOIDCProviders(ctx context.Context) {
	oauthConfigs = make(map[string]OidcProvider)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		value := parts[1]
		key_parts := strings.SplitN(key, "_", 3)
		if key_parts[0] == "OIDC" && key_parts[2] == "URL" {
			oidc_name := key_parts[1]
			oidc_url := value
			client_id := os.Getenv("OIDC_" + oidc_name + "_CLIENT")
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
				RedirectURL: "http://127.0.0.1:8080/auth/google/callback",
				Endpoint: provider.Endpoint(),
				Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
			}
			oauthConfigs[oidc_name] = OidcProvider{
				Verifier: verifier,
				Config:   config,
			}

			log.Println("Registered provider for oidc name " + oidc_name)

		}
	}

}

func RegisterLoginHandlers(router *gin.RouterGroup){
	router.GET("/login/:oidc", loginFlow)
}


func loginFlow(c *gin.Context) {
	// oidc_name := c.Param("oidc")
	// oidc_url := os.Getenv("OIDC_"+oidc_name+"_URL")
	// client_id := os.Getenv("OIDC_"+oidc_name+"_CLIENT")


}
