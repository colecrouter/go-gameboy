package generators

import (
	"testing"

	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/conditions"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
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
		{"Jump when condition true", 0x0100, 0x1234, true, 0x1235},
		{"Skip when condition false", 0x0100, 0x1234, false, 0x0104},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC

			// Set immediate data
			cpu.Memory[tt.initialPC+1], cpu.Memory[tt.initialPC+2] = helpers.FromRegisterPair(tt.jumpAddress)

			// Set condition
			cpu.Flags().Zero = tt.condition

			cpu.Execute(Jump(conditions.Z))

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
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
		{"Relative jump forward when true", 0x0100, 10, true, 0x010D},
		{"Relative jump backward when true", 0x0100, -5, true, 0x00FE},
		{"Skip when condition false", 0x0100, 10, false, 0x0103},
		{"Skip when condition false (neg)", 0x0100, -5, false, 0x0103},
		{"Zero offset jump when true", 0x0100, 0, true, 0x0103},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC

			// Set condition
			cpu.Flags().Zero = tt.condition

			// Set immediate data
			cpu.Memory[tt.initialPC+1] = uint8(tt.offset)

			cpu.Execute(JumpRelative(conditions.Z))

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
		expectedPC uint16
	}{
		{"Return when condition true", 0x0100, 0x1234, true, 0x1234},
		{"Skip when condition false", 0x0100, 0x1234, false, 0x0102},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC
			cpu.Registers().SP = 0xFFF0

			// Write return address to stack
			cpu.Memory[cpu.Registers().SP+1], cpu.Memory[cpu.Registers().SP] = helpers.FromRegisterPair(tt.stackAddress)

			// Set condition
			cpu.Flags().Zero = tt.condition

			cpu.Execute(Return(conditions.Z))

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
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
		{"RST 0x00", 0x0100, 0xFFF0, 0x0000, 0x00, 0xFFEE},
		{"RST 0x08", 0x0200, 0xFFF0, 0x0008, 0x08, 0xFFEE},
		{"RST 0x10", 0x0300, 0xFFF0, 0x0010, 0x10, 0xFFEE},
		{"RST 0x18", 0x0400, 0xFFF0, 0x0018, 0x18, 0xFFEE},
		{"RST 0x20", 0x0500, 0xFFF0, 0x0020, 0x20, 0xFFEE},
		{"RST 0x28", 0x0600, 0xFFF0, 0x0028, 0x28, 0xFFEE},
		{"RST 0x30", 0x0700, 0xFFF0, 0x0030, 0x30, 0xFFEE},
		{"RST 0x38", 0x0800, 0xFFF0, 0x0038, 0x38, 0xFFEE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newMockCPU()
			cpu.Registers().PC = tt.initialPC
			cpu.Registers().SP = tt.initialSP

			cpu.Execute(ResetPC(tt.rstAddress))

			assert.Equal(t, tt.expectedPC, cpu.Registers().PC, "unexpected PC value")
			assert.Equal(t, tt.expectedSP, cpu.Registers().SP, "unexpected SP value")

			// Check that return address was pushed correctly
			pushedLow := cpu.Memory[cpu.Registers().SP]
			pushedHigh := cpu.Memory[cpu.Registers().SP+1]
			returnAddr := helpers.ToRegisterPair(pushedHigh, pushedLow)
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

		// Set immediate data.
		cpu.Memory[cpu.Registers().PC+1] = 0x34
		cpu.Memory[cpu.Registers().PC+2] = 0x12

		cpu.Execute(Call(conditions.Always))

		expectedRetAddr := uint16(0x0103)
		expectedPC := uint16(0x1235)
		expectedSP := uint16(0xFFEE)

		assert.Equal(t, expectedPC, cpu.Registers().PC, "unexpected PC value")
		assert.Equal(t, expectedSP, cpu.Registers().SP, "unexpected SP value")

		// Check that the return address was pushed correctly.
		pushedLow := cpu.Memory[cpu.Registers().SP]
		pushedHigh := cpu.Memory[cpu.Registers().SP+1]
		returnAddr := helpers.ToRegisterPair(pushedHigh, pushedLow)
		assert.Equal(t, expectedRetAddr, returnAddr, "incorrect return address pushed to stack")
	})

	// Test call when condition is false.
	t.Run("Call when condition false", func(t *testing.T) {
		cpu := newMockCPU()
		// Setup initial PC and SP.
		cpu.Registers().PC = 0x0100
		cpu.Registers().SP = 0xFFF0

		cpu.Flags().Zero = false

		cpu.Execute(Call(conditions.Z))

		// Should leave PC and SP unchanged.
		assert.Equal(t, uint16(0x0104), cpu.Registers().PC, "unexpected PC value")
		assert.Equal(t, uint16(0xFFF0), cpu.Registers().SP, "SP should remain unchanged")
	})
}
