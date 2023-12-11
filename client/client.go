package client

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/encoding/prototext"
	pbproto "google.golang.org/protobuf/proto"
)

const (
	defaultReconnectInterval = 1 * time.Second
	defaultSendBufferSize    = 100
)

type dialer func() (net.Conn, error)
type dialerFactory func(addr string) dialer

func makeTCPDialer(addr string) dialer {
	return func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
}

func makeUDPDialer(addr string) dialer {
	return func() (net.Conn, error) {
		return kcp.Dial(addr)
	}
}

// Client interacts with the game server, including establishing connection
// to server, reading data from server and writing data to server etc.
type Client struct {
	dialer dialer
	codec  proto.Codec
	conn   atomic.Value

	ctx    context.Context
	cancel context.CancelFunc

	mu          sync.Mutex
	callbacks   []OnMessageCallback
	recovering  atomic.Bool
	reconnectCh chan struct{}
	requestCh   chan *proto.Message
}

func NewTCPClient(addr string) *Client {
	return newClient(makeTCPDialer(addr))
}

func NewUDPClient(addr string) *Client {
	return newClient(makeUDPDialer(addr))
}

func newClient(dial dialer) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		dialer:      dial,
		ctx:         ctx,
		cancel:      cancel,
		requestCh:   make(chan *proto.Message, defaultSendBufferSize),
		reconnectCh: make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	conn, err := c.dialer()
	if err != nil {
		return err
	}

	go c.handleConnection(conn)
	return nil
}

func (c *Client) Close() {
	if conn, ok := c.conn.Load().(net.Conn); ok {
		conn.Close()
	}

	c.cancel()
}

func (c *Client) handleConnection(conn net.Conn) {
	c.conn.Store(conn)
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()

	go c.read(ctx, conn)
	go c.write(ctx, conn)

	for { // Start failure recover loop
		select {
		case <-ctx.Done():
			return
		case <-c.reconnectCh:
			if conn, err := c.reconnect(ctx); err == nil {
				c.recovering.Store(false)
				go c.handleConnection(conn)
			}
		}
	}
}

func (c *Client) read(ctx context.Context, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := c.codec.Decode(conn)
		if err == nil {
			logrus.WithFields(logrus.Fields{
				"serverAddr": conn.RemoteAddr(),
				"protocol":   conn.RemoteAddr().Network(),
			}).WithField("message", prototext.Format(msg)).
				Debug("Client read new proto message from server")

			c.notifyOnMessage(msg)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"serverAddr": conn.RemoteAddr(),
			"protocol":   conn.RemoteAddr().Network(),
		}).WithError(err).Debug("Client failed to read from server")

		c.failureRecover()
		return
	}
}

func (c *Client) write(ctx context.Context, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.requestCh:
			err := c.codec.Encode(msg, conn)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"serverAddr": conn.RemoteAddr(),
					"protocol":   conn.RemoteAddr().Network(),
				}).WithError(err).Debug("Client failed to write to server")

				c.failureRecover()
				return
			}
		}
	}
}

func (c *Client) reconnect(ctx context.Context) (conn net.Conn, err error) {
	timer := time.NewTimer(defaultReconnectInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			conn, err = c.dialer()
			if err == nil {
				return
			}
		}

		timer.Reset(defaultReconnectInterval)
	}
}

func (c *Client) failureRecover() {
	if !c.recovering.CompareAndSwap(false, true) {
		// Already in recovering
		return
	}

	// Close the old connection
	if conn, ok := c.conn.Load().(net.Conn); ok {
		conn.Close()
	}

	// Try to reconnect the server
	if len(c.reconnectCh) == 0 {
		c.reconnectCh <- struct{}{}
	}
}

func (c *Client) send(m pbproto.Message) error {
	msg, err := proto.NewRequestMessage(m)
	if err != nil {
		return err
	}

	select {
	case c.requestCh <- msg:
		return nil
	default:
		return errors.New("send buffer is full")
	}
}

type OnMessageCallback func(msg *proto.Message)

func (c *Client) notifyOnMessage(msg *proto.Message) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.callbacks {
		c.callbacks[i](msg)
	}
}

func (c *Client) OnMessage(cb OnMessageCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.callbacks = append(c.callbacks, cb)
}

func (c *Client) Info() error {
	return c.send(&proto.InfoRequest{})
}

func (c *Client) Login(username, password string) error {
	return c.send(&proto.LoginRequest{
		Username: username,
		Password: password,
	})
}

func (c *Client) Logout() error {
	return c.send(&proto.LogoutRequest{})
}

func (c *Client) GenerateRandomNickname(sex, culture int32) error {
	return c.send(&proto.GenerateRandomNicknameRequest{
		Sex: sex, Culture: culture,
	})
}
