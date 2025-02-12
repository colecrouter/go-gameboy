package lr35902

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
)

// LR35902 is the original GameBoy CPU
type LR35902 struct {
	initialized bool
	registers   struct {
		a, b, c, d, e, h, l uint8
		sp, pc              uint16
	}
	flags  Flags
	bus    *memory.Bus
	io     *registers.Registers
	cb     bool
	ime    bool
	lastPC uint16
}

// Step executes the next instruction in the CPU's memory.
// Returns the number of clock cycles the instruction took.
func (c *LR35902) Step() int {
	if !c.initialized {
		panic("CPU not initialized")
	}

	// Get instruction
	opcode := c.bus.Read(c.registers.pc)

	// Run instruction
	var instruction instruction
	if c.cb {
		instruction = cbInstructions[opcode]
		c.cb = false
	} else {
		instruction = instructions[opcode]
	}

	var mnemonic string
	if c.bus.Read(c.registers.pc-1) == 0xCB {
		mnemonic = getCBMnemonic(opcode)
	} else {
		mnemonic = mnemonics[opcode]
	}

	_ = mnemonic

	op := instruction.op
	cycles := instruction.c
	increment := instruction.p

	if c.registers.pc == 0x100 {
		c.ime = true
	}

	c.lastPC = c.registers.pc
	if op == nil {
		fmt.Printf("Unimplemented instruction: 0x%02X\r\n", opcode)
		c.registers.pc++
	} else {
		op(c)
		c.registers.pc += uint16(increment)
	}

	return cycles
}

func NewLR35902(bus *memory.Bus, ioRegisters *registers.Registers) *LR35902 {
	if bus == nil {
		panic("Bus is nil")
	}

	cpu := &LR35902{initialized: true}

	cpu.bus = bus
	cpu.io = ioRegisters

	// Initialize registers to default values
	cpu.registers.b = 0x00
	cpu.registers.c = 0x13
	cpu.registers.d = 0x00
	cpu.registers.e = 0xD8
	cpu.registers.h = 0x01
	cpu.registers.l = 0x4D
	cpu.registers.a = 0x01

	// Initialize flags to default values
	cpu.flags.Write(0xB0)

	return cpu
}

func (c *LR35902) PC() uint16 {
	return c.registers.pc
}
