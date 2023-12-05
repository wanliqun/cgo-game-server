package command

import (
	"context"
	"errors"

	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
	pbprotol "google.golang.org/protobuf/proto"
)

var (
	_ Command = (*LoginCommand)(nil)
	_ Command = (*LogoutCommand)(nil)
)

type Command interface {
	Execute(context.Context) (pbprotol.Message, error)
}

type LoginCommand struct {
	reqeuest      *proto.LoginRequest
	playerService *service.PlayerService
}

func (cmd *LoginCommand) Execute(ctx context.Context) (pbprotol.Message, error) {
	session := ctx.Value(server.CtxKeySession).(*server.Session)
	_, err := cmd.playerService.Login(cmd.reqeuest, session)
	return nil, err
}

type LogoutCommand struct {
	playerService *service.PlayerService
}

func (cmd *LogoutCommand) Execute(ctx context.Context) (pbprotol.Message, error) {
	cv := ctx.Value(service.CtxKeyPlayer)
	if cv == nil {
		return nil, errors.New("not logined yet")
	}

	player := cv.(*service.Player)
	cmd.playerService.Kickoff(player)

	return nil, nil
}
