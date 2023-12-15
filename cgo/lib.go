package cgo

// #cgo LDFLAGS: -L.. -lnamegen
// #include "./cpp/lib-bridge.h"
// #include <stdio.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"

	"github.com/wanliqun/cgo-game-server/common"
)

func Init(resourcePath string) {
	cstr := C.CString(resourcePath)
	C.LIB_Load(cstr)
	C.free(unsafe.Pointer(cstr))
}

type CGOFakeNameGenerator struct{}

func (g *CGOFakeNameGenerator) Generate(sex common.Gender, cult common.Culture) (str string) {
	cstr := C.LIB_GetName(C.int(sex), C.int(cult))
	if cstr == nil {
		return
	}

	str = C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return str
}
