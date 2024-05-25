# Faraway

```
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge
```

## Word of Wisdom

### Dev & Build

```sh
# Prepare
go install go.uber.org/mock/mockgen@latest
go generate ./...

# Run tests
go test ./...
go test -race ./...

# Run linter
golangci-lint run ./...

# Build server
go build -o server cmd/server/main.go

# Build client
go build -o client cmd/client/main.go

# Run server
CONFIG=conf/server.yaml go run cmd/server/main.go

# Run client
CONFIG=conf/client.yaml go run cmd/client/main.go
```

Check the [documentation](https://golangci-lint.run/welcome/install/#local-installation) if you don't know how to install **golangci-lint**.

### Server
Server expects the following environment variables:
  - **LOG_LEVEL** - Logging level (*error*, *warning*, *info*, *debug*), default is *info*
  - **HOST** - host to listen on
  - **PORT** - TCP port to listen

Configuration file is located at `conf/server.yaml`.

Build & run:
```sh
docker build . -f Dockerfile.server -t faraway-wow-server
docker run -p 9090:9090 faraway-wow-server
```

### Client

## Network protocol

In this case, JSON is chosen for simplicity in the test project. However, it is not very suitable for real production usage with a TCP server because it is not optimal in size, not the fastest in terms of marshalling, dispatching, and additionally, it does not provide information about message length, which is crucial in network data transmission.

For a real production scenario, something like protobuf, [msgpack](https://msgpack.org), or a custom binary Tag-Length-Value protocol would be more appropriate. The choice should depend on the planned client-server interaction, client specificity (for example, mobile clients might require traffic minimization), the project's ecosystem, and available diagnostic tools (it might be necessary, for example, to implement a plugin for Wireshark).

Also, in a real production environment, encryption between the client and server would be required. Since there was no such requirement in the assignment, I omitted this aspect for simplicity, and all interaction occurs as plain text.


