package server

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/xtaci/kcp-go/v5"
)

type ContextKey string

type HandlerFunc func(context.Context, *proto.Message) *proto.Message

type Server struct {
	listener   net.Listener
	msgHandler HandlerFunc
	sessionMgr *SessionManager
	codec      *proto.Codec

	context context.Context
	cancel  context.CancelFunc
}

func NewServer(
	codec *proto.Codec, sessionMgr *SessionManager, msgHandler HandlerFunc) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		codec:      codec,
		msgHandler: msgHandler,
		sessionMgr: sessionMgr,
		context:    ctx,
		cancel:     cancel,
	}
}

func (s *Server) ListenTCP(endpoint string) (err error) {
	if s.listener != nil {
		return errors.New("server already listened")
	}

	s.listener, err = net.Listen("tcp", endpoint)
	return err
}

func (s *Server) ListenUDP(endpoint string) (err error) {
	if s.listener != nil {
		return errors.New("server already listened")
	}

	s.listener, err = kcp.Listen(endpoint)
	return err
}

func (s *Server) Start() error {
	if s.listener == nil {
		return errors.New("server not listened yet")
	}

	go func() { // Accept new connections.
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				logrus.WithError(err).
					WithField("listenAddr", s.listener.Addr()).
					Debug("Server failed to accept new client connection")
				continue
			}

			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	logger := logrus.WithFields(logrus.Fields{
		"listenAddr": s.listener.Addr(),
		"remoteAddr": conn.RemoteAddr(),
	})
	logger.Debug("Client connection established")

	session := NewSession(conn)
	defer session.Close()

	s.sessionMgr.Add(session)
	defer s.sessionMgr.Remove(session)

	for {
		msg, err := s.codec.Decode(conn)
		if err != nil {
			logger.WithError(err).
				Debug("Server codec failed to decode proto message")
			break
		}

		ctx := context.WithValue(s.context, CtxKeySession, session)
		resp := s.msgHandler(ctx, msg)

		if err := s.codec.Encode(resp, conn); err != nil {
			logger.WithError(err).
				Debug("Server codec failed to encode proto message")
			break
		}

		session.Refresh()
	}

	logger.Debug("Client connection closed")
}

func (s *Server) Stop() error {
	defer func() {
		s.listener = nil

		s.sessionMgr.CloseAll()
		s.sessionMgr.Clear()
	}()

	return s.listener.Close()
}
