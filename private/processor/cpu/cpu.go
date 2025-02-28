package cpu

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

type CPU interface {
	Flags() *flags.Flags
	GetImmediate8() uint8
	GetImmediate16() uint16
	Registers() *registers.Registers
	Read(uint16) uint8
	Write(uint16, uint8)
	Read16(uint16) (uint8, uint8)
	Write16(uint16, uint16)
	Clock()
	IncrementPC()
	Halt()
	Stop()
	EI()
	EIWithDelay()
	PrefixCB()
}

// toRegisterPair returns a 16-bit register pair from two 8-bit registers
// If you want BC, pass B and C in that order
func ToRegisterPair(high, low uint8) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// cpu.FromRegisterPair returns two 8-bit registers from a 16-bit register pair
// It returns it in low, high order. E.g. if you have BC (CB in little endian) it will return B, C
func FromRegisterPair(val uint16) (high uint8, low uint8) {
	return uint8(val >> 8), uint8(val)
}
