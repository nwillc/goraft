#!/usr/bin/env bash

for member in one two three four five; do
  ./goraft -member ${member} &
done
