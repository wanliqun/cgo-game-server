package service

import (
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/config"
	"github.com/wanliqun/cgo-game-server/server"
)

type Factory struct {
	Player    *PlayerService
	Auxiliary *AuxiliaryService
}

func NewFactory(
	conf *config.Config,
	sessionMgr *server.SessionManager,
	monickerGenerator common.MonickerGenerator) *Factory {

	playerSvc := NewPlayerService(conf, sessionMgr)
	auxSvc := NewAuxiliaryService(conf, monickerGenerator, playerSvc, sessionMgr)
	return &Factory{Player: playerSvc, Auxiliary: auxSvc}
}
