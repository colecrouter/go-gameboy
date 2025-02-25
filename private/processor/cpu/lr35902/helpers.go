package lr35902

// Helpers
func (c *LR35902) getImmediate8() uint8 {
	val := c.bus.Read(c.Registers.PC + 1)

	<-c.clock        // Use an additional m-cycle to read the immediate value
	c.Registers.PC++ // Increment the program counter to the next instruction

	return val
}

func (c *LR35902) getImmediate16() uint16 {
	high, low := c.bus.Read16(c.Registers.PC + 1)

	<-c.clock           // Use an additional m-cycle to read the immediate value
	c.Registers.PC += 2 // Increment the program counter to the next instruction

	return toRegisterPair(high, low)
}

// toRegisterPair returns a 16-bit register pair from two 8-bit registers
// If you want BC, pass B and C in that order
func toRegisterPair(high, low uint8) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// fromRegisterPair returns two 8-bit registers from a 16-bit register pair
// It returns it in low, high order. E.g. if you have BC (CB in little endian) it will return B, C
func fromRegisterPair(val uint16) (high uint8, low uint8) {
	return uint8(val >> 8), uint8(val)
}
