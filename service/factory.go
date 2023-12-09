package service

import (
	"github.com/wanliqun/cgo-game-server/common"
	"github.com/wanliqun/cgo-game-server/server"
)

type Factory struct {
	Player    *PlayerService
	Auxiliary *AuxiliaryService
}

func NewFactory(
	sessionMgr *server.SessionManager, monickerGenerator common.MonickerGenerator) *Factory {
	return &Factory{
		Player:    NewPlayerService(sessionMgr),
		Auxiliary: NewAuxiliaryService(monickerGenerator),
	}
}
