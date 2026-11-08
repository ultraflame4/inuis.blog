# Configuration

All configuration is done via **environment variables**.

## Database

| Variable | Description |
|---|---|
| `DATABASE_DRIVER` | Database driver (e.g. `postgres`, `sqlite`) |
| `DATABASE_URL` | Database connection string (DSN). See [GORM docs](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL) for format details. |

### Example (Postgres)

```env
DATABASE_DRIVER="postgres"
DATABASE_URL="host=localhost user=postgres password=... dbname=inuis port=5432 sslmode=disable"
```

## Server Initialization

| Variable | Default | Description |
|---|---|---|
| `AUTO_INIT` | `true` | When `true`, the server auto-migrates the database and initializes admin users on startup. Set to `false` in multi-server deployments. |

## Identity Providers

Identity providers are configured via environment variables prefixed with `OIDC_`. See [Identity Providers](identity-providers.md) for full details.

## Robots.txt

There are currently no plans to make `robots.txt` configurable through the application. Use a reverse proxy (e.g. nginx, Caddy) or a CDN (e.g. Cloudflare) to control crawling behaviour.