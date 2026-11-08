curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to ./tools
export PATH="$PATH:./tools"
./tools/just build