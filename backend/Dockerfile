FROM golang:1.25.5-alpine

COPY ./ /app
WORKDIR /app

RUN go mod download \
    && go build -o server ./cmd/server

EXPOSE 7001
ENTRYPOINT ["/app/server"]
