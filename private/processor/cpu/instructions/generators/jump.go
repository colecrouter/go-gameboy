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

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			return &[]shared.MicroOp{
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					c.Registers().PC = helpers.ToRegisterPair(ctx.Z, ctx.W)
					return nil
				},
				Idle,
			}
		} else {
			return nil
		}
	}, NextPC)
}

func JumpRelative(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			// Read the immediate value into Z
			c.Registers().PC++
			ctx.Z = c.Read(c.Registers().PC)
			return nil
		},
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			// Get current PC
			pc := c.Registers().PC

			// Calculate new address
			newAddr := uint16(int16(pc) + int16(int8(ctx.Z)))

			ctx.W, ctx.Z = helpers.FromRegisterPair(newAddr)

			return &[]shared.MicroOp{
				func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
					c.Registers().PC = uint16(helpers.ToRegisterPair(ctx.W, ctx.Z))
					return nil
				},
			}
		} else {
			return nil
		}
	}, NextPC)
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
	var ops []shared.MicroOp

	if condition != conditions.Always {
		ops = []shared.MicroOp{
			Idle,
		}
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
				NextPC,
			}
		}
	})
}

func ReturnInterrupt() []shared.MicroOp {
	return []shared.MicroOp{
		StackIntoZ,
		StackIntoW,
		JumpToWZSetIME,
		NextPC,
	}
}
func Call(condition conditions.Condition) []shared.MicroOp {
	ops := []shared.MicroOp{
		Immediate8IntoZ,
		Immediate8IntoW,
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		if condition.Test(c.Flags()) {
			// Decrement SP
			c.Registers().SP--

			return &[]shared.MicroOp{
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
			return nil
		}
	}, NextPC)
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
