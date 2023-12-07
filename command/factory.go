package command

import (
	"errors"

	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/service"
	pbprotol "google.golang.org/protobuf/proto"
)

type Factory struct {
	svcFactory *service.Factory
}

func (f *Factory) CreateCommand(msg pbprotol.Message) (Command, error) {
	switch v := msg.(type) {
	case *proto.LoginRequest:
		return NewLoginCommand(v, f.svcFactory.Player), nil
	case *proto.LogoutRequest:
		return NewLogoutCommand(f.svcFactory.Player), nil
	case *proto.InfoRequest:
		return NewInfoCommand(f.svcFactory.Auxiliary), nil
	case *proto.GenerateRandomNicknameRequest:
		return NewGenerateRandomNicknameCommand(v, f.svcFactory.Auxiliary), nil
	default:
		return nil, errors.New("unknown command request")
	}
}
