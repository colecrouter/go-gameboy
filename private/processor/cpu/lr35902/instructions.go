package lr35902

type instruction struct {
	op func(c *LR35902) // operation
	c  int              // cycles
	p  int              // pc advance
}

var instructions = [0x100]instruction{
	// NOP
	0x00: {c: 4, p: 1, op: func(c *LR35902) {}},
	// LD BC,d16
	0x01: {c: 12, p: 3, op: func(c *LR35902) {
		c.load16(&c.Registers.b, &c.Registers.c, c.getImmediate16())
	}},
	// LD (BC),A
	0x02: {c: 8, p: 1, op: func(c *LR35902) {
		address := toRegisterPair(c.Registers.b, c.Registers.c)
		c.bus.Write(address, c.Registers.a)
	}},
	// INC BC
	0x03: {c: 8, p: 1, op: func(c *LR35902) {
		c.inc16(&c.Registers.b, &c.Registers.c)
	}},
	// INC B
	0x04: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.b)
	}},
	// DEC B
	0x05: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.b)
	}},
	// LD B,d8
	0x06: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.getImmediate8())
	}},
	// RLCA: rotate A left circularly (ignore previous carry)
	0x07: {c: 4, p: 1, op: func(c *LR35902) {
		c.rotate(&c.Registers.a, true, false, false)
	}},
	// LD (a16),SP
	0x08: {c: 20, p: 3, op: func(c *LR35902) {
		address := c.getImmediate16()
		bytes := uint(c.Registers.sp)
		c.bus.Write(address, uint8(bytes&0xFF))
		c.bus.Write(address+1, uint8(bytes>>8))
	}},
	// ADD HL,BC
	0x09: {c: 8, p: 1, op: func(c *LR35902) {
		c.add16(&c.Registers.h, &c.Registers.l, c.Registers.b, c.Registers.c)
	}},
	// LD A,(BC)
	0x0A: {c: 8, p: 1, op: func(c *LR35902) {
		address := uint16(c.Registers.b)<<8 | uint16(c.Registers.c)
		c.Registers.a = c.bus.Read(address)
	}},
	// DEC BC
	0x0B: {c: 8, p: 1, op: func(c *LR35902) {
		c.dec16(&c.Registers.b, &c.Registers.c)
	}},
	// INC C
	0x0C: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.c)
	}},
	// DEC C
	0x0D: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.c)
	}},
	// LD C,d8
	0x0E: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.getImmediate8())
	}},
	// RRCA: rotate A right circularly (ignore previous carry)
	0x0F: {c: 4, p: 1, op: func(c *LR35902) {
		c.rotate(&c.Registers.a, false, false, false)
	}},
	// STOP 0
	0x10: {c: 4, p: 2, op: func(c *LR35902) {
		// NOP
	}},
	// LD DE,d16
	0x11: {c: 12, p: 3, op: func(c *LR35902) {
		c.load16(&c.Registers.d, &c.Registers.e, c.getImmediate16())
	}},
	// LD (DE),A
	0x12: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.d, c.Registers.e)
		c.bus.Write(addr, c.Registers.a)
	}},
	// INC DE
	0x13: {c: 8, p: 1, op: func(c *LR35902) {
		c.inc16(&c.Registers.d, &c.Registers.e)
	}},
	// INC D
	0x14: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.d)
	}},
	// DEC D
	0x15: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.d)
	}},
	// LD D,d8
	0x16: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.getImmediate8())
	}},
	// RLA: rotate A left through carry (use previous carry)
	0x17: {c: 4, p: 1, op: func(c *LR35902) {
		c.rotate(&c.Registers.a, true, true, false)
	}},
	// JR r8
	0x18: {c: 12, p: 0, op: func(c *LR35902) {
		c.jumpRelative(int8(c.getImmediate8()), true)
	}},
	// ADD HL,DE
	0x19: {c: 8, p: 1, op: func(c *LR35902) {
		c.add16(&c.Registers.h, &c.Registers.l, c.Registers.d, c.Registers.e)
	}},
	// LD A,(DE)
	0x1A: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.d, c.Registers.e)
		c.Registers.a = c.bus.Read(addr)
	}},
	// DEC DE
	0x1B: {c: 8, p: 1, op: func(c *LR35902) {
		c.dec16(&c.Registers.d, &c.Registers.e)
	}},
	// INC E
	0x1C: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.e)
	}},
	// DEC E
	0x1D: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.e)
	}},
	// LD E,d8
	0x1E: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.getImmediate8())
	}},
	// RRA: rotate A right through carry (use previous carry)
	0x1F: {c: 4, p: 1, op: func(c *LR35902) {
		c.rotate(&c.Registers.a, false, true, false)
	}},

	// JR NZ,r8
	0x20: {c: 12, p: 0, op: func(c *LR35902) {
		c.jumpRelative(int8(c.getImmediate8()), !c.flags.Zero)
	}},
	// LD HL,d16
	0x21: {c: 12, p: 3, op: func(c *LR35902) {
		c.load16(&c.Registers.h, &c.Registers.l, c.getImmediate16())
	}},
	// LD (HL+),A
	0x22: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l) // updated order
		c.bus.Write(addr, c.Registers.a)
		c.inc16(&c.Registers.h, &c.Registers.l) // updated order
	}},
	// INC HL
	0x23: {c: 8, p: 1, op: func(c *LR35902) {
		c.inc16(&c.Registers.h, &c.Registers.l)
	}},
	// INC H
	0x24: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.h)
	}},
	// DEC H
	0x25: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.h)
	}},
	// LD H,d8
	0x26: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.getImmediate8())
	}},
	// DAA
	0x27: {c: 4, p: 1, op: func(c *LR35902) {
		c.decimalAdjust()
	}},
	// JR Z,r8
	0x28: {c: 12, p: 0, op: func(c *LR35902) {
		c.jumpRelative(int8(c.getImmediate8()), c.flags.Zero)
	}},
	// ADD HL,HL
	0x29: {c: 8, p: 1, op: func(c *LR35902) {
		c.add16(&c.Registers.h, &c.Registers.l, c.Registers.h, c.Registers.l)
	}},
	// LD A,(HL+)
	0x2A: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l) // updated order
		c.Registers.a = c.bus.Read(addr)
		c.inc16(&c.Registers.h, &c.Registers.l)
	}},
	// DEC HL
	0x2B: {c: 8, p: 1, op: func(c *LR35902) {
		c.dec16(&c.Registers.h, &c.Registers.l)
	}},
	// INC L
	0x2C: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.l)
	}},
	// DEC L
	0x2D: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.l)
	}},
	// LD L,d8
	0x2E: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.getImmediate8())
	}},
	// CPL
	0x2F: {c: 4, p: 1, op: func(c *LR35902) {
		c.Registers.a = ^c.Registers.a
		c.setFlags(Leave, Set, Set, Leave)
	}},
	// JR NC,r8
	0x30: {c: 12, p: 0, op: func(c *LR35902) {
		c.jumpRelative(int8(c.getImmediate8()), !c.flags.Carry)
	}},
	// LD SP,d16
	0x31: {c: 12, p: 3, op: func(c *LR35902) {
		c.Registers.sp = c.getImmediate16()
	}},
	// LD (HL-),A
	0x32: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l) // updated order
		c.bus.Write(addr, c.Registers.a)
		c.dec16(&c.Registers.h, &c.Registers.l) // updated order
	}},
	// INC SP
	0x33: {c: 8, p: 1, op: func(c *LR35902) {
		c.Registers.sp++
	}},
	// INC (HL)
	0x34: {c: 12, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.inc8(&val)
		c.bus.Write(addr, val)
	}},
	// DEC (HL)
	0x35: {c: 12, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.dec8(&val)
		c.bus.Write(addr, val)
	}},
	// LD (HL),d8
	0x36: {c: 12, p: 2, op: func(c *LR35902) {
		val := c.getImmediate8()
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, val)
	}},
	// SCF
	0x37: {c: 4, p: 1, op: func(c *LR35902) {
		c.setFlags(Leave, Reset, Reset, Set)
	}},
	// JR C,r8
	0x38: {c: 12, p: 0, op: func(c *LR35902) {
		c.jumpRelative(int8(c.getImmediate8()), c.flags.Carry)
	}},
	// ADD HL,SP
	0x39: {c: 8, p: 1, op: func(c *LR35902) {
		high, low := fromRegisterPair(c.Registers.sp)
		c.add16(&c.Registers.h, &c.Registers.l, high, low)
	}},
	// LD A,(HL-)
	0x3A: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l) // updated order
		c.Registers.a = c.bus.Read(addr)
		c.dec16(&c.Registers.h, &c.Registers.l)
	}},
	// DEC SP
	0x3B: {c: 8, p: 1, op: func(c *LR35902) {
		c.Registers.sp--
	}},
	// INC A
	0x3C: {c: 4, p: 1, op: func(c *LR35902) {
		c.inc8(&c.Registers.a)
	}},
	// DEC A
	0x3D: {c: 4, p: 1, op: func(c *LR35902) {
		c.dec8(&c.Registers.a)
	}},
	// LD A,d8
	0x3E: {c: 8, p: 2, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.getImmediate8())
	}},
	// CCF
	0x3F: {c: 4, p: 1, op: func(c *LR35902) {
		carry := Reset
		if !c.flags.Carry {
			carry = Set
		}
		c.setFlags(Leave, Reset, Reset, carry)
	}},
	// LD B,B
	0x40: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.b)
	}},
	// LD B,C
	0x41: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.c)
	}},
	// LD B,D
	0x42: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.d)
	}},
	// LD B,E
	0x43: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.e)
	}},
	// LD B,H
	0x44: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.h)
	}},
	// LD B,L
	0x45: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.l)
	}},
	// LD B,(HL)
	0x46: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.b = c.bus.Read(addr)
	}},
	// LD B,A
	0x47: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.b, c.Registers.a)
	}},
	// LD C,B
	0x48: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.b)
	}},
	// LD C,C
	0x49: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.c)
	}},
	// LD C,D
	0x4A: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.d)
	}},
	// LD C,E
	0x4B: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.e)
	}},
	// LD C,H
	0x4C: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.h)
	}},
	// LD C,L
	0x4D: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.l)
	}},
	// LD C,(HL)
	0x4E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.c = c.bus.Read(addr)
	}},
	// LD C,A
	0x4F: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.c, c.Registers.a)
	}},
	// LD D,B
	0x50: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.b)
	}},
	// LD D,C
	0x51: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.c)
	}},
	// LD D,D
	0x52: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.d)
	}},
	// LD D,E
	0x53: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.e)
	}},
	// LD D,H
	0x54: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.h)
	}},
	// LD D,L
	0x55: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.l)
	}},
	// LD D,(HL)
	0x56: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.d = c.bus.Read(addr)
	}},
	// LD D,A
	0x57: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.d, c.Registers.a)
	}},
	// LD E,B
	0x58: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.b)
	}},
	// LD E,C
	0x59: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.c)
	}},
	// LD E,D
	0x5A: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.d)
	}},
	// LD E,E
	0x5B: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.e)
	}},
	// LD E,H
	0x5C: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.h)
	}},
	// LD E,L
	0x5D: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.l)
	}},
	// LD E,(HL)
	0x5E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.e = c.bus.Read(addr)
	}},
	// LD E,A
	0x5F: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.e, c.Registers.a)
	}},
	// LD H,B
	0x60: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.b)
	}},
	// LD H,C
	0x61: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.c)
	}},
	// LD H,D
	0x62: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.d)
	}},
	// LD H,E
	0x63: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.e)
	}},
	// LD H,H
	0x64: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.h)
	}},
	// LD H,L
	0x65: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.l)
	}},
	// LD H,(HL)
	0x66: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.h = c.bus.Read(addr)
	}},
	// LD H,A
	0x67: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.h, c.Registers.a)
	}},
	// LD L,B
	0x68: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.b)
	}},
	// LD L,C
	0x69: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.c)
	}},
	// LD L,D
	0x6A: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.d)
	}},
	// LD L,E
	0x6B: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.e)
	}},
	// LD L,H
	0x6C: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.h)
	}},
	// LD L,L
	0x6D: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.l)
	}},
	// LD L,(HL)
	0x6E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.l = c.bus.Read(addr)
	}},
	// LD L,A
	0x6F: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.l, c.Registers.a)
	}},
	// LD (HL),B
	0x70: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.b)
	}},
	// LD (HL),C
	0x71: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.c)
	}},
	// LD (HL),D
	0x72: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.d)
	}},
	// LD (HL),E
	0x73: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.e)
	}},
	// LD (HL),H
	0x74: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.h)
	}},
	// LD (HL),L
	0x75: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.l)
	}},
	// HALT
	0x76: {c: 4, p: 1, op: func(c *LR35902) {
		c.halted = true
	}},
	// LD (HL),A
	0x77: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.bus.Write(addr, c.Registers.a)
	}},
	// LD A,B
	0x78: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.b)
	}},
	// LD A,C
	0x79: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.c)
	}},
	// LD A,D
	0x7A: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.d)
	}},
	// LD A,E
	0x7B: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.e)
	}},
	// LD A,H
	0x7C: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.h)
	}},
	// LD A,L
	0x7D: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.l)
	}},
	// LD A,(HL)
	0x7E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		c.Registers.a = c.bus.Read(addr)
	}},
	// LD A,A
	0x7F: {c: 4, p: 1, op: func(c *LR35902) {
		c.load8(&c.Registers.a, c.Registers.a)
	}},
	// ADD A,B
	0x80: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.b)
	}},
	// ADD A,C
	0x81: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.c)
	}},
	// ADD A,D
	0x82: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.d)
	}},
	// ADD A,E
	0x83: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.e)
	}},
	// ADD A,H
	0x84: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.h)
	}},
	// ADD A,L
	0x85: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.l)
	}},
	// ADD A,(HL)
	0x86: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.add8(&c.Registers.a, val)
	}},
	// ADD A,A
	0x87: {c: 4, p: 1, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.Registers.a)
	}},
	// ADC A,B
	0x88: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.b)
	}},
	// ADC A,C
	0x89: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.c)
	}},
	// ADC A,D
	0x8A: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.d)
	}},
	// ADC A,E
	0x8B: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.e)
	}},
	// ADC A,H
	0x8C: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.h)
	}},
	// ADC A,L
	0x8D: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.l)
	}},
	// ADC A,(HL)
	0x8E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.addc8(&c.Registers.a, val)
	}},
	// ADC A,A
	0x8F: {c: 4, p: 1, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.Registers.a)
	}},
	// SUB B
	0x90: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.b)
	}},
	// SUB C
	0x91: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.c)
	}},
	// SUB D
	0x92: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.d)
	}},
	// SUB E
	0x93: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.e)
	}},
	// SUB H
	0x94: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.h)
	}},
	// SUB L
	0x95: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.l)
	}},
	// SUB (HL)
	0x96: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.sub8(&c.Registers.a, val)
	}},
	// SUB A
	0x97: {c: 4, p: 1, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.Registers.a)
	}},
	// SBC A,B
	0x98: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.b)
	}},
	// SBC A,C
	0x99: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.c)
	}},
	// SBC A,D
	0x9A: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.d)
	}},
	// SBC A,E
	0x9B: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.e)
	}},
	// SBC A,H
	0x9C: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.h)
	}},
	// SBC A,L
	0x9D: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.l)
	}},
	// SBC A,(HL)
	0x9E: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.subc8(&c.Registers.a, val)
	}},
	// SBC A,A
	0x9F: {c: 4, p: 1, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.Registers.a)
	}},
	// AND B
	0xA0: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.b)
	}},
	// AND C
	0xA1: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.c)
	}},
	// AND D
	0xA2: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.d)
	}},
	// AND E
	0xA3: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.e)
	}},
	// AND H
	0xA4: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.h)
	}},
	// AND L
	0xA5: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.l)
	}},
	// AND (HL)
	0xA6: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.and8(&c.Registers.a, val)
	}},
	// AND A
	0xA7: {c: 4, p: 1, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.Registers.a)
	}},
	// XOR B
	0xA8: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.b)
	}},
	// XOR C
	0xA9: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.c)
	}},
	// XOR D
	0xAA: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.d)
	}},
	// XOR E
	0xAB: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.e)
	}},
	// XOR H
	0xAC: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.h)
	}},
	// XOR L
	0xAD: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.l)
	}},
	// XOR (HL)
	0xAE: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.xor8(&c.Registers.a, val)
	}},
	// XOR A
	0xAF: {c: 4, p: 1, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.Registers.a)
	}},
	// OR B
	0xB0: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.b)
	}},
	// OR C
	0xB1: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.c)
	}},
	// OR D
	0xB2: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.d)
	}},
	// OR E
	0xB3: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.e)
	}},
	// OR H
	0xB4: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.h)
	}},
	// OR L
	0xB5: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.l)
	}},
	// OR (HL)
	0xB6: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.or8(&c.Registers.a, val)
	}},
	// OR A
	0xB7: {c: 4, p: 1, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.Registers.a)
	}},
	// CP B
	0xB8: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.b)
	}},
	// CP C
	0xB9: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.c)
	}},
	// CP D
	0xBA: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.d)
	}},
	// CP E
	0xBB: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.e)
	}},
	// CP H
	0xBC: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.h)
	}},
	// CP L
	0xBD: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.l)
	}},
	// CP (HL)
	0xBE: {c: 8, p: 1, op: func(c *LR35902) {
		addr := toRegisterPair(c.Registers.h, c.Registers.l)
		val := c.bus.Read(addr)
		c.cp8(c.Registers.a, val)
		a := 1
		_ = a
	}},
	// CP A
	0xBF: {c: 4, p: 1, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.Registers.a)
	}},
	// RET NZ
	0xC0: {c: 8, p: 0, op: func(c *LR35902) {
		c.ret(!c.flags.Zero)
	}},
	// POP BC
	0xC1: {c: 12, p: 1, op: func(c *LR35902) {
		c.pop16(&c.Registers.b, &c.Registers.c)
	}},
	// JP NZ,a16
	0xC2: {c: 12, p: 0, op: func(c *LR35902) {
		c.jump(c.getImmediate16(), !c.flags.Zero)
	}},
	// JP a16
	0xC3: {c: 16, p: 0, op: func(c *LR35902) {
		c.jump(c.getImmediate16(), true)
	}},
	// CALL NZ,a16
	0xC4: {c: 12, p: 0, op: func(c *LR35902) {
		c.call(c.getImmediate16(), !c.flags.Zero)
	}},
	// PUSH BC
	0xC5: {c: 16, p: 1, op: func(c *LR35902) {
		c.push16(c.Registers.b, c.Registers.c)
	}},
	// ADD A,d8
	0xC6: {c: 8, p: 2, op: func(c *LR35902) {
		c.add8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 00H
	0xC7: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x00)
	}},
	// RET Z
	0xC8: {c: 8, p: 0, op: func(c *LR35902) {
		c.ret(c.flags.Zero)
	}},
	// RET
	0xC9: {c: 16, p: 0, op: func(c *LR35902) {
		c.ret(true)
	}},
	// JP Z,a16
	0xCA: {c: 12, p: 0, op: func(c *LR35902) {
		c.jump(c.getImmediate16(), c.flags.Zero)
	}},
	// PREFIX CB
	0xCB: {c: 4, p: 1, op: func(c *LR35902) {
		c.cb = true
	}},
	// CALL Z,a16
	0xCC: {c: 12, p: 0, op: func(c *LR35902) {
		c.call(c.getImmediate16(), c.flags.Zero)
	}},
	// CALL a16
	0xCD: {c: 24, p: 0, op: func(c *LR35902) {
		c.call(c.getImmediate16(), true)
	}},
	// ADC A,d8
	0xCE: {c: 8, p: 2, op: func(c *LR35902) {
		c.addc8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 08H
	0xCF: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x08)
	}},
	// RET NC
	0xD0: {c: 8, p: 0, op: func(c *LR35902) {
		c.ret(!c.flags.Carry)
	}},
	// POP DE
	0xD1: {c: 12, p: 1, op: func(c *LR35902) {
		c.pop16(&c.Registers.d, &c.Registers.e)
	}},
	// JP NC,a16
	0xD2: {c: 12, p: 0, op: func(c *LR35902) {
		c.jump(c.getImmediate16(), !c.flags.Carry)
	}},
	// INVALID
	// CALL NC,a16
	0xD4: {c: 12, p: 0, op: func(c *LR35902) {
		c.call(c.getImmediate16(), !c.flags.Carry)
	}},
	// PUSH DE
	0xD5: {c: 16, p: 1, op: func(c *LR35902) {
		c.push16(c.Registers.d, c.Registers.e)
	}},
	// SUB d8
	0xD6: {c: 8, p: 2, op: func(c *LR35902) {
		c.sub8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 10H
	0xD7: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x10)
	}},
	// RET C
	0xD8: {c: 8, p: 0, op: func(c *LR35902) {
		c.ret(c.flags.Carry)
	}},
	// RETI
	0xD9: {c: 16, p: 0, op: func(c *LR35902) {
		c.ret(true)
		c.ime = true
	}},
	// JP C,a16
	0xDA: {c: 12, p: 0, op: func(c *LR35902) {
		c.jump(c.getImmediate16(), c.flags.Carry)
	}},
	// INVALID
	// CALL C,a16
	0xDC: {c: 12, p: 0, op: func(c *LR35902) {
		c.call(c.getImmediate16(), c.flags.Carry)
	}},
	// INVALID
	// SBC A,d8
	0xDE: {c: 8, p: 2, op: func(c *LR35902) {
		c.subc8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 18H
	0xDF: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x18)
	}},
	// LDH (a8),A
	0xE0: {c: 12, p: 2, op: func(c *LR35902) {
		addr := 0xFF00 + uint16(c.getImmediate8())
		c.load8Mem(c.Registers.a, addr)
	}},
	// POP HL
	0xE1: {c: 12, p: 1, op: func(c *LR35902) {
		c.pop16(&c.Registers.h, &c.Registers.l)
	}},
	// LD (C),A
	0xE2: {c: 8, p: 1, op: func(c *LR35902) {
		addr := 0xFF00 + uint16(c.Registers.c)
		c.load8Mem(c.Registers.a, addr)
	}},
	// INVALID
	// INVALID
	// PUSH HL
	0xE5: {c: 16, p: 1, op: func(c *LR35902) {
		c.push16(c.Registers.h, c.Registers.l)
	}},
	// AND d8
	0xE6: {c: 8, p: 2, op: func(c *LR35902) {
		c.and8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 20H
	0xE7: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x20)
	}},
	// ADD SP,r8
	0xE8: {c: 16, p: 2, op: func(c *LR35902) {
		c.addSPr8()
	}},
	// JP (HL)
	0xE9: {c: 4, p: 0, op: func(c *LR35902) {
		c.Registers.PC = toRegisterPair(c.Registers.h, c.Registers.l)
	}},
	// LD (a16),A
	0xEA: {c: 16, p: 3, op: func(c *LR35902) {
		addr := c.getImmediate16()
		c.bus.Write(addr, c.Registers.a)
	}},
	// INVALID
	// INVALID
	// INVALID
	// XOR d8
	0xEE: {c: 8, p: 2, op: func(c *LR35902) {
		c.xor8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 28H
	0xEF: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x28)
	}},
	// LDH A,(a8)
	0xF0: {c: 12, p: 2, op: func(c *LR35902) {
		addr := 0xFF00 + uint16(c.getImmediate8())
		c.Registers.a = c.bus.Read(addr)
	}},
	// POP AF
	0xF1: {c: 12, p: 1, op: func(c *LR35902) {
		c.popAF()
	}},
	// LD A,(C)
	0xF2: {c: 8, p: 1, op: func(c *LR35902) {
		addr := 0xFF00 + uint16(c.Registers.c)
		c.Registers.a = c.bus.Read(addr)
	}},
	// DI
	0xF3: {c: 4, p: 1, op: func(c *LR35902) {
		// c.ime = false
	}},
	// INVALID
	// PUSH AF
	0xF5: {c: 16, p: 1, op: func(c *LR35902) {
		reg := c.flags.Read()
		c.push16(c.Registers.a, reg)
	}},
	// OR d8
	0xF6: {c: 8, p: 2, op: func(c *LR35902) {
		c.or8(&c.Registers.a, c.getImmediate8())
	}},
	// RST 30H
	0xF7: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x30)
	}},
	// LD HL,SP+r8
	0xF8: {c: 12, p: 2, op: func(c *LR35902) {
		immediate := int8(c.getImmediate8())
		c.loadHLSPOffset(immediate)
	}},
	// LD SP,HL
	0xF9: {c: 8, p: 1, op: func(c *LR35902) {
		c.Registers.sp = toRegisterPair(c.Registers.h, c.Registers.l)
	}},
	// LD A,(a16)
	0xFA: {c: 16, p: 3, op: func(c *LR35902) {
		addr := c.getImmediate16()
		c.Registers.a = c.bus.Read(addr)
	}},
	// EI
	0xFB: {c: 4, p: 1, op: func(c *LR35902) {
		c.eiDelay = 2
	}},
	// INVALID
	// INVALID
	// CP d8
	0xFE: {c: 8, p: 2, op: func(c *LR35902) {
		c.cp8(c.Registers.a, c.getImmediate8())
	}},
	// RST 38H
	0xFF: {c: 16, p: 0, op: func(c *LR35902) {
		c.rst(0x38)
	}},
}
