#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd -P)"

set -eou pipefail

cd "${SCRIPT_DIR}/../pkg"

go build

for member in one two three four five; do
  ./goraft -member ${member} -log-level debug &
done
