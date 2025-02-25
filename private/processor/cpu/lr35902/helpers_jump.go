package lr35902

// Jump
func (c *LR35902) jump(addr uint16, condition bool) {
	if condition {
		c.Registers.PC = addr - 1
		return
	}

	c.Registers.PC += 2 // 3-byte instruction
}
func (c *LR35902) jumpRelative(offset int8, condition bool) {
	if condition {
		c.Registers.PC = uint16(int32(c.Registers.PC)+2+int32(offset)) - 1
		return
	}

	c.Registers.PC += 1 // 2-byte instruction
}

// Subroutines
func (c *LR35902) ret(condition bool) {
	if !condition {
		c.Registers.PC++
		<-c.clock
		return
	}

	// Pop the return address in little endian order.
	high, low := c.Read16(c.Registers.SP)
	c.Registers.SP += 2

	addr := toRegisterPair(high, low)
	c.Registers.PC = addr - 1 // Assign popped address directly
}

func (c *LR35902) call(addr uint16, condition bool) {
	retAddr := c.Registers.PC + 1 // 3-byte instruction
	if !condition {
		c.Registers.PC = retAddr
		return
	}

	// Decrement SP by 2 and push return address using Write16.
	c.Registers.SP -= 2
	c.bus.Write16(c.Registers.SP, retAddr)

	c.Registers.PC = addr - 1 // Assign new address directly
}

func (c *LR35902) rst(addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.Registers.PC + 1

	c.Registers.SP -= 2
	c.bus.Write16(c.Registers.SP, retAddr)

	c.Registers.PC = addr - 1
}
