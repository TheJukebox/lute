FROM golang:1.25.5-alpine

COPY ./ /app
WORKDIR /app

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -o server ./cmd/server

EXPOSE 7001
ENTRYPOINT ["/app/server"]
