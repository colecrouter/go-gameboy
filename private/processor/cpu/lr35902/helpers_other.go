package lr35902

// Other
func (c *LR35902) decimalAdjust() {
	// Copied from https://stackoverflow.com/a/57837042/9731890
	var t uint8

	if c.flags.HalfCarry || (c.registers.a&0x0F) > 9 {
		t++
	}

	if c.flags.Carry || c.registers.a > 0x99 {
		t += 2
		c.flags.Carry = true
	}

	// Builds final H flag
	if c.flags.Subtract && !c.flags.HalfCarry {
		c.flags.HalfCarry = false
	} else {
		if c.flags.Subtract && c.flags.HalfCarry {
			c.flags.HalfCarry = ((c.registers.a & 0x0F) < 6)
		} else {
			c.flags.HalfCarry = ((c.registers.a & 0x0F) >= 0x0A)
		}
	}

	switch t {
	case 1:
		if c.flags.Subtract {
			c.registers.a -= 6
		} else {
			c.registers.a += 6
		}
	case 2:
		if c.flags.Subtract {
			c.registers.a -= 0x60
		} else {
			c.registers.a += 0x60
		}
	case 3:
		if c.flags.Subtract {
			c.registers.a -= 0x66
		} else {
			c.registers.a += 0x66
		}
	}

	if c.registers.a == 0 {
		c.flags.Zero = true
	} else {
		c.flags.Zero = false
	}
}
