package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
)

/*
RLA
┌────────────────────┐
│ ┌──┐  ┌─────────┐  │
└─│CY│<─│7<──────0│<─┘
  └──┘  └─────────┘
             A

RLCA
      ┌──────────────┐
┌──┐  │ ┌─────────┐  │
│CY│<─┴─│7<──────0│<─┘
└──┘    └─────────┘
             A
*/

func Rotate(r operands.Operand[uint8], left, useCarryBit, updateZ bool) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		val := r.Read(c)

		// Extract the bit that will be rotated out.
		// For left rotate, that is the MSB; for right rotate, the LSB.
		var carriedOut uint8
		if left {
			carriedOut = (val & 0x80) >> 7
		} else {
			carriedOut = val & 1
		}

		// Determine the input bit.
		// For instructions that use the external carry, use c.flags.Carry;
		// otherwise, use the bit that got shifted out.
		var carryIn uint8
		if useCarryBit {
			if c.Flags().Carry {
				carryIn = 1
			} else {
				carryIn = 0
			}
		} else {
			carryIn = carriedOut
		}

		// Perform the rotation.
		if left {
			val = (val << 1) | carryIn
		} else {
			val = (val >> 1) | (carryIn << 7)
		}

		// flags.Set the new Carry flag based on the bit that was rotated out.
		var newCarryFlag = flags.Reset
		if carriedOut == 1 {
			newCarryFlag = flags.Set
		}

		// Update the Zero flag only if requested (CB rotates update Z).
		var newZFlag flags.FlagState
		if updateZ {
			if val == 0 {
				newZFlag = flags.Set
			} else {
				newZFlag = flags.Reset
			}
		} else {
			newZFlag = flags.Reset // Always clear Z when !updateZ.
		}

		// Write the result back to the register.
		r.Write(c, val)

		// The N and H flags are reset on rotate instructions.
		c.Flags().Set(newZFlag, flags.Reset, flags.Reset, newCarryFlag)

		switch r.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		}

		c.Registers().PC++
		return nil
	})
}
func Shift(r operands.Operand[uint8], left, arithmeticRight bool) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		p := r.Read(c)
		original := p

		// Determine what the new Carry flag should be.
		var newCarry bool

		if left {
			// For left shifts, the high bit goes into carry.
			newCarry = (original & 0x80) != 0
			p = original << 1
		} else {
			// For right shifts, the low bit is carried out.
			newCarry = (original & 0x01) != 0
			if arithmeticRight {
				// For arithmetic right shifts, preserve the original MSB.
				msb := original & 0x80
				p = (original >> 1) | msb
			} else {
				p = original >> 1
			}
		}

		// flags.Set the Zero flag if the result is zero.
		flagZero := (p == 0)

		// Now, convert boolean conditions into flag behaviors.
		// For c.Flags().Set(zero, N, H, carry):
		var zeroBehavior, nBehavior, hBehavior, carryBehavior flags.FlagState

		if flagZero {
			zeroBehavior = flags.Set
		} else {
			zeroBehavior = flags.Reset
		}

		// The N and H flags are always cleared after these shifts.
		nBehavior = flags.Reset
		hBehavior = flags.Reset

		if newCarry {
			carryBehavior = flags.Set
		} else {
			carryBehavior = flags.Reset
		}

		// Write the result back to the register.
		r.Write(c, p)

		// Update the flags.
		c.Flags().Set(zeroBehavior, nBehavior, hBehavior, carryBehavior)

		switch r.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		}

		c.Registers().PC++
		return nil
	})
}
func Swap(r operands.Operand[uint8]) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		p := r.Read(c)

		// Swap the upper and lower nibbles
		p = (p&0xf)<<4 | p>>4

		zero := flags.Reset
		if p == 0 {
			zero = flags.Set
		}

		r.Write(c, p)

		c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)

		switch r.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		}

		c.Registers().PC++
		return nil
	})
}

// Comparing
func Compare(r1, r2 operands.Operand[uint8]) []shared.MicroOp {
	ops := []shared.MicroOp{}
	switch val := r2.(type) {
	case *operands.ImmediateOperand8:
		ops = append(ops, Immediate8IntoZ)
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		a := r1.Read(c)
		b := r2.Read(c)
		zero := flags.Reset
		carry := flags.Reset
		hc := flags.Reset

		if a < b {
			carry = flags.Set
		}
		if a&0xF < b&0xF {
			hc = flags.Set
		}
		if a == b {
			zero = flags.Set
		}

		c.Flags().Set(zero, flags.Set, hc, carry)

		c.Registers().PC++

		return nil
	})
}
func ReadBit(r operands.Operand[uint8], b uint8) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		val := r.Read(c)

		zero := flags.Reset
		if val&(1<<b) == 0 {
			zero = flags.Set
		}
		c.Flags().Set(zero, flags.Reset, flags.Set, flags.Leave)

		c.Registers().PC++

		return nil
	})
}
func ComplementAcc() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			a := c.Registers().A
			a = ^a
			c.Registers().A = a
			c.Flags().Set(flags.Leave, flags.Set, flags.Set, flags.Leave)

			c.Registers().PC++

			return nil
		},
	}
}
func ComplementCarry() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			carry := c.Flags().Carry
			c.Flags().Set(flags.Leave, flags.Reset, flags.Set, flags.Leave)
			c.Flags().Carry = !carry

			c.Registers().PC++

			return nil
		},
	}
}
func SetCarry() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.Flags().Set(flags.Leave, flags.Reset, flags.Reset, flags.Set)

			c.Registers().PC++

			return nil
		},
	}
}

// Bit manipulation
func ResetBit(r operands.Operand[uint8], b uint8) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		p := r.Read(c)

		p &= ^(1 << b)

		r.Write(c, p)

		switch r.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		}

		c.Registers().PC++
		return nil
	})
}
func SetBit(r operands.Operand[uint8], b uint8) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := r.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		p := r.Read(c)

		p |= 1 << b

		r.Write(c, p)

		switch r.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		}

		c.Registers().PC++
		return nil
	})
}
