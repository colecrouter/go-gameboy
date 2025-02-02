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
	flags         Flags
	bus           memory.Device
	io            *registers.Registers
	done, doClock chan struct{}
	clocking      chan uint8
}

// Clock emulates a clock cycle on the CPU
func (c *LR35902) Clock() {
	if !c.initialized {
		panic("CPU not initialized")
	}

	// Get instruction
	opcode := c.bus.Read(c.registers.pc)

	// Run instruction
	switch opcode {
	// ...existing cases...
	case 0x88: // ADC A,B
		// Changed: call addc8 to include carry input and let the global PC increment handle advancing.
		c.addc8(&c.registers.a, c.registers.b)
	// ...existing cases...
	default:
		instruction, ok := instructions[opcode]
		if ok {
			mnemonic := mnemonics[opcode]
			_ = mnemonic
			instruction.op(c)
		} else {
			panic(fmt.Sprintf("Unknown opcode: 0x%X", opcode))
		}
	}

	// Skip if instruction is still running
	// select {
	// case <-c.done:
	// case <-c.doClock:
	// 	c.clocking <- op
	// 	return
	// }

	// Increment PC
	c.registers.pc++

}

func NewLR35902(bus *memory.Bus, ioRegisters *registers.Registers) *LR35902 {
	if bus == nil {
		panic("Bus is nil")
	}

	cpu := &LR35902{initialized: true}

	// cpu.registers.pc = 0x0100

	// cpu.registers.PC = 0x0100
	cpu.done = make(chan struct{})
	cpu.clocking = make(chan uint8)
	cpu.doClock = make(chan struct{}, 1)
	cpu.bus = bus
	cpu.io = ioRegisters

	// Initialize registers to default values
	cpu.registers.b = 0x00
	cpu.registers.c = 0x13
	cpu.registers.d = 0x84
	cpu.registers.e = 0xD2
	cpu.registers.h = 0x01
	cpu.registers.l = 0x4D
	cpu.registers.a = 0x11

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
