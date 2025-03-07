package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJump(t *testing.T) {
	tests := []struct {
		name        string
		initialPC   uint16
		jumpAddress uint16
		condition   bool
		// For taken jump, expected becomes jumpAddress; false leaves PC unchanged.
		expectedPC uint16
	}{
		{"Jump when condition true", 0x0100, 0x1234, true, 0x1233},
		{"Skip when condition false", 0x0100, 0x1234, false, 0x0100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC

			jump(cpu, tt.jumpAddress, tt.condition)

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
			// No flag assertions in jump; leave unchanged.
		})
	}
}

func TestJumpRelative(t *testing.T) {
	tests := []struct {
		name       string
		initialPC  uint16
		offset     int8
		condition  bool
		expectedPC uint16
	}{
		{"Relative jump forward when true", 0x0100, 10, true, 0x010A},
		{"Relative jump backward when true", 0x0100, -5, true, 0x00FB},
		{"Skip when condition false", 0x0100, 10, false, 0x0100},
		{"Skip when condition false (neg)", 0x0100, -5, false, 0x0100},
		{"Zero offset jump when true", 0x0100, 0, true, 0x0100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC

			jumpRelative(cpu, tt.offset, tt.condition)

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
		})
	}
}

func TestRet(t *testing.T) {
	tests := []struct {
		name         string
		initialPC    uint16
		stackAddress uint16
		condition    bool
		// For taken return: expected becomes the value from stack; false leaves PC unchanged.
		expectedPC  uint16
		clockCalled bool
	}{
		{"Return when condition true", 0x0100, 0x1234, true, 0x1233, false},
		{"Skip when condition false", 0x0100, 0x1234, false, 0x0100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC
			cpu.Registers().SP = 0xFFF0

			// Write return address to stack
			cpu.Memory[cpu.Registers().SP] = uint8(tt.stackAddress & 0xFF)
			cpu.Memory[cpu.Registers().SP+1] = uint8((tt.stackAddress >> 8) & 0xFF)

			ret(cpu, tt.condition)

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
			assert.Equal(t, tt.clockCalled, cpu.ClockCalled, "unexpected clock call")
			if tt.condition {
				assert.Equal(t, uint16(0xFFF2), cpu.Registers().SP, "SP should be incremented by 2 on return")
			} else {
				assert.Equal(t, uint16(0xFFF0), cpu.Registers().SP, "SP should remain unchanged when condition false")
			}
		})
	}
}

func TestRst(t *testing.T) {
	tests := []struct {
		name       string
		initialPC  uint16
		initialSP  uint16
		rstAddress uint16
		// Updated: expected PC becomes (rstAddress – 1) modulo 16-bit wrap
		expectedPC uint16
		expectedSP uint16
	}{
		{"RST 0x00", 0x0100, 0xFFF0, 0x0000, 0xFFFF, 0xFFEE},
		{"RST 0x08", 0x0200, 0xFFF0, 0x0008, 0x0007, 0xFFEE},
		{"RST 0x10", 0x0300, 0xFFF0, 0x0010, 0x000F, 0xFFEE},
		{"RST 0x18", 0x0400, 0xFFF0, 0x0018, 0x0017, 0xFFEE},
		{"RST 0x20", 0x0500, 0xFFF0, 0x0020, 0x001F, 0xFFEE},
		{"RST 0x28", 0x0600, 0xFFF0, 0x0028, 0x0027, 0xFFEE},
		{"RST 0x30", 0x0700, 0xFFF0, 0x0030, 0x002F, 0xFFEE},
		{"RST 0x38", 0x0800, 0xFFF0, 0x0038, 0x0037, 0xFFEE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC
			cpu.Registers().SP = tt.initialSP

			rst(cpu, tt.rstAddress)

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
			assert.Equal(t, tt.expectedSP, cpu.Registers().SP, "unexpected SP value")

			// Check that return address was pushed correctly
			pushedLow := cpu.Memory[cpu.Registers().SP]
			pushedHigh := cpu.Memory[cpu.Registers().SP+1]
			returnAddr := (uint16(pushedHigh) << 8) | uint16(pushedLow)
			assert.Equal(t, tt.initialPC+1, returnAddr, "incorrect return address pushed to stack")
		})
	}
}

func TestCall(t *testing.T) {
	// Test call when condition is true.
	t.Run("Call when condition true", func(t *testing.T) {
		cpu := newMockCPU()
		// Setup initial PC and SP.
		cpu.Registers().PC = 0x0100
		cpu.Registers().SP = 0xFFF0

		// Use call to branch to 0x1234 when condition true.
		call(cpu, 0x1234, true)

		// Expected return address is initial PC + 1 (we had to read the immediate operand).
		expectedRetAddr := uint16(0x0101)
		// The target is adjusted by subtracting 1: expected PC = target - 1.
		expectedPC := uint16(0x1234 - 1)
		// SP should be decreased by 2.
		expectedSP := uint16(0xFFF0 - 2)

		assert.Equal(t, expectedPC, cpu.Registers().PC, "unexpected PC value")
		assert.Equal(t, expectedSP, cpu.Registers().SP, "unexpected SP value")

		// Check that the return address was pushed correctly.
		pushedLow := cpu.Memory[cpu.Registers().SP]
		pushedHigh := cpu.Memory[cpu.Registers().SP+1]
		returnAddr := (uint16(pushedHigh) << 8) | uint16(pushedLow)
		assert.Equal(t, expectedRetAddr, returnAddr, "incorrect return address pushed to stack")
	})

	// Test call when condition is false.
	t.Run("Call when condition false", func(t *testing.T) {
		cpu := newMockCPU()
		// Setup initial PC and SP.
		cpu.Registers().PC = 0x0100
		cpu.Registers().SP = 0xFFF0

		// When condition is false, call should leave PC and SP unchanged.
		call(cpu, 0x1234, false)
		assert.Equal(t, uint16(0x0100), cpu.Registers().PC, "unexpected PC value")
		assert.Equal(t, uint16(0xFFF0), cpu.Registers().SP, "SP should remain unchanged")
	})
}
