FROM golang:1.23.3-bookworm

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
