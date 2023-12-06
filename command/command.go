package command

import (
	"context"
	"errors"

	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/game/common"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
	pbprotol "google.golang.org/protobuf/proto"
)

var (
	_ Command = (*LoginCommand)(nil)
	_ Command = (*LogoutCommand)(nil)
	_ Command = (*InfoCommand)(nil)
	_ Command = (*GenerateRandomNicknameCommand)(nil)
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

type InfoCommand struct {
	infoService *service.InfoService
}

func (cmd *InfoCommand) Execute(ctx context.Context) (pbprotol.Message, error) {
	srvCfg := config.Shared().Server
	srvStat := cmd.infoService.CollectServerStatus()
	rateMetrics := cmd.infoService.GatherRPCRateMetrics()

	resp := &proto.InfoResponse{
		ServerName:            srvCfg.ServerName,
		MaxPlayerCapacity:     int32(srvCfg.MaxPlayerCapacity),
		MaxConnectionCapacity: int32(srvCfg.MaxConnectionCapacity),

		Metrics:        rateMetrics,
		OnlinePlayers:  int32(srvStat.NumOnlinePlayers),
		TcpConnections: int32(srvStat.NumTCPConnections),
		UdpConnections: int32(srvStat.NumUDPConnections),
	}

	return resp, nil
}

type GenerateRandomNicknameCommand struct {
	request           *proto.GenerateRandomNicknameRequest
	monickerGenerator common.MonickerGenerator
}

func (cmd *GenerateRandomNicknameCommand) Execute(ctx context.Context) (pbprotol.Message, error) {
	nickname := cmd.monickerGenerator.Generate(cmd.request.Sex, cmd.request.Culture)
	return &proto.GenerateRandomNicknameResponse{Nickname: nickname}, nil
}
