package display

import (
	"github.com/colecrouter/gameboy-go/pkg/memory"
)

type Display interface {
	Init()
	Clock()
	Connect(vram *memory.Device)
}
