# Lute

Lute is a self-hosted streaming service, with a focus on concurrent listening,
implemented in [Go](https://go.dev/).

## Contributing

### API

Lute makes use of [gRPC](https://grpc.io/docs/what-is-grpc/core-concepts/) and 
[protocol buffers](https://protobuf.dev/) to implement its API. To contribute
to the API, your development environment will require the Go implementations of
protocol buffers and a protocol buffer compiler, specific to Go.

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