package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd8(t *testing.T) {
	tests := []struct {
		name       string
		initialA   uint8
		valueToAdd uint8
		expectedA  uint8
		expectedZ  bool
		expectedN  bool
		expectedH  bool
		expectedC  bool
	}{
		{"Simple addition", 1, 2, 3, false, false, false, false},
		{"Zero result", 0, 0, 0, true, false, false, false},
		{"Half carry", 0x0F, 0x01, 0x10, false, false, true, false},
		{"Carry", 0xFF, 0x01, 0x00, true, false, true, true},
		{"Half carry and carry", 0xFF, 0x02, 0x01, false, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA

			add8(cpu, &cpu.Registers().A, tt.valueToAdd)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, tt.expectedN, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestAdd16(t *testing.T) {
	tests := []struct {
		name       string
		initialHL  uint16
		valueToAdd uint16
		expectedHL uint16
		expectedN  bool
		expectedH  bool
		expectedC  bool
	}{
		{"Simple addition", 0x0001, 0x0001, 0x0002, false, false, false},
		{"Half carry", 0x0FFF, 0x0001, 0x1000, false, true, false},
		{"Carry", 0xFFFF, 0x0001, 0x0000, false, true, true},
		{"Typical case", 0x1234, 0x4321, 0x5555, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().H, cpu.Registers().L = uint8(tt.initialHL>>8), uint8(tt.initialHL&0xFF)
			highVal, lowVal := uint8(tt.valueToAdd>>8), uint8(tt.valueToAdd&0xFF)

			// Set initial flags to make sure they're properly modified
			cpu.Flags().Zero = true // should remain unchanged

			add16(cpu, &cpu.Registers().H, &cpu.Registers().L, highVal, lowVal)

			actualHL := (uint16(cpu.Registers().H) << 8) | uint16(cpu.Registers().L)
			assert.Equal(t, tt.expectedHL, actualHL, "unexpected HL value")
			assert.Equal(t, true, cpu.Flags().Zero, "Zero flag should remain unchanged")
			assert.Equal(t, tt.expectedN, cpu.Flags().Subtract, "unexpected Subtract flag")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestSub8(t *testing.T) {
	tests := []struct {
		name       string
		initialA   uint8
		valueToSub uint8
		expectedA  uint8
		expectedZ  bool
		expectedH  bool
		expectedC  bool
	}{
		{"Simple subtraction", 5, 3, 2, false, false, false},
		{"Zero result", 3, 3, 0, true, false, false},
		{"Half borrow", 0x10, 0x01, 0x0F, false, true, false},
		{"Borrow", 0, 1, 0xFF, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA

			sub8(cpu, &cpu.Registers().A, tt.valueToSub)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, true, cpu.Flags().Subtract, "Subtract flag should be set")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestAnd8(t *testing.T) {
	tests := []struct {
		name       string
		initialA   uint8
		valueToAnd uint8
		expectedA  uint8
		expectedZ  bool
	}{
		{"Simple AND", 0x55, 0xF0, 0x50, false},
		{"Zero result", 0x55, 0x0A, 0x00, true},
		{"Full mask", 0xFF, 0xFF, 0xFF, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA

			and8(cpu, &cpu.Registers().A, tt.valueToAnd)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, true, cpu.Flags().HalfCarry, "HalfCarry flag should be set")
			assert.Equal(t, false, cpu.Flags().Carry, "Carry flag should be reset")
		})
	}
}

func TestOr8(t *testing.T) {
	tests := []struct {
		name      string
		initialA  uint8
		valueToOr uint8
		expectedA uint8
		expectedZ bool
	}{
		{"Simple OR", 0x55, 0xAA, 0xFF, false},
		{"Zero result", 0x00, 0x00, 0x00, true},
		{"No change", 0xFF, 0x55, 0xFF, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA

			or8(cpu, &cpu.Registers().A, tt.valueToOr)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, false, cpu.Flags().HalfCarry, "HalfCarry flag should be reset")
			assert.Equal(t, false, cpu.Flags().Carry, "Carry flag should be reset")
		})
	}
}

func TestXor8(t *testing.T) {
	tests := []struct {
		name       string
		initialA   uint8
		valueToXor uint8
		expectedA  uint8
		expectedZ  bool
	}{
		{"Simple XOR", 0x55, 0xAA, 0xFF, false},
		{"Zero result", 0x55, 0x55, 0x00, true},
		{"No change", 0xFF, 0x00, 0xFF, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA

			xor8(cpu, &cpu.Registers().A, tt.valueToXor)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, false, cpu.Flags().HalfCarry, "HalfCarry flag should be reset")
			assert.Equal(t, false, cpu.Flags().Carry, "Carry flag should be reset")
		})
	}
}

func TestAddc8(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		valueToAdd   uint8
		initialCarry bool
		expectedA    uint8
		expectedZ    bool
		expectedH    bool
		expectedC    bool
	}{
		{"Add with carry=0", 1, 2, false, 3, false, false, false},
		{"Add with carry=1", 1, 2, true, 4, false, false, false},
		{"Half carry", 0x0F, 0x01, false, 0x10, false, true, false},
		{"Half carry with carry=1", 0x0F, 0x00, true, 0x10, false, true, false},
		{"Full carry", 0xFF, 0x01, false, 0x00, true, true, true},
		{"Full carry with carry=1", 0xFE, 0x01, true, 0x00, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA
			cpu.Flags().Carry = tt.initialCarry

			addc8(cpu, &cpu.Registers().A, tt.valueToAdd)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestSubc8(t *testing.T) {
	tests := []struct {
		name         string
		initialA     uint8
		valueToSub   uint8
		initialCarry bool
		expectedA    uint8
		expectedZ    bool
		expectedH    bool
		expectedC    bool
	}{
		{"Sub with carry=0", 5, 3, false, 2, false, false, false},
		{"Sub with carry=1", 5, 3, true, 1, false, false, false},
		{"Zero result", 3, 3, false, 0, true, false, false},
		{"Zero result with carry=1", 4, 3, true, 0, true, false, false},
		{"Half borrow", 0x10, 0x01, false, 0x0F, false, true, false},
		{"Full borrow", 0, 1, false, 0xFF, false, true, true},
		{"Full borrow with carry=1", 0, 0, true, 0xFF, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.initialA
			cpu.Flags().Carry = tt.initialCarry

			subc8(cpu, &cpu.Registers().A, tt.valueToSub)

			assert.Equal(t, tt.expectedA, cpu.Registers().A, "unexpected A value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, true, cpu.Flags().Subtract, "Subtract flag should be set")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestAddSPr8(t *testing.T) {
	tests := []struct {
		name       string
		initialSP  uint16
		offset     int8
		expectedSP uint16
		expectedH  bool
		expectedC  bool
	}{
		{"Simple addition", 0x1000, 5, 0x1005, false, false},
		{"Half carry", 0x100F, 1, 0x1010, true, false},
		{"Carry", 0x10FF, 1, 0x1100, true, true},
		{"Negative offset", 0x1000, -3, 0x0FFD, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().SP = tt.initialSP
			cpu.Immediate8 = uint8(tt.offset)

			addSPr8(cpu)

			assert.Equal(t, tt.expectedSP, cpu.Registers().SP, "unexpected SP value")
			assert.Equal(t, false, cpu.Flags().Zero, "Zero flag should be reset")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected HalfCarry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}
