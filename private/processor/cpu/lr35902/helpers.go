package lr35902

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

func (c LR35902) ClockAndAck() {
	<-c.clock
	c.clockAck <- struct{}{}
}

// Helpers
func (c *LR35902) GetImmediate8() uint8 {
	c.Clock()
	c.Registers().PC++
	val := c.bus.Read(c.registers.PC)
	c.Ack()

	return val
}

func (c *LR35902) GetImmediate16() uint16 {
	c.Clock()
	c.Registers().PC++
	low := c.bus.Read(c.registers.PC)
	c.Ack()

	c.Clock()
	c.Registers().PC++
	high := c.bus.Read(c.registers.PC)
	c.Ack()

	return cpu.ToRegisterPair(high, low)
}

func (c *LR35902) Read(addr uint16) byte {
	val := c.bus.Read(addr)

	return val
}

// Write writes a byte to the given address
func (c *LR35902) Write(addr uint16, val byte) {
	c.bus.Write(addr, val)
}

// Clock waits for the next clock cycle
func (c *LR35902) Clock() {
	<-c.clock
}

// Ack acknowledges the current clock cycle
func (c *LR35902) Ack() {
	c.clockAck <- struct{}{}
}

// Halt halts the CPU until an interrupt is received
func (c *LR35902) Halt() {
	c.halted = true
}

// Stop halts the CPU until a button is pressed
func (c *LR35902) Stop() {
	panic("not implemented")
}

// EI enables interrupts
func (c *LR35902) EI() {
	c.eiDelay = 0
	c.ime = true
}

// EIWithDelay enables interrupts after a delay
func (c *LR35902) EIWithDelay() {
	c.eiDelay = 2
}

// DI disables interrupts
func (c *LR35902) DI() {
	c.ime = false
}

// PrefixCB sets the CPU to use the CB instruction set
func (c *LR35902) PrefixCB() {
	c.cb = true
}

// Flags returns the CPU's flags
func (c *LR35902) Flags() *flags.Flags {
	return &c.flags
}

// Registers returns the CPU's registers
func (c *LR35902) Registers() *registers.Registers {
	return &c.registers
}
