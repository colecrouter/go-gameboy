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
		A, B, C, D, E, H, L uint8
		SP, PC              uint16
	}
	Flags Flags

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

	ienable := c.ie.Read(0)
	iflag := c.io.InterruptFlag.Read(0)

	// Check for interrupts
	for i := VBlankISR; i <= JoypadISR; i++ {
		var ieBit, ifBit uint8
		ieBit = ienable & (1 << isrOffsets[i])
		ifBit = iflag & (1 << isrOffsets[i])

		if ieBit != 0 && ifBit != 0 {
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

	// if c.Registers.PC >= 0x7000 && c.Registers.PC < 0x9FFF {
	// 	fmt.Printf("")
	// }

	if c.Registers.PC == 0x0a9b {
		// Load stack for debugging
		var stack [63]uint16
		for j := 0; j < len(stack); j++ {
			offset := c.Registers.SP + uint16(j*2)
			if offset > 0xFFFE {
				break
			}
			stack[j] = toRegisterPair(c.bus.Read16(offset))
		}
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

	return cpu
}
