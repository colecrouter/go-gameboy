package lr35902

import (
	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/system"
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
	clock   <-chan struct{}
}

// step executes the next instruction in the CPU's memory.
// Returns the number of T-cycles the instruction took.
func (c *LR35902) MClock() {
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

	var instruction instruction
	var mnemonic string
	var op func(*LR35902)
	var opcode uint8

	// Check for HALT mode
	if c.halted {
		// We need to process a cycle here so that the CPU still runs and checks for interrupts
		// If we return 0, the process will hang, as it will continue to clock empty cycles without stopping or checking for interrupts
		<-c.clock

		// Skip the instruction execution stage
		return
	}

	// Fetch the next instruction
	opcode = c.bus.Read(c.Registers.PC)

	if c.cb {
		instruction = cbInstructions[opcode]
		mnemonic = getCBMnemonic(opcode)
		c.cb = false
	} else {
		instruction = instructions[opcode]
		mnemonic = mnemonics[opcode]
	}

	_ = mnemonic

	op = instruction.op

	c.lastPC = c.Registers.PC

	// Execute instruction
	op(c)

	// We don't clock here, because the fetch stage overlaps with the previous instruction's execute stage

	// Update DI and EI delay
	if c.eiDelay > 0 {
		c.eiDelay--
		if c.eiDelay == 0 {
			c.ime = true
		}
	}

	c.Registers.PC++
	<-c.clock
}

func NewLR35902(broadcaster *system.Broadcaster, bus *memory.Bus, ioRegisters *registers.Registers, ie *registers.Interrupt) *LR35902 {
	cpu := &LR35902{initialized: true}

	if broadcaster != nil {
		cpu.clock = broadcaster.SubscribeM()
	}
	cpu.bus = bus
	cpu.io = ioRegisters
	cpu.ie = ie

	return cpu
}

func (c *LR35902) Run(close <-chan struct{}) {
	if !c.initialized {
		panic("CPU not initialized")
	}

	for {
		select {
		case <-close:
			return
		default:
			c.MClock()
		}
	}
}

func (c *LR35902) Read(addr uint16) byte {
	val := c.bus.Read(addr)
	c.Registers.PC++
	<-c.clock
	return val
}

func (c *LR35902) Read16(addr uint16) (high, low uint8) {
	high, low = c.bus.Read16(addr)
	c.Registers.PC += 2
	<-c.clock
	<-c.clock
	return high, low
}
