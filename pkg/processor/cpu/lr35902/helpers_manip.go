package lr35902

// Register manipulation
func (c *LR35902) inc8(r *uint8) {
	old := *r
	*r++
	zero := Reset
	hc := Reset
	// Half-carry set when lower nibble overflows (0xF -> 0x0)
	if old&0xF == 0xF {
		hc = Set
	}
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Reset, hc, Leave)
}
func (c *LR35902) inc16(high, low *uint8) {
	combined := toRegisterPair(*high, *low)
	combined++
	*high, *low = fromRegisterPair(combined)
}
func (c *LR35902) dec8(r *uint8) {
	old := *r
	*r--
	zero := Reset
	hc := Reset
	// Half-carry set when lower nibble underflows (0x0 -> 0xF)
	if old&0xF == 0x0 {
		hc = Set
	}
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Set, hc, Leave)
}
func (c *LR35902) dec16(high, low *uint8) {
	combined := toRegisterPair(*high, *low)
	combined--
	*high, *low = fromRegisterPair(combined)
}
