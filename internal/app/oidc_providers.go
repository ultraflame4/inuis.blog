package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/ultraflame4/bytesblog/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)


type OidcProvider struct {
	Verifier* oidc.IDTokenVerifier
	Config   oauth2.Config
	Provider* oidc.Provider
	Issuer string
}

var loadedProviders map[string]OidcProvider
var issuer2ProviderName map[string]string

func FindProviderName(issuer_url string) string{
	return issuer2ProviderName[issuer_url]
}

func RegisterOIDCProviders(ctx context.Context) {
	logger := slog.With("logger", "register_oidc")
	server_addr := GetConfig().PublicServerAddr()

	loadedProviders = make(map[string]OidcProvider)
	issuer2ProviderName = make(map[string]string)
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

				logger.Error("Failed to create provider for " + oidc_name + ":\n\t", err.Error())
				continue
			}
			oidcConfig := oidc.Config{
				ClientID: client_id,
			}
			verifier := provider.Verifier(&oidcConfig)
			redirectURL := fmt.Sprintf("http://%s/api/oidc/%s/callback", server_addr, strings.ToLower(oidc_name))
			config := oauth2.Config{
				ClientID:  client_id,
				ClientSecret: os.Getenv("OIDC_" + oidc_name + "_CLIENT_SECRET"),
				RedirectURL: redirectURL,
				Endpoint: provider.Endpoint(),
				Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
			}
			issuer2ProviderName[oidc_url] = oidc_name
			loadedProviders[oidc_name] = OidcProvider{
				Verifier: verifier,
				Config:   config,
				Provider: provider,
				Issuer: oidc_url,
			}

			slog.Info("Registered oidc provider", "name", oidc_name)
			slog.Info(oidc_name,"login",fmt.Sprintf("http://%s/api/oidc/%s/login", server_addr, strings.ToLower(oidc_name)))
			slog.Info(oidc_name,"callback",config.RedirectURL)


		}
	}

}


// Get GetOidcProvider by name. Case insensitive.
func GetOidcProvider(name string) (OidcProvider, bool) {
	provider, ok := loadedProviders[strings.ToUpper(name)]
	return provider, ok
}

// Return list of oidc provider names
func GetProviders()  []string {
	names := make([]string, 0, len(loadedProviders))
	for name, _ := range loadedProviders {
		names = append(names, strings.ToLower(name))
	}
	return names
}


// Remove this OidcProvider from the user
func (p *OidcProvider) RemoveFromUser(db *gorm.DB, ctx context.Context, userID uint) error {
	issuer := p.Issuer
	_, err := gorm.G[models.UserIdentity](db).Where("user_id = ? AND issuer = ?", userID, issuer).Delete(ctx)
	return err
}
