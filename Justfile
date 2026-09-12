set dotenv-load := true
export DATA_DIR := "./temp/data"

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


# Build the client, and go server into a single binary
build:
    just client build
    rm -rf ./embed/dist
    mv ./web/build ./embed/dist
    mkdir -p ./bin
    go build -o ./bin/app .


# Builds and run the binary
build-x: build
    ./bin/app
