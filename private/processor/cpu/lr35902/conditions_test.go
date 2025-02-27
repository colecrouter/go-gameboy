package lr35902

import "github.com/colecrouter/gameboy-go/private/processor/cpu/flags"

// conditionType represents the type of condition that an opcode requires.
// The condition is based on the CPU flags, namely the Zero and Carry flags.
type conditionType int

const (
	CondNone conditionType = iota
	CondNZ
	CondZ
	CondNC
	CondC
)

// conditionMap maps opcodes to their corresponding condition type.
// Only include opcodes with conditional behavior.
var conditionMap = [0x100]conditionType{
	0xC0: CondNZ, // RET NZ
	0xC2: CondNZ, // JP NZ, a16
	0xC4: CondNZ, // CALL NZ, a16
	0xC8: CondZ,  // RET Z
	0xCA: CondZ,  // JP Z, a16
	0xCC: CondZ,  // CALL Z, a16
	0xD0: CondNC, // RET NC
	0xD2: CondNC, // JP NC, a16
	0xD4: CondNC, // CALL NC, a16
	0xD8: CondC,  // RET C
	0xDA: CondC,  // JP C, a16
	0xDC: CondC,  // CALL C, a16
}

// setupConditionalByOpcode configures the CPU flags based on the opcode's condition.
// trigger == true means the condition is met.
func setupConditionalByOpcode(cpu *LR35902, opcode byte, trigger bool) {
	condType := conditionMap[opcode]

	switch condType {
	case CondNZ:
		if trigger {
			// For NZ: condition met → Zero flag must be false.
			cpu.flags.Set(flags.Reset, flags.Leave, flags.Leave, flags.Leave)
		} else {
			// Condition not met: Zero flag is true.
			cpu.flags.Set(flags.Set, flags.Leave, flags.Leave, flags.Leave)
		}
	case CondZ:
		if trigger {
			// For Z: condition met → Zero flag is true.
			cpu.flags.Set(flags.Set, flags.Leave, flags.Leave, flags.Leave)
		} else {
			// Condition not met: Zero flag is false.
			cpu.flags.Set(flags.Reset, flags.Leave, flags.Leave, flags.Leave)
		}
	case CondNC:
		if trigger {
			// For NC: condition met → Carry flag must be false.
			cpu.flags.Set(flags.Leave, flags.Leave, flags.Leave, flags.Reset)
		} else {
			// Condition not met: Carry flag is true.
			cpu.flags.Set(flags.Leave, flags.Leave, flags.Leave, flags.Set)
		}
	case CondC:
		if trigger {
			// For C: condition met → Carry flag is true.
			cpu.flags.Set(flags.Leave, flags.Leave, flags.Leave, flags.Set)
		} else {
			// Condition not met: Carry flag is false.
			cpu.flags.Set(flags.Leave, flags.Leave, flags.Leave, flags.Reset)
		}
	}
}
