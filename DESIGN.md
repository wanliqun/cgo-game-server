# Design

## Scope and Objective

Develop a high-performance and highly available demo game server capable of handling approximately 10,000 concurrent users. The server should support both TCP and UDP connections and utilize cgo for integration with C/C++ components.

## Performance, Security, and Reliability

- Handle 10,000 concurrent users with minimal latency and high throughput.
- Sustain operation for an extended period under high traffic conditions.
- Implement retry mechanisms to automatically reconnect to the server in case of network interruptions.

## C/S Communication Protocol

The game server utilizes the protobuf serialization protocol to exchange data with clients. This protocol enables efficient and structured communication between the server and connected clients.

### Supported protocols

| Command | Behavior |
|--|--|
|INFO|Retrieves server information, including running status etc.|
|LOGIN|Logs the player into the game server.|
|LOGOUT|Logs the player off the game server.|
|GENERATE_RANDOM_NICKNAME|Generates a random nickname based on the specified sex and culture.|

## Assumptions and Constraints

- Leveraging Go routines to efficiently manage concurrent connections.
- Integrating cgo to seamlessly exchange data between Go and C/C++.

## Acceptance Criteria and Success Metrics

- Maintain service uptime (no crash, no memory leak) for at least one hour under extreme traffic conditions (~10,000 concurrent users).
- Automatically reconnect to the server in the event of abnormal connections.

## High-Level Design

### Main Components and Responsibilities

1. Server

- Application Main
Manages application life cycle, dependency injection, and configuration loading.

- TCP/UDP Server
Establishes connections with clients, including connection establishment, message codec and handling.

- Session Manager
Manages established sessions, periodically performing timeout checks to terminate inactive connections.

- Protocol Encodec/Decodec
Encodes or decodes messages from network packet data for application layer to handle.

- Middleware Pipeline
Decouples the core business logic of the game server from cross-cutting concerns.

- Command
Encapsulates the application's game logic for handling specific requests from the client. 

2. CGO

- C++ Dynamic Link Library
Compile the [C++ Random Name Generator](https://github.com/dasmig/name-generator) into a dynamic `.so` library, which will provide access to the name generation functionality from the Go application.

- Go Stub Struct
Generate stub objects to call methods from the C++ dynamic link library, which will act as bridges between Go and C++, allowing the Go code to interact seamlessly with the C++ library's functions and data structures.

3. Client

- Connection Management 
Establishes TCP/UDP connection to the remote server and handles connection lifecycle.

- Reconnection Mechanism
Implements an exponential backoff retry mechanism to automatically reconnect to the server in case of connection loss or abnormal conditions.

- API Interface
Provides a wrapped API interface that abstracts away the underlying network communication and error handling, allowing users to interact with the server in a straightforward manner.

4. Metrics

- Collects performance metrics, including QPS (queries per second) and response latency, to measure the server's performance and identify potential bottlenecks.

5. Testing 

- Stress testing
 Simulates a high volume of client requests to assess the server's performance and stability under load.

### Additional Considerations

The server can be extended by adding more features for future improvement, such as:

- Security
Implement appropriate security measures to protect the server from unauthorized access and malicious attacks.

- Scalability 
Design the system to be scalable to handle a growing number of users and requests.

- Performance
Optimize the system for performance to minimize latency and maximize throughput.

- Reliability
Implement mechanisms to ensure the reliability of the system, including fault tolerance and disaster recovery.

- Persistence
Persist the counter and set data to a database so that they are not lost when the server is restarted.

- Testing
Employ more comprehensive testing strategy such as unit testing, integration testing and load testing for more robustness and reliability.

- Logging and Monitoring
Implement comprehensive logging and monitoring mechanisms to track system activity, identify potential issues, and gather performance metrics. 

## Detail Design

### Protobuf Definition

```protobuf
enum Sex {
  MALE = 0;
  FEMALE = 1;
}

enum Culture {
  AMERICAN = 0;
  ARGENTINIAN = 1;
  AUSTRALIAN = 2;
  BRAZILIAN = 3;
  BRITISH = 4;
  BULGARIAN = 5;
  CANADIAN = 6; 
  CHINESE = 7;
  DANISH = 8;
  FINNISH = 9;
  FRENCH = 10;
  GERMAN = 11;
  KAZAKH = 12;
  MEXICAN = 13;
  NORWEGIAN = 14;
  POLISH = 15;
  PORTUGUESE = 16;
  RUSSIAN = 17;
  SPANISH = 18;
  SWEDISH = 19;
  TURKISH = 20;
  UKRAINIAN = 21;
}

enum StatusCode {
  OK = 0;
  UNKNOWN_ERROR = 1;
  INVALID_PARAMETER = 2;
  RESOURCE_NOT_FOUND = 3;
  PERMISSION_DENIED = 4;
}

// The message type enum
enum MessageType {
  INFO = 0; // INFO command
  LOGIN = 1; // LOGIN command
  LOGOUT = 2; // LOGOUT command
  GENERATE_RANDOM_NICKNAME = 3; // GENERATE_RANDOM_NICKNAME command
}

// Define the main request and response messages

message LoginRequest {
  optional string username = 1;
  optional string password = 2;
}

message LoginResponse { }

message LogoutRequest {}

message LogoutResponse { }

message InfoRequest {}

message InfoResponse {
  optional string server_name = 1; // server name
  optional int32 num_online_players = 2; // number of online players
  optional int32 max_player_capacity = 3; // max player capacity
  optional int32 max_connection_capacity = 4; // max connection capacity
  map<string, string> metrics = 5; // metrics status as key-value pairs
}

message GenerateRandomNicknameRequest {
  optional Sex sex = 1; // Sex
  optional Culture culture = 2; // Culture
}

message GenerateRandomNicknameResponse {
  optional string nickname = 1; // Nickname
}

// Message for encapsulating different request types
message Request {
  MessageType type = 1;
  oneof body {
    InfoRequest info = 2;
    LoginRequest login = 3;
    LogoutRequest logout = 4;
    GenerateRandomNicknameRequest generate_random_nickname = 5;
  }
}

// Message for conveying response status information
message Status {
  StatusCode code = 1; // Status code
  optional string message = 2; // Descriptive message
}

// Message for encapsulating different response types
message Response {
  MessageType type = 1;
  oneof body {
    Status status = 2;
    InfoResponse info = 3;
    LoginResponse login = 4;
    LogoutResponse logout = 5;
    GenerateRandomNicknameResponse generate_random_nickname = 6;
  }
}

message Message {
  oneof body {
    Request request = 1;
    Response response = 2;
  }
}
```

### Application

```go
type Application struct {
    tcpServer, udpServer *server.Server
}

func (app *Application) Run() {...}
func (app *Application) Close() {...}
```

### Server

```go
// HandlerFunc type
type HandlerFunc func(context.Cotnext, *proto.Request) *proto.Response

type Server struct {
    codec Codec
    listener net.Listener
    msgHandler HandlerFunc
    sessionMgr *SessionManager
}

func NewServer(
    protoType ProtocolType, codec Codec, 
    sessionMgr *SessionManager,  msgHandler HandlerFunc, 
) (*Server, errror) { ... }

func (s *Server) Start() error { ... }
func (s *Server) Stop() error { ... }
```

### Codec

```go
// Codec serializes and deserializes data between protocol message and 
// underlying network package. 
type Codec interface {
    Encode(*proto.Message}) ([]byte, error)
    Decode([]byte) (*proto.Message, error)
}
```

### SessionManager

```go
type Session struct {
    Conn net.conn
    SessionID string
    LastActive time.Time
}

type SessionManager struct {
    sessions map[string]*Session // Session ID => Session Object
}

func (m *SessionManager) Add(s *Session) {...}
func (m *SessionManager) Remove(s *Session) {...}
func (m *SessionManager) CloseAll() { ... }
func (m *SessionManager) Start() { ... }
func (m *SessionManager) Stop() { ... }
```

### Middlewares

```go
type MiddlewareFunc func(HandlerFunc) HandlerFunc
func MiddlewareChain(middlewares ...MiddlewareFunc) HandlerFunc { ... }

func PanicRecover(HandlerFunc) HandlerFunc {...}
func Authenticater(HandlerFunc) HandlerFunc {...}
func Logger(HandlerFunc) HandlerFunc {...}
func Validator(next HandlerFunc) HandlerFunc { ... }
func MetricsGather(next HandlerFunc) HandlerFunc { ... }
func CommandExecutor(next HandlerFunc) HandlerFunc { ... }
```

### Service

```go
type Player struct {
    Username string
    Session *Session
}

type PlayerService struct {
    userPlayers map[string]*Player // username => *Player
    sessionPlayers map[string]*Player // session => *Player
}

// Event handler for session termination event.
func (s *PlayerService) OnSessionTerminatedEvent(e *SessionTerminatedEvent) { ... }

func (s *PlayerService) Add(p *Player) { ... }
func (s *PlayerService) Remove(p *Player) { ... }
func (s *PlayerService) GetByUser(username string) *Player { ... }
func (s *PlayerService) GetBySession(session *Session) *Player { ... }
```

### Command

```go
type CommandFactory struct {}
func (f *CommandFactory) CreateCommand (msg *Message) (any, error){ ... }

type Command interface {
    Execute(ctx context.Context) (any, error)
}
```

### Client

```go
type Client struct {
    codec Codec
    conn net.Conn
    dialer func(endpoint string) (net.Conn, error)
}

func NewClient(protoType ProtocolType, endpoint string, codec Codec) (*Client, error)
func (c *Client) Close() error { ... }

func (c *Client) Info() (*proto.InfoResponse, error)
func (c *Client) Login(username, password string) (*proto.LoginResponse, error)
func (c *Client) Logout() (*proto.LogoutResponse, error)
func (c *Client) GenerateRandomNickname(sex, culture int) (*GenerateRandomNicknameResponse, error)
```

### Auxillary Types

```go
// Config represents the configuration parameters for the game server.
type Config struct {
    // ServerName is the name of the game server.
    ServerName string
    // ServerPassword is the password required to connect to the server.
    ServerPassword string
    // TCPEndpoint is the TCP endpoint on which the server listens.
    TCPEndpoint string
    // UDPEndpoint is the UDP endpoint on which the server listens.
    UDPEndpoint string
    // MaxPlayerCapacity is the maximum number of players that the server
    // can accommodate simultaneously.
    MaxPlayerCapacity int
    // MaxConnectionCapacity is the maximum number of concurrent connections
    // that the server can handle.
    MaxConnectionCapacity int
    // LogLevel is the logging level for the server. 
    // Valid values are "debug", "info", "warning", "error", and "fatal".
    LogLevel string
}
```

### Metrics

- Gauges
	- `concurrent_players`: Measures the number of players currently connected in the game. 
	- `connected_tcp_connections`: Measures the number of established TCP connections to the game server. 
	- `connected_udp_connections`: Measures the number of established UDP connections to the game server. 

- Timers
	- `overall_qps`: Measures the overall query per second (QPS) rate across all commands processed by the server.
	- `overall_latency`: Measures the overall latency distribution for all commands processed by the server.
	- `generate_random_nickname_qps`: Measures the QPS rate specifically for the `GENERATE_RANDOM_NICKNAME` command. 
	- `generate_random_nickname_latency`: Measures the latency distribution for the `GENERATE_RANDOM_NICKNAME` command. 
