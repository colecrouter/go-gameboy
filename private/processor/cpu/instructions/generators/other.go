package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
)

func DecimalAdjust() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			a := c.Registers().A
			subtract := c.Flags().Subtract
			halfCarry := c.Flags().HalfCarry
			carry := c.Flags().Carry
			offset := uint8(0)

			if !subtract {
				if halfCarry || (a&0x0F) > 0x09 {
					offset |= 0x06
				}
				if carry || a > 0x99 {
					offset |= 0x60
				}
				a += offset
			} else {
				if halfCarry {
					offset |= 0x06
				}
				if carry {
					offset |= 0x60
				}
				a -= offset
			}

			c.Registers().A = a
			c.Flags().Zero = (a == 0)
			c.Flags().HalfCarry = false
			c.Flags().Carry = ((offset & 0x60) != 0)

			return nil
		},
	}
}

func Stop() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.Stop()

			return nil
		},
	}
}

func Halt() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.Halt()

			return nil
		},
	}
}

func Nop() []shared.MicroOp {
	return []shared.MicroOp{
		Idle,
	}
}

func EIWithDelay() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.EIWithDelay()

			c.Registers().PC++

			return nil
		},
	}
}

func DisableInterrupts() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.DI()

			c.Registers().PC++

			return nil
		},
	}
}

func PrefixCB() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.PrefixCB()

			c.Registers().PC++

			return nil
		},
	}
}
