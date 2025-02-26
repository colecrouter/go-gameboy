package lr35902

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
)

// Helpers
func (c *LR35902) GetImmediate8() uint8 {
	val := c.Read(c.registers.PC + 1)

	<-c.clock        // Use an additional m-cycle to read the immediate value
	c.registers.PC++ // Increment the program counter to the next instruction

	return val
}

func (c *LR35902) GetImmediate16() uint16 {
	high, low := c.Read16(c.registers.PC + 1)

	<-c.clock           // Use an additional m-cycle to read the immediate value
	c.registers.PC += 2 // Increment the program counter to the next instruction

	return cpu.ToRegisterPair(high, low)
}

func (c *LR35902) Read(addr uint16) byte {
	val := c.Read(addr)
	c.registers.PC++
	<-c.clock
	return val
}

func (c *LR35902) Read16(addr uint16) (high, low uint8) {
	high, low = c.Read16(addr)
	c.registers.PC += 2
	<-c.clock
	<-c.clock
	return high, low
}

// Write writes a byte to the given address
func (c *LR35902) Write(addr uint16, val byte) {
	c.Write(addr, val)
	c.registers.PC++
	<-c.clock
}

// Write16 writes a 16-bit value to the given address
func (c *LR35902) Write16(addr uint16, val uint16) {
	c.Write16(addr, val)
	c.registers.PC += 2
	<-c.clock
	<-c.clock
}

// Clock increments the program counter and waits for the next clock cycle
func (c *LR35902) Clock() {
	c.registers.PC++
	<-c.clock
}

// Halt halts the CPU until an interrupt is received
func (c *LR35902) Halt() {
	c.halted = true
}

// Stop halts the CPU until a button is pressed
func (c *LR35902) Stop() {
	c.Stop()
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
