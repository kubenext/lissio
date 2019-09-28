#!/bin/sh
# generate golang for protobuf

protoc -I$GOPATH/src/github.com/kubenext/lissio/vendor -I$GOPATH/src/github.com/kubenext/lissio -I. --go_out=plugins=grpc:. dashboard.proto
