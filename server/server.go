package server

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/xtaci/kcp-go/v5"
)

const (
	defaultShutdownTimeout = 5 * time.Second
)

const (
	ServerStatusInitial int32 = iota
	ServerStatusStarted
	ServerStatusStopped
)

var (
	errServerClosed         = errors.New("server closed")
	errProtocolNotSupported = errors.New("only TCP/UDP are supported")
)

type Server struct {
	*ConnectionHandler              // Connection handler
	listener           net.Listener // Net listener
	status             atomic.Int32 // Server status
}

func NewTCPServer(addr string, ch *ConnectionHandler) (srv *Server, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{ConnectionHandler: ch, listener: l}, nil
}

func NewUDPServer(addr string, ch *ConnectionHandler) (srv *Server, err error) {
	l, err := kcp.Listen(addr)
	if err != nil {
		return nil, err
	}

	return &Server{ConnectionHandler: ch, listener: l}, nil
}

// Serve always returns a non-nil error and closes l.
// After closed, the returned error is errServerClosed.
func (srv *Server) Serve() error {
	if srv.status.Load() == ServerStatusStopped {
		return errServerClosed
	}

	if !srv.status.CompareAndSwap(ServerStatusInitial, ServerStatusStarted) {
		return nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"endpoint": srv.listener.Addr(),
		"protocol": srv.listener.Addr().Network(),
	})
	logger.Info("Server listened endpoint started serving")

	defer srv.listener.Close()
	for {
		// TODO: Enforce max connections capacity in case of server overload.
		conn, err := srv.listener.Accept()
		if err == nil {
			go srv.Handle(srv.listener, conn)
			continue
		}

		if srv.status.Load() == ServerStatusStopped {
			return errServerClosed
		}

		logger.WithError(err).Debug("Server failed to accept connection")
		return err
	}
}

func (srv *Server) Close() error {
	if srv.status.Load() == ServerStatusStopped {
		return errServerClosed
	}

	if !srv.status.CompareAndSwap(ServerStatusStarted, ServerStatusStopped) {
		return nil
	}

	// Close listener
	srv.listener.Close()

	// Terminate all connections.
	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return srv.SessManager.TerminateAll(ctx)
}

type ConnectionHandler struct {
	Handler     HandlerFunc     // Connection handler
	SessManager *SessionManager // Session manager
	Codec       *proto.Codec    // Protocol codec
}

func NewConnectionHandler(
	h HandlerFunc, mgr *SessionManager, codec *proto.Codec) *ConnectionHandler {
	return &ConnectionHandler{
		Handler: h, SessManager: mgr, Codec: codec,
	}
}

// handleConnection handles new accepted connection from net listener.
func (ch *ConnectionHandler) Handle(l net.Listener, conn net.Conn) {
	logger := logrus.WithFields(logrus.Fields{
		"protocol":   l.Addr().Network(),
		"listenAddr": l.Addr(),
		"remoteAddr": conn.RemoteAddr(),
	})
	logger.Debug("New connection established")

	session := NewSession(conn)
	ch.SessManager.Add(session)
	defer ch.SessManager.Terminate(session)

	for {
		msg, err := ch.Codec.Decode(conn)
		if err != nil {
			logger.WithError(err).
				Debug("Codec failed to decode proto message")
			break
		}

		ctx := NewContextFromSession(context.Background(), session)
		resp := ch.Handler(ctx, NewMessage(msg))

		if err := ch.Codec.Encode(resp.ProtoMessage(), conn); err != nil {
			logger.WithError(err).
				Debug("Codec failed to encode proto message")
			break
		}

		session.Refresh()
	}

	logger.Debug("Connection termiated")
}
