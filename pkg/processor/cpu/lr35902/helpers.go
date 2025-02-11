package lr35902

import (
	"math/bits"
)

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
func (c *LR35902) rst(addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.registers.pc + 1

	c.registers.sp -= 2
	c.bus.Write16(c.registers.sp, retAddr)

	c.registers.pc = addr - 1
}

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

// Jump
func (c *LR35902) jump(addr uint16, condition bool) {
	if condition {
		// Subtract 1 to account for PC increment after instruction fetch.
		c.registers.pc = addr - 1
	}
}
func (c *LR35902) jumpRelative(offset int8, condition bool) {
	if condition {
		c.registers.pc = uint16(int32(c.registers.pc) + int32(offset))
	}
}

// Helpers
func (c *LR35902) getImmediate8() uint8 {
	val := c.bus.Read(c.registers.pc + 1)
	c.registers.pc++
	return val
}

func (c *LR35902) getImmediate16() uint16 {
	high, low := c.bus.Read16(c.registers.pc + 1)
	c.registers.pc += 2
	return toRegisterPair(high, low)
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

// Subroutines
func (c *LR35902) ret(condition bool) {
	if !condition {
		return
	}

	// Pop the address from the stack in little endian order: low byte then high byte.
	high, low := c.bus.Read16(c.registers.sp)
	c.registers.sp += 2

	addr := toRegisterPair(high, low)

	c.registers.pc = addr - 1
}

func (c *LR35902) call(addr uint16, condition bool) {
	if !condition {
		return
	}

	retAddr := c.registers.pc + 1

	// Decrement SP by 2 and push return address using Write16.
	c.registers.sp -= 2
	c.bus.Write16(c.registers.sp, retAddr)

	c.registers.pc = addr - 1
}

// Other
func (c *LR35902) decimalAdjust() {
	// Copied from https://stackoverflow.com/a/57837042/9731890
	var t uint8

	if c.flags.HalfCarry || (c.registers.a&0x0F) > 9 {
		t++
	}

	if c.flags.Carry || c.registers.a > 0x99 {
		t += 2
		c.flags.Carry = true
	}

	// Builds final H flag
	if c.flags.Subtract && !c.flags.HalfCarry {
		c.flags.HalfCarry = false
	} else {
		if c.flags.Subtract && c.flags.HalfCarry {
			c.flags.HalfCarry = ((c.registers.a & 0x0F) < 6)
		} else {
			c.flags.HalfCarry = ((c.registers.a & 0x0F) >= 0x0A)
		}
	}

	switch t {
	case 1:
		if c.flags.Subtract {
			c.registers.a -= 6
		} else {
			c.registers.a += 6
		}
	case 2:
		if c.flags.Subtract {
			c.registers.a -= 0x60
		} else {
			c.registers.a += 0x60
		}
	case 3:
		if c.flags.Subtract {
			c.registers.a -= 0x66
		} else {
			c.registers.a += 0x66
		}
	}

	if c.registers.a == 0 {
		c.flags.Zero = true
	} else {
		c.flags.Zero = false
	}
}

// registerPair returns a 16-bit register pair from two 8-bit registers
// If you want BC, pass B and C in that order
func toRegisterPair(high, low uint8) uint16 {
	return uint16(high)<<8 | uint16(low)
}

// fromRegisterPair returns two 8-bit registers from a 16-bit register pair
// It returns it in low, high order. E.g. if you have BC (CB in little endian) it will return B, C
func fromRegisterPair(val uint16) (high uint8, low uint8) {
	return uint8(val >> 8), uint8(val)
}
