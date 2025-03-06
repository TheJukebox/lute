#! /bin/sh
cd /usr/src/lute
go mod download
go build -v -o /usr/local/bin/lute ./cmd/server/main.go
lute
