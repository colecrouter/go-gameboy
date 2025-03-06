package lr35902

import (
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

func (c *LR35902) isr(isr ISR) {
	// STAT interrupt is special
	if isr == LCDSTATISR {
		switch c.io.LCDStatus.PPUMode {
		case io.HBlank:
			if !c.io.LCDStatus.Mode0Interrupt {
				return
			}
		case io.VBlank:
			if !c.io.LCDStatus.Mode1Interrupt {
				return
			}
		case io.OAMScan:
			if !c.io.LCDStatus.Mode2Interrupt {
				return
			}
		case io.Drawing:
			if !c.io.LCDStatus.LYCInterrupt {
				return
			}
		}
	}

	// Two additional m-cycles
	c.ClockAndAck()
	c.ClockAndAck()

	// Push PC onto stack
	// This consumes an additional 2 m-cycles

	highPC, lowPC := cpu.FromRegisterPair(c.registers.PC)
	c.Clock()
	c.Registers().SP--
	c.Write(c.registers.SP, highPC)
	c.Ack()

	c.Clock()
	c.Registers().SP--
	c.Write(c.registers.SP, lowPC)
	c.Ack()

	// Jump to ISR
	// PC won't be incremented, so don't -1
	c.registers.PC = isrAddresses[isr]

	// One last m-cycle for the write(?)
	c.ClockAndAck()

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
