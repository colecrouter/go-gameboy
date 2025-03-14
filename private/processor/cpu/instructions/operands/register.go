package operands

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

type Register uint

const (
	A Register = iota
	B
	C
	D
	E
	H
	L
	F
)

type RegisterOperand struct {
	Register Register
}

func (r *RegisterOperand) Read(c cpu.CPU) uint8 {
	switch r.Register {
	case A:
		return c.Registers().A
	case B:
		return c.Registers().B
	case C:
		return c.Registers().C
	case D:
		return c.Registers().D
	case E:
		return c.Registers().E
	case H:
		return c.Registers().H
	case L:
		return c.Registers().L
	case F:
		return c.Flags().Read()
	default:
		panic("Invalid register")
	}
}

func (r *RegisterOperand) Write(c cpu.CPU, val uint8) {
	switch r.Register {
	case A:
		c.Registers().A = val
	case B:
		c.Registers().B = val
	case C:
		c.Registers().C = val
	case D:
		c.Registers().D = val
	case E:
		c.Registers().E = val
	case H:
		c.Registers().H = val
	case L:
		c.Registers().L = val
	case F:
		c.Flags().Write(val)
	default:
		panic("Invalid register")
	}
}
