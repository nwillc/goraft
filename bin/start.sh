#!/usr/bin/env bash

go build ./...

for member in one two three four five; do
  ./goraft -member ${member} &
done
