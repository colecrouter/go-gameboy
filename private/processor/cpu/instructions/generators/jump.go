package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/conditions"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

// Jump
func Jump(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		Immediate8IntoZ,
		Immediate8IntoW,
	}

	ops = append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			return &[]shared.MicroOp{
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					c.Registers().PC = helpers.ToRegisterPair(ctx.W, ctx.Z)
					return nil
				},
				Idle,
			}
		} else {
			return &[]shared.MicroOp{
				Idle,
			}
		}
	})

	return ops
}

func JumpRelative(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		Immediate8IntoZ,
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			// Get separate bits from PC
			high, low := helpers.FromRegisterPair(c.Registers().PC)

			// Cast to signed 8-bit integer, store back in Z
			lowSigned := int8(low)
			ctx.Z = uint8(int16(high) + int16(lowSigned))

			// Get the carry from the 7th bit
			carry := (low & 0x80) != 0

			// Get the sign of Z
			sign := (ctx.Z & 0x80) != 0

			// Calculate the adjustment
			var adj int8
			if carry && !sign {
				adj = 1
			} else if !carry && sign {
				adj = -1
			} else {
				adj = 0
			}

			// Set W to the high byte of PC plus the adjustment
			ctx.W = high + uint8(adj)

			return &[]shared.MicroOp{
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					c.Registers().PC = uint16(int32(c.Registers().PC)+int32(ctx.Z)) + 1
					return nil
				},
			}
		} else {
			return &[]shared.MicroOp{
				Idle,
			}
		}
	})
}

func JumpHL() []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Set PC to HL
			c.Registers().PC = helpers.ToRegisterPair(c.Registers().H, c.Registers().L)

			return nil
		},
	}
}

// Subroutines
func Return(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		Idle,
	}

	if condition != conditions.Always {
		// I think this is supposed to be a cycle to evaluate the condition
		// Not sure
		ops = append(ops, Idle)
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			return &[]shared.MicroOp{
				StackIntoZ,
				StackIntoW,
				JumpToWZ,
			}
		} else {
			return &[]shared.MicroOp{
				Idle,
			}
		}
	})
}

func ReturnInterrupt() []shared.MicroOp {
	return []shared.MicroOp{
		StackIntoZ,
		StackIntoW,
		JumpToWZSetIME,
	}
}
func Call(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		Immediate8IntoZ,
		Immediate8IntoW,
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			return &[]shared.MicroOp{
				// Decrement SP
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					c.Registers().SP--
					return nil
				},
				// Push MSB of PC onto the stack
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					high, _ := helpers.FromRegisterPair(c.Registers().PC)
					c.Write(c.Registers().SP, high)

					// Decrement SP
					c.Registers().SP--
					return nil
				},
				// Push LSB of PC onto the stack
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					_, low := helpers.FromRegisterPair(c.Registers().PC)
					c.Write(c.Registers().SP, low)

					// Set PC to the address
					c.Registers().PC = helpers.ToRegisterPair(ctx.W, ctx.Z)
					return nil
				},
			}
		} else {
			return &[]shared.MicroOp{
				Idle,
			}
		}
	})
}
func ResetPC(addr uint16) []shared.MicroOp {
	return []shared.MicroOp{
		// Decrement SP
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			c.Registers().SP--

			return nil
		},
		// Push MSB of PC onto the stack
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			high, _ := helpers.FromRegisterPair(c.Registers().PC)
			c.Write(c.Registers().SP, high)

			// Decrement SP
			c.Registers().SP--

			return nil
		},
		// Push LSB of PC onto the stack
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			_, low := helpers.FromRegisterPair(c.Registers().PC)
			c.Write(c.Registers().SP, low)

			// Set PC to the address
			c.Registers().PC = addr

			return nil
		},
		Idle,
	}
}
