package cgo_test

// #cgo LDFLAGS: -L.. -L. -lnamegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wanliqun/cgo-game-server/cgo"
	"github.com/wanliqun/cgo-game-server/common"
)

const (
	reourceDir = "../resources"
)

func TestCGOMonickerGenerator(t *testing.T) {
	cgo.Init(reourceDir)

	g := cgo.CGOFakeNameGenerator{}
	name := g.Generate(common.Male, common.CHINESE)
	assert.NotEmpty(t, name, "generated nickname shouldn't be empty")
}
