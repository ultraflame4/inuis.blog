# API Reference

All API routes are prefixed with `/api`.

## Overview

| Group | Description |
|---|---|
| [App Meta](#app-meta) | Application metadata |
| [OIDC Auth](#oidc-authentication) | Login flow via external providers |
| [Current User (`/api/me`)](#current-user-session-apime) | Authenticated user actions (posts, profile, sessions) |
| [Public Posts](#public-posts) | Read published posts |
| [Public Users](#public-user-routes) | Browse a user's published posts |

## Objects

### Session Object

```json
{
  "user_id": "test_user",
  "username": "Test User",
  "permissions": []
}
```

| Field | Type | Description |
|---|---|---|
| `user_id` | string | Unique user identifier |
| `username` | string | Display name |
| `permissions` | string[] | UI-facing permission flags; actual authorization is enforced server-side |

---

## App Meta

`GET /api/meta`

Returns application-level metadata.

---

## OIDC Authentication

OIDC authentication is handled externally by configured providers.

See [Identity Providers](identity-providers.md)

### Routes

`GET /api/oidc/providers`
: List available OIDC login providers.

`GET /api/oidc/:oidc/login`
: Start the OAuth login flow for the given provider. Provider names are case-insensitive (e.g. `google` or `GOOGLE`).

`GET /api/oidc/:oidc/callback`
: OAuth callback endpoint — the external provider redirects here after authentication.

---

## Current User Session (`/api/me`)

All routes in this group require an authenticated session.

### Profile & Sessions

`GET /api/me/`
: Get current user's session info. Returns a [Session object](#session-object).

`PUT /api/me/username`
: Update current user's display name.

`GET /api/me/logout`
: Log out the current session only.

`GET /api/me/logout/all`
: Log out all sessions for the current user.

### OIDC Providers

`GET /api/me/oidc-providers`
: List OIDC providers connected to the current user.

`GET /api/me/oidc/:oidc/disconnect`
: Disconnect a linked OIDC provider.

### Posts

`POST /api/me/posts/save`
: Create or update a post (saved as a draft).

`GET /api/me/posts`
: List all posts belonging to the current user.

`GET /api/me/posts/:id`
: Get a specific post owned by the current user.

`DELETE /api/me/posts/:id`
: Delete a post owned by the current user.

`POST /api/me/posts/:id/publish`
: Publish a draft post.

---

## Public Posts

`GET /api/posts/:id`
: Get a published post by ID.

---

## Public User Routes

`GET /api/user/:user_id/posts`
: List a specific user's published posts.

---

## Authentication

### Login Flow (MFA)

**Step 1 —** `POST /login`

Client submits credentials. If MFA is not configured, a session is created immediately. If MFA is required, the server responds with:

```json
{
  "mfa_required": true,
  "challenge_id": "abcdefg..."
}
```

The `challenge_id` is a JWT containing a `started_at` timestamp, allowing the server to enforce dynamic expiry.

**Step 2 —** `POST /login/2fa`

Client submits the MFA code along with the challenge ID. On success, a session is created.

### Notes

- Account recovery and communication require a mail server and are not yet implemented; reliance on external OAuth providers is preferred.
- `permissions` in the session object are informational for the UI — actual authorization checks happen server-side.