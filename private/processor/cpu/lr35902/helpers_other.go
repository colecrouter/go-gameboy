package lr35902

// Other
func (c *LR35902) decimalAdjust() {
	a := c.registers.a

	var carry uint8
	var halfCarry uint8
	var zero uint8
	if c.flags.Carry {
		carry = 1
	}
	if c.flags.HalfCarry {
		halfCarry = 1
	}
	if c.flags.Zero {
		zero = 1
	}

	if !c.flags.Subtract {
		if c.flags.HalfCarry || (a&0xF) > 9 {
			a += 0x06
		}
		if c.flags.Carry || a > 0x9F {
			a += 0x60
		}
	} else {
		if c.flags.HalfCarry {
			a = (a - 6) & 0xFF
		}
		if c.flags.Carry {
			a -= 0x60
		}
	}

	c.flags.Write(c.flags.Read() & ^(halfCarry | zero))

	if (int(a) & 0x100) == 0x100 {
		c.flags.Write(c.flags.Read() | carry)
	}

	a &= 0xFF

	if a == 0 {
		c.flags.Write(c.flags.Read() | zero)
	}

	c.registers.a = a

	// a := c.registers.a
	// subtract := c.flags.Subtract
	// halfCarry := c.flags.HalfCarry
	// carry := c.flags.Carry
	// offset := uint8(0)

	// if !subtract {
	// 	if halfCarry || (a&0x0F) > 0x09 {
	// 		offset |= 0x06
	// 	}
	// 	if carry || a > 0x99 {
	// 		offset |= 0x60
	// 	}
	// 	a += offset
	// } else {
	// 	if halfCarry {
	// 		offset |= 0x06
	// 	}
	// 	if carry {
	// 		offset |= 0x60
	// 	}
	// 	a -= offset
	// }

	// c.registers.a = a
	// c.flags.Zero = (a == 0)
	// c.flags.HalfCarry = false
	// c.flags.Carry = ((offset & 0x60) != 0)
}
