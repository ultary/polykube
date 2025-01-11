#!/bin/sh

# Protocol Buffer Compiler Installation
#  * https://grpc.io/docs/protoc-installation/
#
# ==== macos ===
# brew update
# brew install protobuf
# brew upgrade protobuf
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.2
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
# go install github.com/wasilibs/go-protoc-gen-grpc/cmd/protoc-gen-grpc_python@v1.69.0

protoc \
  --proto_path ./kluster \
  --go_opt paths=source_relative \
  --go_out ./kluster \
  --go-grpc_opt paths=source_relative \
  --go-grpc_out ./kluster \
  --python_out ./console \
  --pyi_out ./console \
  --grpc_python_out ./console \
  api/grpc/v1/kluster.proto \
  api/grpc/v1/system.proto
