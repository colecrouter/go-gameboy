package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

// Arithmetic
func Add(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			source := op1.Read(c)
			dest := op2.Read(c)

			zero := flags.Reset
			carry := flags.Reset
			hc := flags.Reset

			sum := uint16(source) + uint16(dest)
			if sum > 0xFF {
				carry = flags.Set
			}
			if ((source)&0xF)+(dest&0xF) > 0xF {
				hc = flags.Set
			}

			result := uint8(sum)
			if result == 0 {
				zero = flags.Set
			}
			op1.Write(c, result)
			c.Flags().Set(zero, flags.Reset, hc, carry)

			c.Registers().PC++

			return nil
		},
	}
}
func Add16(op1, op2 *operands.RegisterPairOperand) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {

			// Do low bit first
			sum := op1.Read(c) + op2.Read(c)
			_, low := helpers.FromRegisterPair(sum)

			existingHigh, _ := helpers.FromRegisterPair(op1.Read(c))

			combined := helpers.ToRegisterPair(existingHigh, low)

			// TODO flags
			op1.Write(c, combined)

			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Do high bit
			sum := op1.Read(c) + op2.Read(c)
			high, _ := helpers.FromRegisterPair(sum)

			_, existingLow := helpers.FromRegisterPair(op1.Read(c))

			combined := helpers.ToRegisterPair(high, existingLow)

			// TODO flags
			op1.Write(c, combined)

			return nil
		},
		Idle,
	}
}
func Sub(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			source := op1.Read(c)
			dest := op2.Read(c)

			zero := flags.Reset
			carry := flags.Reset
			hc := flags.Reset

			diff := int16(source) - int16(dest)
			if diff < 0 {
				carry = flags.Set
				diff += 256
			}
			if (source & 0xF) < (dest & 0xF) {
				hc = flags.Set
			}
			result := uint8(diff)
			if result == 0 {
				zero = flags.Set
			}
			op1.Write(c, result)
			c.Flags().Set(zero, flags.Set, hc, carry)

			c.Registers().PC++

			return nil
		},
	}
}
func And(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			val := op1.Read(c)
			r := op2.Read(c)

			r &= val
			zero := flags.Reset
			if r == 0 {
				zero = flags.Set
			}
			c.Flags().Set(zero, flags.Reset, flags.Set, flags.Reset)

			op2.Write(c, r)

			c.Registers().PC++

			return nil
		},
	}
}
func AddSPImmediate() []shared.MicroOp {
	return []shared.MicroOp{
		Immediate8IntoZ,
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			operand := int8(ctx.Z)
			zero := flags.Reset
			hc := flags.Reset
			carry := flags.Reset

			// Compute half-carry and carry flags using only the lower nibble/byte.
			if (c.Registers().SP&0xF)+(uint16(uint8(operand))&0xF) > 0xF {
				hc = flags.Set
			}
			if (c.Registers().SP&0xFF)+uint16(uint8(operand)) > 0xFF {
				carry = flags.Set
			}
			c.Flags().Set(zero, flags.Reset, hc, carry)

			return nil
		},
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			operand := int8(ctx.Z)
			result := c.Registers().SP + uint16(operand)
			c.Registers().SP = result

			return nil
		},
		Idle,
	}
}
func Or(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			val1 := op1.Read(c)
			val2 := op2.Read(c)
			result := val1 | val2
			zero := flags.Reset
			if result == 0 {
				zero = flags.Set
			}
			c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)
			op1.Write(c, result)
			c.Registers().PC++
			return nil
		},
	}
}
func Xor(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			val1 := op1.Read(c)
			val2 := op2.Read(c)
			result := val1 ^ val2
			zero := flags.Reset
			if result == 0 {
				zero = flags.Set
			}
			c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)
			op1.Write(c, result)
			c.Registers().PC++
			return nil
		},
	}
}
func AddC(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			val1 := op1.Read(c)
			val2 := op2.Read(c)
			zero := flags.Reset
			carry := flags.Reset
			hc := flags.Reset
			var carryIn uint16 = 0
			if c.Flags().Carry {
				carryIn = 1
			}
			if ((val1 & 0xF) + (val2 & 0xF) + uint8(carryIn)) > 0xF {
				hc = flags.Set
			}
			sum := uint16(val1) + uint16(val2) + carryIn
			if sum > 0xFF {
				carry = flags.Set
			}
			result := uint8(sum)
			if result == 0 {
				zero = flags.Set
			}
			op1.Write(c, result)
			c.Flags().Set(zero, flags.Reset, hc, carry)
			c.Registers().PC++
			return nil
		},
	}
}
func SubC(op1, op2 operands.Operand[uint8]) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			val1 := op1.Read(c)
			val2 := op2.Read(c)
			carryIn := uint8(0)
			if c.Flags().Carry {
				carryIn = 1
			}
			diff := int16(val1) - int16(val2) - int16(carryIn)
			carry := flags.Reset
			if diff < 0 {
				carry = flags.Set
				diff += 256
			}
			hc := flags.Reset
			if (val1 & 0xF) < ((val2 & 0xF) + carryIn) {
				hc = flags.Set
			}
			result := uint8(diff)
			zero := flags.Reset
			if result == 0 {
				zero = flags.Set
			}
			op1.Write(c, result)
			c.Flags().Set(zero, flags.Set, hc, carry)
			c.Registers().PC++
			return nil
		},
	}
}
