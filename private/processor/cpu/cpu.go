package cpu

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

type CPU interface {
	Flags() *flags.Flags
	Registers() *registers.Registers
	Read(uint16) uint8
	Write(uint16, uint8)
	Halt()
	Stop()
	EI()
	EIWithDelay()
	DI()
	PrefixCB()
}
