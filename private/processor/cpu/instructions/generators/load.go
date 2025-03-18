package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

type incDec uint

const (
	None incDec = iota
	Inc
	Dec
)

// Memory access
func LoadInc(dest, val operands.Operand[uint8], incDest, incVal incDec) []shared.MicroOp {
	ops := []shared.MicroOp{}

	// Check if we need to take another step to load the value
	switch v := val.(type) {
	case *operands.ImmediateOperand8:
		ops = append(ops, Immediate8IntoZ)
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
		ops = append(ops, IndirectImmediateIntoW)
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(v.Indirectable))
	}

	switch d := dest.(type) {
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
		ops = append(ops, IndirectImmediateIntoW)
		ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Read the ZW value into the Z register
			c.Write(helpers.ToRegisterPair(ctx.Z, ctx.W), d.Read(c))
			return nil
		})
	}

	// Do the operation
	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		switch val.(type) {
		case *operands.RegisterOperand:
			dest.Write(c, val.Read(c))
		default:
			dest.Write(c, ctx.Z)
		}

		// Handle case for HL+ and HL- instructions
		switch incDest {
		case Inc:
			val.Write(c, val.Read(c)+1)
		case Dec:
			val.Write(c, val.Read(c)-1)
		}

		return nil
	}, NextPC)

	return ops
}
func Load(dest, val operands.Operand[uint8]) []shared.MicroOp {
	return LoadInc(dest, val, None, None)
}
func Load16[T operands.OperandSize](dest operands.Operand[T], val operands.Operand[uint16]) []shared.MicroOp {
	ops := []shared.MicroOp{}

	switch val.(type) {
	case *operands.ImmediateOperand16:
		ops = append(ops, Immediate8IntoZ)
		ops = append(ops, Immediate8IntoW)
	}
	switch any(dest).(type) {
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
		ops = append(ops, IndirectImmediateIntoW)
	}

	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		switch dest := any(dest).(type) {
		case *operands.RegisterPairOperand:
			dest.Write(c, val.Read(c))
		case *operands.ImmediateIndirectOperand:
			switch val.(type) {
			case *operands.ImmediateOperand16:
				high, low := helpers.FromRegisterPair(val.Read(c))
				c.Write(helpers.ToRegisterPair(ctx.Z, ctx.W), low)
				c.Write(helpers.ToRegisterPair(ctx.Z, ctx.W)+1, high)
			}
		}

		return nil
	})

	// 16-bit loads always need an extra cycle, not sure why
	ops = append(ops, NextPC)

	return ops
}
func LoadHigh(dest operands.Operand[uint8], val operands.Operand[uint8]) []shared.MicroOp {
	ops := []shared.MicroOp{}

	switch val := val.(type) {
	case *operands.RegisterOperand:
		ops = append(ops, RegisterIntoZ(val.Register))
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
	}

	switch dest.(type) {
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
		ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.Write(0xFF00+uint16(ctx.Z), val.Read(c))
			return nil
		},
			NextPC)
	case *operands.IndirectOperand:
		ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			dest.Write(c, val.Read(c))
			c.Registers().PC++
			return nil
		})
	}

	switch val.(type) {
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			ctx.Z = c.Read(0xFF00 + uint16(ctx.Z))
			return nil
		},
			func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
				dest.Write(c, ctx.Z)
				c.Registers().PC++
				return nil
			})

	case *operands.IndirectOperand:
		ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			dest.Write(c, val.Read(c))
			c.Registers().PC++
			return nil
		})
	}

	return ops
}
func LoadHLSPOffset() []shared.MicroOp {
	return []shared.MicroOp{
		Immediate8IntoZ,
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			sp := c.Registers().SP

			// Get LSB of SP + offset
			_, lsb := helpers.FromRegisterPair(sp)
			offset := int8(ctx.Z)
			result := uint8(int16(lsb) + int16(offset)) // Sign-extend offset.
			c.Registers().L = result

			// Compute flags.
			zero := flags.Reset
			hc := flags.Reset
			carry := flags.Reset
			if offset >= 0 {
				if ((sp & 0xF) + (uint16(offset) & 0xF)) > 0xF {
					hc = flags.Set
				}
				if ((sp & 0xFF) + (uint16(offset) & 0xFF)) > 0xFF {
					carry = flags.Set
				}
			} else {
				if (sp & 0xF) < (uint16(-offset) & 0xF) {
					hc = flags.Set
				}
				if (sp & 0xFF) < (uint16(-offset) & 0xFF) {
					carry = flags.Set
				}
			}
			c.Flags().Set(zero, flags.Reset, hc, carry)
			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Compute high-byte adjustment.
			offset := int8(ctx.Z)
			var highAdjustment int16
			if offset >= 0 {
				if c.Flags().Carry {
					highAdjustment = 1
				} else {
					highAdjustment = 0
				}
			} else {
				if c.Flags().Carry {
					highAdjustment = -1
				} else {
					highAdjustment = 0
				}
			}

			msb, _ := helpers.FromRegisterPair(c.Registers().SP)
			ctx.W = uint8(int16(msb) + highAdjustment)
			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Set H register.
			c.Registers().H = ctx.W
			c.Registers().PC++
			return nil
		},
	}
}

func Push(op operands.Operand[uint16]) []shared.MicroOp {
	// First, decrement the stack pointer
	ops := []shared.MicroOp{func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		c.Registers().SP--
		return nil
	}}

	// Push the low byte onto the stack, then decrement the stack pointer
	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		_, low := helpers.FromRegisterPair(op.Read(c))
		c.Write(c.Registers().SP, low)
		c.Registers().SP--
		return nil
	})

	// Push the high byte onto the stack
	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		high, _ := helpers.FromRegisterPair(op.Read(c))
		c.Write(c.Registers().SP, high)
		return nil
	})

	// Extra cycle
	ops = append(ops, NextPC)

	return ops
}

func Pop(op operands.Operand[uint16]) []shared.MicroOp {
	// First, read the high byte from the stack
	ops := []shared.MicroOp{
		StackIntoZ,
		StackIntoW,
	}

	// Write the value to the register pair
	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		op.Write(c, helpers.ToRegisterPair(ctx.W, ctx.Z))

		c.Registers().PC++
		return nil
	})

	return ops
}
