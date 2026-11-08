# User Management

## Signup Behaviour

- Users can only sign in via an external [Identity Provider](04-identity-providers.md).
- Server owners can enable **invite-only mode**, which requires invite codes for all signups.
- If no admin users exist, invite codes are **always required**, regardless of the invite-only setting.

## Invite Codes

### Automatic Admin Code

On startup, the server checks for admin users. If none are found:

1. Invite codes become mandatory for signup.
2. A **single-use admin invite code** is generated (or retrieved from the database).
3. Use this code to create the first admin user.

### Single-Use Codes

All invite codes are single-use by default. Once redeemed, a code cannot be used again.

### Invite-Only Mode

When enabled, every new signup requires a valid invite code — even after admin users exist. This is useful for private or team blogs.