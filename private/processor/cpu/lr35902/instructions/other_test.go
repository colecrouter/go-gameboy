package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimalAdjust(t *testing.T) {
	tests := []struct {
		name      string
		initialA  uint8
		subtract  bool
		halfCarry bool
		carry     bool
		expectedA uint8
		expectedZ bool
		expectedC bool
	}{
		// Addition cases (subtract = false)
		{"Add: No adjustment needed", 0x12, false, false, false, 0x12, false, false},
		{"Add: Lower digit > 9", 0x3A, false, false, false, 0x40, false, false},
		{"Add: HC set", 0x12, false, true, false, 0x18, false, false},
		{"Add: Upper digit > 9", 0xA2, false, false, false, 0x02, false, true},
		{"Add: Carry set", 0x42, false, false, true, 0xA2, false, true},
		{"Add: Both digits need adjustment", 0x9A, false, false, false, 0x00, true, true},
		{"Add: HC & upper digit > 9", 0xAD, false, true, false, 0x13, false, true},

		// Subtraction cases (subtract = true)
		{"Sub: No adjustment needed", 0x12, true, false, false, 0x12, false, false},
		{"Sub: HC set", 0x12, true, true, false, 0x0C, false, false},
		{"Sub: Carry set", 0x12, true, false, true, 0xB2, false, true},
		{"Sub: Both HC & C set", 0x12, true, true, true, 0xAC, false, true},
		{"Sub: Result zero", 0x00, true, false, false, 0x00, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA
			cpu.Flags().Subtract = tt.subtract
			cpu.Flags().HalfCarry = tt.halfCarry
			cpu.Flags().Carry = tt.carry

			decimalAdjust(cpu)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.subtract, cpu.Flags().Subtract, "Subtract flag should be preserved")
			assert.Equal(t, false, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}
