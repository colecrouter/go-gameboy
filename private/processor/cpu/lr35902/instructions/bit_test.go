package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	// Test cases for left rotation
	t.Run("RLCA - rotate A left (no carry)", func(t *testing.T) {
		cpu := newMockCPU()
		cpu.Registers().A = 0x80 // 10000000

		// Rotate left, no carry bit, don't update Z flag
		rotate(cpu, &cpu.Registers().A, true, false, false)

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
		rotate(cpu, &cpu.Registers().A, true, true, false)

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
		rotate(cpu, &cpu.Registers().A, false, false, false)

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
		rotate(cpu, &cpu.Registers().A, false, true, false)

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
		rotate(cpu, &cpu.Registers().B, true, false, true)

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
		rotate(cpu, &cpu.Registers().B, true, false, true)

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
		shift(cpu, &cpu.Registers().A, true, false)

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
		shift(cpu, &cpu.Registers().A, false, false)

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
		shift(cpu, &cpu.Registers().A, false, true)

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

			swap(cpu, &cpu.Registers().B)

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
			cp8(cpu, tt.aValue, tt.compareValue)

			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			// CP always sets N
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Flags().Carry = true // Initial carry is set

			bit(cpu, tt.bitPosition, tt.value)

			assert.Equal(t, tt.expectedZero, cpu.Flags().Zero, "unexpected Zero flag")
			assert.False(t, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.True(t, cpu.Flags().HalfCarry, "Half-carry flag should be set")
			assert.True(t, cpu.Flags().Carry, "Carry flag should be preserved")
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			value := tt.initialValue

			res(cpu, tt.bitPosition, &value)

			assert.Equal(t, tt.expectedValue, value, "unexpected result value")
			// Flags are not affected by RES instruction
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			value := tt.initialValue

			set(cpu, tt.bitPosition, &value)

			assert.Equal(t, tt.expectedValue, value, "unexpected result value")
			// Flags are not affected by SET instruction
		})
	}
}
