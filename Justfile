set dotenv-load := true
export DATA_DIR := "./data"

_:
    just --list

[parallel]
dev: _dev-server-proxied (client "dev")

[env("PROXIED_ADDR", "localhost:5173")]
_dev-server-proxied:
    go run .

server:
    go run .

[working-directory("./web")]
client *args:
    npm run {{args}}


configure: _configure_client

[working-directory("./web")]
_configure_client:
    npm install


build:
    just client build
    rm -rf ./dist
    mv ./web/build ./dist