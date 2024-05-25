# Faraway

```
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge
```

## Drawbacks and solutions

### Network protocol

In this case, JSON is chosen for simplicity in the test project. The existing JSON protocol can be enhanced by adding a message tag and incorporating static error codes. However, it is not very suitable for real production usage with a TCP server because it is not optimal in size, not the fastest in terms of marshalling, dispatching, and additionally, it does not provide information about message length, which is crucial in network data transmission.

For a real production scenario, something like protobuf, [msgpack](https://msgpack.org), or a custom binary Tag-Length-Value protocol would be more appropriate. The choice should depend on the planned client-server interaction, client specificity (for example, mobile clients might require traffic minimization), the project's ecosystem, and available diagnostic tools (it might be necessary, for example, to implement a plugin for Wireshark).

Also, in a real production environment, encryption between the client and server would be required. Since there was no such requirement in the assignment, I omitted this aspect for simplicity, and all interaction occurs as plain text.

### Synchronous client and server

To avoid overengineering, both the server and client were implemented to be synchronous (a specific sequence of messages is expected). With an asynchronous client and server, it would have required complicating the logic by adding something like an FSM and a message dispatcher.

### Proof of work algorithm

The Proof of Work algorithm is a somewhat like interpretation of [Hashcash](https://en.wikipedia.org/wiki/Hashcash), which is used in anti-spam filters, DDoS protectors and in the Bitcoin network.

It can be described by the following pseudocode:
```
number_of_low_bits = rate * factor
mask = (1 << number_of_low_bits) - 1
result = (sha256(challenge + solution) & mask) == 0
```

Obviously, for real protection, the complexity of such an algorithm will not be sufficient, and a more complex challenge should be introduced. As an option, you can also consider another cryptographic function, such as [Argon2](https://en.wikipedia.org/wiki/Argon2), or some other [KDF](https://en.wikipedia.org/wiki/Key_derivation_function), since they are typically designed to be resource-intensive for brute-force attacks.

### Network interaction
There is a wide scope for improving the network code - for example, by implementing reconnection after a disconnect or by handling large amounts of data that do not fit into the buffer all at once.

### Logging

For simplicity, a global [zap](https://github.com/uber-go/zap) logger was used, but it's not the best solution. It's better to pass the logger explicitly or through the context.

However, the logging approach heavily depends on the service, the number of microservices and their distribution, the load of service, criticality of logs, and the service deployment environment.

## Word of Wisdom

## Run

This command builds everything you need and runs server and 10 clients:
```sh
docker-compose build && docker-compose up
```

Logs:
```
docker-compose logs server
docker-compose logs client
```


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

# Store deps to vendor 
go mod tidy
go mod vendor

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
docker run -e HOST=0.0.0.0 -e PORT=9090 -p 9090:9090 faraway-wow-server
```

### Client
Client expects the following environment variables:
  - **LOG_LEVEL** - Logging level (*error*, *warning*, *info*, *debug*), default is *info*
  - **HOST** - host to connect to
  - **PORT** - TCP port to connect to

Configuration file is located at `conf/client.yaml`.

Build & run:
```sh
docker build . -f Dockerfile.client -t faraway-wow-client
docker run -e HOST=<server listen addr> -e PORT=9090 faraway-wow-client
```
