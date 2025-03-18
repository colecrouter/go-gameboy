package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
)

// Increment increments an 8-bit register.
func Increment(op operands.Operand[uint8]) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := op.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		old := op.Read(c)
		newVal := old + 1
		zero := flags.Reset
		hc := flags.Reset
		if old&0xF == 0xF {
			hc = flags.Set
		}
		if newVal == 0 {
			zero = flags.Set
		}
		op.Write(c, newVal)
		c.Flags().Set(zero, flags.Reset, hc, flags.Leave)

		switch op.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		default:
			c.Registers().PC++
			return nil
		}
	})
}

// Increment16 increments a 16-bit register pair.
func Increment16(op *operands.RegisterPairOperand) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			combined := op.Read(c)
			combined++
			op.Write(c, combined)

			return nil
		},
		NextPC,
	}
}

// Decrement decrements an 8-bit register.
func Decrement(op operands.Operand[uint8]) []shared.MicroOp {
	var ops []shared.MicroOp

	switch val := op.(type) {
	case *operands.IndirectOperand:
		ops = append(ops, IndirectIntoZ(val.Indirectable))
	}

	return append(ops, func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
		old := op.Read(c)
		newVal := old - 1
		zero := flags.Reset
		hc := flags.Reset
		if old&0xF == 0x0 {
			hc = flags.Set
		}
		if newVal == 0 {
			zero = flags.Set
		}
		op.Write(c, newVal)
		c.Flags().Set(zero, flags.Set, hc, flags.Leave)

		switch op.(type) {
		case *operands.IndirectOperand:
			return &[]shared.MicroOp{NextPC}
		default:
			c.Registers().PC++
			return nil
		}
	})
}

// Decrement16 decrements a 16-bit register pair.
func Decrement16(op *operands.RegisterPairOperand) []shared.MicroOp {
	return []shared.MicroOp{
		func(c cpu.CPU, ctx *shared.Context) *[]shared.MicroOp {
			combined := op.Read(c)
			combined--
			op.Write(c, combined)

			return nil
		},
		NextPC,
	}
}
