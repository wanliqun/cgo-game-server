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
|GENERATE_RANDOM_NICKNAME|Generates a random nickname based on specified gender and culture.|

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
Implements a retry mechanism to automatically reconnect to the server in case of connection loss or abnormal conditions.

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
Persist data to a database so that they are not lost when the server is restarted.

- Testing
Employ more comprehensive testing strategy such as unit testing, integration testing and load testing for more robustness and reliability.

- Logging and Monitoring
Implement comprehensive logging and monitoring mechanisms to track system activity, identify potential issues, and gather performance metrics. 

## Detail Design

### Protobuf Definition

```protobuf

// The message type enum
enum MessageType {
  INFO = 0; // INFO command
  LOGIN = 1; // LOGIN command
  LOGOUT = 2; // LOGOUT command
  GENERATE_RANDOM_NICKNAME = 3; // GENERATE_RANDOM_NICKNAME command
}

// Login in
message LoginRequest {
  string username = 1;
  string password = 2;
}
message LoginResponse { }

// Log out
message LogoutRequest {}
message LogoutResponse { }

// Info
message InfoRequest {}
message InfoResponse {
  string server_name = 1; // server name
  string uptime = 2; // server uptime
  int32 online_players = 3; // number of online players
  int32 total_connections = 4; // total number of network connections
  map<string, string> metrics = 5; // overall metrics as key-value pairs
}

// Generate random nickname
message GenerateRandomNicknameRequest {
  int32 sex = 1; // Sex
  int32 culture = 2; // Culture
}
message GenerateRandomNicknameResponse {
  string nickname = 1; // Nickname
}

// Message for conveying response status information
message Status {
  int32 code = 1; // Status code
  string message = 2; // Descriptive message
}
```

### Application

```go
type Application struct {
  sessionMgr *server.SessionManager
  udpServer  *server.Server
  tcpServer  *server.Server
  restServer *rest.Server
}

func (app *Application) Run() {...}
func (app *Application) Close() {...}
```

### Server

```go
type Server struct {
  *ConnectionHandler              // Connection handler
  listener           net.Listener // Net listener
  status             atomic.Int32 // Server status
}

func NewTCPServer(addr string, ch *ConnectionHandler) (srv *Server, err error) {...}
func NewUDPServer(addr string, ch *ConnectionHandler) (srv *Server, err error) {...}
func (srv *Server) Serve() error {...}
func (srv *Server) Close() error {...}
```

### Codec

```go
type Codec struct {}
func (c *Codec) Encode(msg *Message, w io.Writer) error {...}
func (c *Codec) Decode(r io.Reader) (*Message, error){...}
```

### SessionManager

```go
type Session struct {
  ID         string   // Session ID
  Conn       net.Conn // Underlying network connection
  lastActive int64    // Last active timestamp
}

func (s *Session) Refresh() { ... }
func (s *Session) Close() { ... }

type SessionManager struct {
  sessions map[string]*Session // Session ID => Session Object
}

func (m *SessionManager) Add(s *Session) {...}
func (m *SessionManager) Count() int {...}
func (m *SessionManager) Terminate(sess *Session) error {...}
func (m *SessionManager) TerminateAll(ctx context.Context) error { ... }
func (m *SessionManager) Start() { ... }
func (m *SessionManager) Stop() { ... }
```

### Middlewares

```go
type MiddlewareFunc func(server.HandlerFunc) server.HandlerFunc
func MiddlewareChain(handler server.HandlerFunc, middlewares ...MiddlewareFunc) (server.HandlerFunc, error) {...}

func PanicRecover(next server.HandlerFunc) server.HandlerFunc {...}
func MsgValidator(next server.HandlerFunc) server.HandlerFunc {...}
func Logger(next server.HandlerFunc) server.HandlerFunc {...}
func Authenticator(s *service.PlayerService) MiddlewareFunc { ... }
func Metrics(next server.HandlerFunc) server.HandlerFunc  { ... }
```

### Service

```go
type Player struct {
  Username string
  Session *Session
}

type PlayerService struct {
  mu          sync.Mutex
  config      *config.Config
  usrPlayers  map[string]*Player // username=>Player
  sessPlayers map[string]*Player // session=>Player
  sessionMgr  *server.SessionManager
}

func (s *PlayerService) Add(p *Player) { ... }
func (s *PlayerService) Kickoff(p *Player) {...}
func (s *PlayerService) GetByUser(username string) *Player { ... }
func (s *PlayerService) GetBySession(sessionID string) *Player { ... }
func (s *PlayerService) Login(req *proto.LoginRequest, session *server.Session) (*Player, error) {...}
func (s *PlayerService) OnSessionTerminatedEvent(e *SessionTerminatedEvent) { ... }
```

### Command

```go
type Command interface {
  Execute(context.Context) (pbproto.Message, error)
}

type Executor struct {
  svcFactory *service.Factory
}

func (e *Executor) Execute(ctx context.Context, msg *server.Message) *server.Message { ... }
```

### Client

```go
type Client struct {
  dialer func() (net.Conn, error)
  codec  proto.Codec
  conn  net.Conn
  ...
}

func NewTCPClient(addr string) *Client {...}
func NewUDPClient(addr string) *Client {...}

func (c *Client) Connect() error {...}
func (c *Client) Close() {...}

type OnMessageCallback func(msg *proto.Message)
func (c *Client) OnMessage(cb OnMessageCallback) {...}

func (c *Client) Info() error {...}
func (c *Client) Login(username, password string) error {...}
func (c *Client) Logout() error {...}
func (c *Client) GenerateRandomNickname(sex, culture int32) error {...}
```

### Auxillary Types

```go
type LogConfig struct {
  Level      string `default:"info"`
  ForceColor bool   `default:"true"`
}

type ServerConfig struct {
  Name                  string `default:"cgo_game_server"`
  Password              string `default:"helloworld"`
  TCPEndpoint           string `default:":8765"`
  UDPEndpoint           string `default:":8765"`
  HTTPEndpoint          string `default:":8787"`
  MaxPlayerCapacity     int    `default:"10000"`
  MaxConnectionCapacity int    `default:"15000"`
}

type CGOConfig struct {
  Enabled     bool   `default:":false"`
  ResourceDir string `default:"./resources"`
}

type Config struct {
  Log    LogConfig
  Server ServerConfig
  CGO    CGOConfig
}
```

### Metrics

- `rpc.rate.overall`
  - Measures the overall query per second (QPS) rate across all commands processed by the server.
  - Measures the overall latency distribution for all commands processed by the server.
- `rpc.rate.${command}.[success|error]`
  - Measures the QPS rate specifically for the RPC command. 
  - Measures the latency distribution for the RPC command. 