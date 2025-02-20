package lr35902

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
)

// LR35902 is the original GameBoy CPU
type LR35902 struct {
	initialized bool
	Registers   struct {
		a, b, c, d, e, h, l uint8
		sp, PC              uint16
	}
	flags   Flags
	bus     *memory.Bus
	io      *registers.Registers
	ie      *registers.Interrupt
	cb      bool
	ime     bool
	eiDelay int
	lastPC  uint16
	halted  bool
}

// Step executes the next instruction in the CPU's memory.
// Returns the number of clock cycles the instruction took.
func (c *LR35902) Step() int {
	if !c.initialized {
		panic("CPU not initialized")
	}

	// Check for interrupts
	for i := VBlankISR; i <= JoypadISR; i++ {
		// TODO gross
		if c.ie.Read(0)&(1<<isrOffsets[i]) != 0 && c.io.InterruptFlag.Read(0)&(1<<isrOffsets[i]) != 0 {
			if c.ime {
				c.isr(i)

				// Clear interrupt flag
				c.io.InterruptFlag.Write(0, c.io.InterruptFlag.Read(0)&^(1<<isrOffsets[i]))
			}

			// Cancel HALT mode
			c.halted = false
		}
	}

	// Check for HALT mode
	if c.halted {
		// We need to return 4 cycles here so that the CPU still runs and checks for interrupts
		// If we return 0, the process will hang, as it will continue to clock empty cycles without stopping or checking for interrupts
		return 4
	}

	// Get instruction
	opcode := c.bus.Read(c.Registers.PC)

	// Run instruction
	var instruction instruction
	var mnemonic string
	if c.cb {
		instruction = cbInstructions[opcode]
		mnemonic = getCBMnemonic(opcode)
		c.cb = false
	} else {
		instruction = instructions[opcode]
		mnemonic = mnemonics[opcode]
	}

	if c.Registers.PC == 0xc2f6 {
		fmt.Printf("")
	}

	_ = mnemonic

	op := instruction.op
	cycles := instruction.c
	increment := instruction.p

	c.lastPC = c.Registers.PC
	if op == nil {
		fmt.Printf("Unimplemented instruction: 0x%02X\r\n", opcode)
		c.Registers.PC++
	} else {
		op(c)
		c.Registers.PC += uint16(increment)
	}

	// Delay EI effect: Decrement counter at the very end of the instruction.
	if c.eiDelay > 0 {
		c.eiDelay--
		if c.eiDelay == 0 {
			c.ime = true
			// (Optional log: fmt.Printf("IME enabled at PC: 0x%04X\n", c.registers.pc))
		}
	}

	return cycles
}

func NewLR35902(bus *memory.Bus, ioRegisters *registers.Registers, ie *registers.Interrupt) *LR35902 {
	cpu := &LR35902{initialized: true}

	cpu.bus = bus
	cpu.io = ioRegisters
	cpu.ie = ie

	// Initialize registers to default values
	cpu.Registers.b = 0x00
	cpu.Registers.c = 0x13
	cpu.Registers.d = 0x00
	cpu.Registers.e = 0xD8
	cpu.Registers.h = 0x01
	cpu.Registers.l = 0x4D
	cpu.Registers.a = 0x01

	// Initialize flags to default values
	cpu.flags.Write(0xB0)

	return cpu
}
