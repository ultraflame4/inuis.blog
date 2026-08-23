package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teenyblog/internal/app"
	"github.com/teenyblog/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// Set secure cookies
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

// Generate PKCE challenge from verifier
func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// Handlers
func (*ApiHandler) StartOauthLoginFlow(c *gin.Context) {
	oidc_name := strings.ToUpper(c.Param("oidc"))
	config, ok := app.GetOidcProvider(oidc_name)
	if !ok {
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

	if len(config.Config.ClientSecret) == 0 {
		// No client secret. Default to pkce
		pkceVerifier, err := app.RandString(32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate PKCE verifier: " + err.Error(),
			})
			return
		}
		setCallbackCookie(c, "pkceVerifier", pkceVerifier)
		challenge := pkceChallenge(pkceVerifier)
		c.Redirect(http.StatusFound,
			config.Config.AuthCodeURL(state,
			oauth2.SetAuthURLParam("code_challenge", challenge),
				oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			))

	}
	setCallbackCookie(c, "state", state)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate PKCE verifier: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, config.Config.AuthCodeURL(state))

}


func (a *ApiHandler) AcceptOAuthCallback(c *gin.Context) {
	oidc_name := strings.ToUpper(c.Param("oidc"))
	conf, ok := app.GetOidcProvider(oidc_name)
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Unknown OIDC Provider",
		})
		return
	}

	ctx := context.Background()

	
	var oauth2Token *oauth2.Token
	if verifier , err := c.Cookie("pkceVerifier"); err == nil {
		token, err := conf.Config.Exchange(
			ctx,
			c.Query("code"),
			oauth2.SetAuthURLParam("code_verifier", verifier),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to exchange token: " + err.Error(),
			})
			return
		}
		oauth2Token = token
	} else {
		state, err := c.Cookie("state")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "state not found",
			})
			return
		}
		_ = state
	
		token, err := conf.Config.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to exchange token: " + err.Error(),
			})
			return
		}
		oauth2Token = token
	}
	

	userInfo, err := conf.Provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get userinfo: " + err.Error(),
		})
		return
	}

	// ---- Get / Create user


	
	// Find existing identities with this email.
	existing, err := gorm.G[models.UserIdentity](a.db).Where("email = ?", userInfo.Email).Find(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError);
		return
	}


	var user models.User;
	if len(existing) == 0 {
		slog.Info("Creating new user for email: " + userInfo.Email)
		// Create new user
		email_parts := strings.SplitN(userInfo.Email, "@", 2)
		discriminator_hash := sha256.Sum256([]byte(email_parts[1]))
		discriminator := base64.RawURLEncoding.EncodeToString(discriminator_hash[:])[:8]
		// Create user
		user = models.User{
			Username: email_parts[0] + "_" + discriminator,
		}
		if err := a.db.Transaction(func (tx *gorm.DB) error {
			if err := gorm.G[models.User](tx).Create(c.Request.Context(), &user); err != nil {
				return err
			}
			// Attach identity
			identity := models.UserIdentity{
				Subject: userInfo.Subject,
				Email:   &userInfo.Email,
				Issuer:  conf.Issuer,
				UserID:  user.ID,
			}
			if err := gorm.G[models.UserIdentity](tx).Create(c.Request.Context(), &identity); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		slog.Info("Created user: " + user.Username)

	} else {
		// Fetch identity
		identity, err := gorm.G[models.UserIdentity](a.db).
			Preload("User", nil).
			Where("subject = ? AND issuer = ?", userInfo.Subject, conf.Issuer).
			First(c.Request.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// User with this email exists. But does not have this identity registered.
				// Should reject oauth, ask user to login with correct identity first.
				// then connect provider when logged in
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"error": "User with email " + userInfo.Email + " already exists, but does not have this identity registered!",
				})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to fetch user identity: " + err.Error(),
				})
			}
			return
		}
		// Continue to get user object
		slog.Info("Found user identity: %s %s (sub: %s)\n", identity.Email, identity.Issuer, identity.Subject)
		user = identity.User
	}


	// Revoke any existings session from same browser / cookie
	if existing_session, err := c.Cookie("session_token"); err == nil {
		slog.Debug("Revoking existing session with token: " + existing_session)

		if err := app.RevokeSession(a.db, ctx, existing_session); err != nil {
			slog.Error("Failed to revoke existing session: " + err.Error())
		}
	}
	// Create user session
	slog.Debug("Creating session for user: " + user.Username)
	session, err := app.CreateSession(a.db, c.Request.Context(), user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session: " + err.Error(),
		})
		return
	}

	// resp := struct {
	// 	OAuth2Token *oauth2.Token
	// 	UserInfo    *oidc.UserInfo
	// 	User        *models.User
	// 	Session     *models.UserSession
	// }{oauth2Token, userInfo, &user, session}

	c.SetCookie("session_token",string(session.Token), app.SessionMaxAge(), "", "", false, true);

	// c.JSON(http.StatusOK, resp)
	c.Redirect(http.StatusFound, "/me")

}


func (a *ApiHandler) ListLoginProviders(c *gin.Context) {
	providers := app.GetProviders()
	c.JSON(http.StatusOK, providers)
}