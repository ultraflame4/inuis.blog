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

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/ultraflame4/inuis/internal/app"
	"github.com/ultraflame4/inuis/models"
	"github.com/ultraflame4/inuis/routers/middlewares"
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

// Retrieves session from cookie, aborts with error if session is non-existent or invalid.
//
// You probably want `meta.GetSession` instead which is protected on the middleware level.
//
// Should only use this when your endpoint has features that *may* require an authenticated session.
// Such as:
// - Elevation
// - Action(s) on protected resource
//
// Even then, you should try to seperate the protected / unprotected functionality into seperate
// endpoints to prevent confusion.
//
// This method really only exists for Oauth callback where a single only a single endpoint is used
// for callbacks by Identity Providers. (Having multiple makes it inconvenient, especially if we add more feautres / intents).
//
// Returns (user session, aborted)
func (a *ApiHandler) checkSession(c *gin.Context) (*models.UserSession, bool) {
	sid, _ := c.Cookie(middlewares.COOKIE_SESSION_TOKEN)
	s := app.GetSession(a.db, c.Request.Context(), sid)
	if s == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User not logged in.",
		})
		return nil, true
	}
	return s, false
}

var oauth2_SelectAccount = oauth2.SetAuthURLParam("prompt", "select_account")

// User wants to add a new identity
const IntentAddIdentity = "add_identity"
const IntentElevate = "elevate"
const IntentSignup = "signup"

// Handlers
func (a *ApiHandler) StartOauthLoginFlow(c *gin.Context) {
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

	// If intent is to signup, do some prechecks before we start oauth flow
	if intent == IntentSignup {
		if fail := a.handleAllowSignupsPrechecks(inviteCode, c); fail {
			return
		}
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
				oauth2_SelectAccount,
			))
		return
	} else {
		// Not pkce flow. remove set pkce cookie to false
		// Use to prevent conflicts when switching provider flows.
		setCallbackCookie(c, "pkceVerifier", "false")
		c.Redirect(http.StatusFound, config.Config.AuthCodeURL(state, oauth2_SelectAccount))
		return
	}
	// Nothing below here. Requests are fully handled in if-else statement
}

func (a *ApiHandler) handleAllowSignupsPrechecks(inviteCode string, c *gin.Context) bool {
	adminless := a.config.IsAdminless(a.db, c.Request.Context())
	require_invite := a.config.RequireInviteCode || adminless
	if a.config.ALLOW_SIGNUP || adminless {
		// Check if invite code exists before proceeding with oauth flow
		has_invite := len(inviteCode) != 0
		if has_invite {
			// Check if invite code is valid
			if valid := app.VerifyInviteCode(a.db, c.Request.Context(), inviteCode); !valid {
				c.Redirect(http.StatusFound, "/failure?reason=bad_invite")
				return true
			}
		} else if require_invite {
			c.Redirect(http.StatusFound, "/failure?reason=require_invite")
			return true
		}
	} else {
		c.Redirect(http.StatusFound, "/failure?reason=signup_disable")
		return true
	}
	return false
}

func (a *ApiHandler) AcceptOAuthCallback(c *gin.Context) {
	oidc_name := strings.ToUpper(c.Param("oidc"))
	provider, ok := app.GetOidcProvider(oidc_name)
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
		token, err := provider.Config.Exchange(
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

		token, err := provider.Config.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to exchange token: " + err.Error(),
			})
			return
		}
		oauth2Token = token
	}

	idtoken, err := provider.Provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get userinfo: " + err.Error(),
		})
		return
	}

	intent, _ := c.Cookie("intent")

	var user models.User

	// --- Active Session flows ---
	// - These intents require an active session
	switch intent {
	case IntentAddIdentity:
		a.handleIntentAddIdentity(c, idtoken, provider, oidc_name)
		return
	case IntentElevate:
		a.handleIntentElevate(c, idtoken, provider)
		return
	}

	// --- Login Session + Signup Flow ---
	// - When logging in or signing up
	// - Flow creates new active session

	// Find existing identities with this email.
	existing, err := app.FindUserIdentitiesByEmail(idtoken.Email, a.db, c.Request.Context())
	existing_issuers := make([]string, 0, len(existing))
	for _, e := range existing {
		existing_issuers = append(existing_issuers, e.Issuer)
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// No identities with this email found.
	// Require signup flow.
	if len(existing) == 0 {
		// Require signup intent!
		if intent != IntentSignup {
			c.Redirect(http.StatusFound, fmt.Sprintf("/failure?reason=signup_required&email=%s", idtoken.Email))
			return
		}

		new_user, ok := a.createUser(c, idtoken, provider)
		if !ok {
			return
		}
		user = new_user

	} else {
		// Fetch identity using subject and issuer
		identity, err := app.FindUserIdentity(idtoken.Subject, provider.Issuer, a.db, c.Request.Context())
		if err != nil {
			// User exists but we cannot verify that the identity they logged in with is valid / connected to the account.
			a.handleIdentityError(err, existing_issuers, provider, existing, c, idtoken)
			return
		}
		// Continue to get user object
		slog.Info("Found user identity", "email", identity.Email, "issuer", identity.Issuer, "sub", identity.Subject)
		user = identity.User
	}

	// Revoke any existings session from same browser / cookie
	session, aborted := a.loginUserSession(ctx, c, user)
	if aborted {
		return
	}

	c.SetCookie("session_token", string(session.Token), a.config.MaxSessionAge, "", "", false, true)
	c.Redirect(http.StatusFound, "/me")

}

// Logins and creates the user ssession
func (a *ApiHandler) loginUserSession(ctx context.Context, c *gin.Context, user models.User) (*models.UserSession, bool) {
	if existing_session, err := c.Cookie("session_token"); err == nil {
		slog.Debug("Revoking existing session with token: " + existing_session)
		if err := app.RevokeSessionToken(a.db, ctx, existing_session); err != nil {
			slog.Warn("Failed to revoke existing session: " + err.Error())
		}
	}
	// Create user session
	slog.Debug("Creating session for user: " + user.Username)
	session, err := app.CreateSession(a.db, c.Request.Context(), user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session: " + err.Error(),
		})
		return nil, true
	}
	return session, false
}

// Handle errors when fetching existing user's identity
func (*ApiHandler) handleIdentityError(err error, existing_issuers []string, conf app.OidcProvider, existing []models.UserIdentity, c *gin.Context, userInfo *oidc.UserInfo) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// existing user has identity from the provider
		// but our records does not have it recorded.
		//
		// User likely did not connect the identity with us.

		// Check if this identity provider is registed with user
		// If not found, return "bad_identity" error.
		ident_idx := slices.Index(existing_issuers, conf.Issuer)
		if ident_idx == -1 {
			// List known identities for existing user
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
		// Edge case here, the user has this identity registered but with as different subject.
		// This can happen when the server changes the OIDC client token, or secret
		// Should reject oauth, and return "bad_subject" error.

		// Get user's identity, and show the expected subject value. And also include a list of providers
		// the user can log in with.
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
		// Generic catch all error
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user identity: " + err.Error(),
		})
	}
	return
}

func (a *ApiHandler) createUser(c *gin.Context, userInfo *oidc.UserInfo, conf app.OidcProvider) (models.User, bool) {
	user := models.User{}

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
			return user, false
		}
		inviteCode := app.ConsumeInviteCode(a.db, c.Request.Context(), inviteCode_s)
		if inviteCode == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "Invite code has expired or is invalid.",
				"invite_code": inviteCode_s,
			})
			return user, false
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
		return user, false
	}
	user = new_user
	slog.Info("Created user: " + user.Username)
	return user, true
}

func (a *ApiHandler) handleIntentElevate(c *gin.Context, userInfo *oidc.UserInfo, conf app.OidcProvider) {
	// Cant elevate if the user is not logged in.
	session, aborted := a.checkSession(c)
	if aborted {
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

func (a *ApiHandler) handleIntentAddIdentity(c *gin.Context, userInfo *oidc.UserInfo, conf app.OidcProvider, oidc_name string) {
	// Cant add identity if the user is not logged in.
	session, aborted := a.checkSession(c)
	if aborted {
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
}
