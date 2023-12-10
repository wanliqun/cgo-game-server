package command

import (
	"context"

	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
	pbproto "google.golang.org/protobuf/proto"
)

var (
	_ Command = (*LoginCommand)(nil)
	_ Command = (*LogoutCommand)(nil)
	_ Command = (*InfoCommand)(nil)
	_ Command = (*GenerateRandomNicknameCommand)(nil)
)

type Command interface {
	Execute(context.Context) (pbproto.Message, error)
}

type LoginCommand struct {
	reqeuest      *proto.LoginRequest
	playerService *service.PlayerService
}

func NewLoginCommand(
	reqeuest *proto.LoginRequest, playerService *service.PlayerService) *LoginCommand {
	return &LoginCommand{
		reqeuest:      reqeuest,
		playerService: playerService,
	}
}

func (cmd *LoginCommand) Execute(ctx context.Context) (pbproto.Message, error) {
	session := ctx.Value(server.CtxKeySession).(*server.Session)
	_, err := cmd.playerService.Login(cmd.reqeuest, session)

	return nil, err
}

type LogoutCommand struct {
	playerService *service.PlayerService
}

func NewLogoutCommand(playerService *service.PlayerService) *LogoutCommand {
	return &LogoutCommand{playerService: playerService}
}

func (cmd *LogoutCommand) Execute(ctx context.Context) (pbproto.Message, error) {
	player, _ := service.PlayerFromContext(ctx)
	cmd.playerService.Kickoff(player)

	return nil, nil
}

type InfoCommand struct {
	axService *service.AuxiliaryService
}

func NewInfoCommand(axService *service.AuxiliaryService) *InfoCommand {
	return &InfoCommand{axService: axService}
}

func (cmd *InfoCommand) Execute(ctx context.Context) (pbproto.Message, error) {
	srvCfg := config.Shared().Server
	srvStat := cmd.axService.CollectServerStatus()
	rateMetrics := cmd.axService.GatherRPCRateMetrics()

	return &proto.InfoResponse{
		ServerName:            srvCfg.Name,
		MaxPlayerCapacity:     int32(srvCfg.MaxPlayerCapacity),
		MaxConnectionCapacity: int32(srvCfg.MaxConnectionCapacity),

		Metrics:        rateMetrics,
		OnlinePlayers:  int32(srvStat.NumOnlinePlayers),
		TcpConnections: int32(srvStat.NumTCPConnections),
		UdpConnections: int32(srvStat.NumUDPConnections),
	}, nil
}

type GenerateRandomNicknameCommand struct {
	request   *proto.GenerateRandomNicknameRequest
	axService *service.AuxiliaryService
}

func NewGenerateRandomNicknameCommand(
	request *proto.GenerateRandomNicknameRequest,
	axService *service.AuxiliaryService) *GenerateRandomNicknameCommand {
	return &GenerateRandomNicknameCommand{
		request:   request,
		axService: axService,
	}
}

func (cmd *GenerateRandomNicknameCommand) Execute(ctx context.Context) (pbproto.Message, error) {
	nickname := cmd.axService.Generate(cmd.request.Sex, cmd.request.Culture)
	return &proto.GenerateRandomNicknameResponse{Nickname: nickname}, nil
}
