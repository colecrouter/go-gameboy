package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

/*
	This file contains predicate functions that are meant to be used while
	constructing instruction generators. They return procedures that abstract
	away common patterns, such as getting an immediate value into a register.

	Their goal is to make the instruction generators more readable.
*/

// Immediate8IntoZ reads the immediate 8-bit value into the Z register.
func Immediate8IntoZ(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Read the value into Z
	ctx.Z = c.Read(c.Registers().PC)

	// Increment PC
	c.Registers().PC++

	return nil
}

// Immediate8IntoW reads the immediate 8-bit value into the W register.
func Immediate8IntoW(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Read the value into W
	ctx.W = c.Read(c.Registers().PC)

	// Increment PC
	c.Registers().PC++

	return nil
}

// RegisterIntoZ reads the value in the specified register into the Z register.
func RegisterIntoZ(reg operands.Register) func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	return func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		// Read the value in the specified register into Z
		switch reg {
		case operands.A:
			ctx.Z = c.Registers().A
		case operands.B:
			ctx.Z = c.Registers().B
		case operands.C:
			ctx.Z = c.Registers().C
		case operands.D:
			ctx.Z = c.Registers().D
		case operands.E:
			ctx.Z = c.Registers().E
		case operands.H:
			ctx.Z = c.Registers().H
		case operands.L:
			ctx.Z = c.Registers().L
		case operands.F:
			ctx.Z = c.Flags().Read()
		default:
			panic("Invalid register")
		}

		return nil
	}
}

// RegisterIntoW reads the value in the specified register into the W register.
func RegisterIntoW(reg operands.Register) func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	return func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		// Read the value in the specified register into W
		switch reg {
		case operands.A:
			ctx.W = c.Registers().A
		case operands.B:
			ctx.W = c.Registers().B
		case operands.C:
			ctx.W = c.Registers().C
		case operands.D:
			ctx.W = c.Registers().D
		case operands.E:
			ctx.W = c.Registers().E
		case operands.H:
			ctx.W = c.Registers().H
		case operands.L:
			ctx.W = c.Registers().L
		case operands.F:
			ctx.W = c.Flags().Read()
		default:
			panic("Invalid register")
		}

		return nil
	}
}

// JumpToWZ jumps to the address in WZ. WZ is little endian, so Z should be fetched first, then W.
func JumpToWZ(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Jump to the address in WZ
	c.Registers().PC = helpers.ToRegisterPair(ctx.W, ctx.Z)

	return nil
}

// JumpToWZSetIME jumps to the address in WZ and enables interrupts.
func JumpToWZSetIME(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Jump to the address in WZ
	c.Registers().PC = helpers.ToRegisterPair(ctx.W, ctx.Z)

	// Enable interrupts
	c.EI()

	return nil
}

// StackIntoZ reads the value at the top of the stack into the Z register, then increments SP.
func StackIntoZ(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Read the value at the top of the stack into Z
	ctx.Z = c.Read(c.Registers().SP)

	// Increment SP
	c.Registers().SP++

	return nil
}

// StackIntoW reads the value at the top of the stack into the W register, then increments SP.
func StackIntoW(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Read the value at the top of the stack into W
	ctx.W = c.Read(c.Registers().SP)

	// Increment SP
	c.Registers().SP++

	return nil
}

// NextPC increments the program counter without doing anything.
func NextPC(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Do nothing
	c.Registers().PC++

	return nil
}

// Idle does nothing. It is meant to be used in specific cases (such as jump instructions)
// where a clock cycle is needed, but PC should not be incremented. In most cases, NextPC
// should be used instead.
func Idle(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	return nil
}

// IndirectIntoZ reads the value at the address in the specified register pair into the Z register.
func IndirectIntoZ(reg operands.Indirectable) func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	return func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		// Read the value at the address in the specified register pair into Z
		switch reg {
		case operands.AF_:
			ctx.Z = c.Read(helpers.ToRegisterPair(c.Registers().A, c.Flags().Read()))
		case operands.BC_:
			ctx.Z = c.Read(helpers.ToRegisterPair(c.Registers().B, c.Registers().C))
		case operands.DE_:
			ctx.Z = c.Read(helpers.ToRegisterPair(c.Registers().D, c.Registers().E))
		case operands.HL_:
			ctx.Z = c.Read(helpers.ToRegisterPair(c.Registers().H, c.Registers().L))
		case operands.SP_:
			ctx.Z = c.Read(c.Registers().SP)
		case operands.A_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().A))
		case operands.B_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().B))
		case operands.C_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().C))
		case operands.D_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().D))
		case operands.E_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().E))
		case operands.H_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().H))
		case operands.L_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Registers().L))
		case operands.F_:
			ctx.Z = c.Read(helpers.ToRegisterPair(0xFF, c.Flags().Read()))
		default:
			panic("Invalid register pair")
		}

		return nil
	}
}

func IndirectFromZ(reg operands.Indirectable) func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	return func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		// Write the value in Z to the address in the specified register pair
		switch reg {
		case operands.AF_:
			c.Write(helpers.ToRegisterPair(c.Registers().A, c.Flags().Read()), ctx.Z)
		case operands.BC_:
			c.Write(helpers.ToRegisterPair(c.Registers().B, c.Registers().C), ctx.Z)
		case operands.DE_:
			c.Write(helpers.ToRegisterPair(c.Registers().D, c.Registers().E), ctx.Z)
		case operands.HL_:
			c.Write(helpers.ToRegisterPair(c.Registers().H, c.Registers().L), ctx.Z)
		case operands.SP_:
			c.Write(c.Registers().SP, ctx.Z)
		case operands.A_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().A), ctx.Z)
		case operands.B_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().B), ctx.Z)
		case operands.C_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().C), ctx.Z)
		case operands.D_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().D), ctx.Z)
		case operands.E_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().E), ctx.Z)
		case operands.H_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().H), ctx.Z)
		case operands.L_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Registers().L), ctx.Z)
		case operands.F_:
			c.Write(helpers.ToRegisterPair(0xFF, c.Flags().Read()), ctx.Z)
		default:
			panic("Invalid register pair")
		}

		return nil
	}
}

// IndirectImmediateIntoZ reads the value at the immediate address into the Z register.
func IndirectImmediateIntoZ(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Increment PC
	c.Registers().PC++

	// Read the value at the immediate address into Z
	ctx.Z = c.Read(helpers.ToRegisterPair(c.Read(c.Registers().PC), c.Read(c.Registers().PC+1)))

	return nil
}

// IndirectImmediateIntoW reads the value at the immediate address into the W register.
func IndirectImmediateIntoW(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Increment PC
	c.Registers().PC++

	// Read the value at the immediate address into W
	ctx.W = c.Read(helpers.ToRegisterPair(c.Read(c.Registers().PC), c.Read(c.Registers().PC+1)))

	return nil
}

// IndirectImmediateFromZ writes the value in Z to the immediate address.
func IndirectImmediateFromZ(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Increment PC
	c.Registers().PC++

	// Write the value in Z to the immediate address
	c.Write(helpers.ToRegisterPair(c.Read(c.Registers().PC), c.Read(c.Registers().PC+1)), ctx.Z)

	return nil
}

// IndirectImmediateFromW writes the value in W to the immediate address.
func IndirectImmediateFromW(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
	// Increment PC
	c.Registers().PC++

	// Write the value in W to the immediate address
	c.Write(helpers.ToRegisterPair(c.Read(c.Registers().PC), c.Read(c.Registers().PC+1)), ctx.W)

	return nil
}
