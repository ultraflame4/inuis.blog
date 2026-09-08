package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/bytesblog/internal/app"
	"github.com/ultraflame4/bytesblog/models"
	"github.com/ultraflame4/bytesblog/routers/middlewares"
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
	setCallbackCookie(c, "state", state)

	// -- Store query options in cookies
	// You might be wondering, why are we not only setting the cookies if query exists.
	// Because:
	// 1. If query don't exists, we need to clear the cookie.
	// 2. So why not, just set the cookies anyways regardless if query exists or not
	// 3. If query doesn't exist, the values will be "" (empty). We can check for that on the callback.

	inviteCode := c.Query("invite_code")
	setCallbackCookie(c, "inviteCode", inviteCode)

	intent := c.Query("intent")
	setCallbackCookie(c, "intent", intent)
	// -- end

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
		return
	} else {
		// Not pkce flow. remove set pkce cookie to false
		// Use to prevent conflicts when switching provider flows.
		setCallbackCookie(c, "pkceVerifier", "false")

		c.Redirect(http.StatusFound, config.Config.AuthCodeURL(state))
		return
	}
	// Nothing below here. Requests are fully handled in if-else statement
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

	state, err := c.Cookie("state")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "state not found",
		})
		return
	}
	if state != c.Query("state") {
		// slog.Warn(fmt.Sprintf("Expected state %s got state %s", state, c.Query("state")))
		setCallbackCookie(c, "state", "-")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "state mismatched!",
		})
		return
	}

	ctx := context.Background()

	var oauth2Token *oauth2.Token
	if verifier, err := c.Cookie("pkceVerifier"); err == nil && verifier != "false" {
		token, err := conf.Config.Exchange(
			ctx,
			c.Query("code"),
			oauth2.SetAuthURLParam("code_verifier", verifier),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "(PKCE) Failed to exchange token: " + err.Error(),
			})
			return
		}
		oauth2Token = token
	} else {

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
	existing_issuers := make([]string, 0, len(existing))
	for _, e := range existing {
		existing_issuers = append(existing_issuers, e.Issuer)
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	intent, _ := c.Cookie("intent")
	// User wants to add a new identity
	const IntentAddIdentity = "add_identity"
	const IntentElevate = "elevate"
	const IntentSignup = "signup"

	var user models.User

	switch intent {
	case IntentAddIdentity:
		// Cant add identity if the user is not logged in.
		sid, _ := c.Cookie(middlewares.COOKIE_SESSION_TOKEN)
		session := app.GetSession(a.db, c.Request.Context(), sid)
		if session == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not logged in.",
			})
			return
		}

		user := session.User

		// Check if the is already in use!
		existing, err := app.FindUserIdentity(userInfo.Subject, conf.Issuer, a.db, c.Request.Context())
		if err == nil {
			// If existing identity already belongs to user, ignroe and  just return success
			if existing.UserID == user.ID {
				// Ignore, do nothing
			} else {
				// Otherwise redirect to error (Identity in use!)
				c.Redirect(http.StatusFound,
					fmt.Sprintf("/failure?reason=add_identity_in_use&provider=%s&email=%s", oidc_name, userInfo.Email),
				)
				return
			}
		} else {
			// Attach identity
			identity := models.UserIdentity{
				Subject: userInfo.Subject,
				Email:   &userInfo.Email,
				Issuer:  conf.Issuer,
				UserID:  user.ID,
			}
			if err := gorm.G[models.UserIdentity](a.db).Create(c.Request.Context(), &identity); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to add identity. Error: " + err.Error(),
				})
				return
			}
		}

		c.Redirect(http.StatusFound, "/success?reason=add_identity&provider="+oidc_name)
		return

	case IntentElevate:
		// Cant elevate if the user is not logged in.
		sid, _ := c.Cookie(middlewares.COOKIE_SESSION_TOKEN)
		session := app.GetSession(a.db, c.Request.Context(), sid)
		if session == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not logged in.",
			})
			return
		}
		// Check if the identity user elevated with matches the user session!
		identity, err := app.FindUserIdentity(userInfo.Subject, conf.Issuer, a.db, c.Request.Context())
		if err != nil || identity.UserID != session.UserID {
			c.Redirect(http.StatusFound, "/failure?reason=elevate_fail")
			return
		}
		session.ElevateNow(a.db, c.Request.Context())
		c.Redirect(http.StatusFound, "/success?reason=session_elevated")
		return
	}

	// No identities with this email found. Create a new user.
	if len(existing) == 0 {

		// Require signup intent!
		if intent != IntentSignup {
			c.Redirect(http.StatusFound, fmt.Sprintf("/failure?reason=signup_required&email=%s", userInfo.Email))
			return
		}

		// Create empty user
		user = models.User{}

		// Check for invite code requirements
		// - Get invite code from cookie
		// - Can safely ignore error, inciteCode_s should be empty if not present. (or has been reset by login)
		inviteCode_s, _ := c.Cookie("inviteCode")
		inviteCode_present := len(inviteCode_s) != 0

		var grantRoles string
		// - Do invite code checks if: 1. required, 2. invite code is present.
		if a.config.RequireInviteCode || inviteCode_present {
			if !inviteCode_present {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invite code required.",
				})
				return
			}
			inviteCode := app.ConsumeInviteCode(a.db, c.Request.Context(), inviteCode_s)
			if inviteCode == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":       "Invite code has expired or is invalid.",
					"invite_code": inviteCode_s,
				})
				return
			}
			slog.Info("Invite code consumed: " + inviteCode_s)
			grantRoles = *inviteCode.GrantRoles
		}

		slog.Info("Creating new user for email: " + userInfo.Email)
		// Create new user
		email_parts := strings.SplitN(userInfo.Email, "@", 2)
		discriminator_hash := sha256.Sum256([]byte(email_parts[1]))
		discriminator := base64.RawURLEncoding.EncodeToString(discriminator_hash[:])[:8]
		// Create user
		username := email_parts[0] + "_" + discriminator

		identity := models.UserIdentity{
			Subject: userInfo.Subject,
			Email:   &userInfo.Email,
			Issuer:  conf.Issuer,
		}

		new_user, err := app.CreateUser(username, identity, grantRoles, a.db, c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user = new_user
		slog.Info("Created user: " + user.Username)

	} else {
		// Fetch identity
		identity, err := app.FindUserIdentity(userInfo.Subject, conf.Issuer, a.db, c.Request.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Here we already know there is a user with this email.

				// There is a technically an edge case here where the user has identity registered, but not with the same subject
				// AKA their subject for the identity has changed somehow. (eg. IDP changed but using same domain)
				// We are going to ignore this edge case and treat it as a "bad_identity" error.

				// Try get expected identity for this user via issuer, from the list of user's existing identities
				ident_idx := slices.Index(existing_issuers, conf.Issuer)
				// If not found, return "bad_identity" error.
				if ident_idx == -1 {
					var allowed_providers strings.Builder
					for idx, value := range existing {
						provider_name := strings.ToLower(app.FindProviderName(value.Issuer))
						if idx != 0 {
							allowed_providers.WriteByte(',')
						}
						allowed_providers.WriteString(provider_name)
					}

					// Identity not registered.
					// Should reject oauth, ask user to login with correct identity first, then connect provider when logged in
					c.Redirect(http.StatusFound, fmt.Sprintf("/failure?reason=bad_identity&allowed=%s&email=%s", allowed_providers.String(), userInfo.Email))
					return
				}
				// Edge case here, the user has this identity registered but with a different subject.
				// Should reject oauth, and return "bad_subject" error.

				// Retrieve the user's existing identity for this issuer
				identity := existing[ident_idx]

				var allowed_providers strings.Builder
				for idx, value := range existing {
					// Skip the current provider, as it's subject field is now non-working
					if value.Issuer == conf.Issuer {
						continue
					}
					provider_name := strings.ToLower(app.FindProviderName(value.Issuer))
					if idx != 0 {
						allowed_providers.WriteByte(',')
					}
					allowed_providers.WriteString(provider_name)
				}

				c.Redirect(http.StatusFound, fmt.Sprintf("/failure?reason=bad_subject&expected=%s&received=%s&email=%s&allowed=%s", identity.Subject, userInfo.Subject, userInfo.Email, allowed_providers.String()))
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

		if err := app.RevokeSessionToken(a.db, ctx, existing_session); err != nil {
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

	c.SetCookie("session_token", string(session.Token), a.config.MaxSessionAge, "", "", false, true)

	// c.JSON(http.StatusOK, resp)
	c.Redirect(http.StatusFound, "/me")

}

func (a *ApiHandler) ListLoginProviders(c *gin.Context) {
	providers := app.GetProviders()
	c.JSON(http.StatusOK, providers)
}
