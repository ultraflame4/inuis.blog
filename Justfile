set dotenv-load := true
export DATA_DIR := "./data"

_:
    just --list

[parallel]
dev: _dev-server-proxied dev-client

[env("PROXIED_ADDR", "localhost:5173")]
_dev-server-proxied:
    go run .

dev-server:
    go run .

[working-directory("./web")]
dev-client:
    npm run dev
