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
	switch val.(type) {
	case *operands.ImmediateOperand8:
		ops = append(ops, Immediate8IntoZ)
	case *operands.ImmediateIndirectOperand:
		ops = append(ops, IndirectImmediateIntoZ)
	case *operands.IndirectOperand:
		v := val.(*operands.IndirectOperand)
		ops = append(ops, IndirectIntoZ(v.Indirectable))
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

		// PC
		switch val.(type) {
		case *operands.IndirectOperand:
			// Nothing, we'll increment the PC later
		default:
			c.Registers().PC++
		}

		return nil
	})

	// Indirect instructions need an extra step to write the value to memory
	switch dest.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, Idle)
	}

	// TODO LD (nn), A

	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		c.Registers().PC++
		return nil
	})

	return ops
}
func Load(dest, val operands.Operand[uint8]) []shared.MicroOp {
	return LoadInc(dest, val, None, None)
}
func Load16(rp, val operands.Operand[uint16]) []shared.MicroOp {
	ops := []shared.MicroOp{}

	switch val.(type) {
	case *operands.ImmediateOperand16:
		ops = append(ops, Immediate8IntoZ)
		ops = append(ops, Immediate8IntoW)
	}

	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		val.Write(c, helpers.ToRegisterPair(ctx.Z, ctx.W))

		return nil
	})

	// 16-bit loads always need an extra cycle, not sure why
	ops = append(ops, Idle)

	return ops
}

func LoadHLSPOffset(offset operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		Immediate8IntoZ,
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			sp := c.Registers().SP

			// Get LSB of SP + offset
			_, lsb := helpers.FromRegisterPair(sp)
			offset := int8(ctx.Z)
			result := uint8(int16(lsb) + int16(offset)) // Sign-extend offset correctly.

			// Set L register
			c.Registers().L = result

			// Compute flags using the raw byte value.
			unsignedOffset := uint16(uint8(offset))
			var hc, carry = flags.Reset, flags.Reset

			if ((sp & 0xF) + (unsignedOffset & 0xF)) > 0xF {
				hc = flags.Set
			}
			if ((sp & 0xFF) + (unsignedOffset & 0xFF)) > 0xFF {
				carry = flags.Set
			}

			// Store result in HL and update flags: Z and N reset.
			c.Flags().Set(flags.Reset, flags.Reset, hc, carry)

			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {

			// Need to adjust the sign of the offset
			adjust := 0x00
			if ctx.Z > 0x7F {
				adjust = 0xFF
			}

			// I don't really know what's going on here, some sort of complement?
			carry := 0
			if c.Flags().Carry {
				carry = 1
			}
			msb, _ := helpers.FromRegisterPair(c.Registers().SP)
			ctx.W = uint8(int16(msb) + int16(adjust) + int16(carry))

			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Set H register
			c.Registers().H = ctx.W

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
	ops = append(ops, Idle)

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
		op.Write(c, helpers.ToRegisterPair(ctx.Z, ctx.W))
		return nil
	})

	// Extra cycle
	ops = append(ops, Idle)

	return ops
}
