package lr35902

import "math/bits"

// Arithmetic
func (c *LR35902) add8(r *uint8, val uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	sum := uint16(*r) + uint16(val)
	if sum > 0xFF {
		carry = Set
	}
	if ((*r)&0xF)+(val&0xF) > 0xF {
		hc = Set
	}

	*r = uint8(sum)
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Reset, hc, carry)
}
func (c *LR35902) add16(lowDest, highDest *uint8, lowVal, highVal uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	sum, car := bits.Add(uint(*highDest), uint(highVal), 0)
	if car > 0 {
		carry = Set
	}

	if (*highDest&0xf)+(highVal&0xf) > 0xf {
		hc = Set
	}

	*highDest = uint8(sum)

	sum, car = bits.Add(uint(*lowDest), uint(lowVal), 0)
	if car > 0 {
		carry = Set
	}

	if (*lowDest&0xf)+(lowVal&0xf) > 0xf {
		hc = Set
	}

	*lowDest = uint8(sum)

	if *highDest == 0 && *lowDest == 0 {
		zero = Set
	}

	c.setFlags(zero, Reset, hc, carry)
}
func (c *LR35902) sub8(r *uint8, val uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	sum, car := bits.Sub(uint(*r), uint(val), 0)
	if car > 0 {
		carry = Set
	}

	// Fix half-carry: set if lower nibble of *r is less than lower nibble of val.
	if (*r & 0xF) < (val & 0xF) {
		hc = Set
	}

	*r = uint8(sum)

	if *r == 0 {
		zero = Set
	}

	c.setFlags(zero, Set, hc, carry)
}
func (c *LR35902) and8(r *uint8, val uint8) {
	*r &= val
	zero := Reset
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Reset, Set, Reset)
}
func (c *LR35902) or8(r *uint8, val uint8) {
	*r |= val
	zero := Reset
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Reset, Reset, Reset)
}
func (c *LR35902) xor8(r *uint8, val uint8) {
	*r ^= val
	zero := Reset
	if *r == 0 {
		zero = Set
	}
	c.setFlags(zero, Reset, Reset, Reset)
}
func (c *LR35902) addc8(r *uint8, val uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	var carryIn uint16 = 0
	if c.flags.Carry {
		carryIn = 1
	}

	// Compute half-carry using 8-bit arithmetic.
	if ((*r)&0xF)+(val&0xF)+uint8(carryIn) > 0xF {
		hc = Set
	}

	sum := uint16(*r) + uint16(val) + carryIn
	if sum > 0xFF {
		carry = Set
	}

	*r = uint8(sum)
	if *r == 0 {
		zero = Set
	}

	c.setFlags(zero, Reset, hc, carry)
}
func (c *LR35902) subc8(r *uint8, val uint8) {
	zero := Reset
	carry := Reset
	hc := Reset

	// Adjust val if carry flag is set
	if c.flags.Carry {
		val++
	}

	sum, car := bits.Sub(uint(*r), uint(val), 0)
	if car > 0 {
		carry = Set
	}

	// Fix half-carry: set if lower nibble of *r is less than lower nibble of val.
	if (*r & 0xF) < (val & 0xF) {
		hc = Set
	}

	*r = uint8(sum)

	if *r == 0 {
		zero = Set
	}

	c.setFlags(zero, Set, hc, carry)
}
