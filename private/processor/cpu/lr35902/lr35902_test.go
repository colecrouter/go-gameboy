package lr35902

import (
	"testing"

	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

func TestByteLengths(t *testing.T) {
	// For instructions we can't test, add them to the blacklist.
	blacklist := map[uint8]bool{
		0x18: true, // JR r8
		0xC3: true, // JP a16
		0xC7: true, // RST 00H
		0xC9: true, // RET
		0xCF: true, // RST 08H
		0xCD: true, // CALL a16
		0xD7: true, // RST 10H
		0xD9: true, // RETI
		0xDF: true, // RST 18H
		0xE7: true, // RST 20H
		0xE9: true, // JP (HL)
		0xEF: true, // RST 28H
		0xF7: true, // RST 30H
	}

	// Add invalid ops to blacklist.
	for i, l := range instrLengths {
		if l == 0 {
			blacklist[uint8(i)] = true
		}
	}

	for i := range uint8(0xFF) {
		mnemonic := mnemonics[i]

		if blacklist[i] {
			continue
		}

		t.Run(mnemonic, func(t *testing.T) {
			cpu, mem, _ := newTestCPU()
			cpu.registers.PC = 0
			mem.Write(0, uint8(i))

			// Setup so that conditional instructions don't jump
			setupConditionalByOpcode(cpu, uint8(i), false)

			// Preload clock ticks equal to the expected instruction length and close the channel.
			ticks := instrLengths[i]
			manualClock := make(chan struct{}, ticks)
			for j := 0; j < int(ticks); j++ {
				manualClock <- struct{}{}
			}
			close(manualClock)
			cpu.clock = manualClock
			cpu.clockAck = make(chan struct{}, ticks+10) // New ack channel

			cpu.MClock()

			if int(cpu.registers.PC) != instrLengths[i] {
				t.Errorf("PC: got %d, want %d", cpu.registers.PC, instrLengths[i])
			} else {
				t.Logf("PC: got %d", cpu.registers.PC)
			}
		})
	}
}

func TestCyclesUnconditional(t *testing.T) {
	for i := range uint8(0xFF) {
		// Only run for unconditional instructions.
		if instrCycles[i] != instrCyclesCond[i] || instrCycles[i] == 0 {
			continue
		}
		mnemonic := mnemonics[i]
		t.Run(mnemonic, func(t *testing.T) {
			runCyclesTest(t, uint8(i), instrCycles[i], false, nil)
		})
	}
}

// TestCyclesConditional tests opcodes with conditional cycle differences.
func TestCyclesConditional(t *testing.T) {
	// For each opcode where the unconditional cycles differ from the conditional ones,
	// and one exists in the conditionMap, run subtests.
	for i := range uint8(0xFF) {
		// Both of the below checks *should* effectively be the same
		// Might as well keep them in in case of mistakes

		// Skip if timings are identical.
		if instrCycles[i] == instrCyclesCond[i] || instrCyclesCond[i] == 0 {
			continue
		}
		// Skip if this opcode doesn't have a condition type.
		if cond := conditionMap[uint8(i)]; cond == CondNone {
			continue
		}

		mnemonic := mnemonics[i]
		t.Run(mnemonic, func(t *testing.T) {
			runCyclesTest(t, uint8(i), instrCyclesCond[i], true, nil)
		})
	}
}

// TestCyclesCB tests the cycle timings for CB-prefixed instructions.
func TestCyclesCB(t *testing.T) {
	for i := range uint8(0xFF) {
		mnemonic := getCBMnemonic(uint8(i))
		t.Run(mnemonic, func(t *testing.T) {
			adjust := func(cpu *LR35902) {
				cpu.cb = true
			}
			// For CB opcodes, preload ticks = instrCyclesCB[i] - 1 and expect PC to advance by that amount.
			ticks := instrCyclesCB[i] - 1
			runCyclesTest(t, uint8(i), ticks, false, adjust)
		})
	}
}

func TestJPHL(t *testing.T) {
	c, mem, _ := newTestCPU()
	jumpTarget := uint16(0x1234)
	c.registers.H, c.registers.L = helpers.FromRegisterPair(jumpTarget)
	mem.Write(0, 0xE9) // JP (HL) opcode

	// Preload clock with enough ticks (assume 1 tick is sufficient)
	manualClock := make(chan struct{}, 1)
	manualClock <- struct{}{}
	close(manualClock)
	c.clock = manualClock
	c.clockAck = make(chan struct{}, 10) // New ack channel

	c.MClock()

	if c.registers.PC != jumpTarget {
		t.Errorf("JP (HL) failed: got PC %d, want %d", c.registers.PC, jumpTarget)
	}
}
