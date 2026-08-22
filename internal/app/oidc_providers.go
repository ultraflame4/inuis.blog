package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)


type OidcProvider struct {
	Verifier* oidc.IDTokenVerifier
	Config   oauth2.Config
	Provider* oidc.Provider
	Issuer string
}

var loadedProviders map[string]OidcProvider
func RegisterOIDCProviders(ctx context.Context) {
	server_addr := PublicServerAddr()

	loadedProviders = make(map[string]OidcProvider)
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

				slog.Error("Failed to create provider for " + oidc_name + ":\n\t", err.Error())
				continue
			}
			oidcConfig := oidc.Config{
				ClientID: client_id,
			}
			verifier := provider.Verifier(&oidcConfig)
			config := oauth2.Config{
				ClientID:  client_id,
				ClientSecret: os.Getenv("OIDC_" + oidc_name + "_CLIENT_SECRET"),
				RedirectURL: fmt.Sprintf("http://%s/api/oauth2/%s/callback", server_addr, strings.ToLower(oidc_name)) ,
				Endpoint: provider.Endpoint(),
				Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
			}
			loadedProviders[oidc_name] = OidcProvider{
				Verifier: verifier,
				Config:   config,
				Provider: provider,
				Issuer: oidc_url,
			}

			slog.Info("Registered provider for oidc name " + oidc_name)
			slog.Info("Login url for oidc '" + oidc_name + "' is " + fmt.Sprintf("http://%s/api/oauth2/%s/login", server_addr, strings.ToLower(oidc_name)))
			slog.Info("Callback url for oidc '" + oidc_name + "' is " + config.RedirectURL)


		}
	}

}


func GetOidcProvider(name string) (OidcProvider, bool) {
	provider, ok := loadedProviders[name]
	return provider, ok
}
