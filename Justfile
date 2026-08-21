set dotenv-load := true
export DATA_DIR := "./data"

_:
    just --list

[parallel]
dev: dev-server dev-client

dev-server:
    go run . 

[working-directory("./client")]
dev-client:
    npm run dev