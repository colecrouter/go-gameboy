package lr35902

import (
	"testing"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/system"
)

func TestByteLengths(t *testing.T) {
	// Create a new LR35902 CPU
	bus := &memory.Bus{}
	ir := &io.Interrupt{}
	ie := &io.Interrupt{}
	ioreg := io.NewRegisters(nil, bus, ir)
	mem := &memory.Memory{Buffer: make([]uint8, 0x10000)}
	bus.AddDevice(0, 0xFFFF, mem)
	broadcaster := system.NewBroadcaster()
	cpu := NewLR35902(broadcaster, bus, ioreg, ie)

	go system.ClockGenerator(broadcaster, 4)

	for i := 0; i < 0x100; i++ {
		// Lookup mnemonic
		mnemonic := mnemonics[i]

		t.Run(mnemonic, func(t *testing.T) {
			// Reset PC
			cpu.registers.PC = 0

			// Load instruction
			mem.Write(0, uint8(i))

			// Execute instruction
			cpu.MClock()

			// Check PC
			// +1 because the PC is incremented after the instruction is fetched
			if int(cpu.registers.PC) != instrLengths[i] {
				t.Errorf("PC: got %d, want %d", cpu.registers.PC, instrLengths[i])
			}
		})
	}
}

// TestCyclesUnconditional uses the new helper to deduplicate boilerplate.
func TestCyclesUnconditional(t *testing.T) {
	for i := range uint8(0xFF) {
		// Only run for unconditional instructions.
		if instrCycles[i] != instrCyclesCond[i] {
			continue
		}
		mnemonic := mnemonics[i]
		t.Run(mnemonic, func(t *testing.T) {
			runCyclesTest(t, uint8(i), instrCycles[i], instrLengths[i], false, nil)
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
		if instrCycles[i] == instrCyclesCond[i] {
			continue
		}
		// Skip if this opcode doesn't have a condition type.
		if cond := conditionMap[uint8(i)]; cond == CondNone {
			continue
		}

		mnemonic := mnemonics[i]
		t.Run(mnemonic, func(t *testing.T) {
			runCyclesTest(t, uint8(i), instrCyclesCond[i], instrLengths[i], true, nil)
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
			runCyclesTest(t, uint8(i), ticks, ticks, false, adjust)
		})
	}
}
