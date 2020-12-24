#!/bin/sh

protoc -I api/raftapi/ api/raftapi/raftapi.proto --go_out=plugins=grpc:api/raftapi/
