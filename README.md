# cgo-game-server

 This is a demo project showcasing a high-performance TCP/UDP game server built with Go and cgo. It includes a server, simulator, and stress testing tool for development and evaluation. It also provides a RESTful service that exposes two RPC methods to view the server status and metrics.

## Quick Start

### Clone the Repository

```bash
git clone --recurse-submodules https://github.com/wanliqun/cgo-game-server.git
```

This downloads the required submodule ( [C++ Random Name Generator](https://github.com/dasmig/name-generator.git) ) automatically.

### Build the C++ Library (cgo)

**Requirements:**
- automake
- clang++

```bash
make cgo
```

This generates a `.so` dynamic link library used by cgo.

### Run the Application:

Install Go (v1.21+) and edit the configuration file (config/config.yml) for customization, then run:

```bash
go run main.go --help
```

This displays usage instructions and available commands.

- Start Server (TCP, UDP, RESTful):
```bash
go run main.go server
```

- Start Simulator (client for debugging):
```bash
go run main.go simulator
```

- Stress Testing:

Establish Mass Connections and Measure TPS:
```bash
go run main.go loadrunner -w 100 -r 100 -d 1h
```

> -w: Number of concurrent workers (go-routines)
> -r: Robots (clients) allocated per worker
> -d: Duration of the test


### Metrics

`cgo-game-server` starts a RESTful service which provides two RPC methods to view the server status and metrics. eg., You can open the following URL in a web browser to check the server running statistics if you are running as predefined configurations: 
- Server status: http://127.0.0.1:8787/status 
- RPC metrics: http://127.0.0.1:8787/metrics

## Acknowledgement

The following projects are referenced:

- [Handling 1M TCP connections in Go](https://github.com/smallnest/1m-go-tcp-server)
- [Handling 1M websockets connections in Go](https://github.com/eranyanay/1m-go-websockets)