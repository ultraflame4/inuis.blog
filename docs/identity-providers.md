# Identity Providers

Users authenticate exclusively through external Identity Providers via OpenID Connect (OIDC). The server scans environment variables starting with `OIDC_` to discover and configure providers automatically.

## Adding a Provider

For each provider, set three environment variables:

| Variable | Required | Description |
|---|---|---|
| `OIDC_<NAME>_URL` | Yes | Base URL of the OIDC provider. The server appends `/.well-known/openid-configuration` automatically. |
| `OIDC_<NAME>_CLIENT_ID` | Yes | OAuth client ID from the provider. |
| `OIDC_<NAME>_CLIENT_SECRET` | No | OAuth client secret. If omitted, the **PKCE flow** is used instead. |

> [!Important]
> Use the base URL only — do **not** include `/.well-known/openid-configuration` in the URL.

`<NAME>` will also act as the id for the provider (e.g. `GOOGLE`, `GITHUB`). Users will log in at `/login/<NAME>` (case-insensitive).

### Example: Google

```bash
OIDC_GOOGLE_URL="https://accounts.google.com"
OIDC_GOOGLE_CLIENT_ID="...your google client id..."
OIDC_GOOGLE_CLIENT_SECRET="...your google client secret..."
```

Users can then log in at `/login/google` or `/login/GOOGLE`.
