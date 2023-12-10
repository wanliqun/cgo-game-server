package command

import (
	"context"
	"errors"

	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
)

var (
	errMsgTypeNotSupported = errors.New("message type not supported")
)

type Executor struct {
	svcFactory *service.Factory
}

func NewExecutor(svcFactory *service.Factory) *Executor {
	return &Executor{svcFactory: svcFactory}
}

func (e *Executor) Execute(ctx context.Context, msg *server.Message) *server.Message {
	var cmd Command
	req := msg.GetRequest()

	switch {
	case req.GetInfo() != nil:
		cmd = NewInfoCommand(e.svcFactory.Auxiliary)
	case req.GetLogin() != nil:
		cmd = NewLoginCommand(req.GetLogin(), e.svcFactory.Player)
	case req.GetLogout() != nil:
		cmd = NewLogoutCommand(e.svcFactory.Player)
	case req.GetGenerateRandomNickname() != nil:
		v := req.GetGenerateRandomNickname()
		cmd = NewGenerateRandomNicknameCommand(v, e.svcFactory.Auxiliary)
	// TODO: extend for more message types support
	default:
		err := server.NewBadRequestError(errMsgTypeNotSupported)
		return server.NewMessageWithError(err)
	}

	pbmsg, err := cmd.Execute(ctx)
	if err != nil || pbmsg == nil {
		return server.NewMessageWithError(err)
	}

	return server.NewMessage(proto.NewResponseMessage(pbmsg))
}
