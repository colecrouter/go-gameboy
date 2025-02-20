package lr35902

// Jump
func (c *LR35902) jump(addr uint16, condition bool) {
	if condition {
		c.Registers.PC = addr
		return
	}

	c.Registers.PC += 3
}
func (c *LR35902) jumpRelative(offset int8, condition bool) {
	if condition {
		c.Registers.PC = uint16(int32(c.Registers.PC) + 2 + int32(offset))
		return
	}

	c.Registers.PC += 2
}

// Subroutines
func (c *LR35902) ret(condition bool) {
	if !condition {
		c.Registers.PC += 1
		return
	}

	// Pop the return address in little endian order.
	high, low := c.bus.Read16(c.Registers.sp)
	c.Registers.sp += 2

	addr := toRegisterPair(high, low)
	c.Registers.PC = addr // Assign popped address directly
}

func (c *LR35902) call(addr uint16, condition bool) {
	if !condition {
		c.Registers.PC += 3
		return
	}

	retAddr := c.Registers.PC + 3

	// Decrement SP by 2 and push return address using Write16.
	c.Registers.sp -= 2
	c.bus.Write16(c.Registers.sp, retAddr)

	c.Registers.PC = addr
}

func (c *LR35902) rst(addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.Registers.PC + 1

	c.Registers.sp -= 2
	c.bus.Write16(c.Registers.sp, retAddr)

	c.Registers.PC = addr
}
