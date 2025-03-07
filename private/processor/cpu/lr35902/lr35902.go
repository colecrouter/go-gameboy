package lr35902

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/logging"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/registers"
	"github.com/colecrouter/gameboy-go/private/system"
)

// LR35902 is the original GameBoy CPU
type LR35902 struct {
	initialized bool
	registers   registers.Registers
	flags       flags.Flags
	bus         *memory.Bus
	io          *io.Registers
	ie          *io.Interrupt
	cb          bool
	ime         bool
	eiDelay     int
	lastPC      uint16
	halted      bool
	clock       <-chan struct{}

	logger logging.Logger
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

	var instruction instructions.Instruction
	var mnemonic string
	var op func(cpu.CPU)
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
	// We don't clock here, because the fetch stage overlaps with the previous instruction's execute stage
	opcode = c.bus.Read(c.registers.PC)

	if c.cb {
		instruction = instructions.CBInstructions[opcode]
		mnemonic = getCBMnemonic(opcode)
		c.cb = false
	} else {
		instruction = instructions.Instructions[opcode]
		mnemonic = mnemonics[opcode]
	}

	_ = mnemonic

	op = instruction

	c.lastPC = c.registers.PC

	if c.registers.PC == 0x29d0 {
		if c.registers.A&0b1111 != 0b1111 {
			fmt.Printf("")
		}
	}

	if c.registers.PC == 0x03a0 {
		fmt.Printf("")
	}

	// Execute instruction
	op(c)

	// Update DI and EI delay
	if c.eiDelay > 0 {
		c.eiDelay--
		if c.eiDelay == 0 {
			c.ime = true
		}
	}

	c.registers.PC++
	<-c.clock
}

func NewLR35902(broadcaster *system.Broadcaster, bus *memory.Bus, ioRegisters *io.Registers, ie *io.Interrupt) *LR35902 {
	cpu := &LR35902{initialized: true}

	if broadcaster != nil {
		cpu.clock = broadcaster.Subscribe(system.MRisingEdge)
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

func (c *LR35902) printStack() []uint16 {
	var stack [63]uint16
	var j int
	for j = range len(stack) {
		offset := uint32(c.registers.SP) + uint32(j*2)
		if offset > 0xFFFE {
			break
		}
		stack[j] = cpu.ToRegisterPair(c.bus.Read16(uint16(offset)))
	}

	return stack[0:j]
}
