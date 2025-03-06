package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

type Instruction func(c cpu.CPU) // Operation

var Instructions = [0x100]Instruction{
	// NOP
	0x00: func(c cpu.CPU) {},
	// LD BC,d16
	0x01: func(c cpu.CPU) { load16(c, &c.Registers().B, &c.Registers().C, c.GetImmediate16()) },
	// LD (BC),A
	0x02: func(c cpu.CPU) {
		c.Clock()
		address := cpu.ToRegisterPair(c.Registers().B, c.Registers().C)
		c.Ack()
		c.Write(address, c.Registers().A)
	},
	// INC BC
	0x03: func(c cpu.CPU) { inc16(c, &c.Registers().B, &c.Registers().C) },
	// INC B
	0x04: func(c cpu.CPU) { inc8(c, &c.Registers().B) },
	// DEC B
	0x05: func(c cpu.CPU) { dec8(c, &c.Registers().B) },
	// LD B,d8
	0x06: func(c cpu.CPU) { load8(c, &c.Registers().B, c.GetImmediate8()) },
	// RLCA: rotate A left circularly (ignore previous carry)
	0x07: func(c cpu.CPU) { rotate(c, &c.Registers().A, true, false, false) },
	// LD (a16),SP
	0x08: func(c cpu.CPU) {
		address := c.GetImmediate16()
		c.Clock()
		bytes := uint(c.Registers().SP)
		c.Ack()

		c.Clock()
		c.Write(address, uint8(bytes&0xFF))
		c.Ack()

		c.Write(address+1, uint8(bytes>>8))
	},
	// ADD HL,BC
	0x09: func(c cpu.CPU) { add16(c, &c.Registers().H, &c.Registers().L, c.Registers().B, c.Registers().C) },
	// LD A,(BC)
	0x0A: func(c cpu.CPU) {
		c.Clock()
		address := cpu.ToRegisterPair(c.Registers().B, c.Registers().C)
		c.Ack()
		c.Registers().A = c.Read(address)
	},
	// DEC BC
	0x0B: func(c cpu.CPU) { dec16(c, &c.Registers().B, &c.Registers().C) },
	// INC C
	0x0C: func(c cpu.CPU) { inc8(c, &c.Registers().C) },
	// DEC C
	0x0D: func(c cpu.CPU) { dec8(c, &c.Registers().C) },
	// LD C,d8
	0x0E: func(c cpu.CPU) { load8(c, &c.Registers().C, c.GetImmediate8()) },
	// RRCA: rotate A right circularly (ignore previous carry)
	0x0F: func(c cpu.CPU) { rotate(c, &c.Registers().A, false, false, false) },
	// STOP 0
	0x10: func(c cpu.CPU) {},
	// LD DE,d16
	0x11: func(c cpu.CPU) { load16(c, &c.Registers().D, &c.Registers().E, c.GetImmediate16()) },
	// LD (DE),A
	0x12: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().D, c.Registers().E)
		c.Ack()
		c.Write(addr, c.Registers().A)
	},
	// INC DE
	0x13: func(c cpu.CPU) { inc16(c, &c.Registers().D, &c.Registers().E) },
	// INC D
	0x14: func(c cpu.CPU) { inc8(c, &c.Registers().D) },
	// DEC D
	0x15: func(c cpu.CPU) { dec8(c, &c.Registers().D) },
	// LD D,d8
	0x16: func(c cpu.CPU) { load8(c, &c.Registers().D, c.GetImmediate8()) },
	// RLA: rotate A left through carry (use previous carry)
	0x17: func(c cpu.CPU) { rotate(c, &c.Registers().A, true, true, false) },
	// JR r8
	0x18: func(c cpu.CPU) { jumpRelative(c, int8(c.GetImmediate8()), true) },
	// ADD HL,DE
	0x19: func(c cpu.CPU) { add16(c, &c.Registers().H, &c.Registers().L, c.Registers().D, c.Registers().E) },
	// LD A,(DE)
	0x1A: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().D, c.Registers().E)
		c.Ack()
		c.Registers().A = c.Read(addr)
	},
	// DEC DE
	0x1B: func(c cpu.CPU) { dec16(c, &c.Registers().D, &c.Registers().E) },
	// INC E
	0x1C: func(c cpu.CPU) { inc8(c, &c.Registers().E) },
	// DEC E
	0x1D: func(c cpu.CPU) { dec8(c, &c.Registers().E) },
	// LD E,d8
	0x1E: func(c cpu.CPU) { load8(c, &c.Registers().E, c.GetImmediate8()) },
	// RRA: rotate A right through carry (use previous carry)
	0x1F: func(c cpu.CPU) { rotate(c, &c.Registers().A, false, true, false) },

	// JR NZ,r8
	0x20: func(c cpu.CPU) { jumpRelative(c, int8(c.GetImmediate8()), !c.Flags().Zero) },
	// LD HL,d16
	0x21: func(c cpu.CPU) { load16(c, &c.Registers().H, &c.Registers().L, c.GetImmediate16()) },
	// LD (HL+),A
	0x22: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().A)
		inc16(c, &c.Registers().H, &c.Registers().L)
	},
	// INC HL
	0x23: func(c cpu.CPU) { inc16(c, &c.Registers().H, &c.Registers().L) },
	// INC H
	0x24: func(c cpu.CPU) { inc8(c, &c.Registers().H) },
	// DEC H
	0x25: func(c cpu.CPU) { dec8(c, &c.Registers().H) },
	// LD H,d8
	0x26: func(c cpu.CPU) { load8(c, &c.Registers().H, c.GetImmediate8()) },
	// DAA
	0x27: func(c cpu.CPU) { decimalAdjust(c) },
	// JR Z,r8
	0x28: func(c cpu.CPU) { jumpRelative(c, int8(c.GetImmediate8()), c.Flags().Zero) },
	// ADD HL,HL
	0x29: func(c cpu.CPU) { add16(c, &c.Registers().H, &c.Registers().L, c.Registers().H, c.Registers().L) },
	// LD A,(HL+)
	0x2A: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Registers().A = c.Read(addr)
		inc16(c, &c.Registers().H, &c.Registers().L)
	},
	// DEC HL
	0x2B: func(c cpu.CPU) { dec16(c, &c.Registers().H, &c.Registers().L) },
	// INC L
	0x2C: func(c cpu.CPU) { inc8(c, &c.Registers().L) },
	// DEC L
	0x2D: func(c cpu.CPU) { dec8(c, &c.Registers().L) },
	// LD L,d8
	0x2E: func(c cpu.CPU) { load8(c, &c.Registers().L, c.GetImmediate8()) },
	// CPL
	0x2F: func(c cpu.CPU) {
		c.Registers().A = ^c.Registers().A
		c.Flags().Set(flags.Leave, flags.Set, flags.Set, flags.Leave)
	},
	// JR NC,r8
	0x30: func(c cpu.CPU) { jumpRelative(c, int8(c.GetImmediate8()), !c.Flags().Carry) },
	// LD SP,d16
	0x31: func(c cpu.CPU) { c.Registers().SP = c.GetImmediate16() },
	// LD (HL-),A
	0x32: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Write(addr, c.Registers().A)
		dec16(c, &c.Registers().H, &c.Registers().L)
	},
	// INC SP
	0x33: func(c cpu.CPU) {
		c.Clock()
		c.Registers().SP++
		c.Ack()
	},
	// INC (HL)
	0x34: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)

		c.Clock()
		val := c.Read(addr)
		c.Ack()

		c.Clock()
		inc8(c, &val)
		c.Ack()

		c.Write(addr, val)
	},
	// DEC (HL)
	0x35: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Clock()
		val := c.Read(addr)
		c.Ack()
		c.Clock()
		dec8(c, &val)
		c.Ack()
		c.Write(addr, val)
	},
	// LD (HL),d8
	0x36: func(c cpu.CPU) {
		val := c.GetImmediate8()
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, val)
	},
	// SCF
	0x37: func(c cpu.CPU) { c.Flags().Set(flags.Leave, flags.Reset, flags.Reset, flags.Set) },
	// JR C,r8
	0x38: func(c cpu.CPU) { jumpRelative(c, int8(c.GetImmediate8()), c.Flags().Carry) },
	// ADD HL,SP
	0x39: func(c cpu.CPU) {
		high, low := cpu.FromRegisterPair(c.Registers().SP)
		add16(c, &c.Registers().H, &c.Registers().L, high, low)
	},
	// LD A,(HL-)
	0x3A: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L) // updated order
		c.Registers().A = c.Read(addr)
		dec16(c, &c.Registers().H, &c.Registers().L)
	},
	// DEC SP
	0x3B: func(c cpu.CPU) {
		c.Clock()
		c.Registers().SP--
		c.Ack()
	},
	// INC A
	0x3C: func(c cpu.CPU) { inc8(c, &c.Registers().A) },
	// DEC A
	0x3D: func(c cpu.CPU) { dec8(c, &c.Registers().A) },
	// LD A,d8
	0x3E: func(c cpu.CPU) { load8(c, &c.Registers().A, c.GetImmediate8()) },
	// CCF
	0x3F: func(c cpu.CPU) {
		carry := flags.Reset
		if !c.Flags().Carry {
			carry = flags.Set
		}
		c.Flags().Set(flags.Leave, flags.Reset, flags.Reset, carry)
	},
	// LD B,B
	0x40: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().B) },
	// LD B,C
	0x41: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().C) },
	// LD B,D
	0x42: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().D) },
	// LD B,E
	0x43: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().E) },
	// LD B,H
	0x44: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().H) },
	// LD B,L
	0x45: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().L) },
	// LD B,(HL)
	0x46: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().B = c.Read(addr)
	},
	// LD B,A
	0x47: func(c cpu.CPU) { load8(c, &c.Registers().B, c.Registers().A) },
	// LD C,B
	0x48: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().B) },
	// LD C,C
	0x49: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().C) },
	// LD C,D
	0x4A: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().D) },
	// LD C,E
	0x4B: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().E) },
	// LD C,H
	0x4C: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().H) },
	// LD C,L
	0x4D: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().L) },
	// LD C,(HL)
	0x4E: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().C = c.Read(addr)
	},
	// LD C,A
	0x4F: func(c cpu.CPU) { load8(c, &c.Registers().C, c.Registers().A) },
	// LD D,B
	0x50: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().B) },
	// LD D,C
	0x51: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().C) },
	// LD D,D
	0x52: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().D) },
	// LD D,E
	0x53: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().E) },
	// LD D,H
	0x54: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().H) },
	// LD D,L
	0x55: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().L) },
	// LD D,(HL)
	0x56: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().D = c.Read(addr)
	},
	// LD D,A
	0x57: func(c cpu.CPU) { load8(c, &c.Registers().D, c.Registers().A) },
	// LD E,B
	0x58: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().B) },
	// LD E,C
	0x59: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().C) },
	// LD E,D
	0x5A: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().D) },
	// LD E,E
	0x5B: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().E) },
	// LD E,H
	0x5C: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().H) },
	// LD E,L
	0x5D: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().L) },
	// LD E,(HL)
	0x5E: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().E = c.Read(addr)
	},
	// LD E,A
	0x5F: func(c cpu.CPU) { load8(c, &c.Registers().E, c.Registers().A) },
	// LD H,B
	0x60: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().B) },
	// LD H,C
	0x61: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().C) },
	// LD H,D
	0x62: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().D) },
	// LD H,E
	0x63: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().E) },
	// LD H,H
	0x64: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().H) },
	// LD H,L
	0x65: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().L) },
	// LD H,(HL)
	0x66: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().H = c.Read(addr)
	},
	// LD H,A
	0x67: func(c cpu.CPU) { load8(c, &c.Registers().H, c.Registers().A) },
	// LD L,B
	0x68: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().B) },
	// LD L,C
	0x69: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().C) },
	// LD L,D
	0x6A: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().D) },
	// LD L,E
	0x6B: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().E) },
	// LD L,H
	0x6C: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().H) },
	// LD L,L
	0x6D: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().L) },
	// LD L,(HL)
	0x6E: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().L = c.Read(addr)
	},
	// LD L,A
	0x6F: func(c cpu.CPU) { load8(c, &c.Registers().L, c.Registers().A) },
	// LD (HL),B
	0x70: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().B)
	},
	// LD (HL),C
	0x71: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().C)
	},
	// LD (HL),D
	0x72: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().D)
	},
	// LD (HL),E
	0x73: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().E)
	},
	// LD (HL),H
	0x74: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().H)
	},
	// LD (HL),L
	0x75: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().L)
	},
	// HALT
	0x76: func(c cpu.CPU) { c.Halt() },
	// LD (HL),A
	0x77: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Write(addr, c.Registers().A)
	},
	// LD A,B
	0x78: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().B) },
	// LD A,C
	0x79: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().C) },
	// LD A,D
	0x7A: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().D) },
	// LD A,E
	0x7B: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().E) },
	// LD A,H
	0x7C: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().H) },
	// LD A,L
	0x7D: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().L) },
	// LD A,(HL)
	0x7E: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		c.Registers().A = c.Read(addr)
	},
	// LD A,A
	0x7F: func(c cpu.CPU) { load8(c, &c.Registers().A, c.Registers().A) },
	// ADD A,B
	0x80: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().B) },
	// ADD A,C
	0x81: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().C) },
	// ADD A,D
	0x82: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().D) },
	// ADD A,E
	0x83: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().E) },
	// ADD A,H
	0x84: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().H) },
	// ADD A,L
	0x85: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().L) },
	// ADD A,(HL)
	0x86: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Clock()
		val := c.Read(addr)
		c.Ack()
		add8(c, &c.Registers().A, val)
	},
	// ADD A
	0x87: func(c cpu.CPU) { add8(c, &c.Registers().A, c.Registers().A) },
	// ADC B
	0x88: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().B) },
	// ADC C
	0x89: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().C) },
	// ADC D
	0x8A: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().D) },
	// ADC E
	0x8B: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().E) },
	// ADC H
	0x8C: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().H) },
	// ADC L
	0x8D: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().L) },
	// ADC (HL)
	0x8E: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Clock()
		val := c.Read(addr)
		c.Ack()
		addc8(c, &c.Registers().A, val)
	},
	// ADC A
	0x8F: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.Registers().A) },
	// SUB B
	0x90: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().B) },
	// SUB C
	0x91: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().C) },
	// SUB D
	0x92: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().D) },
	// SUB E
	0x93: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().E) },
	// SUB H
	0x94: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().H) },
	// SUB L
	0x95: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().L) },
	// SUB (HL)
	0x96: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		val := c.Read(addr)
		sub8(c, &c.Registers().A, val)
	},
	// SUB A
	0x97: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.Registers().A) },
	// SBC B
	0x98: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().B) },
	// SBC C
	0x99: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().C) },
	// SBC D
	0x9A: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().D) },
	// SBC E
	0x9B: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().E) },
	// SBC H
	0x9C: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().H) },
	// SBC L
	0x9D: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().L) },
	// SBC (HL)
	0x9E: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		c.Ack()
		subc8(c, &c.Registers().A, val)
	},
	// SBC A
	0x9F: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.Registers().A) },
	// AND B
	0xA0: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().B) },
	// AND C
	0xA1: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().C) },
	// AND D
	0xA2: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().D) },
	// AND E
	0xA3: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().E) },
	// AND H
	0xA4: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().H) },
	// AND L
	0xA5: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().L) },
	// AND (HL)
	0xA6: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		c.Ack()

		and8(c, &c.Registers().A, val)
	},
	// AND A
	0xA7: func(c cpu.CPU) { and8(c, &c.Registers().A, c.Registers().A) },
	// XOR B
	0xA8: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().B) },
	// XOR C
	0xA9: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().C) },
	// XOR D
	0xAA: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().D) },
	// XOR E
	0xAB: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().E) },
	// XOR H
	0xAC: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().H) },
	// XOR L
	0xAD: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().L) },
	// XOR (HL)
	0xAE: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		val := c.Read(addr)
		c.Ack()
		xor8(c, &c.Registers().A, val)
	},
	// XOR A
	0xAF: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.Registers().A) },
	// OR B
	0xB0: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().B) },
	// OR C
	0xB1: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().C) },
	// OR D
	0xB2: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().D) },
	// OR E
	0xB3: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().E) },
	// OR H
	0xB4: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().H) },
	// OR L
	0xB5: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().L) },
	// OR (HL)
	0xB6: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()
		val := c.Read(addr)
		or8(c, &c.Registers().A, val)
	},
	// OR A
	0xB7: func(c cpu.CPU) { or8(c, &c.Registers().A, c.Registers().A) },
	// CP B
	0xB8: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().B) },
	// CP C
	0xB9: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().C) },
	// CP D
	0xBA: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().D) },
	// CP E
	0xBB: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().E) },
	// CP H
	0xBC: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().H) },
	// CP L
	0xBD: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().L) },
	// CP (HL)
	0xBE: func(c cpu.CPU) {
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Clock()
		val := c.Read(addr)
		c.Ack()
		cp8(c, c.Registers().A, val)
		a := 1
		_ = a
	},
	// CP A
	0xBF: func(c cpu.CPU) { cp8(c, c.Registers().A, c.Registers().A) },
	// RET NZ
	0xC0: func(c cpu.CPU) { ret(c, !c.Flags().Zero) },
	// POP BC
	0xC1: func(c cpu.CPU) { pop16(c, &c.Registers().B, &c.Registers().C) },
	// JP NZ,a16
	0xC2: func(c cpu.CPU) { jump(c, c.GetImmediate16(), !c.Flags().Zero) },
	// JP a16
	0xC3: func(c cpu.CPU) { jump(c, c.GetImmediate16(), true) },
	// CALL NZ,a16
	0xC4: func(c cpu.CPU) { call(c, c.GetImmediate16(), !c.Flags().Zero) },
	// PUSH BC
	0xC5: func(c cpu.CPU) { push16(c, c.Registers().B, c.Registers().C) },
	// ADD A,d8
	0xC6: func(c cpu.CPU) { add8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 00H
	0xC7: func(c cpu.CPU) { call(c, 0x00, true) },
	// RET Z
	0xC8: func(c cpu.CPU) { ret(c, c.Flags().Zero) },
	// RET
	0xC9: func(c cpu.CPU) { ret(c, true) },
	// JP Z,a16
	0xCA: func(c cpu.CPU) { jump(c, c.GetImmediate16(), c.Flags().Zero) },
	// PREFIX CB
	0xCB: func(c cpu.CPU) { c.PrefixCB() },
	// CALL Z,a16
	0xCC: func(c cpu.CPU) { call(c, c.GetImmediate16(), c.Flags().Zero) },
	// CALL a16
	0xCD: func(c cpu.CPU) { call(c, c.GetImmediate16(), true) },
	// ADC A,d8
	0xCE: func(c cpu.CPU) { addc8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 08H
	0xCF: func(c cpu.CPU) { call(c, 0x08, true) },
	// RET NC
	0xD0: func(c cpu.CPU) { ret(c, !c.Flags().Carry) },
	// POP DE
	0xD1: func(c cpu.CPU) { pop16(c, &c.Registers().D, &c.Registers().E) },
	// JP NC,a16
	0xD2: func(c cpu.CPU) { jump(c, c.GetImmediate16(), !c.Flags().Carry) },
	// INVALID
	// CALL NC,a16
	0xD4: func(c cpu.CPU) { call(c, c.GetImmediate16(), !c.Flags().Carry) },
	// PUSH DE
	0xD5: func(c cpu.CPU) { push16(c, c.Registers().D, c.Registers().E) },
	// SUB d8
	0xD6: func(c cpu.CPU) { sub8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 10H
	0xD7: func(c cpu.CPU) { call(c, 0x10, true) },
	// RET C
	0xD8: func(c cpu.CPU) { ret(c, c.Flags().Carry) },
	// RETI
	0xD9: func(c cpu.CPU) {
		ret(c, true)
		c.EI()
	},
	// JP C,a16
	0xDA: func(c cpu.CPU) { jump(c, c.GetImmediate16(), c.Flags().Carry) },
	// INVALID
	// CALL C,a16
	0xDC: func(c cpu.CPU) { call(c, c.GetImmediate16(), c.Flags().Carry) },
	// INVALID
	// SBC d8
	0xDE: func(c cpu.CPU) { subc8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 18H
	0xDF: func(c cpu.CPU) { call(c, 0x18, true) },
	// LDH (a8),A
	0xE0: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.GetImmediate8())
		load8Mem(c, c.Registers().A, addr)
	},
	// POP HL
	0xE1: func(c cpu.CPU) { pop16(c, &c.Registers().H, &c.Registers().L) },
	// LD (C),A
	0xE2: func(c cpu.CPU) {
		addr := 0xFF00 + uint16(c.Registers().C)
		load8Mem(c, c.Registers().A, addr)
	},
	// INVALID
	// INVALID
	// PUSH HL
	0xE5: func(c cpu.CPU) { push16(c, c.Registers().H, c.Registers().L) },
	// AND d8
	0xE6: func(c cpu.CPU) { and8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 20H
	0xE7: func(c cpu.CPU) { call(c, 0x20, true) },
	// ADD SP,r8
	0xE8: func(c cpu.CPU) { addSPr8(c) },
	// JP (HL)
	0xE9: func(c cpu.CPU) { c.Registers().PC = cpu.ToRegisterPair(c.Registers().H, c.Registers().L) - 1 }, // <insert cussing here>
	// LD (a16),A
	0xEA: func(c cpu.CPU) {
		addr := c.GetImmediate16()
		c.Clock()
		c.Write(addr, c.Registers().A)
		c.Ack()
	},
	// INVALID
	// INVALID
	// INVALID
	// XOR d8
	0xEE: func(c cpu.CPU) { xor8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 28H
	0xEF: func(c cpu.CPU) { call(c, 0x28, true) },
	// LDH A,(a8)
	0xF0: func(c cpu.CPU) {
		addr := uint16(c.GetImmediate8())

		c.Clock()
		addr += 0xFF00
		c.Ack()

		c.Registers().A = c.Read(addr)
	},
	// POP AF
	0xF1: func(c cpu.CPU) { popAF(c) },
	// LD A,(C)
	0xF2: func(c cpu.CPU) {
		c.Clock()
		addr := 0xFF00 + uint16(c.Registers().C)
		c.Ack()
		c.Registers().A = c.Read(addr)
	},
	// DI
	0xF3: func(c cpu.CPU) {
		c.DI()
	},
	// INVALID
	// PUSH AF
	0xF5: func(c cpu.CPU) { push16(c, c.Registers().A, c.Flags().Read()) },
	// OR d8
	0xF6: func(c cpu.CPU) { or8(c, &c.Registers().A, c.GetImmediate8()) },
	// RST 30H
	0xF7: func(c cpu.CPU) { call(c, 0x30, true) },
	// LD HL,SP+r8
	0xF8: func(c cpu.CPU) {
		immediate := int8(c.GetImmediate8())
		loadHLSPOffset(c, immediate)
	},
	// LD SP,HL
	0xF9: func(c cpu.CPU) {
		c.Clock()
		addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
		c.Ack()

		c.Registers().SP = addr
	},
	// LD A,(a16)
	0xFA: func(c cpu.CPU) {
		addr := c.GetImmediate16()
		c.Clock()
		c.Registers().A = c.Read(addr)
		c.Ack()
	},
	// EI
	0xFB: func(c cpu.CPU) { c.EIWithDelay() },
	// INVALID
	// INVALID
	// CP d8
	0xFE: func(c cpu.CPU) { cp8(c, c.Registers().A, c.GetImmediate8()) },
	// RST 38H
	0xFF: func(c cpu.CPU) { call(c, 0x38, true) },
}
