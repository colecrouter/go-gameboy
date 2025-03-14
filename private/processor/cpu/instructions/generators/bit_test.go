package generators

import (
	"testing"

	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/enums"
	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	// Test cases for left rotation
	t.Run("RLCA - rotate A left (no carry)", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x80 // 10000000

		// Rotate left, no carry bit, don't update Z flag
		// Rotate(cpu, &cpu.Registers().A, true, false, false)
		cpu.Execute(Rotate(A, true, false, false))

		assert.Equal(t, uint8(0x01), cpu.Registers().A, "A should rotate left with bit 7 to bit 0")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 7 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	t.Run("RLA - rotate A left through carry", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x80 // 10000000
		cpu.Flags().Carry = true // Initial carry is set

		// Rotate left, use carry bit, don't update Z flag
		cpu.Execute(Rotate(A, true, true, false))

		assert.Equal(t, uint8(0x01), cpu.Registers().A, "A should rotate left with carry into bit 0")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 7 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	// Test cases for right rotation
	t.Run("RRCA - rotate A right (no carry)", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x01 // 00000001

		// Rotate right, no carry bit, don't update Z flag
		cpu.Execute(Rotate(A, false, false, false))

		assert.Equal(t, uint8(0x80), cpu.Registers().A, "A should rotate right with bit 0 to bit 7")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 0 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	t.Run("RRA - rotate A right through carry", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x01 // 00000001
		cpu.Flags().Carry = true // Initial carry is set

		// Rotate right, use carry bit, don't update Z flag
		cpu.Execute(Rotate(A, false, true, false))

		assert.Equal(t, uint8(0x80), cpu.Registers().A, "A should rotate right with carry into bit 7")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 0 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	// CB-prefixed rotations (with Z flag update)
	t.Run("RLC B - CB prefixed rotate left", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().B = 0x80 // 10000000

		// Rotate left, no carry bit, but update Z flag (CB prefix behavior)
		cpu.Execute(Rotate(B, true, false, true))

		assert.Equal(t, uint8(0x01), cpu.Registers().B, "B should rotate left")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset (result not zero)")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	t.Run("RLC B - CB prefixed rotate left with zero result", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().B = 0x00 // 00000000

		// Rotate left, no carry bit, update Z flag (CB prefix behavior)
		cpu.Execute(Rotate(B, true, false, true))

		assert.Equal(t, uint8(0x00), cpu.Registers().B, "B should remain zero")
		assert.False(t, cpu.Flags().Carry, "Carry flag should be reset (bit 7 was 0)")
		assert.True(t, cpu.Flags().Zero, "Zero flag should be set (result is zero)")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})
}

func TestShift(t *testing.T) {
	t.Run("SLA - shift left arithmetic", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x80 // 10000000

		// Shift left, not arithmetic right (parameter is ignored for left shifts)
		cpu.Execute(Shift(A, true, false))

		assert.Equal(t, uint8(0x00), cpu.Registers().A, "A should be shifted left with 0 into bit 0")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 7 was 1)")
		assert.True(t, cpu.Flags().Zero, "Zero flag should be set (result is zero)")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	t.Run("SRL - shift right logical", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x81 // 10000001

		// Shift right, not arithmetic (logical)
		cpu.Execute(Shift(A, false, false))

		assert.Equal(t, uint8(0x40), cpu.Registers().A, "A should be shifted right logically")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 0 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset (result not zero)")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})

	t.Run("SRA - shift right arithmetic", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x81 // 10000001

		// Shift right, arithmetic (preserves sign bit)
		cpu.Execute(Shift(A, false, true))

		assert.Equal(t, uint8(0xC0), cpu.Registers().A, "A should be shifted right with MSB preserved")
		assert.True(t, cpu.Flags().Carry, "Carry flag should be set (bit 0 was 1)")
		assert.False(t, cpu.Flags().Zero, "Zero flag should be reset (result not zero)")
		assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
		assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
	})
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name         string
		value        uint8
		expected     uint8
		expectedZero bool
	}{
		{"Swap non-zero", 0x12, 0x21, false},
		{"Swap zero", 0x00, 0x00, true},
		{"Swap mixed nibbles", 0xAB, 0xBA, false},
		{"Swap with zero nibble", 0x0F, 0xF0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().B = tt.value

			cpu.Execute(Swap(B))

			assert.Equal(t, tt.expected, cpu.Registers().B, "unexpected result value")
			assert.Equal(t, tt.expectedZero, cpu.Flags().Zero, "unexpected Zero flag")
			assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.False(t, cpu.Flags().HalfCarry, "Half-carry flag should be reset")
			assert.False(t, cpu.Flags().Carry, "Carry flag should be reset")
		})
	}
}

func TestCP8(t *testing.T) {
	tests := []struct {
		name         string
		aValue       uint8
		compareValue uint8
		expectedZ    bool // should be true only when values are equal
		expectedH    bool
		expectedC    bool
	}{
		{"Equal values", 0x42, 0x42, true, false, false},
		{"A > compare value", 0xFF, 0x01, false, false, false},
		// Updated: When A < compared value, Zero flag should be false.
		{"A < compare value", 0x01, 0xFF, false, true, true},
		// Updated: When result is non-zero, Zero flag must be false.
		{"Half-carry check", 0x20, 0x21, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().A = tt.aValue
			cpu.Registers().B = tt.compareValue
			// Execute Compare using the compare value only
			cpu.Execute(Compare(A, B))

			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			// CP always sets Subtract flag
			assert.True(t, cpu.Flags().Subtract, "Subtract flag should be set")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected Half-carry flag")
			assert.Equal(t, tt.expectedC, cpu.Flags().Carry, "unexpected Carry flag")
		})
	}
}

func TestBit(t *testing.T) {
	tests := []struct {
		name         string
		value        uint8
		bitPosition  uint8
		expectedZero bool
	}{
		{"Test bit 0 set", 0x01, 0, false},
		{"Test bit 0 reset", 0xFE, 0, true},
		{"Test bit 7 set", 0x80, 7, false},
		{"Test bit 7 reset", 0x7F, 7, true},
		{"Test middle bit set", 0x08, 3, false},
		{"Test middle bit reset", 0xF7, 3, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newMockCPU()
			// Use register A for testing.
			cpu.Registers().A = tc.value
			// Clear flags to avoid previous side-effects.
			cpu.Flags().Zero = false
			cpu.Flags().Carry = false

			cpu.Execute(ReadBit(A, tc.bitPosition))

			assert.Equal(t, tc.expectedZero, cpu.Flags().Zero, "unexpected Zero flag")
		})
	}
}

func TestRes(t *testing.T) {
	tests := []struct {
		name          string
		initialValue  uint8
		bitPosition   uint8
		expectedValue uint8
	}{
		{"Reset bit 0", 0xFF, 0, 0xFE},
		{"Reset bit 7", 0xFF, 7, 0x7F},
		{"Reset already reset bit", 0xF7, 3, 0xF7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newMockCPU()
			// Use register A as the operand.
			cpu.Registers().A = tc.initialValue
			cpu.Execute(ResetBit(A, tc.bitPosition))
			assert.Equal(t, tc.expectedValue, cpu.Registers().A, "unexpected result value")
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name          string
		initialValue  uint8
		bitPosition   uint8
		expectedValue uint8
	}{
		{"Set bit 0", 0x00, 0, 0x01},
		{"Set bit 7", 0x00, 7, 0x80},
		{"Set already set bit", 0x08, 3, 0x08},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newMockCPU()
			// Use register A as the operand.
			cpu.Registers().A = tc.initialValue
			cpu.Execute(SetBit(A, tc.bitPosition))
			assert.Equal(t, tc.expectedValue, cpu.Registers().A, "unexpected result value")
		})
	}
}
