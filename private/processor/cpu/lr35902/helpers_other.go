package lr35902

// Other
func (c *LR35902) decimalAdjust() {
	a := c.registers.a
	subtract := c.flags.Subtract
	halfCarry := c.flags.HalfCarry
	carry := c.flags.Carry
	offset := uint8(0)

	if !subtract {
		if halfCarry || (a&0x0F) > 0x09 {
			offset |= 0x06
		}
		if carry || a > 0x99 {
			offset |= 0x60
		}
		a += offset
	} else {
		if halfCarry {
			offset |= 0x06
		}
		if carry {
			offset |= 0x60
		}
		a -= offset
	}

	c.registers.a = a
	c.flags.Zero = (a == 0)
	c.flags.HalfCarry = false
	c.flags.Carry = ((offset & 0x60) != 0)
}
