package conditions

import "github.com/colecrouter/gameboy-go/private/processor/cpu/flags"

type Condition uint

const (
	Always Condition = iota
	NZ
	Z
	NC
	C
)

func (c Condition) Test(flags *flags.Flags) bool {
	switch c {
	case NZ:
		return !flags.Zero
	case Z:
		return flags.Zero
	case NC:
		return !flags.Carry
	case C:
		return flags.Carry
	}
	return false
}
