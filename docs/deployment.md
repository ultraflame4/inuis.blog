# Deployment

## Single Server (Default)

`inuis` was built with a single-server deployment in mind. On every startup it will:

1. Auto-migrate the database schema.
2. Initialize admin users if none exist.

No extra configuration is needed — just set your database variables and run.

## Multi-Server / Edge Deployment

When running multiple backend instances behind a load balancer, you'll want to disable `AUTO_INIT` to prevent race conditions between servers and uncontrolled database migrations.

Set `AUTO_INIT` to `false`:

```env
AUTO_INIT=false
```

### Identity Provider Consistency

To prevent weird behaviour across instances, ensure all `OIDC_*` environment variables are **identical** on every server.
Mismatched provider configs may break the OAuth callback flow.