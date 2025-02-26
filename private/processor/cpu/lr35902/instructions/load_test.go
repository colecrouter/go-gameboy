package instructions

import (
	"testing"

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

			load8(cpu, &r, tt.value)

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
			var high, low uint8

			load16(cpu, &high, &low, tt.value)

			assert.Equal(t, tt.expectedH, high, "unexpected high byte value")
			assert.Equal(t, tt.expectedL, low, "unexpected low byte value")
		})
	}
}

func TestLoadHLSPOffset(t *testing.T) {
	tests := []struct {
		name       string
		initialSP  uint16
		offset     int8
		expectedHL uint16
		expectedZ  bool
		expectedN  bool
		expectedH  bool
		expectedC  bool
	}{
		{"Positive offset", 0x1000, 5, 0x1005, false, false, false, false},
		{"Half carry", 0x0FFF, 1, 0x1000, false, false, true, true},
		{"Carry", 0xFFFF, 1, 0x0000, false, false, true, true},
		// a borrow occurs from both the lower nibble and lower byte,
		// so both HalfCarry and Carry should be set.
		{"Negative offset", 0x1000, -5, 0x0FFB, false, false, true, true},
		{"Zero offset", 0x0500, 0, 0x0500, false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = tt.initialSP

			loadHLSPOffset(cpu, tt.offset)

			actualHL := (uint16(cpu.Registers().H) << 8) | uint16(cpu.Registers().L)
			assert.Equal(t, tt.expectedHL, actualHL, "unexpected HL value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.expectedN, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
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
			cpu.Memory[cpu.Registers().SP] = uint8(tt.stackValue & 0xFF)          // Low byte
			cpu.Memory[cpu.Registers().SP+1] = uint8((tt.stackValue >> 8) & 0xFF) // High byte

			var high, low uint8
			pop16(cpu, &high, &low)

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

			push16(cpu, tt.highValue, tt.lowValue)

			assert.Equal(t, initialSP-2, cpu.Registers().SP, "SP should be decremented by 2")
			assert.Equal(t, tt.lowValue, cpu.Memory[cpu.Registers().SP], "unexpected low byte on stack")
			assert.Equal(t, tt.highValue, cpu.Memory[cpu.Registers().SP+1], "unexpected high byte on stack")
		})
	}
}

func TestLoad8Mem(t *testing.T) {
	tests := []struct {
		name    string
		value   uint8
		address uint16
	}{
		{"Store to zero address", 0x42, 0x0000},
		{"Store to high address", 0x7F, 0xFF80},
		{"Store zero value", 0x00, 0x1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()

			load8Mem(cpu, tt.value, tt.address)

			assert.Equal(t, tt.value, cpu.Memory[tt.address], "value not stored correctly at address")
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

			popAF(cpu)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.expectedN, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
			assert.Equal(t, initialSP+2, cpu.Registers().SP, "SP should be incremented by 2")
		})
	}
}
