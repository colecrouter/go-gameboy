package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

// MockCPU implements a minimal CPU for testing, with internal registers and flags.
type MockCPU struct {
	// Internal state similar to the real CPU.
	regs        registers.Registers
	flgs        flags.Flags
	Memory      []byte
	ClockCalled bool

	// Immediate operands (for testing GetImmediate8/16).
	Immediate8  uint8
	Immediate16 uint16

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
		flgs:   flags.Flags{},
		Memory: make([]byte, 0x10000),
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

// Clock simulates waiting for a clock cycle.
func (m *MockCPU) Clock() {
	m.ClockCalled = true
}

// Ack acknowledges a clock cycle.
func (m *MockCPU) Ack() {
	m.ClockCalled = false
}

// Read returns the byte stored in memory at the given address.
func (m *MockCPU) Read(addr uint16) uint8 {
	return m.Memory[addr]
}

// Write stores a byte in memory at the given address.
func (m *MockCPU) Write(addr uint16, val uint8) {
	m.Memory[addr] = val
}

// Read16 returns two bytes from memory in little-endian order.
func (m *MockCPU) Read16(addr uint16) (uint8, uint8) {
	return m.Memory[addr+1], m.Memory[addr]
}

// Write16 writes a 16-bit value to memory in little-endian order.
func (m *MockCPU) Write16(addr uint16, val uint16) {
	m.Memory[addr] = uint8(val & 0xFF)
	m.Memory[addr+1] = uint8(val >> 8)
}

// GetImmediate8 returns the immediate 8-bit operand.
func (m *MockCPU) GetImmediate8() uint8 {
	return m.Immediate8
}

// GetImmediate16 returns the immediate 16-bit operand.
func (m *MockCPU) GetImmediate16() uint16 {
	return m.Immediate16
}

// FromRegisterPair splits a 16-bit value into two 8-bit registers.
func (m *MockCPU) FromRegisterPair(val uint16) (uint8, uint8) {
	return uint8(val >> 8), uint8(val & 0xFF)
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
