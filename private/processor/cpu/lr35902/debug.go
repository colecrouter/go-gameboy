package lr35902

import "fmt"

var mnemonics = [0x100]string{
	0x00: "NOP",
	0x01: "LD BC,d16",
	0x02: "LD (BC),A",
	0x03: "INC BC",
	0x04: "INC B",
	0x05: "DEC B",
	0x06: "LD B,d8",
	0x07: "RLCA",
	0x08: "LD (a16),SP",
	0x09: "ADD HL,BC",
	0x0A: "LD A,(BC)",
	0x0B: "DEC BC",
	0x0C: "INC C",
	0x0D: "DEC C",
	0x0E: "LD C,d8",
	0x0F: "RRCA",
	0x10: "STOP 0",
	0x11: "LD DE,d16",
	0x12: "LD (DE),A",
	0x13: "INC DE",
	0x14: "INC D",
	0x15: "DEC D",
	0x16: "LD D,d8",
	0x17: "RLA",
	0x18: "JR r8",
	0x19: "ADD HL,DE",
	0x1A: "LD A,(DE)",
	0x1B: "DEC DE",
	0x1C: "INC E",
	0x1D: "DEC E",
	0x1E: "LD E,d8",
	0x1F: "RRA",
	0x20: "JR NZ,r8",
	0x21: "LD HL,d16",
	0x22: "LD (HL+),A",
	0x23: "INC HL",
	0x24: "INC H",
	0x25: "DEC H",
	0x26: "LD H,d8",
	0x27: "DAA",
	0x28: "JR Z,r8",
	0x29: "ADD HL,HL",
	0x2A: "LD A,(HL+)",
	0x2B: "DEC HL",
	0x2C: "INC L",
	0x2D: "DEC L",
	0x2E: "LD L,d8",
	0x2F: "CPL",
	0x30: "JR NC,r8",
	0x31: "LD SP,d16",
	0x32: "LD (HL-),A",
	0x33: "INC SP",
	0x34: "INC (HL)",
	0x35: "DEC (HL)",
	0x36: "LD (HL),d8",
	0x37: "SCF",
	0x38: "JR C,r8",
	0x39: "ADD HL,SP",
	0x3A: "LD A,(HL-)",
	0x3B: "DEC SP",
	0x3C: "INC A",
	0x3D: "DEC A",
	0x3E: "LD A,d8",
	0x3F: "CCF",
	0x40: "LD B,B",
	0x41: "LD B,C",
	0x42: "LD B,D",
	0x43: "LD B,E",
	0x44: "LD B,H",
	0x45: "LD B,L",
	0x46: "LD B,(HL)",
	0x47: "LD B,A",
	0x48: "LD C,B",
	0x49: "LD C,C",
	0x4A: "LD C,D",
	0x4B: "LD C,E",
	0x4C: "LD C,H",
	0x4D: "LD C,L",
	0x4E: "LD C,(HL)",
	0x4F: "LD C,A",
	0x50: "LD D,B",
	0x51: "LD D,C",
	0x52: "LD D,D",
	0x53: "LD D,E",
	0x54: "LD D,H",
	0x55: "LD D,L",
	0x56: "LD D,(HL)",
	0x57: "LD D,A",
	0x58: "LD E,B",
	0x59: "LD E,C",
	0x5A: "LD E,D",
	0x5B: "LD E,E",
	0x5C: "LD E,H",
	0x5D: "LD E,L",
	0x5E: "LD E,(HL)",
	0x5F: "LD E,A",
	0x60: "LD H,B",
	0x61: "LD H,C",
	0x62: "LD H,D",
	0x63: "LD H,E",
	0x64: "LD H,H",
	0x65: "LD H,L",
	0x66: "LD H,(HL)",
	0x67: "LD H,A",
	0x68: "LD L,B",
	0x69: "LD L,C",
	0x6A: "LD L,D",
	0x6B: "LD L,E",
	0x6C: "LD L,H",
	0x6D: "LD L,L",
	0x6E: "LD L,(HL)",
	0x6F: "LD L,A",
	0x70: "LD (HL),B",
	0x71: "LD (HL),C",
	0x72: "LD (HL),D",
	0x73: "LD (HL),E",
	0x74: "LD (HL),H",
	0x75: "LD (HL),L",
	0x76: "HALT",
	0x77: "LD (HL),A",
	0x78: "LD A,B",
	0x79: "LD A,C",
	0x7A: "LD A,D",
	0x7B: "LD A,E",
	0x7C: "LD A,H",
	0x7D: "LD A,L",
	0x7E: "LD A,(HL)",
	0x7F: "LD A,A",
	0x80: "ADD A,B",
	0x81: "ADD A,C",
	0x82: "ADD A,D",
	0x83: "ADD A,E",
	0x84: "ADD A,H",
	0x85: "ADD A,L",
	0x86: "ADD A,(HL)",
	0x87: "ADD A,A",
	0x88: "ADC A,B",
	0x89: "ADC A,C",
	0x8A: "ADC A,D",
	0x8B: "ADC A,E",
	0x8C: "ADC A,H",
	0x8D: "ADC A,L",
	0x8E: "ADC A,(HL)",
	0x8F: "ADC A,A",
	0x90: "SUB B",
	0x91: "SUB C",
	0x92: "SUB D",
	0x93: "SUB E",
	0x94: "SUB H",
	0x95: "SUB L",
	0x96: "SUB (HL)",
	0x97: "SUB A",
	0x98: "SBC A,B",
	0x99: "SBC A,C",
	0x9A: "SBC A,D",
	0x9B: "SBC A,E",
	0x9C: "SBC A,H",
	0x9D: "SBC A,L",
	0x9E: "SBC A,(HL)",
	0x9F: "SBC A,A",
	0xA0: "AND B",
	0xA1: "AND C",
	0xA2: "AND D",
	0xA3: "AND E",
	0xA4: "AND H",
	0xA5: "AND L",
	0xA6: "AND (HL)",
	0xA7: "AND A",
	0xA8: "XOR B",
	0xA9: "XOR C",
	0xAA: "XOR D",
	0xAB: "XOR E",
	0xAC: "XOR H",
	0xAD: "XOR L",
	0xAE: "XOR (HL)",
	0xAF: "XOR A",
	0xB0: "OR B",
	0xB1: "OR C",
	0xB2: "OR D",
	0xB3: "OR E",
	0xB4: "OR H",
	0xB5: "OR L",
	0xB6: "OR (HL)",
	0xB7: "OR A",
	0xB8: "CP B",
	0xB9: "CP C",
	0xBA: "CP D",
	0xBB: "CP E",
	0xBC: "CP H",
	0xBD: "CP L",
	0xBE: "CP (HL)",
	0xBF: "CP A",
	0xC0: "RET NZ",
	0xC1: "POP BC",
	0xC2: "JP NZ,a16",
	0xC3: "JP a16",
	0xC4: "CALL NZ,a16",
	0xC5: "PUSH BC",
	0xC6: "ADD A,d8",
	0xC7: "RST 00H",
	0xC8: "RET Z",
	0xC9: "RET",
	0xCA: "JP Z,a16",
	0xCB: "PREFIX CB",
	0xCC: "CALL Z,a16",
	0xCD: "CALL a16",
	0xCE: "ADC A,d8",
	0xCF: "RST 08H",
	0xD0: "RET NC",
	0xD1: "POP DE",
	0xD2: "JP NC,a16",
	0xD3: "INVALID",
	0xD4: "CALL NC,a16",
	0xD5: "PUSH DE",
	0xD6: "SUB d8",
	0xD7: "RST 10H",
	0xD8: "RET C",
	0xD9: "RETI",
	0xDA: "JP C,a16",
	0xDB: "INVALID",
	0xDC: "CALL C,a16",
	0xDD: "INVALID",
	0xDE: "SBC A,d8",
	0xDF: "RST 18H",
	0xE0: "LDH (a8),A",
	0xE1: "POP HL",
	0xE2: "LD (C),A",
	0xE3: "INVALID",
	0xE4: "INVALID",
	0xE5: "PUSH HL",
	0xE6: "AND d8",
	0xE7: "RST 20H",
	0xE8: "ADD SP,r8",
	0xE9: "JP (HL)",
	0xEA: "LD (a16),A",
	0xEB: "INVALID",
	0xEC: "INVALID",
	0xED: "INVALID",
	0xEE: "XOR d8",
	0xEF: "RST 28H",
	0xF0: "LDH A,(a8)",
	0xF1: "POP AF",
	0xF2: "LD A,(C)",
	0xF3: "DI",
	0xF4: "INVALID",
	0xF5: "PUSH AF",
	0xF6: "OR d8",
	0xF7: "RST 30H",
	0xF8: "LD HL,SP+r8",
	0xF9: "LD SP,HL",
	0xFA: "LD A,(a16)",
	0xFB: "EI",
	0xFC: "INVALID",
	0xFD: "INVALID",
	0xFE: "CP d8",
	0xFF: "RST 38H",
}

// getCBMnemonic returns a string mnemonic for a given CB opcode.
func getCBMnemonic(op byte) string {
	// Get target register based on lower 3 bits
	reg := []string{"B", "C", "D", "E", "H", "L", "(HL)", "A"}[op&0x07]

	// Determine instruction based on upper bits
	// First half
	if (op & 0x0F) < 0x08 {
		switch op >> 4 {
		case 0x0:
			return fmt.Sprintf("RLC %s", reg)
		case 0x1:
			return fmt.Sprintf("RL %s", reg)
		case 0x2:
			return fmt.Sprintf("SLA %s", reg)
		case 0x3:
			return fmt.Sprintf("SWAP %s", reg)
		case 0x4:
			return fmt.Sprintf("BIT 0,%s", reg)
		case 0x5:
			return fmt.Sprintf("BIT 2,%s", reg)
		case 0x6:
			return fmt.Sprintf("BIT 4,%s", reg)
		case 0x7:
			return fmt.Sprintf("BIT 6,%s", reg)
		case 0x8:
			return fmt.Sprintf("RES 0,%s", reg)
		case 0x9:
			return fmt.Sprintf("RES 2,%s", reg)
		case 0xA:
			return fmt.Sprintf("RES 4,%s", reg)
		case 0xB:
			return fmt.Sprintf("RES 6,%s", reg)
		case 0xC:
			return fmt.Sprintf("SET 0,%s", reg)
		case 0xD:
			return fmt.Sprintf("SET 2,%s", reg)
		case 0xE:
			return fmt.Sprintf("SET 4,%s", reg)
		case 0xF:
			return fmt.Sprintf("SET 6,%s", reg)
		}
	} else {
		// Second half
		switch op >> 4 {
		case 0x0:
			return fmt.Sprintf("RRC %s", reg)
		case 0x1:
			return fmt.Sprintf("RR %s", reg)
		case 0x2:
			return fmt.Sprintf("SRA %s", reg)
		case 0x3:
			return fmt.Sprintf("SRL %s", reg)
		case 0x4:
			return fmt.Sprintf("BIT 1,%s", reg)
		case 0x5:
			return fmt.Sprintf("BIT 3,%s", reg)
		case 0x6:
			return fmt.Sprintf("BIT 5,%s", reg)
		case 0x7:
			return fmt.Sprintf("BIT 7,%s", reg)
		case 0x8:
			return fmt.Sprintf("RES 1,%s", reg)
		case 0x9:
			return fmt.Sprintf("RES 3,%s", reg)
		case 0xA:
			return fmt.Sprintf("RES 5,%s", reg)
		case 0xB:
			return fmt.Sprintf("RES 7,%s", reg)
		case 0xC:
			return fmt.Sprintf("SET 1,%s", reg)
		case 0xD:
			return fmt.Sprintf("SET 3,%s", reg)
		case 0xE:
			return fmt.Sprintf("SET 5,%s", reg)
		case 0xF:
			return fmt.Sprintf("SET 7,%s", reg)
		}
	}
	return "Unknown CB"
}

// PrintRegisters prints the current state of the CPU registers and flags.
func (c *LR35902) PrintRegisters() {
	fmt.Printf("\r\nRegisters:\r\n")
	fmt.Printf("A: 0x%02X  B: 0x%02X  C: 0x%02X  D: 0x%02X  E: 0x%02X\r\n", c.registers.A, c.registers.B, c.registers.C, c.registers.D, c.registers.E)
	fmt.Printf("H: 0x%02X  L: 0x%02X\r\n", c.registers.H, c.registers.L)
	fmt.Printf("SP: 0x%04X  PC: 0x%04X\r\n", c.registers.SP, c.registers.PC)
	fmt.Printf("Flags: Z=%t  N=%t  H=%t  C=%t\r\n", c.Flags().Zero, c.Flags().Subtract, c.Flags().HalfCarry, c.Flags().Carry)
}
