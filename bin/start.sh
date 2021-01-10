#!/usr/bin/env bash

set -eou pipefail

(cd pkg ; go build)

for member in one two three four five; do
  ./pkg/goraft -member ${member} &
done
