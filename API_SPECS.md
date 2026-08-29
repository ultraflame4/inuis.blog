# Overview
This are specifications for a backend of a supersimple
blogging engine / app.


# OIDC Authentication 
All OIDC authentication is handled by the external OIDC provider.
All providers must support PKCE as it is the only supported flow.

## Configuring Providers
OIDC providers are configured via environment variables.

You need to set the following environment variables to configure an OIDC provider:
- `OIDC_<PROVIDER_NAME>_URL`: URL of the OIDC provider. This is used to get the oidc address.
  - Example: `https://accounts.google.com` will resolve to `https://accounts.google.com/.well-known/openid-configuration`
  - DO NOT PUT the DISCOVERY address! Use the portion WITHOUT `/.well-known/openid_configuration`!
- `OIDC_<PROVIDER_NAME>_CLIENT_ID`: Client ID for OIDC provider
- `OIDC_<PROVIDER_NAME>_CLIENT_SECRET`: Client secret for OIDC provider
where `<PROVIDER_NAME>` is the name of the provider (e.g. `GOOGLE`).

User can then login via `/login/<PROVIDER_NAME>`

Example for Google:
```bash
OIDC_GOOGLE_URL="https://accounts.google.com"
OIDC_GOOGLE_CLIENT="..."
OIDC_GOOGLE_SECRET="..."
```
User can then login via `/login/google` or `/login/GOOGLE` (case-insensitive)

# Objects
## Users

### Session Object
What a session object might look like:
```json
{
  "user_id": "test_user",
  "username": "Test User",
  "permissions": [
    // list of user permissions for UI.
    // Actual checks is done in backend.
  ]
}
```

## Posts

# API Routes

All routes are prefixed with `/api`.

## App Meta
`GET /api/meta` — Get application metadata

## OIDC Authentication
`GET /api/oidc/providers` — List available OIDC login providers
`GET /api/oidc/:oidc/login` — Start OAuth login flow for a provider
`GET /api/oidc/:oidc/callback` — Accept OAuth callback from a provider

## Current User Session (`/api/me`)
All routes in this group require an authenticated session.

`GET /api/me/` — Get current user session info
`PUT /api/me/username` — Update current user's username
`GET /api/me/logout` — Log out current session
`GET /api/me/logout/all` — Log out all sessions
`GET /api/me/oidc-providers` — List connected OIDC providers for current user
`GET /api/me/oidc/:oidc/disconnect` — Disconnect an OIDC provider
`POST /api/me/posts/save` — Create or update a post (draft)
`GET /api/me/posts` — List current user's posts
`GET /api/me/posts/:id` — Get a specific post owned by current user
`DELETE /api/me/posts/:id` — Delete a post owned by current user
`POST /api/me/posts/:id/publish` — Publish a draft post

## Public Posts
`GET /api/posts/:id` — Get a published post by ID

## Public User Routes
`GET /api/user/:user_id/posts` — List a specific user's published posts

# Login & User Identity

## Notes

## Email Notifications.
User account resets + communication requires minimally an email server.
Hence might just rely on external services and use OAuth for user authentication.
(i.e bluesky, github, google etc.)

## Login Flow (2FA)
Proposed 2FA login flow

1. `/login` -> Client submits username and passowrd
  1. Server verifies username & password
  2. Server checks for 2FA (if no 2FA, just create session)
  3. If 2FA, return 
     ```json
     {
        mfa_required: true
        challenge_id: "abcdefg..." // Likely to be a JWT token with a started_at timestamp field. So server can dynamically change expiry.
     }
     ```
2. `/login/2fa` -> Client submits login with challenge id
  1. Server authenticates and creates session
