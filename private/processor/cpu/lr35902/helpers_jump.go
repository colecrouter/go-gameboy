package lr35902

// Jump
func (c *LR35902) jump(addr uint16, condition bool) {
	if condition {
		c.registers.pc = addr
		return
	}

	c.registers.pc += 3
}
func (c *LR35902) jumpRelative(offset int8, condition bool) {
	if condition {
		c.registers.pc = uint16(int32(c.registers.pc) + 2 + int32(offset))
		return
	}

	c.registers.pc += 2
}

// Subroutines
func (c *LR35902) ret(condition bool) {
	if !condition {
		c.registers.pc += 1
		return
	}

	// Pop the return address in little endian order.
	high, low := c.bus.Read16(c.registers.sp)
	c.registers.sp += 2

	addr := toRegisterPair(high, low)
	c.registers.pc = addr // Assign popped address directly
}

func (c *LR35902) call(addr uint16, condition bool) {
	if !condition {
		c.registers.pc += 3
		return
	}

	retAddr := c.registers.pc + 3

	// Decrement SP by 2 and push return address using Write16.
	c.registers.sp -= 2
	c.bus.Write16(c.registers.sp, retAddr)

	c.registers.pc = addr
}

func (c *LR35902) rst(addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.registers.pc + 1

	c.registers.sp -= 2
	c.bus.Write16(c.registers.sp, retAddr)

	c.registers.pc = addr
}
