#!/usr/bin/env bash

tempfile=$(mktemp /tmp/rollcall.XXXXXX)
trap "rm -rf ${tempfile}" EXIT

function pingMember() {
    local member=$1
    go run cmd/client/ping/main.go -member ${member} 2>&1 | grep role
    echo >> "${tempfile}"
}

for member in one two three four five; do
  pingMember $member &
done

while [ "$(wc -l "${tempfile}" 2> /dev/null | awk '{ print $1 }')" != "5" ]; do
    sleep 0.1
done

