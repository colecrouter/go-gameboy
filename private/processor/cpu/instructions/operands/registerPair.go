package operands

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

type RegisterPair uint

const (
	AF RegisterPair = iota
	BC
	DE
	HL
	SP
)

type RegisterPairOperand struct {
	RegisterPair RegisterPair
}

func (rp *RegisterPairOperand) Read(c cpu.CPU) uint16 {
	switch rp.RegisterPair {
	case AF:
		return helpers.ToRegisterPair(c.Registers().A, c.Flags().Read())
	case BC:
		return helpers.ToRegisterPair(c.Registers().B, c.Registers().C)
	case DE:
		return helpers.ToRegisterPair(c.Registers().D, c.Registers().E)
	case HL:
		return helpers.ToRegisterPair(c.Registers().H, c.Registers().L)
	case SP:
		return c.Registers().SP
	}
	return 0
}

func (rp *RegisterPairOperand) Write(c cpu.CPU, value uint16) {
	switch rp.RegisterPair {
	case AF:
		high, low := helpers.FromRegisterPair(value)
		c.Registers().A = high
		c.Flags().Write(low)
	case BC:
		c.Registers().B, c.Registers().C = helpers.FromRegisterPair(value)
	case DE:
		c.Registers().D, c.Registers().E = helpers.FromRegisterPair(value)
	case HL:
		c.Registers().H, c.Registers().L = helpers.FromRegisterPair(value)
	case SP:
		c.Registers().SP = value
	}
}

func ToRegisterPair(high, low Register) RegisterPair {
	switch high {
	case A:
		switch low {
		case F:
			return AF
		}
	case B:
		switch low {
		case C:
			return BC
		}
	case D:
		switch low {
		case E:
			return DE
		}
	case H:
		switch low {
		case L:
			return HL
		}
	}
	return 0
}

func FromRegisterPair(rp RegisterPair) (Register, Register) {
	switch rp {
	case AF:
		return A, F
	case BC:
		return B, C
	case DE:
		return D, E
	case HL:
		return H, L
	}
	return 0, 0
}
