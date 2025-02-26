package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

type Instruction struct {
	OP func(c cpu.CPU) // Operation
	P  uint            // PC advance
}

var Instructions = [0x100]Instruction{
	// NOP
	0x00: {P: 1, OP: func(c cpu.CPU) {}},
	// LD BC,d16
	0x01: {P: 3, OP: func(c cpu.CPU) {
		load16(c, &c.Registers().B, &c.Registers().C, c.GetImmediate16())
	}},
	// LD (BC),A
	0x02: {P: 1, OP: func(c cpu.CPU) {
		address := cpu.ToRegisterPair(c.Registers().B, c.Registers().C)
		c.Write(address, c.Registers().A)
	}},
	// INC BC
	0x03: {P: 1, OP: func(c cpu.CPU) {
		inc16(c, &c.Registers().B, &c.Registers().C)
	}},
	// INC B
	0x04: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().B)
	}},
	// DEC B
	0x05: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().B)
	}},
	// LD B,d8
	0x06: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.GetImmediate8())
	}},
	// RLCA: rotate A left circularly (ignore previous carry)
	0x07: {P: 1, OP: func(c cpu.CPU) {
		rotate(c, &c.Registers().A, true, false, false)
	}},
	// LD (a16),SP
	0x08: {P: 3, OP: func(c cpu.CPU) {
		address := c.GetImmediate16()
		bytes := uint(c.Registers().SP)
		c.Write(address, uint8(bytes&0xFF))
		c.Write(address+1, uint8(bytes>>8))
	}},
	// ADD HL,BC
	0x09: {P: 1, OP: func(c cpu.CPU) {
		add16(c, &c.Registers().H, &c.Registers().L, c.Registers().B, c.Registers().C)
	}},
	// LD A,(BC)
	0x0A: {P: 1, OP: func(c cpu.CPU) {
		address := uint16(c.Registers().B)<<8 | uint16(c.Registers().C)
		c.Registers().A = c.Read(address)
	}},
	// DEC BC
	0x0B: {P: 1, OP: func(c cpu.CPU) {
		dec16(c, &c.Registers().B, &c.Registers().C)
	}},
	// INC C
	0x0C: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().C)
	}},
	// DEC C
	0x0D: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().C)
	}},
	// LD C,d8
	0x0E: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.GetImmediate8())
	}},
	// RRCA: rotate A right circularly (ignore previous carry)
	0x0F: {P: 1, OP: func(c cpu.CPU) {
		rotate(c, &c.Registers().A, false, false, false)
	}},
	// STOP 0
	0x10: {P: 2, OP: func(c cpu.CPU) {
		// NOP
	}},
	// LD DE,d16
	0x11: {P: 3, OP: func(c cpu.CPU) {
		load16(c, &c.Registers().D, &c.Registers().E, c.GetImmediate16())
	}},
	// LD (DE),A
	0x12: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().D, c.Registers().E)
		c.Write(addr, c.Registers().A)
	}},
	// INC DE
	0x13: {P: 1, OP: func(c cpu.CPU) {
		inc16(c, &c.Registers().D, &c.Registers().E)
	}},
	// INC D
	0x14: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().D)
	}},
	// DEC D
	0x15: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().D)
	}},
	// LD D,d8
	0x16: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.GetImmediate8())
	}},
	// RLA: rotate A left through carry (use previous carry)
	0x17: {P: 1, OP: func(c cpu.CPU) {
		rotate(c, &c.Registers().A, true, true, false)
	}},
	// JR r8
	0x18: {P: 0, OP: func(c cpu.CPU) {
		jumpRelative(c, int8(c.GetImmediate8()), true)
	}},
	// ADD HL,DE
	0x19: {P: 1, OP: func(c cpu.CPU) {
		add16(c, &c.Registers().H, &c.Registers().L, c.Registers().D, c.Registers().E)
	}},
	// LD A,(DE)
	0x1A: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().D, c.Registers().E)
		c.Registers().A = c.Read(addr)
	}},
	// DEC DE
	0x1B: {P: 1, OP: func(c cpu.CPU) {
		dec16(c, &c.Registers().D, &c.Registers().E)
	}},
	// INC E
	0x1C: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().E)
	}},
	// DEC E
	0x1D: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().E)
	}},
	// LD E,d8
	0x1E: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.GetImmediate8())
	}},
	// RRA: rotate A right through carry (use previous carry)
	0x1F: {P: 1, OP: func(c cpu.CPU) {
		rotate(c, &c.Registers().A, false, true, false)
	}},

	// JR NZ,r8
	0x20: {P: 0, OP: func(c cpu.CPU) {
		jumpRelative(c, int8(c.GetImmediate8()), !c.Flags().Zero)
	}},
	// LD HL,d16
	0x21: {P: 3, OP: func(c cpu.CPU) {
		load16(c, &c.Registers().H, &c.Registers().L, c.GetImmediate16())
	}},
	// LD (HL+),A
	0x22: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Write(addr, c.Registers().A)
		inc16(c, &c.Registers().H, &c.Registers().L) // updated order
	}},
	// INC HL
	0x23: {P: 1, OP: func(c cpu.CPU) {
		inc16(c, &c.Registers().H, &c.Registers().L)
	}},
	// INC H
	0x24: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().H)
	}},
	// DEC H
	0x25: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().H)
	}},
	// LD H,d8
	0x26: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.GetImmediate8())
	}},
	// DAA
	0x27: {P: 1, OP: func(c cpu.CPU) {
		decimalAdjust(c)
	}},
	// JR Z,r8
	0x28: {P: 0, OP: func(c cpu.CPU) {
		jumpRelative(c, int8(c.GetImmediate8()), c.Flags().Zero)
	}},
	// ADD HL,HL
	0x29: {P: 1, OP: func(c cpu.CPU) {
		add16(c, &c.Registers().H, &c.Registers().L, c.Registers().H, c.Registers().L)
	}},
	// LD A,(HL+)
	0x2A: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Registers().A = c.Read(addr)
		inc16(c, &c.Registers().H, &c.Registers().L)
	}},
	// DEC HL
	0x2B: {P: 1, OP: func(c cpu.CPU) {
		dec16(c, &c.Registers().H, &c.Registers().L)
	}},
	// INC L
	0x2C: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().L)
	}},
	// DEC L
	0x2D: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().L)
	}},
	// LD L,d8
	0x2E: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.GetImmediate8())
	}},
	// CPL
	0x2F: {P: 1, OP: func(c cpu.CPU) {
		c.Registers().A = ^c.Registers().A
		c.Flags().Set(flags.Leave, flags.Set, flags.Set, flags.Leave)
	}},
	// JR NC,r8
	0x30: {P: 0, OP: func(c cpu.CPU) {
		jumpRelative(c, int8(c.GetImmediate8()), !c.Flags().Carry)
	}},
	// LD SP,d16
	0x31: {P: 3, OP: func(c cpu.CPU) {
		c.Registers().SP = c.GetImmediate16()
	}},
	// LD (HL-),A
	0x32: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Write(addr, c.Registers().A)
		dec16(c, &c.Registers().H, &c.Registers().L) // updated order
	}},
	// INC SP
	0x33: {P: 1, OP: func(c cpu.CPU) {
		c.Registers().SP++
	}},
	// INC (HL)
	0x34: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		inc8(c, &val)
		c.Write(addr, val)
	}},
	// DEC (HL)
	0x35: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		dec8(c, &val)
		c.Write(addr, val)
	}},
	// LD (HL),d8
	0x36: {P: 2, OP: func(c cpu.CPU) {
		val := c.GetImmediate8()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, val)
	}},
	// SCF
	0x37: {P: 1, OP: func(c cpu.CPU) {
		c.Flags().Set(flags.Leave, flags.Reset, flags.Reset, flags.Set)
	}},
	// JR C,r8
	0x38: {P: 0, OP: func(c cpu.CPU) {
		jumpRelative(c, int8(c.GetImmediate8()), c.Flags().Carry)
	}},
	// ADD HL,SP
	0x39: {P: 1, OP: func(c cpu.CPU) {
		high, low := cpu.FromRegisterPair(c.Registers().SP)
		add16(c, &c.Registers().H, &c.Registers().L, high, low)
	}},
	// LD A,(HL-)
	0x3A: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Registers().A = c.Read(addr)
		dec16(c, &c.Registers().H, &c.Registers().L)
	}},
	// DEC SP
	0x3B: {P: 1, OP: func(c cpu.CPU) {
		c.Registers().SP--
	}},
	// INC A
	0x3C: {P: 1, OP: func(c cpu.CPU) {
		inc8(c, &c.Registers().A)
	}},
	// DEC A
	0x3D: {P: 1, OP: func(c cpu.CPU) {
		dec8(c, &c.Registers().A)
	}},
	// LD A,d8
	0x3E: {P: 2, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// CCF
	0x3F: {P: 1, OP: func(c cpu.CPU) {
		carry := flags.Reset
		if !c.Flags().Carry {
			carry = flags.Set
		}
		c.Flags().Set(flags.Leave, flags.Reset, flags.Reset, carry)
	}},
	// LD B,B
	0x40: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().B)
	}},
	// LD B,C
	0x41: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().C)
	}},
	// LD B,D
	0x42: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().D)
	}},
	// LD B,E
	0x43: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().E)
	}},
	// LD B,H
	0x44: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().H)
	}},
	// LD B,L
	0x45: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().L)
	}},
	// LD B,(HL)
	0x46: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().B = c.Read(addr)
	}},
	// LD B,A
	0x47: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().B, c.Registers().A)
	}},
	// LD C,B
	0x48: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().B)
	}},
	// LD C,C
	0x49: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().C)
	}},
	// LD C,D
	0x4A: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().D)
	}},
	// LD C,E
	0x4B: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().E)
	}},
	// LD C,H
	0x4C: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().H)
	}},
	// LD C,L
	0x4D: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().L)
	}},
	// LD C,(HL)
	0x4E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().C = c.Read(addr)
	}},
	// LD C,A
	0x4F: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().C, c.Registers().A)
	}},
	// LD D,B
	0x50: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().B)
	}},
	// LD D,C
	0x51: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().C)
	}},
	// LD D,D
	0x52: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().D)
	}},
	// LD D,E
	0x53: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().E)
	}},
	// LD D,H
	0x54: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().H)
	}},
	// LD D,L
	0x55: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().L)
	}},
	// LD D,(HL)
	0x56: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().D = c.Read(addr)
	}},
	// LD D,A
	0x57: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().D, c.Registers().A)
	}},
	// LD E,B
	0x58: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().B)
	}},
	// LD E,C
	0x59: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().C)
	}},
	// LD E,D
	0x5A: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().D)
	}},
	// LD E,E
	0x5B: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().E)
	}},
	// LD E,H
	0x5C: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().H)
	}},
	// LD E,L
	0x5D: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().L)
	}},
	// LD E,(HL)
	0x5E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().E = c.Read(addr)
	}},
	// LD E,A
	0x5F: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().E, c.Registers().A)
	}},
	// LD H,B
	0x60: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().B)
	}},
	// LD H,C
	0x61: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().C)
	}},
	// LD H,D
	0x62: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().D)
	}},
	// LD H,E
	0x63: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().E)
	}},
	// LD H,H
	0x64: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().H)
	}},
	// LD H,L
	0x65: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().L)
	}},
	// LD H,(HL)
	0x66: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().H = c.Read(addr)
	}},
	// LD H,A
	0x67: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().H, c.Registers().A)
	}},
	// LD L,B
	0x68: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().B)
	}},
	// LD L,C
	0x69: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().C)
	}},
	// LD L,D
	0x6A: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().D)
	}},
	// LD L,E
	0x6B: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().E)
	}},
	// LD L,H
	0x6C: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().H)
	}},
	// LD L,L
	0x6D: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().L)
	}},
	// LD L,(HL)
	0x6E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().L = c.Read(addr)
	}},
	// LD L,A
	0x6F: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().L, c.Registers().A)
	}},
	// LD (HL),B
	0x70: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().B)
	}},
	// LD (HL),C
	0x71: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().C)
	}},
	// LD (HL),D
	0x72: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().D)
	}},
	// LD (HL),E
	0x73: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().E)
	}},
	// LD (HL),H
	0x74: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().H)
	}},
	// LD (HL),L
	0x75: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().L)
	}},
	// HALT
	0x76: {P: 1, OP: func(c cpu.CPU) {
		c.Halt()
	}},
	// LD (HL),A
	0x77: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().A)
	}},
	// LD A,B
	0x78: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().B)
	}},
	// LD A,C
	0x79: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().C)
	}},
	// LD A,D
	0x7A: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().D)
	}},
	// LD A,E
	0x7B: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().E)
	}},
	// LD A,H
	0x7C: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().H)
	}},
	// LD A,L
	0x7D: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().L)
	}},
	// LD A,(HL)
	0x7E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Registers().A = c.Read(addr)
	}},
	// LD A,A
	0x7F: {P: 1, OP: func(c cpu.CPU) {
		load8(c, &c.Registers().A, c.Registers().A)
	}},
	// ADD A,B
	0x80: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().B)
	}},
	// ADD A,C
	0x81: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().C)
	}},
	// ADD A,D
	0x82: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().D)
	}},
	// ADD A,E
	0x83: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().E)
	}},
	// ADD A,H
	0x84: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().H)
	}},
	// ADD A,L
	0x85: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().L)
	}},
	// ADD A,(HL)
	0x86: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		add8(c, &c.Registers().A, val)
	}},
	// ADD A,A
	0x87: {P: 1, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.Registers().A)
	}},
	// ADC A,B
	0x88: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().B)
	}},
	// ADC A,C
	0x89: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().C)
	}},
	// ADC A,D
	0x8A: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().D)
	}},
	// ADC A,E
	0x8B: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().E)
	}},
	// ADC A,H
	0x8C: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().H)
	}},
	// ADC A,L
	0x8D: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().L)
	}},
	// ADC A,(HL)
	0x8E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		addc8(c, &c.Registers().A, val)
	}},
	// ADC A,A
	0x8F: {P: 1, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.Registers().A)
	}},
	// SUB B
	0x90: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().B)
	}},
	// SUB C
	0x91: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().C)
	}},
	// SUB D
	0x92: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().D)
	}},
	// SUB E
	0x93: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().E)
	}},
	// SUB H
	0x94: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().H)
	}},
	// SUB L
	0x95: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().L)
	}},
	// SUB (HL)
	0x96: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		sub8(c, &c.Registers().A, val)
	}},
	// SUB A
	0x97: {P: 1, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.Registers().A)
	}},
	// SBC A,B
	0x98: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().B)
	}},
	// SBC A,C
	0x99: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().C)
	}},
	// SBC A,D
	0x9A: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().D)
	}},
	// SBC A,E
	0x9B: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().E)
	}},
	// SBC A,H
	0x9C: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().H)
	}},
	// SBC A,L
	0x9D: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().L)
	}},
	// SBC A,(HL)
	0x9E: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		subc8(c, &c.Registers().A, val)
	}},
	// SBC A,A
	0x9F: {P: 1, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.Registers().A)
	}},
	// AND B
	0xA0: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().B)
	}},
	// AND C
	0xA1: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().C)
	}},
	// AND D
	0xA2: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().D)
	}},
	// AND E
	0xA3: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().E)
	}},
	// AND H
	0xA4: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().H)
	}},
	// AND L
	0xA5: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().L)
	}},
	// AND (HL)
	0xA6: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		and8(c, &c.Registers().A, val)
	}},
	// AND A
	0xA7: {P: 1, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.Registers().A)
	}},
	// XOR B
	0xA8: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().B)
	}},
	// XOR C
	0xA9: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().C)
	}},
	// XOR D
	0xAA: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().D)
	}},
	// XOR E
	0xAB: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().E)
	}},
	// XOR H
	0xAC: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().H)
	}},
	// XOR L
	0xAD: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().L)
	}},
	// XOR (HL)
	0xAE: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		xor8(c, &c.Registers().A, val)
	}},
	// XOR A
	0xAF: {P: 1, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.Registers().A)
	}},
	// OR B
	0xB0: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().B)
	}},
	// OR C
	0xB1: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().C)
	}},
	// OR D
	0xB2: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().D)
	}},
	// OR E
	0xB3: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().E)
	}},
	// OR H
	0xB4: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().H)
	}},
	// OR L
	0xB5: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().L)
	}},
	// OR (HL)
	0xB6: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		or8(c, &c.Registers().A, val)
	}},
	// OR A
	0xB7: {P: 1, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.Registers().A)
	}},
	// CP B
	0xB8: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().B)
	}},
	// CP C
	0xB9: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().C)
	}},
	// CP D
	0xBA: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().D)
	}},
	// CP E
	0xBB: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().E)
	}},
	// CP H
	0xBC: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().H)
	}},
	// CP L
	0xBD: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().L)
	}},
	// CP (HL)
	0xBE: {P: 1, OP: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		cp8(c, c.Registers().A, val)
		a := 1
		_ = a
	}},
	// CP A
	0xBF: {P: 1, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.Registers().A)
	}},
	// RET NZ
	0xC0: {P: 0, OP: func(c cpu.CPU) {
		ret(c, !c.Flags().Zero)
	}},
	// POP BC
	0xC1: {P: 1, OP: func(c cpu.CPU) {
		pop16(c, &c.Registers().B, &c.Registers().C)
	}},
	// JP NZ,a16
	0xC2: {P: 0, OP: func(c cpu.CPU) {
		jump(c, c.GetImmediate16(), !c.Flags().Zero)
	}},
	// JP a16
	0xC3: {P: 0, OP: func(c cpu.CPU) {
		jump(c, c.GetImmediate16(), true)
	}},
	// CALL NZ,a16
	0xC4: {P: 0, OP: func(c cpu.CPU) {
		call(c, c.GetImmediate16(), !c.Flags().Zero)
	}},
	// PUSH BC
	0xC5: {P: 1, OP: func(c cpu.CPU) {
		push16(c, c.Registers().B, c.Registers().C)
	}},
	// ADD A,d8
	0xC6: {P: 2, OP: func(c cpu.CPU) {
		add8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 00H
	0xC7: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x00)
	}},
	// RET Z
	0xC8: {P: 0, OP: func(c cpu.CPU) {
		ret(c, c.Flags().Zero)
	}},
	// RET
	0xC9: {P: 0, OP: func(c cpu.CPU) {
		ret(c, true)
	}},
	// JP Z,a16
	0xCA: {P: 0, OP: func(c cpu.CPU) {
		jump(c, c.GetImmediate16(), c.Flags().Zero)
	}},
	// PREFIX CB
	0xCB: {P: 1, OP: func(c cpu.CPU) {
		c.PrefixCB()
	}},
	// CALL Z,a16
	0xCC: {P: 0, OP: func(c cpu.CPU) {
		call(c, c.GetImmediate16(), c.Flags().Zero)
	}},
	// CALL a16
	0xCD: {P: 0, OP: func(c cpu.CPU) {
		call(c, c.GetImmediate16(), true)
	}},
	// ADC A,d8
	0xCE: {P: 2, OP: func(c cpu.CPU) {
		addc8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 08H
	0xCF: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x08)
	}},
	// RET NC
	0xD0: {P: 0, OP: func(c cpu.CPU) {
		ret(c, !c.Flags().Carry)
	}},
	// POP DE
	0xD1: {P: 1, OP: func(c cpu.CPU) {
		pop16(c, &c.Registers().D, &c.Registers().E)
	}},
	// JP NC,a16
	0xD2: {P: 0, OP: func(c cpu.CPU) {
		jump(c, c.GetImmediate16(), !c.Flags().Carry)
	}},
	// INVALID
	// CALL NC,a16
	0xD4: {P: 0, OP: func(c cpu.CPU) {
		call(c, c.GetImmediate16(), !c.Flags().Carry)
	}},
	// PUSH DE
	0xD5: {P: 1, OP: func(c cpu.CPU) {
		push16(c, c.Registers().D, c.Registers().E)
	}},
	// SUB d8
	0xD6: {P: 2, OP: func(c cpu.CPU) {
		sub8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 10H
	0xD7: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x10)
	}},
	// RET C
	0xD8: {P: 0, OP: func(c cpu.CPU) {
		ret(c, c.Flags().Carry)
	}},
	// RETI
	0xD9: {P: 0, OP: func(c cpu.CPU) {
		ret(c, true)
		c.EI()
	}},
	// JP C,a16
	0xDA: {P: 0, OP: func(c cpu.CPU) {
		jump(c, c.GetImmediate16(), c.Flags().Carry)
	}},
	// INVALID
	// CALL C,a16
	0xDC: {P: 0, OP: func(c cpu.CPU) {
		call(c, c.GetImmediate16(), c.Flags().Carry)
	}},
	// INVALID
	// SBC A,d8
	0xDE: {P: 2, OP: func(c cpu.CPU) {
		subc8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 18H
	0xDF: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x18)
	}},
	// LDH (a8),A
	0xE0: {P: 2, OP: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.GetImmediate8())
		load8Mem(c, c.Registers().A, addr)
	}},
	// POP HL
	0xE1: {P: 1, OP: func(c cpu.CPU) {
		pop16(c, &c.Registers().H, &c.Registers().L)
	}},
	// LD (C),A
	0xE2: {P: 1, OP: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.Registers().C)
		load8Mem(c, c.Registers().A, addr)
	}},
	// INVALID
	// INVALID
	// PUSH HL
	0xE5: {P: 1, OP: func(c cpu.CPU) {
		push16(c, c.Registers().H, c.Registers().L)
	}},
	// AND d8
	0xE6: {P: 2, OP: func(c cpu.CPU) {
		and8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 20H
	0xE7: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x20)
	}},
	// ADD SP,r8
	0xE8: {P: 2, OP: func(c cpu.CPU) {
		addSPr8(c)
	}},
	// JP (HL)
	0xE9: {P: 0, OP: func(c cpu.CPU) {
		c.Registers().PC = cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
	}},
	// LD (a16),A
	0xEA: {P: 3, OP: func(c cpu.CPU) {
		addr := c.GetImmediate16()
		c.Write(addr, c.Registers().A)
	}},
	// INVALID
	// INVALID
	// INVALID
	// XOR d8
	0xEE: {P: 2, OP: func(c cpu.CPU) {
		xor8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 28H
	0xEF: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x28)
	}},
	// LDH A,(a8)
	0xF0: {P: 2, OP: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.GetImmediate8())
		c.Registers().A = c.Read(addr)
	}},
	// POP AF
	0xF1: {P: 1, OP: func(c cpu.CPU) {
		popAF(c)
	}},
	// LD A,(C)
	0xF2: {P: 1, OP: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.Registers().C)
		c.Registers().A = c.Read(addr)
	}},
	// DI
	0xF3: {P: 1, OP: func(c cpu.CPU) {
		// c.ime = false
	}},
	// INVALID
	// PUSH AF
	0xF5: {P: 1, OP: func(c cpu.CPU) {
		reg := c.Flags().Read()
		push16(c, c.Registers().A, reg)
	}},
	// OR d8
	0xF6: {P: 2, OP: func(c cpu.CPU) {
		or8(c, &c.Registers().A, c.GetImmediate8())
	}},
	// RST 30H
	0xF7: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x30)
	}},
	// LD HL,SP+r8
	0xF8: {P: 2, OP: func(c cpu.CPU) {
		immediate := int8(c.GetImmediate8())
		loadHLSPOffset(c, immediate)
	}},
	// LD SP,HL
	0xF9: {P: 1, OP: func(c cpu.CPU) {
		c.Registers().SP = cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
	}},
	// LD A,(a16)
	0xFA: {P: 3, OP: func(c cpu.CPU) {
		addr := c.GetImmediate16()
		c.Registers().A = c.Read(addr)
	}},
	// EI
	0xFB: {P: 1, OP: func(c cpu.CPU) {
		c.EIWithDelay()
	}},
	// INVALID
	// INVALID
	// CP d8
	0xFE: {P: 2, OP: func(c cpu.CPU) {
		cp8(c, c.Registers().A, c.GetImmediate8())
	}},
	// RST 38H
	0xFF: {P: 0, OP: func(c cpu.CPU) {
		rst(c, 0x38)
	}},
}
