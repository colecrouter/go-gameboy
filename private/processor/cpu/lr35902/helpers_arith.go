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
func (c *LR35902) add16(highDest, lowDest *uint8, highVal, lowVal uint8) {
	carry := Reset
	hc := Reset

	firstVal := toRegisterPair(highVal, lowVal)
	secondVal := toRegisterPair(*highDest, *lowDest)

	sum := firstVal + secondVal

	if sum < firstVal || sum < secondVal {
		carry = Set
	}
	if (firstVal&0xFFF)+(secondVal&0xFFF) > 0xFFF {
		hc = Set
	}

	*highDest, *lowDest = fromRegisterPair(sum)

	c.setFlags(Leave, Reset, hc, carry)
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
func (c *LR35902) addSPr8() {
	operand := int8(c.getImmediate8())
	result := c.registers.sp + uint16(int16(operand))

	// Compute half-carry and carry flags using only the lower nibble/byte.
	hc := Reset
	if (c.registers.sp&0xF)+(uint16(uint8(operand))&0xF) > 0xF {
		hc = Set
	}
	carry := Reset
	if (c.registers.sp&0xFF)+uint16(uint8(operand)) > 0xFF {
		carry = Set
	}

	// Write result to SP (wraps naturally to 16 bits) and update flags, Z and N are reset.
	c.registers.sp = result
	c.setFlags(Reset, Reset, hc, carry)
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
