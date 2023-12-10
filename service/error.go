package service

import (
	"errors"

	"github.com/wanliqun/cgo-game-server/server"
)

const (
	StatusInvalidPassword = iota + 1000
)

var (
	errInvalidPassword = &server.StatusError{
		Code: StatusInvalidPassword,
		Err:  errors.New("invalid password"),
	}
)
