package lr35902

// Memory access
func (c *LR35902) load8(r *uint8, val uint8) {
	*r = val
}

func (c *LR35902) load16(high, low *uint8, val uint16) {
	*high, *low = fromRegisterPair(val)
	// Flags not affected for plain 16-bit loads.
}

func (c *LR35902) loadHLSPOffset(offset int8) {
	result := c.Registers.sp + uint16(int16(offset))
	var hc, carry = Reset, Reset

	if (c.Registers.sp&0xF)+(uint16(uint8(offset))&0xF) > 0xF {
		hc = Set
	}
	if (c.Registers.sp&0xFF)+(uint16(uint8(offset))) > 0xFF {
		carry = Set
	}

	// Load computed result into HL and update flags: Z and N reset.
	c.Registers.h, c.Registers.l = fromRegisterPair(result)
	c.setFlags(Reset, Reset, hc, carry)
}

func (c *LR35902) pop16(high, low *uint8) {
	*high, *low = c.bus.Read16(c.Registers.sp)
	c.Registers.sp += 2
}
func (c *LR35902) push16(high, low uint8) {
	c.Registers.sp -= 2
	c.bus.Write16(c.Registers.sp, toRegisterPair(high, low))
}
func (c *LR35902) load8Mem(r uint8, addr uint16) {
	// For LDH (n), A
	c.bus.Write(addr, r)
}

// popAF pops register AF from the stack and updates A and flag fields.
func (c *LR35902) popAF() {
	high, low := c.bus.Read16(c.Registers.sp)
	c.Registers.sp += 2
	c.Registers.a = high
	// Update flags: bit7: Z, bit6: N, bit5: H, bit4: C (lower 4 bits ignored)
	c.flags.Zero = (low & 0x80) != 0
	c.flags.Subtract = (low & 0x40) != 0
	c.flags.HalfCarry = (low & 0x20) != 0
	c.flags.Carry = (low & 0x10) != 0
}
