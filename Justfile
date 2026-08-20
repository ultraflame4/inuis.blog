set dotenv-load := true
export DATA_DIR := "./data"

_:
    just --list

dev:
    go run .
