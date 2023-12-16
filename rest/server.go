package rest

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/service"
)

type Server struct {
	*http.Server
	listener net.Listener // Net listener
}

func NewServer(endpoint string, svcFactory *service.Factory) (*Server, error) {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: ln,
		Server: &http.Server{
			Addr:        endpoint,
			ReadTimeout: 1 * time.Minute,
			Handler:     newRouter(svcFactory),
		},
	}, nil
}

func (s *Server) Serve() error {
	logrus.WithFields(logrus.Fields{
		"endpoint": s.listener.Addr(),
		"protocol": "http",
	}).Info("Server listened endpoint started serving")

	err := s.Server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	s.listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
