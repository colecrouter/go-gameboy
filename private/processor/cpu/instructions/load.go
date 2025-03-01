package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

// Memory access
func load8(c cpu.CPU, r *uint8, val uint8) {
	*r = val
}

func load16(c cpu.CPU, high, low *uint8, val uint16) {
	*high, *low = cpu.FromRegisterPair(val)
	// Flags not affected for plain 16-bit loads.
}

func loadHLSPOffset(c cpu.CPU, offset int8) {
	sp := c.Registers().SP
	// Sign-extend offset correctly.
	result := sp + uint16(int16(offset))

	// Compute flags using the raw byte value.
	unsignedOffset := uint16(uint8(offset))
	var hc, carry = flags.Reset, flags.Reset

	if ((sp & 0xF) + (unsignedOffset & 0xF)) > 0xF {
		hc = flags.Set
	}
	if ((sp & 0xFF) + (unsignedOffset & 0xFF)) > 0xFF {
		carry = flags.Set
	}

	// Store result in HL and update flags: Z and N reset.
	c.Registers().H, c.Registers().L = cpu.FromRegisterPair(result)
	c.Flags().Set(flags.Reset, flags.Reset, hc, carry)
}

func pop16(c cpu.CPU, high, low *uint8) {
	*high, *low = c.Read16(c.Registers().SP)
	c.Registers().SP += 2
}
func push16(c cpu.CPU, high, low uint8) {
	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, cpu.ToRegisterPair(high, low))
}
func load8Mem(c cpu.CPU, r uint8, addr uint16) {
	// For LDH (n), A
	c.Write(addr, r)
}

// popAF pops register AF from the stack and updates A and flag fields.
func popAF(c cpu.CPU) {
	high, low := c.Read16(c.Registers().SP)
	c.Registers().SP += 2
	c.Registers().A = high
	// Update flags: bit7: Z, bit6: N, bit5: H, bit4: C (lower 4 bits ignored)
	c.Flags().Zero = (low & 0x80) != 0
	c.Flags().Subtract = (low & 0x40) != 0
	c.Flags().HalfCarry = (low & 0x20) != 0
	c.Flags().Carry = (low & 0x10) != 0
}
