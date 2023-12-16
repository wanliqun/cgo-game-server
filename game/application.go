package game

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/wanliqun/cgo-game-server/cgo"
	"github.com/wanliqun/cgo-game-server/command"
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/middlewares"
	"github.com/wanliqun/cgo-game-server/proto"
	"github.com/wanliqun/cgo-game-server/rest"
	"github.com/wanliqun/cgo-game-server/server"
	"github.com/wanliqun/cgo-game-server/service"
	"github.com/wanliqun/cgo-game-server/util"
)

type Application struct {
	conf       *config.Config
	sessionMgr *server.SessionManager
	udpServer  *server.Server
	tcpServer  *server.Server
	restServer *rest.Server
}

func NewApplication(configYaml string) (*Application, error) {
	cfg, err := config.NewConfig(configYaml)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load config")
	}

	if err := config.InitLogger(&cfg.Log); err != nil {
		return nil, errors.WithMessage(err, "failed to init logger")
	}

	var monickerGenerator common.MonickerGenerator
	if cfg.CGO.Enabled {
		cgo.Init(cfg.CGO.ResourceDir)
		monickerGenerator = &cgo.CGOFakeNameGenerator{}
	} else {
		monickerGenerator = &common.GoFakerNameGenerator{}
	}

	sessionMgr := server.NewSessionManager()

	svcFactory := service.NewFactory(cfg, sessionMgr, monickerGenerator)
	cmdExecutor := command.NewExecutor(svcFactory)

	msgHandler, err := middlewares.MiddlewareChain(
		cmdExecutor.Execute,
		middlewares.PanicRecover,
		middlewares.Logger,
		middlewares.MsgValidator,
		middlewares.Authenticator(svcFactory.Player),
		middlewares.Metrics,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to build middleware chain")
	}

	codec := proto.NewCodec()
	connHandler := server.NewConnectionHandler(msgHandler, sessionMgr, codec)

	udpServer, err := server.NewUDPServer(cfg.Server.UDPEndpoint, connHandler)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to new UDP server")
	}

	tcpServer, err := server.NewTCPServer(cfg.Server.TCPEndpoint, connHandler)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to new TCP server")
	}

	restServer, err := rest.NewServer(cfg.Server.HTTPEndpoint, svcFactory)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to new RESTful server")
	}

	return &Application{
		conf:       cfg,
		sessionMgr: sessionMgr,
		udpServer:  udpServer,
		tcpServer:  tcpServer,
		restServer: restServer,
	}, nil
}

func (app *Application) Run() {
	go app.sessionMgr.Start()
	go app.udpServer.Serve()
	go app.tcpServer.Serve()
	go app.restServer.Serve()

	util.GracefulShutdown(&sync.WaitGroup{}, app.Close)
}

func (app *Application) Close() {
	app.sessionMgr.Stop()
	app.udpServer.Close()
	app.tcpServer.Close()
	app.restServer.Close()
}
