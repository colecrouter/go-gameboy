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
	ime   bool
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

	if c.registers.pc == 0x02a0 {
		fmt.Println("Breakpoint")
	}

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
	fmt.Printf("\r\nRegisters:\r\n")
	fmt.Printf("A: 0x%02X  B: 0x%02X  C: 0x%02X  D: 0x%02X  E: 0x%02X\r\n", c.registers.a, c.registers.b, c.registers.c, c.registers.d, c.registers.e)
	fmt.Printf("H: 0x%02X  L: 0x%02X\r\n", c.registers.h, c.registers.l)
	fmt.Printf("SP: 0x%04X  PC: 0x%04X\r\n", c.registers.sp, c.registers.pc)
	fmt.Printf("Flags: Z=%t  N=%t  H=%t  C=%t\r\n", c.flags.Zero, c.flags.Subtract, c.flags.HalfCarry, c.flags.Carry)
}

func (c *LR35902) PC() uint16 {
	return c.registers.pc
}

func (c *LR35902) ISR(isr ISR) {
	// If master interrupt enable is disabled, return
	if !c.ime {
		return
	}

	// If the interrupt is disabled, return
	if c.io.InterruptEnable.Read()<<isrOffsets[isr] == 0 {
		return
	}

	// STAT interrupt is special
	if isr == LCDSTATISR {
		switch c.io.LCDStatus.PPUMode {
		case registers.HBlank:
			if !c.io.LCDStatus.Mode0Interrupt {
				return
			}
		case registers.VBlank:
			if !c.io.LCDStatus.Mode1Interrupt {
				return
			}
		case registers.OAMScan:
			if !c.io.LCDStatus.Mode2Interrupt {
				return
			}
		case registers.Drawing:
			if !c.io.LCDStatus.LYCInterrupt {
				return
			}
		}
	}

	// Push PC onto stack
	c.registers.sp -= 2
	c.bus.Write(c.registers.sp, uint8(c.registers.pc>>8))
	c.bus.Write(c.registers.sp+1, uint8(c.registers.pc))

	// Jump to ISR
	c.registers.pc = isrAddresses[isr]

	// Disable interrupts
	c.ime = false
}

type ISR int

const (
	VBlankISR ISR = iota
	LCDSTATISR
	TimerISR
	SerialISR
	JoypadISR
)

var isrAddresses = [5]uint16{
	VBlankISR:  0x0040,
	LCDSTATISR: 0x0048,
	TimerISR:   0x0050,
	SerialISR:  0x0058,
	JoypadISR:  0x0060,
}

var isrOffsets = [5]uint8{
	VBlankISR:  0,
	LCDSTATISR: 1,
	TimerISR:   2,
	SerialISR:  3,
	JoypadISR:  4,
}
