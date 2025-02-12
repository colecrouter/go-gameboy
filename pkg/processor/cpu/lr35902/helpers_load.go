package lr35902

// Memory access
func (c *LR35902) load8(r *uint8, val uint8) {
	*r = val
}
func (c *LR35902) load16(high, low *uint8, val uint16) {
	*high, *low = fromRegisterPair(val)
}
func (c *LR35902) pop16(high, low *uint8) {
	*high, *low = c.bus.Read16(c.registers.sp)
	c.registers.sp += 2
}
func (c *LR35902) push16(high, low uint8) {
	c.registers.sp -= 2
	c.bus.Write16(c.registers.sp, toRegisterPair(high, low))
}
func (c *LR35902) load8Mem(r uint8, addr uint16) {
	// For LDH (n), A
	c.bus.Write(addr, r)
}
