#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd -P)"

cd "${SCRIPT_DIR}/../pkg" || exit

go run cmd/client/append/main.go -member $1 -config-file ../config.json -value $2
