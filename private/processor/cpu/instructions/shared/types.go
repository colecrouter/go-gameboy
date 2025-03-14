package shared

import "github.com/colecrouter/gameboy-go/private/processor/cpu"

// MicroOp is a function that performs a single operation on the CPU.
type MicroOp func(c cpu.CPU, ctx *Context) *[]MicroOp

// Instruction is a sequence of micro-operations that make up a single instruction.
type Instruction []MicroOp
