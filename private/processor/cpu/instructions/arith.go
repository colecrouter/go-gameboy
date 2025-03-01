package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

// Arithmetic
func add8(c cpu.CPU, r *uint8, val uint8) {
	zero := flags.Reset
	carry := flags.Reset
	hc := flags.Reset

	sum := uint16(*r) + uint16(val)
	if sum > 0xFF {
		carry = flags.Set
	}
	if ((*r)&0xF)+(val&0xF) > 0xF {
		hc = flags.Set
	}

	*r = uint8(sum)
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Reset, hc, carry)
}
func add16(c cpu.CPU, highDest, lowDest *uint8, highVal, lowVal uint8) {
	carry := flags.Reset
	hc := flags.Reset

	firstVal := cpu.ToRegisterPair(highVal, lowVal)
	secondVal := cpu.ToRegisterPair(*highDest, *lowDest)

	sum := firstVal + secondVal

	if sum < firstVal || sum < secondVal {
		carry = flags.Set
	}
	if (firstVal&0xFFF)+(secondVal&0xFFF) > 0xFFF {
		hc = flags.Set
	}

	*highDest, *lowDest = cpu.FromRegisterPair(sum)

	c.Flags().Set(flags.Leave, flags.Reset, hc, carry)
}
func sub8(c cpu.CPU, r *uint8, val uint8) {
	zeroFlag := flags.Reset
	carryFlag := flags.Reset
	halfCarryFlag := flags.Reset

	// The subtraction is performed in a signed context to detect borrow.
	diff := int16(*r) - int16(val)

	if diff < 0 {
		carryFlag = flags.Set
		// Wrap-around for unsigned arithmetic.
		diff += 256
	}

	// Check for half-carry (i.e. borrow from bit 4).
	if (*r & 0xF) < (val & 0xF) {
		halfCarryFlag = flags.Set
	}

	result := uint8(diff)
	if result == 0 {
		zeroFlag = flags.Set
	}

	*r = result

	// In subtraction, the N flag is always set.
	c.Flags().Set(zeroFlag, flags.Set, halfCarryFlag, carryFlag)
}
func and8(c cpu.CPU, r *uint8, val uint8) {
	*r &= val
	zero := flags.Reset
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Reset, flags.Set, flags.Reset)
}
func addSPr8(c cpu.CPU) {
	operand := int8(c.GetImmediate8())
	result := c.Registers().SP + uint16(int16(operand))

	// Compute half-carry and carry flags using only the lower nibble/byte.
	hc := flags.Reset
	if (c.Registers().SP&0xF)+(uint16(uint8(operand))&0xF) > 0xF {
		hc = flags.Set
	}
	carry := flags.Reset
	if (c.Registers().SP&0xFF)+uint16(uint8(operand)) > 0xFF {
		carry = flags.Set
	}

	// Write result to SP (wraps naturally to 16 bits) and update flags, Z and N are reset.
	c.Registers().SP = result
	c.Flags().Set(flags.Reset, flags.Reset, hc, carry)
}
func or8(c cpu.CPU, r *uint8, val uint8) {
	*r |= val
	zero := flags.Reset
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)
}
func xor8(c cpu.CPU, r *uint8, val uint8) {
	*r ^= val
	zero := flags.Reset
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Reset, flags.Reset, flags.Reset)
}
func addc8(c cpu.CPU, r *uint8, val uint8) {
	zero := flags.Reset
	carry := flags.Reset
	hc := flags.Reset

	var carryIn uint16 = 0
	if c.Flags().Carry {
		carryIn = 1
	}

	// Compute half-carry using 8-bit arithmetic.
	if ((*r)&0xF)+(val&0xF)+uint8(carryIn) > 0xF {
		hc = flags.Set
	}

	sum := uint16(*r) + uint16(val) + carryIn
	if sum > 0xFF {
		carry = flags.Set
	}

	*r = uint8(sum)
	if *r == 0 {
		zero = flags.Set
	}

	c.Flags().Set(zero, flags.Reset, hc, carry)
}
func subc8(c cpu.CPU, r *uint8, val uint8) {
	// Get the current carry bit (0 or 1). Adjust according to how your emulator stores flags.
	carryIn := uint8(0)
	if c.Flags().Carry { // assuming you have a helper for reading the carry flag
		carryIn = 1
	}

	// Use int16 arithmetic to account for potential borrow.
	diff := int16(*r) - int16(val) - int16(carryIn)
	carry := flags.Reset
	if diff < 0 {
		carry = flags.Set // Borrow occurred.
		diff += 256       // Wrap-around for 8-bit arithmetic.
	}

	// Check for half-borrow: compare the low nibble of *r with (val + carryIn)'s low nibble.
	hc := flags.Reset
	if ((*r) & 0xF) < ((val & 0xF) + carryIn) {
		hc = flags.Set
	}

	// Determine the result and the Z flag.
	result := uint8(diff)
	zero := flags.Reset
	if result == 0 {
		zero = flags.Set
	}

	*r = result
	// In subtraction with borrow, the N flag is always set.
	c.Flags().Set(zero, flags.Set, hc, carry)
}
