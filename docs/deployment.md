# Deployment


## Docker
Note that the docker examples are untested.
They are there for reference only.

### Dockerfile

Multi-stage build: `just build` handles the frontend and Go compilation, then the binary is copied into a minimal runtime image.

```dockerfile
# Stage 1: Build (Go + Node + just)
FROM golang:1.27.3-trixie AS builder

ENV DEBIAN_FRONTEND=noninteractive

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    nodejs \
    npm \
    just \
    git \
    build-essential \
    pkg-config \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN just build

# Stage 2: Runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /src/bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

### docker-compose.yml

Example docker compose. Edit before using.
```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - # YOUR CONFIG ENV VARS. See configuration.md

    restart: unless-stopped

```

---

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