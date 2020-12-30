#!/usr/bin/env bash

for member in one two three four five; do
  go run cmd/client/shutdown/main.go -member ${member} > /dev/null 2>&1 &
done
