# Getting Started

This guide walks you through compiling, configuring, and deploying **inuis.blog** from scratch. Think of it as your roadmap to a running instance.

> **Note:** There are currently no Docker images or pre-built binaries. You will compile from source and set up your own deployment flow.

---

## Roadmap

1. [Prerequisites](#prerequisites) — what you need installed
2. [Compile](#compile) — clone and build the binary
3. [Configure](#configure) — database, OIDC, and server settings
4. [Deploy](#deploy) — run the binary, reverse proxy, systemd
5. [First Startup](#first-startup) — invite code and admin account

---

## Prerequisites

| Tool | Description |
|---|---|
| [Go](https://go.dev/dl/) 1.21+ | Compiles the server binary |
| [Node.js](https://nodejs.org/) 18+ | Builds the web frontend (embedded into the binary) |
| [Just](https://github.com/casey/just) | Command runner for build shortcuts |
| Git | Clone the repository |


---

## Compile

#### 1. Clone the repository

```bash
git clone https://github.com/ultraflame4/inuis.blog
cd inuis.blog
```

#### 2. Install frontend dependencies

```bash
cd web
npm install
cd ..
```

#### 3. Build

```bash
just build
```

This produces a single binary at `./bin/app`.

The binary embeds the frontend assets — you only need the single `./bin/app` file to deploy.

---

## Configuration

All settings configured via environment variables. 

### Minimal setup (SQLite — easiest for getting started)

```env
DATABASE_DRIVER="sqlite"
DATABASE_URL="./inuis.db"
```

### PostgreSQL setup

```env
DATABASE_DRIVER="postgres"
DATABASE_URL="host=localhost user=postgres password=yourpassword dbname=inuis port=5432 sslmode=disable"
```

### Identity Provider (required for login)

You need at least one OIDC provider so you can sign in. Example with Google:

```env
OIDC_GOOGLE_URL="https://accounts.google.com"
OIDC_GOOGLE_CLIENT_ID="...your google client id..."
OIDC_GOOGLE_CLIENT_SECRET="...your google client secret..."
```

See [Identity Providers](identity-providers.md) for details on setting up other providers and the PKCE flow.

### Full configuration reference

See [Configuration](configuration.md) for all available variables (`AUTO_INIT`, port, etc.).

---

## Deployment

See [Deployment](deployment.md)

---

## First Startup

1. Start the server. The logs will print a **single-use invite code**.
2. Open `http://localhost:8080` (or your domain) in a browser.
3. Sign in with your OIDC provider. You'll be prompted for the invite code.
4. Enter the code — you are now the **admin user**.

> Once an admin exists, future signup behaviour depends on your [invite-only setting](user-manage.md#invite-codes).

---

## Next Steps

- [Configure additional identity providers](identity-providers.md)
- [Set up multi-server deployment](deployment.md)
- [Explore the API](api.md)