package lr35902

import (
	"encoding/binary"
	"log"

	"github.com/colecrouter/gameboy-go/pkg/memory"
)

// LR35902 is the original GameBoy CPU
type LR35902 struct {
	initialized bool
	registers   struct {
		A, B, C, D, E, H, L, Flags uint8
		SP, PC                     uint16
	}
	bus           memory.Device
	done, doClock chan struct{}
	clocking      chan uint8
}

// Set, reset, leave for flags
const (
	set = iota
	reset
	leave
)

func (c *LR35902) setFlags(zero int, subtract int, halfCarry int, carry int) {
	if !c.initialized {
		panic("CPU not initialized")
	}

	flags := []int{zero, subtract, halfCarry, carry}
	for i, v := range flags {
		offset := uint8(7 - i)
		switch v {
		case set:
			c.registers.Flags |= (1 << offset)
		case reset:
			c.registers.Flags |= (0 << offset)
		case leave:
		}
	}
}

// Clock emulates a clock cycle on the CPU
func (c *LR35902) Clock() {
	if !c.initialized {
		panic("CPU not initialized")
	}

	// Get instruction
	op := c.bus.Read(c.registers.PC)
	log.Printf("0x%X: 0x%X\n", c.registers.PC, op)

	// Increment PC
	c.registers.PC++

	// Run instruction
	switch op {
	case 0x00:
		go c.nop()
	case 0x01:
		go c.ldBcD16()
	case 0x02:
		go c.ldpBcA()
	case 0x03:
		go c.incBc()
	case 0x04:
		go c.incB()
	case 0x05:
		go c.decB()
	case 0x06:
		go c.ldBD8()
	case 0x07:
		go c.rlcA()
	case 0x08:
		go c.ldpA16SP()
	case 0x09:
		go c.addHlBc()
	case 0x0A:
		go c.ldApBc()
	case 0x0B:
		go c.decBc()
	case 0x0C:
		go c.incC()
	case 0x0D:
		go c.decC()
	case 0x0E:
		go c.ldCD8()
	case 0x0F:
		go c.rrcA()
	// 0x1x
	case 0x10:
		go c.stop()
	case 0x11:
		go c.ldDeD16()
	case 0x12:
		go c.ldpDeA()
	case 0x13:
		go c.incDe()
	case 0x14:
		go c.incD()
	case 0x15:
		go c.decD()
	case 0x16:
		go c.ldDD8()
	case 0x17:
		go c.rlA()
	case 0x18:
		go c.jrS8()
	case 0x19:
		go c.addHlDe()
	case 0x1A:
		go c.ldApDe()
	case 0x1B:
		go c.decDe()
	case 0x1C:
		go c.incE()
	case 0x1D:
		go c.decE()
	case 0x1E:
		go c.ldED8()
	case 0x1F:
		go c.rrA()
		// 0x2x
	case 0x20:
		go c.jrNzS8()
	case 0x21:
		go c.ldHlD16()
	case 0x22:
		go c.ldpHlA()
	case 0x23:
		go c.incHl()
	case 0x24:
		go c.incH()
	case 0x25:
		go c.decH()
	case 0x26:
		go c.ldHD8()
	case 0x27:
		go c.ddA()
	case 0x28:
		go c.jrZS8()
	case 0x29:
		go c.addHlHl()
	case 0x2A:
		go c.ldApHlp()
	case 0x2B:
		go c.decHl()
	case 0x2C:
		go c.incL()
	case 0x2D:
		go c.decL()
	case 0x2E:
		go c.ldLD8()
	case 0x2F:
		go c.cpl()
	default:
		go func() {
			c.done <- struct{}{}
		}()
	}

	// Skip if instruction is still running
	select {
	case <-c.done:
	case <-c.doClock:
		c.clocking <- op
		return
	}
}

func NewLR35902(bus *memory.Bus) *LR35902 {
	if bus == nil {
		panic("Bus is nil")
	}

	cpu := &LR35902{initialized: true}

	cpu.registers.PC = 0x0100
	cpu.done = make(chan struct{})
	cpu.clocking = make(chan uint8)
	cpu.doClock = make(chan struct{}, 1)
	cpu.bus = bus

	return cpu
}

func toLong(a uint8, b uint8) uint16 {
	return uint16(b)<<8 | uint16(a)
}

func toShort(a uint16) (uint8, uint8) {
	return uint8(a), uint8(a << 8)
}

func uintToLittleEndian16(a uint) uint16 {
	// Convert back to buffer
	buf := make([]byte, 2)
	temp := uint16(a)
	binary.LittleEndian.PutUint16(buf, temp)

	// Convert buffer to uint
	return binary.LittleEndian.Uint16(buf)
}

func littleEndian16ToUint(a uint16) uint {
	// Convert from Little Endian to uint
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, a)
	var bnew = binary.LittleEndian
	temp := bnew.Uint16(buf)

	return uint(temp)
}
