# Configuration

All configuration is done via **environment variables**.

---

## Database

| Variable | Default | Description |
|---|---|---|
| `DATABASE_DRIVER` | `sqlite` | Database driver. Supported: `sqlite`, `postgres`, `mysql`. |
| `DATABASE_URL` | `data.db` | Connection string. For SQLite this is a path relative to `DATA_DIR`. |

### Examples

**SQLite (default):**

```env
DATABASE_DRIVER="sqlite"
DATABASE_URL="inuis.db"
```

**PostgreSQL:**

```env
DATABASE_DRIVER="postgres"
DATABASE_URL="host=localhost user=postgres password=... dbname=inuis port=5432 sslmode=disable"
```

**MySQL:**

```env
DATABASE_DRIVER="mysql"
DATABASE_URL="user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

See [GORM docs](https://gorm.io/docs/connecting_to_the_database.html) for full DSN format details.

---

## Server

| Variable | Default | Description |
|---|---|---|
| `HOST` | `localhost` | Address the server binds to. |
| `PORT` | `8080` | Port the server listens on. |
| `PROXIED_ADDR` | *(none)* | Public-facing address when behind a reverse proxy. Used for OAuth redirect URLs so the browser is sent to the proxy, not the backend. Example: `localhost:5173` (Vite dev proxy). |
| `DATA_DIR` | `/data` | Directory for runtime data. SQLite databases are stored relative to this path. Created automatically if missing. |
| `AUTO_INIT` | `true` | Auto-migrates the database and initializes admin users on startup. Set to `false` in multi-server deployments. |

---

## Sessions

| Variable | Default | Description |
|---|---|---|
| `MAX_SESSION_AGE` | `3600` | Maximum session lifetime in seconds (1 hour). If the value is unparseable, falls back to 4 hours (14400s). |

---

## Branding

| Variable | Default | Description |
|---|---|---|
| `BRAND_NAME` | `inuis.blog` | Brand name displayed in the application UI. |
| `COPYRIGHT_NAME` | *(same as `BRAND_NAME`)* | Name used in copyright notices. Defaults to `BRAND_NAME` if not set. |

---

## Signups & Invites

| Variable | Default | Description |
|---|---|---|
| `ALLOW_SIGNUP` | `true` | When `false`, disables all new account creation. Parsed with `strconv.ParseBool` — if parsing fails, signups remain enabled as a safety measure. |
| `REQUIRE_INVITE_CODE` | `false` | When `true`, all new signups require a valid invite code. Parsed with `strconv.ParseBool` — if parsing fails, invite codes are **required** as a safety measure. |

> **Note:** If no admin users exist, invite codes are always required regardless of `REQUIRE_INVITE_CODE`.

---

## Identity Providers

Identity providers are configured via environment variables prefixed with `OIDC_`. See [Identity Providers](identity-providers.md) for full details.

---

## Robots.txt

There are currently no plans to make `robots.txt` configurable through the application. Use a reverse proxy (e.g. nginx, Caddy) or a CDN (e.g. Cloudflare) to control crawling behaviour.