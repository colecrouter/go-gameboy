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
func (c *LR35902) rotate(r *uint8, left, useCarryBit, updateZ bool) {
	// Extract the bit that will be rotated out.
	// For left rotate, that is the MSB; for right rotate, the LSB.
	var carriedOut uint8
	if left {
		carriedOut = (*r & 0x80) >> 7
	} else {
		carriedOut = *r & 1
	}

	// Determine the input bit.
	// For instructions that use the external carry, use c.flags.Carry;
	// otherwise, use the bit that got shifted out.
	var carryIn uint8
	if useCarryBit {
		if c.Flags.Carry {
			carryIn = 1
		} else {
			carryIn = 0
		}
	} else {
		carryIn = carriedOut
	}

	// Perform the rotation.
	if left {
		*r = (*r << 1) | carryIn
	} else {
		*r = (*r >> 1) | (carryIn << 7)
	}

	// Set the new Carry flag based on the bit that was rotated out.
	var newCarryFlag = Reset
	if carriedOut == 1 {
		newCarryFlag = Set
	}

	// Update the Zero flag only if requested (CB rotates update Z).
	var newZFlag FlagState
	if updateZ {
		if *r == 0 {
			newZFlag = Set
		} else {
			newZFlag = Reset
		}
	} else {
		newZFlag = Reset // Always clear Z when !updateZ.
	}

	// The N and H flags are reset on rotate instructions.
	c.setFlags(newZFlag, Reset, Reset, newCarryFlag)
}
func (c *LR35902) shift(r *uint8, left, arithmeticRight bool) {
	original := *r

	// Determine what the new Carry flag should be.
	var newCarry bool

	if left {
		// For left shifts, the high bit goes into carry.
		newCarry = (original & 0x80) != 0
		*r = original << 1
	} else {
		// For right shifts, the low bit is carried out.
		newCarry = (original & 0x01) != 0
		if arithmeticRight {
			// For arithmetic right shifts, preserve the original MSB.
			msb := original & 0x80
			*r = (original >> 1) | msb
		} else {
			*r = original >> 1
		}
	}

	// Set the Zero flag if the result is zero.
	flagZero := (*r == 0)

	// Now, convert boolean conditions into flag behaviors.
	// For c.setFlags(zero, N, H, carry):
	var zeroBehavior, nBehavior, hBehavior, carryBehavior FlagState

	if flagZero {
		zeroBehavior = Set
	} else {
		zeroBehavior = Reset
	}

	// The N and H flags are always cleared after these shifts.
	nBehavior = Reset
	hBehavior = Reset

	if newCarry {
		carryBehavior = Set
	} else {
		carryBehavior = Reset
	}

	// Update the flags.
	c.setFlags(zeroBehavior, nBehavior, hBehavior, carryBehavior)
}
func (c *LR35902) swap(r *uint8) {
	// Swap the upper and lower nibbles
	*r = (*r&0xf)<<4 | *r>>4

	zero := Reset
	if *r == 0 {
		zero = Set
	}

	c.setFlags(zero, Reset, Reset, Reset)
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

	if r&(1<<b) == 0 {
		zero = Set
	}

	c.setFlags(zero, Reset, Set, Leave)
}

// Bit manipulation
func (c *LR35902) res(b uint8, r *uint8) {
	*r &= ^(1 << b)
}
func (c *LR35902) set(b uint8, r *uint8) {
	*r |= 1 << b
}
