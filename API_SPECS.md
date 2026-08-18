# Overview
This are specifications for a backend of a supersimple
blogging engine / app.


# OIDC Authentication 
All OIDC authentication is handled by the external OIDC provider.
All providers must support PKCE as it is the only supported flow.
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

## Public User Routes
Login
`/user/login`
Login for MFA
`/user/login/2fa`
Create User
`/user/new`
Reset Password 
`/user/reset-password`

## Current User Session
Logout
`/user/logout`
Delete account
`/user/delete`
Query current session
`/user`
Query posts for current session.
`/user/posts`

## Post Actions
GET    - `/posts` (Get Posts, can filter using queryies like: `?user=<user_id>`)
GET    - `/posts/<post_id>` (Get specific Post) (Need permissions)
UPDATE - `/posts/<post_id>` (Update Post)       (Need permissions)
DELETE - `/posts/<post_id>` (Delete Post)       (Need permissions)
POST   - `/posts/new` (Create Post)             (Need permissions)

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
