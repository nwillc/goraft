#!/bin/sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" > /dev/null 2>&1 && pwd -P)"

cd "${SCRIPT_DIR}/.." || exit

go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
asdf reshim golang

#protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#    api/raftapi/raftapi.proto

protoc --proto_path=api/proto/v1 --go_out=pkg/raftapi --go-grpc_out=pkg/raftapi raftapi.proto
