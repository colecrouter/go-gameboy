package lr35902

import (
	"testing"
)

func TestByteLengths(t *testing.T) {
	for i := range uint8(0xFF) {
		mnemonic := mnemonics[i]

		if mnemonic == "INVALID" {
			continue // TODO add more elegant catch case
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
			for j := 0; j < ticks; j++ {
				manualClock <- struct{}{}
			}
			close(manualClock)
			cpu.clock = manualClock

			cpu.MClock()

			if int(cpu.registers.PC) != instrLengths[i] {
				t.Errorf("PC: got %d, want %d", cpu.registers.PC, instrLengths[i])
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
