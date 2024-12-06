FROM golang:1.23.4-bookworm

# Install deps
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
        protobuf-compiler \
        ffmpeg \
        python3 \ 
        python3.11-venv \
        nodejs \
        npm

# install protoc plugins 
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    npm install -g protoc-gen-js && \
    wget https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-linux-x86_64 && \
    mv protoc-gen-grpc-web-1.5.0-linux-x86_64 /usr/local/bin/protoc-gen-grpc-web && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web
