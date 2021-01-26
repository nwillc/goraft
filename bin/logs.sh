#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd -P)"

cd "${SCRIPT_DIR}/../pkg" || exit

for member in one two three four five; do
  go run cmd/client/logs/main.go -member ${member} -config-file ../config.json
done


