package server

import (
	"context"
	"runtime/debug"

	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wanliqun/cgo-game-server/proto"
	pbprotol "google.golang.org/protobuf/proto"
)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

var (
	protoValidator *protovalidate.Validator

	// Build builtin core middleware chain.
	coreMiddlewares = []MiddlewareFunc{panicRecover, msgValidator}
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
func MiddlewareChain(handler HandlerFunc, middlewares ...MiddlewareFunc) (HandlerFunc, error) {
	if len(middlewares) == 0 {
		return nil, errors.New("no middleware provided")
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	// Core middlewares must be called at first.
	for i := len(coreMiddlewares) - 1; i >= 0; i-- {
		handler = coreMiddlewares[i](handler)
	}

	return handler, nil
}

// panicRecover recovers from panic to prevent server crash during message handling.
func panicRecover(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, m pbprotol.Message) (res pbprotol.Message) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			// Print stack info
			debug.PrintStack()

			// Log output context
			sess, _ := SessionFromContext(ctx)
			logrus.WithFields(logrus.Fields{
				"session":    sess,
				"panicError": err,
			}).Error("RPC middleware panic")

			// Resp server error
			res = &proto.Status{
				Code:    proto.StatusCode_UNKNOWN_ERROR,
				Message: "panic crash",
			}
		}()

		return next(ctx, m)
	}
}

// msgValidator validates protobuf messages.
func msgValidator(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, m pbprotol.Message) pbprotol.Message {
		msg := m.(*proto.Message)
		if err := protoValidator.Validate(msg); err != nil {
			return &proto.Status{
				Code:    proto.StatusCode_INVALID_PARAMETER,
				Message: err.Error(),
			}
		}

		return next(ctx, m)
	}
}
