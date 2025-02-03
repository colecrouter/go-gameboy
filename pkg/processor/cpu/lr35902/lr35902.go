package lr35902

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
)

// LR35902 is the original GameBoy CPU
type LR35902 struct {
	initialized bool
	registers   struct {
		a, b, c, d, e, h, l uint8
		sp, pc              uint16
	}
	flags Flags
	bus   memory.Device
	io    *registers.Registers
	cb    bool
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
	var ok bool
	if c.cb {
		instruction, ok = cbInstructions[opcode]
		c.cb = false
	} else {
		instruction, ok = instructions[opcode]
	}
	if !ok {
		panic(fmt.Sprintf("Unknown opcode: 0x%X", opcode))
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

	a := c.io.LCDStatus.LY
	_ = a
	op(c)

	// Increment PC
	c.registers.pc++

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

// PrintRegisters prints the current state of the CPU registers and flags.
func (c *LR35902) PrintRegisters() {
	fmt.Printf("\nRegisters:\n")
	fmt.Printf("A: 0x%02X  B: 0x%02X  C: 0x%02X  D: 0x%02X  E: 0x%02X\n", c.registers.a, c.registers.b, c.registers.c, c.registers.d, c.registers.e)
	fmt.Printf("H: 0x%02X  L: 0x%02X\n", c.registers.h, c.registers.l)
	fmt.Printf("SP: 0x%04X  PC: 0x%04X\n", c.registers.sp, c.registers.pc)
	fmt.Printf("Flags: Z=%t  N=%t  H=%t  C=%t\n", c.flags.Zero, c.flags.Subtract, c.flags.HalfCarry, c.flags.Carry)
}
