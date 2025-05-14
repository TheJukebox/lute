#!/usr/bin/bash

protoc --go_out=gen/stream/ --go_opt=paths=source_relative \
    --go-grpc_out=gen/stream/ --go-grpc_opt=paths=source_relative \
    api/proto/stream.proto 

protoc --go_out=gen/upload/ --go_opt=paths=source_relative \
    --go-grpc_out=gen/upload/ --go-grpc_opt=paths=source_relative \
    api/proto/upload.proto 

npx buf generate 
