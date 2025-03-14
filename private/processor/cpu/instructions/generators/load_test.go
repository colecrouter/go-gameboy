package generators

import (
	"testing"

	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/enums"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
	"github.com/stretchr/testify/assert"
)

func TestLoad8(t *testing.T) {
	tests := []struct {
		name     string
		value    uint8
		expected uint8
	}{
		{"Load zero", 0x00, 0x00},
		{"Load non-zero", 0x42, 0x42},
		{"Load max value", 0xFF, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			var r uint8

			cpu.Execute(Load(A, B))

			assert.Equal(t, tt.expected, r, "unexpected loaded value")
		})
	}
}

func TestLoad16(t *testing.T) {
	tests := []struct {
		name      string
		value     uint16
		expectedH uint8
		expectedL uint8
	}{
		{"Load zero", 0x0000, 0x00, 0x00},
		{"Load typical value", 0x1234, 0x12, 0x34},
		{"Load max value", 0xFFFF, 0xFF, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			high, low := helpers.FromRegisterPair(tt.value)

			cpu.Registers().C = high
			cpu.Registers().B = low

			cpu.Execute(Load16(HL, BC))

			assert.Equal(t, tt.expectedH, cpu.Registers().H, "unexpected high byte value")
			assert.Equal(t, tt.expectedL, cpu.Registers().L, "unexpected low byte value")
		})
	}
}

func TestPop16(t *testing.T) {
	tests := []struct {
		name         string
		stackValue   uint16
		expectedHigh uint8
		expectedLow  uint8
	}{
		{"Pop zero", 0x0000, 0x00, 0x00},
		{"Pop typical value", 0x1234, 0x12, 0x34},
		{"Pop max value", 0xFFFF, 0xFF, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = 0xFFF0
			initialSP := cpu.Registers().SP

			// Write the value to pop in little-endian order
			high, low := helpers.FromRegisterPair(tt.stackValue)
			cpu.Write(initialSP, low)
			cpu.Write(initialSP+1, high)

			cpu.Execute(Pop(HL))

			assert.Equal(t, tt.expectedHigh, high, "unexpected high byte value")
			assert.Equal(t, tt.expectedLow, low, "unexpected low byte value")
			assert.Equal(t, initialSP+2, cpu.Registers().SP, "SP should be incremented by 2")
		})
	}
}

func TestPush16(t *testing.T) {
	tests := []struct {
		name      string
		highValue uint8
		lowValue  uint8
	}{
		{"Push zero", 0x00, 0x00},
		{"Push typical value", 0x12, 0x34},
		{"Push max value", 0xFF, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = 0xFFF0
			initialSP := cpu.Registers().SP

			cpu.Registers().H = tt.highValue
			cpu.Registers().L = tt.lowValue

			cpu.Execute(Push(HL))

			assert.Equal(t, initialSP-2, cpu.Registers().SP, "SP should be decremented by 2")
			assert.Equal(t, tt.lowValue, cpu.Memory[cpu.Registers().SP], "unexpected low byte on stack")
			assert.Equal(t, tt.highValue, cpu.Memory[cpu.Registers().SP+1], "unexpected high byte on stack")
		})
	}
}

func TestPopAF(t *testing.T) {
	tests := []struct {
		name      string
		stackA    uint8
		stackF    uint8
		expectedA uint8
		expectedZ bool
		expectedN bool
		expectedH bool
		expectedC bool
	}{
		{"Pop zero flags", 0x42, 0x00, 0x42, false, false, false, false},
		{"Pop all flags set", 0x24, 0xF0, 0x24, true, true, true, true},
		{"Pop mixed flags", 0xFF, 0xA0, 0xFF, true, false, true, false},
		// Lower 4 bits of F are ignored in GB CPU
		{"Pop with lower bits set", 0x00, 0xFF, 0x00, true, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = 0xFFF0
			initialSP := cpu.Registers().SP

			// Write A and F to stack in little-endian order
			cpu.Memory[cpu.Registers().SP] = tt.stackF   // Low byte (F)
			cpu.Memory[cpu.Registers().SP+1] = tt.stackA // High byte (A)

			cpu.Execute(Pop(AF))

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.expectedN, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
			assert.Equal(t, initialSP+2, cpu.Registers().SP, "SP should be incremented by 2")
		})
	}
}

func TestLoadHLSPOffset(t *testing.T) {
	tests := []struct {
		name          string
		initialSP     uint16
		offset        int8
		expectedHL    uint16
		expectedZero  bool // always false per spec
		expectedSub   bool // always false
		expectedHalf  bool
		expectedCarry bool
	}{
		{
			name:          "Positive offset, no carry",
			initialSP:     0x1000,
			offset:        5,
			expectedHL:    0x1005,
			expectedZero:  false,
			expectedSub:   false,
			expectedHalf:  false,
			expectedCarry: false,
		},
		{
			name:         "Lower nibble overflow sets half-carry",
			initialSP:    0x0FFF,
			offset:       1,
			expectedHL:   0x1000,
			expectedZero: false,
			expectedSub:  false,
			// 0xF+0x1=0x10 → half-carry and carry set.
			expectedHalf:  true,
			expectedCarry: true,
		},
		{
			name:          "Carry from lower byte addition",
			initialSP:     0xFFFF,
			offset:        1,
			expectedHL:    0x0000,
			expectedZero:  false,
			expectedSub:   false,
			expectedHalf:  true,
			expectedCarry: true,
		},
		{
			name:          "Negative offset causes borrow",
			initialSP:     0x1000,
			offset:        -5,
			expectedHL:    0x0FFB,
			expectedZero:  false,
			expectedSub:   false,
			expectedHalf:  true,
			expectedCarry: true,
		},
		{
			name:          "Zero offset: no flags",
			initialSP:     0x0500,
			offset:        0,
			expectedHL:    0x0500,
			expectedZero:  false,
			expectedSub:   false,
			expectedHalf:  false,
			expectedCarry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = tt.initialSP

			cpu.Registers().A = uint8(tt.offset)

			cpu.Execute(LoadHLSPOffset(A))

			actualHL := helpers.ToRegisterPair(cpu.Registers().H, cpu.Registers().L)
			assert.Equal(t, tt.expectedHL, actualHL, "unexpected HL value")
			assert.Equal(t, tt.expectedZero, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.expectedSub, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedHalf, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedCarry, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}
