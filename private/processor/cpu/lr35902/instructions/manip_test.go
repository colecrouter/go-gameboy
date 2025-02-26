package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInc8(t *testing.T) {
	tests := []struct {
		name      string
		initial   uint8
		expected  uint8
		expectedZ bool
		expectedH bool
	}{
		{"Increment from zero", 0x00, 0x01, false, false},
		{"Increment to zero", 0xFF, 0x00, true, true},
		{"Half-carry case", 0x0F, 0x10, false, true},
		{"General case", 0x42, 0x43, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			r := tt.initial

			inc8(cpu, &r)

			assert.Equal(t, tt.expected, r, "unexpected register value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, false, cpu.Flags().Subtract, "Subtract flag should be reset")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected Half-carry flag")
		})
	}
}

func TestInc16(t *testing.T) {
	tests := []struct {
		name         string
		initialHigh  uint8
		initialLow   uint8
		expectedHigh uint8
		expectedLow  uint8
	}{
		{"Increment from zero", 0x00, 0x00, 0x00, 0x01},
		{"Low byte overflow", 0x00, 0xFF, 0x01, 0x00},
		{"Increment to max", 0xFF, 0xFE, 0xFF, 0xFF},
		{"Increment from max", 0xFF, 0xFF, 0x00, 0x00},
		{"General case", 0x12, 0x34, 0x12, 0x35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			high := tt.initialHigh
			low := tt.initialLow

			inc16(cpu, &high, &low)

			assert.Equal(t, tt.expectedHigh, high, "unexpected high byte value")
			assert.Equal(t, tt.expectedLow, low, "unexpected low byte value")
		})
	}
}

func TestDec8(t *testing.T) {
	tests := []struct {
		name      string
		initial   uint8
		expected  uint8
		expectedZ bool
		expectedH bool
	}{
		{"Decrement to zero", 0x01, 0x00, true, false},
		{"Decrement from zero", 0x00, 0xFF, false, true},
		{"Half-carry case", 0x10, 0x0F, false, true},
		{"General case", 0x43, 0x42, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			r := tt.initial

			dec8(cpu, &r)

			assert.Equal(t, tt.expected, r, "unexpected register value")
			assert.Equal(t, tt.expectedZ, cpu.Flags().Zero, "unexpected Zero flag")
			assert.Equal(t, true, cpu.Flags().Subtract, "Subtract flag should be set")
			assert.Equal(t, tt.expectedH, cpu.Flags().HalfCarry, "unexpected Half-carry flag")
		})
	}
}

func TestDec16(t *testing.T) {
	tests := []struct {
		name         string
		initialHigh  uint8
		initialLow   uint8
		expectedHigh uint8
		expectedLow  uint8
	}{
		{"Decrement to zero", 0x00, 0x01, 0x00, 0x00},
		{"Low byte underflow", 0x01, 0x00, 0x00, 0xFF},
		{"Decrement to almost max", 0x00, 0x00, 0xFF, 0xFF},
		{"General case", 0x12, 0x34, 0x12, 0x33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			high := tt.initialHigh
			low := tt.initialLow

			dec16(cpu, &high, &low)

			assert.Equal(t, tt.expectedHigh, high, "unexpected high byte value")
			assert.Equal(t, tt.expectedLow, low, "unexpected low byte value")
		})
	}
}
