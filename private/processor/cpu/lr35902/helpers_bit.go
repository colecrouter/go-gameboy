package lr35902

/*
RLA
┌────────────────────┐
│ ┌──┐  ┌─────────┐  │
└─│CY│<─│7<──────0│<─┘
  └──┘  └─────────┘
             A

RLCA
      ┌──────────────┐
┌──┐  │ ┌─────────┐  │
│CY│<─┴─│7<──────0│<─┘
└──┘    └─────────┘
             A
*/

// Rotates/Shifts
func (c *LR35902) rotate(r *uint8, left bool, useCarryBit bool) {
	var carriedOut uint8
	var carriedIn uint8

	if left {
		carriedOut = *r >> 7
	} else {
		carriedOut = *r & 1
	}

	if useCarryBit {
		if c.flags.Carry {
			carriedIn = 1
		}
	} else {
		carriedIn = carriedOut
	}

	if left {
		*r = *r<<1 | carriedIn
	} else {
		*r = *r>>1 | (carriedIn << 7)
	}

	carryFlag := Reset
	if carriedOut == 1 {
		carryFlag = Set
	}

	c.setFlags(Reset, Reset, Reset, carryFlag)
}
func (c *LR35902) shift(r *uint8, left bool) {
	if left {
		c.flags.Carry = *r&0b10000000 == 0b10000000
		*r <<= 1
	} else {
		c.flags.Carry = *r&0b00000001 == 0b00000001
		*r >>= 1
	}
}
func (c *LR35902) swap(r *uint8) {
	// Swap the upper and lower nibbles
	*r = (*r&0xf)<<4 | *r>>4
}

// Comparing
func (c *LR35902) cp8(r1 uint8, r2 uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	if r1 < r2 {
		carry = Set
	}

	if r1&0xf < r2&0xf {
		hc = Set
	}

	if r1 == r2 {
		zero = Set
	}

	c.setFlags(zero, Set, hc, carry)
}
func (c *LR35902) bit(b uint8, r uint8) {
	zero := Reset
	hc := Set
	var carry FlagState

	if r&(1<<b) == 0 {
		zero = Set
	}

	if b == 7 {
		carry = Set
	} else {
		carry = Reset
	}

	c.setFlags(zero, Reset, hc, carry)
}

// Bit manipulation
func (c *LR35902) res(b uint8, r *uint8) {
	*r &= ^(1 << b)
}
func (c *LR35902) set(b uint8, r *uint8) {
	*r |= 1 << b
}
