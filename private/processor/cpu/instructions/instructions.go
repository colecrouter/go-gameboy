package instructions

import (
	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/enums"
	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/generators"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
)

var Instructions = [0x100]shared.Instruction{
	// NOP
	0x00: Nop(),
	// LD BC,d16
	0x01: Load16(BC, D16),
	// LD (BC),A
	0x02: Load(BC_, A),
	// INC BC
	0x03: Increment16(BC),
	// INC B
	0x04: Increment(B),
	// DEC B
	0x05: Decrement(B),
	// LD B,d8
	0x06: Load(B, D8),
	// RLCA: h.Rotate A left circularly (ignore previous carry)
	0x07: Rotate(A, true, false, false),
	// LD (a16),SP
	0x08: Load16(D16, SP), // ???
	// ADD HL,BC
	0x09: Add16(HL, BC),
	// LD A,(BC)
	0x0A: Load(A, BC_),
	// DEC BC
	0x0B: Decrement16(BC),
	// INC C
	0x0C: Increment(C),
	// DEC C
	0x0D: Decrement(C),
	// LD C,d8
	0x0E: Load(C, D8),
	// RRCA: h.Rotate A right circularly (ignore previous carry)
	0x0F: Rotate(A, false, false, false),
	// STOP 0
	0x10: Stop(),
	// LD DE,d16
	0x11: Load16(DE, D16),
	// LD (DE),A
	0x12: Load(DE_, A),
	// INC DE
	0x13: Increment16(DE),
	// INC D
	0x14: Increment(D),
	// DEC D
	0x15: Decrement(D),
	// LD D,d8
	0x16: Load(D, D8),
	// RLA
	0x17: Rotate(A, true, true, false),
	// JR r8
	0x18: JumpRelative(Always),
	// ADD HL,DE
	0x19: Add16(HL, DE),
	// LD A,(DE)
	0x1A: Load(A, DE_),
	// DEC DE
	0x1B: Decrement16(DE),
	// INC E
	0x1C: Increment(E),
	// DEC E
	0x1D: Decrement(E),
	// LD E,d8
	0x1E: Load(E, D8),
	// RRA
	0x1F: Rotate(A, false, true, false),
	// JR NZ,r8
	0x20: JumpRelative(Zero),
	// LD HL,d16
	0x21: Load16(HL, D16),
	// LD (HL+),A
	0x22: LoadInc(HL_, A, Inc, None),
	// INC HL
	0x23: Increment16(HL),
	// INC H
	0x24: Increment(H),
	// DEC H
	0x25: Decrement(H),
	// LD H,d8
	0x26: Load(H, D8),
	// DAA
	0x27: DecimalAdjust(),
	// JR Z,r8
	0x28: JumpRelative(Zero),
	// ADD HL,HL
	0x29: Add16(HL, HL),
	// LD A,(HL+)
	0x2A: LoadInc(A, HL_, None, Inc),
	// DEC HL
	0x2B: Decrement16(HL),
	// INC L
	0x2C: Increment(L),
	// DEC L
	0x2D: Decrement(L),
	// LD L,d8
	0x2E: Load(L, D8),
	// CPL
	0x2F: ComplementAcc(),
	// JR NC,r8
	0x30: JumpRelative(Carry),
	// LD SP,d16
	0x31: Load16(SP, D16),
	// LD (HL-),A
	0x32: LoadInc(HL_, A, Dec, None),
	// INC SP
	0x33: Increment16(SP),
	// INC (HL)
	0x34: Increment(HL_),
	// DEC (HL)
	0x35: Decrement(HL_),
	// LD (HL),d8
	0x36: Load(HL_, D8),
	// SCF
	0x37: SetCarry(),
	// JR C,r8
	0x38: JumpRelative(Carry),
	// ADD HL,SP
	0x39: Add16(HL, SP),
	// LD A,(HL-)
	0x3A: LoadInc(A, HL_, None, Dec),
	// DEC SP
	0x3B: Decrement16(SP),
	// INC A
	0x3C: Increment(A),
	// DEC A
	0x3D: Decrement(A),
	// LD A,d8
	0x3E: Load(A, D8),
	// CCF
	0x3F: ComplementCarry(),
	// LD B,B
	0x40: Load(B, B),
	// LD B,C
	0x41: Load(B, C),
	// LD B,D
	0x42: Load(B, D),
	// LD B,E
	0x43: Load(B, E),
	// LD B,H
	0x44: Load(B, H),
	// LD B,L
	0x45: Load(B, L),
	// LD B,(HL)
	0x46: Load(B, HL_),
	// LD B,A
	0x47: Load(B, A),
	// LD C,B
	0x48: Load(C, B),
	// LD C,C
	0x49: Load(C, C),
	// LD C,D
	0x4A: Load(C, D),
	// LD C,E
	0x4B: Load(C, E),
	// LD C,H
	0x4C: Load(C, H),
	// LD C,L
	0x4D: Load(C, L),
	// LD C,(HL)
	0x4E: Load(C, HL_),
	// LD C,A
	0x4F: Load(C, A),
	// LD D,B
	0x50: Load(D, B),
	// LD D,C
	0x51: Load(D, C),
	// LD D,D
	0x52: Load(D, D),
	// LD D,E
	0x53: Load(D, E),
	// LD D,H
	0x54: Load(D, H),
	// LD D,L
	0x55: Load(D, L),
	// LD D,(HL)
	0x56: Load(D, HL_),
	// LD D,A
	0x57: Load(D, A),
	// LD E,B
	0x58: Load(E, B),
	// LD E,C
	0x59: Load(E, C),
	// LD E,D
	0x5A: Load(E, D),
	// LD E,E
	0x5B: Load(E, E),
	// LD E,H
	0x5C: Load(E, H),
	// LD E,L
	0x5D: Load(E, L),
	// LD E,(HL)
	0x5E: Load(E, HL_),
	// LD E,A
	0x5F: Load(E, A),
	// LD H,B
	0x60: Load(H, B),
	// LD H,C
	0x61: Load(H, C),
	// LD H,D
	0x62: Load(H, D),
	// LD H,E
	0x63: Load(H, E),
	// LD H,H
	0x64: Load(H, H),
	// LD H,L
	0x65: Load(H, L),
	// LD H,(HL)
	0x66: Load(H, HL_),
	// LD H,A
	0x67: Load(H, A),
	// LD L,B
	0x68: Load(L, B),
	// LD L,C
	0x69: Load(L, C),
	// LD L,D
	0x6A: Load(L, D),
	// LD L,E
	0x6B: Load(L, E),
	// LD L,H
	0x6C: Load(L, H),
	// LD L,L
	0x6D: Load(L, L),
	// LD L,(HL)
	0x6E: Load(L, HL_),
	// LD L,A
	0x6F: Load(L, A),
	// LD (HL),B
	0x70: Load(HL_, B),
	// LD (HL),C
	0x71: Load(HL_, C),
	// LD (HL),D
	0x72: Load(HL_, D),
	// LD (HL),E
	0x73: Load(HL_, E),
	// LD (HL),H
	0x74: Load(HL_, H),
	// LD (HL),L
	0x75: Load(HL_, L),
	// HALT
	0x76: Halt(),
	// LD (HL),A
	0x77: Load(HL_, A),
	// LD A,B
	0x78: Load(A, B),
	// LD A,C
	0x79: Load(A, C),
	// LD A,D
	0x7A: Load(A, D),
	// LD A,E
	0x7B: Load(A, E),
	// LD A,H
	0x7C: Load(A, H),
	// LD A,L
	0x7D: Load(A, L),
	// LD A,(HL)
	0x7E: Load(A, HL_),
	// LD A,A
	0x7F: Load(A, A),
	// ADD A,B
	0x80: Add(A, B),
	// ADD A,C
	0x81: Add(A, C),
	// ADD A,D
	0x82: Add(A, D),
	// ADD A,E
	0x83: Add(A, E),
	// ADD A,H
	0x84: Add(A, H),
	// ADD A,L
	0x85: Add(A, L),
	// ADD A,(HL)
	0x86: Add(A, HL_),
	// ADD A
	0x87: Add(A, A),
	// ADC B
	0x88: AddC(A, B),
	// ADC C
	0x89: AddC(A, C),
	// ADC D
	0x8A: AddC(A, D),
	// ADC E
	0x8B: AddC(A, E),
	// ADC H
	0x8C: AddC(A, H),
	// ADC L
	0x8D: AddC(A, L),
	// ADC (HL)
	0x8E: AddC(A, HL_),
	// ADC A
	0x8F: AddC(A, A),
	// SUB B
	0x90: Sub(A, B),
	// SUB C
	0x91: Sub(A, C),
	// SUB D
	0x92: Sub(A, D),
	// SUB E
	0x93: Sub(A, E),
	// SUB H
	0x94: Sub(A, H),
	// SUB L
	0x95: Sub(A, L),
	// SUB (HL)
	0x96: Sub(A, HL_),
	// SUB A
	0x97: Sub(A, A),
	// SBC B
	0x98: SubC(A, B),
	// SBC C
	0x99: SubC(A, C),
	// SBC D
	0x9A: SubC(A, D),
	// SBC E
	0x9B: SubC(A, E),
	// SBC H
	0x9C: SubC(A, H),
	// SBC L
	0x9D: SubC(A, L),
	// SBC (HL)
	0x9E: Sub(A, HL_),
	// SBC A
	0x9F: SubC(A, A),
	// AND B
	0xA0: And(A, B),
	// AND C
	0xA1: And(A, C),
	// AND D
	0xA2: And(A, D),
	// AND E
	0xA3: And(A, E),
	// AND H
	0xA4: And(A, H),
	// AND L
	0xA5: And(A, L),
	// AND (HL)
	0xA6: And(A, HL_),
	// AND A
	0xA7: And(A, A),
	// Xor B
	0xA8: Xor(A, B),
	// Xor C
	0xA9: Xor(A, C),
	// Xor D
	0xAA: Xor(A, D),
	// Xor E
	0xAB: Xor(A, E),
	// Xor H
	0xAC: Xor(A, H),
	// Xor L
	0xAD: Xor(A, L),
	// Xor (HL)
	0xAE: Xor(A, HL_),
	// Xor A
	0xAF: Xor(A, A),
	// OR B
	0xB0: Or(A, B),
	// OR C
	0xB1: Or(A, C),
	// OR D
	0xB2: Or(A, D),
	// OR E
	0xB3: Or(A, E),
	// OR H
	0xB4: Or(A, H),
	// OR L
	0xB5: Or(A, L),
	// OR (HL)
	0xB6: Or(A, HL_),
	// OR A
	0xB7: Or(A, A),
	// CP B
	0xB8: Compare(A, B),
	// CP C
	0xB9: Compare(A, C),
	// CP D
	0xBA: Compare(A, D),
	// CP E
	0xBB: Compare(A, E),
	// CP H
	0xBC: Compare(A, H),
	// CP L
	0xBD: Compare(A, L),
	// CP (HL)
	0xBE: Compare(A, HL_),
	// CP A
	0xBF: Compare(A, A),
	// RET NZ
	0xC0: Return(NotZero),
	// POP BC
	0xC1: Pop(BC),
	// JP NZ,a16
	0xC2: Jump(NotZero),
	// JP a16
	0xC3: Jump(Always),
	// Call NZ,a16
	0xC4: Call(NotZero),
	// PUSH BC
	0xC5: Push(BC),
	// ADD A,d8
	0xC6: Add(A, D8),
	// ResetPC 00H
	0xC7: ResetPC(0x00),
	// RET Z
	0xC8: Return(Zero),
	// RET
	0xC9: Return(Always),
	// JP Z,a16
	0xCA: Jump(Zero),
	// PREFIX CB
	0xCB: PrefixCB(),
	// Call Z,a16
	0xCC: Call(Zero),
	// Call a16
	0xCD: Call(Always),
	// ADC A,d8
	0xCE: AddC(A, D8),
	// ResetPC 08H
	0xCF: ResetPC(0x08),
	// RET NC
	0xD0: Return(NotCarry),
	// POP DE
	0xD1: Pop(DE),
	// JP NC,a16
	0xD2: Jump(NotCarry),
	// INVALID
	// Call NC,a16
	0xD4: Call(NotCarry),
	// PUSH DE
	0xD5: Push(DE),
	// SUB d8
	0xD6: Sub(A, D8),
	// ResetPC 10H
	0xD7: ResetPC(0x10),
	// RET C
	0xD8: Return(Carry),
	// RETI
	0xD9: ReturnInterrupt(),
	// JP C,a16
	0xDA: Jump(Carry),
	// INVALID
	// Call C,a16
	0xDC: Call(Carry),
	// INVALID
	// SBC d8
	0xDE: SubC(A, D8),
	// ResetPC 18H
	0xDF: ResetPC(0x18),
	// LDH (a8),A
	0xE0: Load(A8_, A),
	// POP HL
	0xE1: Pop(HL),
	// LD (C),A
	0xE2: Load(C_, A),
	// INVALID
	// INVALID
	// PUSH HL
	0xE5: Push(HL),
	// AND d8
	0xE6: And(A, D8),
	// ResetPC 20H
	0xE7: ResetPC(0x20),
	// ADD SP,r8
	0xE8: AddSPImmediate(),
	// JP (HL)
	0xE9: JumpHL(), // <insert cussing here>
	// LD (a16),A
	0xEA: Load(A16_, A),
	// INVALID
	// INVALID
	// INVALID
	// Xor d8
	0xEE: Xor(A, D8),
	// ResetPC 28H
	0xEF: ResetPC(0x28),
	// LDH A,(a8)
	0xF0: Load(A, A8_),
	// POP AF
	0xF1: Pop(AF),
	// LD A,(C)
	0xF2: Load(A, C_),
	// DI
	0xF3: DisableInterrupts(),
	// INVALID
	// PUSH AF
	0xF5: Push(AF),
	// OR d8
	0xF6: Or(A, D8),
	// ResetPC 30H
	0xF7: ResetPC(0x30),
	// LD HL,SP+r8
	0xF8: LoadHLSPOffset(D8),
	// LD SP,HL
	0xF9: Load16(SP, HL),
	// LD A,(a16)
	0xFA: Load(A, A16_),
	// EI
	0xFB: EIWithDelay(),
	// INVALID
	// INVALID
	// CP d8
	0xFE: Compare(A, D8),
	// ResetPC 38H
	0xFF: ResetPC(0x38),
}
