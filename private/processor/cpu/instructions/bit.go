package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

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
func rotate(c cpu.CPU, r *uint8, left, useCarryBit, updateZ bool) {
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
		if c.Flags().Carry {
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

	// flags.Set the new Carry flag based on the bit that was rotated out.
	var newCarryFlag = flags.Reset
	if carriedOut == 1 {
		newCarryFlag = flags.Set
	}

	// Update the Zero flag only if requested (CB rotates update Z).
	var newZFlag flags.FlagState
	if updateZ {
		if *r == 0 {
			newZFlag = flags.Set
		} else {
			newZFlag = flags.Reset
		}
	} else {
		newZFlag = flags.Reset // Always clear Z when !updateZ.
	}

	// The N and H flags are reset on rotate instructions.
	c.Flags().Set(newZFlag, flags.Reset, flags.Reset, newCarryFlag)
}
func shift(c cpu.CPU, r *uint8, left, arithmeticRight bool) {
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

	// flags.Set the Zero flag if the result is zero.
	flagZero := (*r == 0)

	// Now, convert boolean conditions into flag behaviors.
	// For c.Flags().Set(zero, N, H, carry):
	var zeroBehavior, nBehavior, hBehavior, carryBehavior flags.FlagState

	if flagZero {
		zeroBehavior = flags.Set
	} else {
		zeroBehavior = flags.Reset
	}

	// The N and H flags are always cleared after these shifts.
	nBehavior = flags.Reset
	hBehavior = flags.Reset

	if newCarry {
		carryBehavior = flags.Set
	} else {
		carryBehavior = flags.Reset
	}

	// Update the flags.
	c.Flags().Set(zeroBehavior, nBehavior, hBehavior, carryBehavior)
}
func swap(c cpu.CPU, r *uint8) {
	// Swap the upper and lower nibbles
	*r = (*r&0xf)<<4 | *r>>4

	zero := flags.Reset
	if *r == 0 {
		zero = flags.Set
	}

	c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)
}

// Comparing
func cp8(c cpu.CPU, r1 uint8, r2 uint8) {
	zero := flags.Reset
	carry := flags.Reset
	hc := flags.Reset

	if r1 < r2 {
		carry = flags.Set
	}

	if r1&0xf < r2&0xf {
		hc = flags.Set
	}

	if r1 == r2 {
		zero = flags.Set
	}

	c.Flags().Set(zero, flags.Set, hc, carry)
}
func bit(c cpu.CPU, b uint8, r uint8) {
	zero := flags.Reset

	if r&(1<<b) == 0 {
		zero = flags.Set
	}

	c.Flags().Set(zero, flags.Reset, flags.Set, flags.Leave)
}

// Bit manipulation
func res(c cpu.CPU, b uint8, r *uint8) {
	*r &= ^(1 << b)
}
func set(c cpu.CPU, b uint8, r *uint8) {
	*r |= 1 << b
}
