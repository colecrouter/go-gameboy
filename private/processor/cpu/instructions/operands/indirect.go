package operands

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

type Indirectable uint

const (
	AF_ Indirectable = iota
	BC_
	DE_
	HL_
	SP_
	A_
	B_
	C_
	D_
	E_
	H_
	L_
	F_
)

type IndirectOperand struct {
	Indirectable Indirectable
}

func (i *IndirectOperand) Read(c cpu.CPU) uint8 {
	switch i.Indirectable {
	case AF_:
		return c.Read(helpers.ToRegisterPair(c.Registers().A, c.Flags().Read()))
	case BC_:
		return c.Read(helpers.ToRegisterPair(c.Registers().B, c.Registers().C))
	case DE_:
		return c.Read(helpers.ToRegisterPair(c.Registers().D, c.Registers().E))
	case HL_:
		return c.Read(helpers.ToRegisterPair(c.Registers().H, c.Registers().L))
	case SP_:
		return c.Read(c.Registers().SP)
	case A_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().A))
	case B_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().B))
	case C_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().C))
	case D_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().D))
	case E_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().E))
	case H_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().H))
	case L_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Registers().L))
	case F_:
		return c.Read(helpers.ToRegisterPair(0xFF, c.Flags().Read()))
	default:
		panic("Invalid register pair")
	}
}

func (i *IndirectOperand) Write(c cpu.CPU, val uint8) {
	switch i.Indirectable {
	case AF_:
		c.Write(helpers.ToRegisterPair(c.Registers().A, c.Flags().Read()), val)
	case BC_:
		c.Write(helpers.ToRegisterPair(c.Registers().B, c.Registers().C), val)
	case DE_:
		c.Write(helpers.ToRegisterPair(c.Registers().D, c.Registers().E), val)
	case HL_:
		c.Write(helpers.ToRegisterPair(c.Registers().H, c.Registers().L), val)
	case SP_:
		c.Write(c.Registers().SP, val)
	case A_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().A), val)
	case B_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().B), val)
	case C_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().C), val)
	case D_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().D), val)
	case E_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().E), val)
	case H_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().H), val)
	case L_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Registers().L), val)
	case F_:
		c.Write(helpers.ToRegisterPair(0xFF, c.Flags().Read()), val)
	default:
		panic("Invalid register pair")
	}
}
