package generators

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

// MockCPU implements a minimal CPU for testing, with internal registers and flags.
type MockCPU struct {
	// Internal state similar to the real CPU.
	regs        registers.Registers
	flgs        flags.Flags
	Memory      []byte
	ClockCalled bool

	// New clock channels for micro-op synchronization.
	clock    chan struct{}
	clockAck chan struct{}

	// Modes & status flags.
	Halted    bool
	Stopped   bool
	IME       bool
	EIDelayed bool
	CBPrefix  bool
}

// newMockCPU creates a new instance with initialized memory and default registers.
func newMockCPU() *MockCPU {
	return &MockCPU{
		regs: registers.Registers{
			PC: 0x0100, // typical starting PC for GameBoy
			SP: 0xFFFE,
		},
		flgs:     flags.Flags{},
		Memory:   make([]byte, 0x10000),
		clock:    make(chan struct{}, 100),
		clockAck: make(chan struct{}, 100),
	}
}

// Registers returns a pointer to the CPU's internal registers.
func (m *MockCPU) Registers() *registers.Registers {
	return &m.regs
}

// Flags returns a pointer to the CPU's flags.
func (m *MockCPU) Flags() *flags.Flags {
	return &m.flgs
}

// Read returns the byte stored in memory at the given address.
func (m *MockCPU) Read(addr uint16) uint8 {
	return m.Memory[addr]
}

// Write stores a byte in memory at the given address.
func (m *MockCPU) Write(addr uint16, val uint8) {
	m.Memory[addr] = val
}

// Halt puts the CPU in halted mode.
func (m *MockCPU) Halt() {
	m.Halted = true
}

// Stop puts the CPU in stopped mode.
func (m *MockCPU) Stop() {
	m.Stopped = true
}

// EI enables interrupts immediately.
func (m *MockCPU) EI() {
	m.IME = true
}

// EIWithDelay enables interrupts after a delay.
func (m *MockCPU) EIWithDelay() {
	m.EIDelayed = true
}

// DI disables interrupts.
func (m *MockCPU) DI() {
	m.IME = false
}

// PrefixCB sets the CB prefix mode.
func (m *MockCPU) PrefixCB() {
	m.CBPrefix = true
}

// Additional methods for testing.
func (m *MockCPU) Execute(ops shared.Instruction) {
	// Execute instruction
	ctx := &shared.Context{}

	// Simulate the PC increment of the last instruction.
	m.regs.PC++

	for _, op := range ops {
		extra := op(m, ctx)

		if extra != nil {
			for _, e := range *extra {
				e(m, ctx)
			}
		}
	}
}
