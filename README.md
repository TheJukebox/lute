# Lute

Lute is a self-hosted streaming service, with a focus on concurrent listening,
implemented in [Go](https://go.dev/).

## Configuration

You can configure the Lute backend using the `lute.config.json` file or command-line arguments.

Command-line arguments take precedent over the JSON configuration, allowing you to override your default
configuration at runtime:

```bash
lute -h
Usage of lute:
  -debug
        Run Lute in debug mode.
  -grpc int
        The port that Lute should use for gRPC requests. (default 50051)
  -host string
        The hostname or address that the Lute backend should listen on. (default "127.0.0.1")
  -http int
        The port that Lute should use for HTTP requests. (default 8080)
  -pg string
        The hostname or address of the PostgreSQL database. (default "127.0.0.1")
  -pg-port int
        The port of the PostgreSQL database. (default 5432)

lute -host 192.168.100.1 -http 80
```

The JSON configuration in `lute.config.json` sets the default values that Lute will startup with:

```json
{
    "lute": {
        "host": "127.0.0.1",
        "grpc": 50051,
        "http": 8080
    },
    "postgres": {
        "host": "127.0.0.1",
        "port": 5432
    },
    "uploads": "uploads/",
    "debug": false
}
```
todo: write up configuration definition

## Contributing

### API

Lute makes use of [gRPC](https://grpc.io/docs/what-is-grpc/core-concepts/) and 
[protocol buffers](https://protobuf.dev/) to implement its API. To contribute
to the API, your development environment will require a protocol buffer compiler
and Go-specific plugins for protocol buffers and gRPC.

You can [install `protoc`, a protocol buffer compiler](https://grpc.io/docs/protoc-installation/), 
from pre-compiled binaries. Check the documentation for specific instructions 
for your OS.

```bash
# Using apt:
apt install -y protobuf-compiler

# Using Homebrew
brew install protobuf
```

`protoc` requires [plugins](https://github.com/protocolbuffers/protobuf-go) 
to compile protocol buffers into the Go-specific gRPC implementation.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### Compiling Protocol Buffers

You can compile protocol buffers in this project with `protoc`. Generally,
`.proto` files will be found in the [`api/proto`](api/proto) directory.

```bash
protoc --go_out=. --go_opt = paths =source_relative --go-grpc_out=. \
    --go-grpc_opt = paths =source_relative api/proto
```