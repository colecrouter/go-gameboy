package lr35902

import (
	"encoding/binary"
	"math/bits"
)

// ReadFlag is a placeholder implementation to restore compilation; adjust logic as needed.
func (c *LR35902) ReadFlag(flag int) bool {
	return false
}

/*
	Helper functions
*/

func (c *LR35902) inc8(a *uint8) {
	// flags
	zero := reset
	hc := reset

	// Convert from Little Endian to uint
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(*a))
	var bnew = binary.LittleEndian
	temp := bnew.Uint16(buf)

	// Increment
	temp++

	// Convert back to buffer
	binary.LittleEndian.PutUint16(buf, temp)

	// Convert buffer to uint
	temp = binary.LittleEndian.Uint16(buf)

	// Assign
	*a = uint8(temp)

	// half-carry
	if *a&0x10 == 0x10 {
		hc = set
	}

	// zero
	if *a == 0 {
		zero = 1
	}

	c.setFlags(zero, reset, hc, leave)
}

func (c *LR35902) dec8(a *uint8) {
	// flags
	zero := reset
	hc := reset

	// Convert from Little Endian to uint
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(*a))
	var bnew = binary.LittleEndian
	temp := bnew.Uint16(buf)

	// Decrement
	temp--

	// Convert back to buffer
	binary.LittleEndian.PutUint16(buf, temp)

	// Convert buffer to uint
	temp = binary.LittleEndian.Uint16(buf)

	// Assign
	*a = uint8(temp)

	// half-carry
	if *a&0x10 == 0x10 {
		hc = set
	}

	// zero
	if *a == 0 {
		zero = 1
	}

	c.setFlags(zero, reset, hc, leave)
}

func (c *LR35902) dec16(a *uint8, b *uint8) {
	// flags
	zero := reset
	hc := reset

	// Convert from Little Endian to uint
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, toLong(*a, *b))
	var bnew = binary.LittleEndian
	temp := bnew.Uint16(buf)

	// Decrement
	temp--

	// Convert back to buffer
	binary.LittleEndian.PutUint16(buf, temp)

	// Convert buffer to uint
	temp = binary.LittleEndian.Uint16(buf)

	// Assign
	*a, *b = toShort(temp)

	// half-carry
	if *a&0x10 == 0x10 {
		hc = set
	}

	// zero
	if *a == 0 {
		zero = 1
	}

	c.setFlags(zero, reset, hc, leave)
}

func (c *LR35902) ld8(a *uint8, b *uint8) {
	*a = *b
}

func (c *LR35902) ld16(a1 *uint8, a2 *uint8, b1 *uint8, b2 *uint8) {
	*a1 = *b1
	*a2 = *b2
}

func (c *LR35902) ldi8(a *uint8, b *uint8) {
	*a = *b
}

func (c *LR35902) ldi16(a1 *uint8, a2 *uint8, b *uint16) {
	*a1 = uint8(*b)
	*a2 = uint8(*b << 8)
}

func (c *LR35902) add16(a1 *uint8, a2 *uint8, b1 *uint8, b2 *uint8) {
	a := uint16(*a2)<<8 | uint16(*a1)
	b := uint16(*b1)<<8 | uint16(*b2)

	carry := reset
	sum, car := bits.Add(uint(a), uint(b), 0)
	if car > 0 {
		carry = set
	}

	hc := reset
	if (((a & 0xf) + (b & 0xf)) & 0x10) == 0x10 {
		hc = set
	}

	*a1 = uint8(sum)
	*a2 = uint8(sum << 8)

	c.setFlags(leave, reset, hc, carry)
}

func (c *LR35902) add16s1(a1 *uint8, a2 *uint8, b *uint16) {
	a := uint16(*a2)<<8 | uint16(*a1)

	carry := reset
	sum, car := bits.Add(uint(a), uint(*b), 0)
	if car > 0 {
		carry = set
	}

	hc := reset
	if (((a & 0xf) + (*b & 0xf)) & 0x10) == 0x10 {
		hc = set
	}

	*a1 = uint8(sum)
	*a2 = uint8(sum << 8)

	c.setFlags(leave, reset, hc, carry)
}

func (c *LR35902) inc16(a *uint8, b *uint8) {
	zero, hc := reset, reset
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, toLong(*a, *b))
	tmp := binary.LittleEndian.Uint16(buf)
	tmp++
	binary.LittleEndian.PutUint16(buf, tmp)
	newVal := binary.LittleEndian.Uint16(buf)
	*a, *b = toShort(newVal)
	if (*a & 0x0F) == 0 {
		hc = set
	}
	if *a == 0 && *b == 0 {
		zero = set
	}
	c.setFlags(zero, reset, hc, leave)
}

func (c *LR35902) inc16s(a *uint16) {
	// Convert from Little Endian to uint
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, *a)
	var b = binary.LittleEndian
	temp := b.Uint16(buf)

	// Increment
	temp++

	// Convert back to buffer
	binary.LittleEndian.PutUint16(buf, temp)

	// Convert buffer to uint
	*a = binary.LittleEndian.Uint16(buf)
}

/*
	Instructions
*/

func (c *LR35902) nop() {
	c.done <- struct{}{}

}

func (c *LR35902) ldBcD16() {
	c.doClock <- struct{}{}
	a1 := <-c.clocking
	c.doClock <- struct{}{}
	a2 := <-c.clocking
	a := toLong(a1, a2)
	c.ldi16(&c.registers.B, &c.registers.C, &a)
	c.done <- struct{}{}
}

func (c *LR35902) ldpBcA() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.bus.Write(uint16(c.registers.C)<<8|uint16(c.registers.B), c.registers.A)
	c.done <- struct{}{}
}

func (c *LR35902) incBc() {
	c.inc16(&c.registers.B, &c.registers.C)
	c.done <- struct{}{}
}

func (c *LR35902) incB() {
	c.inc8(&c.registers.B)
	c.done <- struct{}{}
}

func (c *LR35902) decB() {
	c.dec8(&c.registers.B)
	c.done <- struct{}{}
}

func (c *LR35902) ldBD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ldi8(&c.registers.B, &a)
	c.done <- struct{}{}
}

func (c *LR35902) rlcA() {
	// TODO fix
	// msb
	carry := leave
	if c.registers.A>>7 == 0 {
		carry = reset
	} else {
		carry = set
	}

	bits.RotateLeft8(c.registers.A, 1)
	c.setFlags(reset, reset, reset, carry)
	c.done <- struct{}{}
}

func (c *LR35902) ldpA16SP() {
	c.doClock <- struct{}{}
	a1 := <-c.clocking
	c.doClock <- struct{}{}
	a2 := <-c.clocking
	c.doClock <- struct{}{}
	<-c.clocking
	c.doClock <- struct{}{}
	<-c.clocking
	a := toLong(a1, a2)
	c.bus.Write(a, uint8(c.registers.SP<<8))
	c.bus.Write(a+1, uint8(c.registers.SP))
	c.done <- struct{}{}
}

func (c *LR35902) addHlBc() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.add16(&c.registers.H, &c.registers.L, &c.registers.B, &c.registers.C)
	c.done <- struct{}{}
}

func (c *LR35902) ldApBc() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.registers.A = c.bus.Read(uint16(c.registers.C)<<8 | uint16(c.registers.B))
	c.done <- struct{}{}
}

func (c *LR35902) decBc() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.dec16(&c.registers.B, &c.registers.C)
	c.done <- struct{}{}
}

func (c *LR35902) incC() {
	c.inc8(&c.registers.C)
	c.done <- struct{}{}
}

func (c *LR35902) decC() {
	c.dec8(&c.registers.C)
	c.done <- struct{}{}
}

func (c *LR35902) ldCD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ld8(&c.registers.C, &a)
	c.done <- struct{}{}
}

func (c *LR35902) rrcA() {
	// msb
	carry := leave
	if c.registers.A>>7 == 0 {
		carry = reset
	} else {
		carry = set
	}

	bits.RotateLeft8(c.registers.A, -1)
	c.setFlags(reset, reset, reset, carry)
	c.done <- struct{}{}
}

// 0x1x
func (c *LR35902) stop() {
	// TODO
	// log.Println("stop")
	c.done <- struct{}{}
}

func (c *LR35902) ldDeD16() {
	c.doClock <- struct{}{}
	a1 := <-c.clocking
	c.doClock <- struct{}{}
	a2 := <-c.clocking
	a := toLong(a1, a2)
	c.ldi16(&c.registers.D, &c.registers.E, &a)
	c.done <- struct{}{}
}

func (c *LR35902) ldpDeA() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.bus.Write(uint16(c.registers.E)<<8|uint16(c.registers.D), c.registers.A)
	c.done <- struct{}{}
}

func (c *LR35902) incDe() {
	c.inc16(&c.registers.D, &c.registers.E)
	c.done <- struct{}{}
}

func (c *LR35902) incD() {
	c.inc8(&c.registers.D)
	c.done <- struct{}{}
}

func (c *LR35902) decD() {
	c.dec8(&c.registers.D)
	c.done <- struct{}{}
}

func (c *LR35902) ldDD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ldi8(&c.registers.D, &a)
	c.done <- struct{}{}
}

func (c *LR35902) rlA() {
	// TODO FIX
	// msb
	carry := leave
	if c.registers.A>>7 == 0 {
		carry = reset
	} else {
		carry = set
	}

	bits.RotateLeft8(c.registers.A, 1)
	c.setFlags(reset, reset, reset, carry)
	c.done <- struct{}{}
}

func (c *LR35902) jrS8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.doClock <- struct{}{}
	<-c.clocking
	c.registers.PC += uint16(a)
	c.done <- struct{}{}
}

func (c *LR35902) addHlDe() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.add16(&c.registers.H, &c.registers.L, &c.registers.D, &c.registers.E)
	c.done <- struct{}{}
}

func (c *LR35902) ldApDe() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.registers.A = c.bus.Read(uint16(c.registers.E)<<8 | uint16(c.registers.D))
	c.done <- struct{}{}
}

func (c *LR35902) decDe() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.dec16(&c.registers.D, &c.registers.E)
	c.done <- struct{}{}
}

func (c *LR35902) incE() {
	c.inc8(&c.registers.E)
	c.done <- struct{}{}
}

func (c *LR35902) decE() {
	c.dec8(&c.registers.E)
	c.done <- struct{}{}
}

func (c *LR35902) ldED8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ld8(&c.registers.E, &a)
	c.done <- struct{}{}
}

func (c *LR35902) rrA() {
	// TODO fix
	// msb
	carry := leave
	if c.registers.A>>7 == 0 {
		carry = reset
	} else {
		carry = set
	}

	bits.RotateLeft8(c.registers.A, -1)
	c.setFlags(reset, reset, reset, carry)
	c.done <- struct{}{}
}

// 0x2x
func (c *LR35902) jrNzS8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	if !c.ReadFlag(0) {
		c.registers.PC += uint16(a)
		<-c.clocking
	}
	c.done <- struct{}{}
}

func (c *LR35902) ldHlD16() {
	c.doClock <- struct{}{}
	a1 := <-c.clocking
	c.doClock <- struct{}{}
	a2 := <-c.clocking
	a := toLong(a1, a2)
	c.ldi16(&c.registers.H, &c.registers.L, &a)
	c.done <- struct{}{}
}

func (c *LR35902) ldpHlA() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.bus.Write(toLong(c.registers.H, c.registers.L), c.registers.A)
	c.inc16(&c.registers.L, &c.registers.H)
	c.done <- struct{}{}
}

func (c *LR35902) incHl() {
	c.inc16(&c.registers.H, &c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) incH() {
	c.inc8(&c.registers.H)
	c.done <- struct{}{}
}

func (c *LR35902) decH() {
	c.dec8(&c.registers.H)
	c.done <- struct{}{}
}

func (c *LR35902) ldHD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ldi8(&c.registers.H, &a)
	c.done <- struct{}{}
}

func (c *LR35902) ddA() {
	// TODO
	c.setFlags(reset, reset, reset, reset)
	c.done <- struct{}{}
}

func (c *LR35902) jrZS8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	if c.ReadFlag(0) {
		c.registers.PC += uint16(a)
		<-c.clocking
	}
	c.done <- struct{}{}
}

func (c *LR35902) addHlHl() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.add16(&c.registers.H, &c.registers.L, &c.registers.H, &c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) ldApHlp() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.registers.A = c.bus.Read(uint16(c.registers.E)<<8 | uint16(c.registers.D))
	c.incHl()
	c.done <- struct{}{}
}

func (c *LR35902) decHl() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.dec16(&c.registers.H, &c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) incL() {
	c.inc8(&c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) decL() {
	c.dec8(&c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) ldLD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	c.ld8(&c.registers.L, &a)
	c.done <- struct{}{}
}

func (c *LR35902) cpl() {
	c.registers.A = ^c.registers.A
	c.setFlags(leave, set, set, leave)
	c.done <- struct{}{}
}

// 0x3x
func (c *LR35902) jrNcS8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	if !c.ReadFlag(3) {
		c.registers.PC += uint16(a)
		<-c.clocking
	}
	c.done <- struct{}{}
}

func (c *LR35902) ldSpD16() {
	c.doClock <- struct{}{}
	a1 := <-c.clocking
	c.doClock <- struct{}{}
	a2 := <-c.clocking
	a := toLong(a1, a2)
	c.registers.SP = a // TODO is this broken?
	c.done <- struct{}{}
}

func (c *LR35902) ldnHlA() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.bus.Write(toLong(c.registers.H, c.registers.L), c.registers.A)
	c.dec16(&c.registers.H, &c.registers.L)
	c.done <- struct{}{}
}

func (c *LR35902) incSp() {
	c.registers.SP++
	c.done <- struct{}{}
}

func (c *LR35902) incpHl() {
	addr := toLong(c.registers.H, c.registers.L)
	value := c.bus.Read(addr)
	c.inc8(&value)
	c.bus.Write(addr, value)
	c.done <- struct{}{}
}

func (c *LR35902) decpHl() {
	addr := toLong(c.registers.H, c.registers.L)
	value := c.bus.Read(addr)
	c.dec8(&value)
	c.bus.Write(addr, value)
	c.done <- struct{}{}
}

func (c *LR35902) ldpHlD8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	addr := c.bus.Read(toLong(c.registers.H, c.registers.L))
	value := c.bus.Read(uint16(addr))
	c.ldi8(&value, &a)
	c.done <- struct{}{}
}

func (c *LR35902) scf() {
	// TODO
	c.setFlags(reset, reset, reset, set)
	c.done <- struct{}{}
}

func (c *LR35902) jrCS8() {
	c.doClock <- struct{}{}
	a := <-c.clocking
	if c.ReadFlag(3) {
		c.registers.PC += uint16(a)
		<-c.clocking
	}
	c.done <- struct{}{}
}

func (c *LR35902) addHlSp() {
	c.doClock <- struct{}{}
	<-c.clocking
	c.add16s1(&c.registers.H, &c.registers.L, &c.registers.SP)
	c.done <- struct{}{}
}
