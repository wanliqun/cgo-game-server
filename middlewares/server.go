package middlewares

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/metrics"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
	"google.golang.org/protobuf/encoding/protojson"
)

type MiddlewareFunc func(server.HandlerFunc) server.HandlerFunc

var (
	protoValidator *protovalidate.Validator

	errPanicCrash   = errors.New("panic crash")
	errAuthRequired = errors.New("authentication required")
)

func init() {
	var err error
	protoValidator, err = protovalidate.New()
	if err != nil {
		panic(errors.WithMessage(err, "failed to initialize proto validator"))
	}
}

// MiddlewareChain builds a server-side middleware chain.
// It takes the main game logic handler as the first argument and chains additional middleware
// functions around it.
func MiddlewareChain(
	handler server.HandlerFunc, middlewares ...MiddlewareFunc) (server.HandlerFunc, error) {
	if len(middlewares) == 0 {
		return nil, errors.New("no middleware provided")
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler, nil
}

// PanicRecover recovers from panic to prevent server crash during message handling.
func PanicRecover(next server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, m *server.Message) (res *server.Message) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			// Print stack info
			debug.PrintStack()

			// Log output context
			sess, _ := server.SessionFromContext(ctx)
			logrus.WithFields(logrus.Fields{
				"session":    sess,
				"panicError": err,
			}).Error("RPC middleware panic")

			// Resp server error
			res = server.NewMessageWithError(errPanicCrash)
		}()

		return next(ctx, m)
	}
}

// MsgValidator validates protobuf messages.
func MsgValidator(next server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, m *server.Message) *server.Message {
		if err := protoValidator.Validate(m); err != nil {
			err = server.NewBadRequestError(err)
			return server.NewMessageWithError(err)
		}

		return next(ctx, m)
	}
}

// Logger logs request, response and handling duration.
func Logger(next server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, m *server.Message) *server.Message {
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			// Skip logging if `debug` level is not enabled.
			return next(ctx, m)
		}

		// Start a timer
		start := time.Now()

		// Log the request
		logrus.WithField("request", m.String()).Debug("Request received")

		// Pass to next handler chain
		resp := next(ctx, m)

		logrus.WithFields(logrus.Fields{
			"response": protojson.Format(resp),
			"elapsed":  time.Since(start),
			"err":      resp.Error,
		}).Debug("Request handled")

		return resp
	}
}

// Authentication middleware.
func Authenticator(s *service.PlayerService) MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, m *server.Message) *server.Message {
			req := m.GetRequest()
			if req.GetLogin() != nil || req.GetInfo() != nil {
				// Non-Auth action required
				return next(ctx, m)
			}

			sess, _ := server.SessionFromContext(ctx)
			player := s.GetBySession(sess.ID)
			if player == nil {
				err := server.NewBadRequestError(errAuthRequired)
				return server.NewMessageWithError(err)
			}

			// Player already logined in
			ctx = service.NewContextFromPlayer(ctx, player)
			return next(ctx, m)
		}
	}
}

// Metrics collects RPC latency and QPS.
func Metrics(next server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, m *server.Message) *server.Message {
		// Start a timer
		start := time.Now()

		// Pass to next handler chain
		resp := next(ctx, m)

		// Rate RPC latency and QPS
		metrics.RPC.Rate(m.Type, resp.Error, start)

		return resp
	}
}
