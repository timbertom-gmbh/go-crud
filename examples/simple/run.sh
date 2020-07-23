#!/bin/sh

protoc --proto_path=:. --twirp_out=$GOPATH/src --go_out=$GOPATH/src examples/simple/simple.proto

go run cmd/generator/main.go \
  -model=SimpleModel \
  -modelpkg=github.com/timbertom-gmbh/go-crud/examples/simple \
  -rpc=Simple \
  -rpcpkg=github.com/timbertom-gmbh/go-crud/examples/simple \
  -out=github.com/timbertom-gmbh/go-crud/examples/simple > examples/simple/result.go