# bytesblog

A tiny, no-nonsense, fast blogging platform written in go.

No trackers, minimal JavaScript and some CSS

This is a minimal, self-hostable blogging platform for developers, writers, and anyone else who wants to share their thoughts, in a fast and simple way.

# Notices
- Users can only login/signup and be authenticated via an external Identity Provider like Google.


# Configuration

All configuration is done via environment variables.

## Adding an Identity Provider
Identity providers are used to authenticate users. 
Any identity provider that supports OpenID Connect can be used.

To add an identity provider, you just need to specify the OIDC config via environment variables.
The server will automatically scan environment variables starting with `OIDC_` to configure OIDC providers.

You need to set the following environment variables to configure an OIDC provider:
- `OIDC_<PROVIDER_NAME>_URL`: URL of the OIDC provider. This is used to get the oidc address.
  - Example: `https://accounts.google.com` will resolve to `https://accounts.google.com/.well-known/openid-configuration`
  - DO NOT PUT the DISCOVERY address! (Use the portion WITHOUT `/.well-known/openid_configuration`)!
- `OIDC_<PROVIDER_NAME>_CLIENT_ID`: Client ID for OIDC provider.
- `OIDC_<PROVIDER_NAME>_CLIENT_SECRET`: Client secret for OIDC provider. If this is not set, the PKCE flow will be used instead.

Example for Google:
```bash
OIDC_GOOGLE_URL="https://accounts.google.com"
OIDC_GOOGLE_CLIENT="...google app client id..."
OIDC_GOOGLE_SECRET="...google app client secret..."
```
