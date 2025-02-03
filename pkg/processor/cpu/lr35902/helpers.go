package lr35902

import (
	"math/bits"
)

// Memory access
func (c *LR35902) load8(r *uint8, val uint8) {
	*r = val
}
func (c *LR35902) load16(high, low *uint8, val uint16) {
	*high = uint8(val >> 8)
	*low = uint8(val)
}
func (c *LR35902) pop16(high, low *uint8) {
	*low = c.bus.Read(c.registers.sp)
	*high = c.bus.Read(c.registers.sp + 1)
	c.registers.sp += 2
}
func (c *LR35902) push16(high, low uint8) {
	c.registers.sp -= 2
	c.bus.Write(c.registers.sp, low)
	c.bus.Write(c.registers.sp+1, high)
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

	if (*r&0xf)-(val&0xf) < 0 {
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

	if c.flags.Carry {
		val++
	}

	sum, car := bits.Sub(uint(*r), uint(val), 0)
	if car > 0 {
		carry = Set
	}

	if (*r&0xf)-(val&0xf) < 0 {
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
	zero := Reset
	hc := Reset

	*r++

	if *r&0x10 == 0x10 {
		hc = Set
	}

	if *r == 0 {
		zero = 1
	}

	c.setFlags(zero, Reset, hc, Leave)
}
func (c *LR35902) inc16(high, low *uint8) {
	combined := toRegisterPair(*high, *low)
	combined++
	*high, *low = fromRegisterPair(combined)
}
func (c *LR35902) dec8(r *uint8) {
	*r--

	hc := Reset
	if *r&0x10 == 0x10 {
		hc = Set
	}

	zero := Reset
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
	retAddr := c.registers.pc + 1 // +1 because the PC is incremented after the instruction is fetched
	c.registers.sp -= 2
	// Push low byte first then high byte
	c.bus.Write(c.registers.sp, uint8(retAddr))
	c.bus.Write(c.registers.sp+1, uint8(retAddr>>8))
	c.registers.pc = addr - 1 // -1 because the PC is incremented after the instruction is fetched
}

// Rotates/Shifts
func (c *LR35902) rotate(r *uint8, left bool, carry bool) {
	var carried uint8

	if left {
		*r = *r<<1 | *r>>7
		if carry {
			carried = *r >> 7
		}
	} else {
		*r = *r>>1 | *r<<7
		if carry {
			carried = *r << 7
		}
	}

	c.setFlags(Reset, Reset, Reset, FlagState(carried))
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
		c.registers.pc += uint16(offset)
	}
}

// Helpers
func (c *LR35902) getImmediate8() uint8 {
	val := c.bus.Read(c.registers.pc + 1)
	c.registers.pc++
	return val
}

func (c *LR35902) getImmediate16() uint16 {
	low := c.bus.Read(c.registers.pc + 1)
	high := c.bus.Read(c.registers.pc + 2)
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

	// Pop the address from the stack: high byte first, then low byte
	high := c.bus.Read(c.registers.sp)
	low := c.bus.Read(c.registers.sp + 1)
	c.registers.sp += 2

	// Subtract 1 so that after the clock cycle increments PC, it equals the intended return address
	c.registers.pc = (uint16(high)<<8 | uint16(low)) - 1
}
func (c *LR35902) call(addr uint16, condition bool) {
	if !condition {
		return
	}

	// Calculate return address (current PC + 3: opcode plus two immediate bytes)
	retAddr := c.registers.pc + 3
	c.registers.sp -= 2
	// Push high byte then low byte of return address
	c.bus.Write(c.registers.sp, uint8(retAddr>>8))
	c.bus.Write(c.registers.sp+1, uint8(retAddr))
	// Adjust target address: subtract 1 to account for PC auto-increment
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

// Interrupts
func (c *LR35902) disableInterrupts() {
	c.io.Interrupts.Joypad = false
	c.io.Interrupts.Serial = false
	c.io.Interrupts.Timer = false
	c.io.Interrupts.LCD = false
	c.io.Interrupts.VBlank = false
}
func (c *LR35902) enableInterrupts() {
	c.io.Interrupts.Joypad = true
	c.io.Interrupts.Serial = true
	c.io.Interrupts.Timer = true
	c.io.Interrupts.LCD = true
	c.io.Interrupts.VBlank = true
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
