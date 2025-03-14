package operands

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

type OperandSize interface {
	uint8 | uint16
}

type Operand[T OperandSize] interface {
	Read(c cpu.CPU) T
	Write(c cpu.CPU, val T)
}
