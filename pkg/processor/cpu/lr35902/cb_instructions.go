package lr35902

import "fmt"

// All of the CB instructions are all the same per first bit
// The only difference is the register that the instruction operates on
// E.g. 0x00-0x07 are all RLC instructions, B, C, D, E, H, L, (HL), A

// Here we'll define two arrays of helper functions to help us construct the final instruction map

type instructionGenerator func(*uint8) func(*LR35902)

var firstHalfCbInstructionsHelper = [16]instructionGenerator{
	// RLC
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.rotate(u, true, true) } },
	// RL
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.rotate(u, true, false) } },
	// SLA
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.shift(u, true) } },
	// SWAP
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.swap(u) } },
	// BIT 0, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(0, *u) } },
	// BIT 2, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(2, *u) } },
	// BIT 4, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(4, *u) } },
	// BIT 6, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(6, *u) } },
	// RES 0, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(0, u) } },
	// RES 2, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(2, u) } },
	// RES 4, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(4, u) } },
	// RES 6, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(6, u) } },
	// SET 0, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(0, u) } },
	// SET 2, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(2, u) } },
	// SET 4, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(4, u) } },
	// SET 6, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(6, u) } },
}

var secondHalfCbInstructionsHelper = [16]instructionGenerator{
	// RRC
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.rotate(u, false, true) } },
	// RR
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.rotate(u, false, false) } },
	// SRA
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.shift(u, false) } },
	// SRL
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.shift(u, true) } },
	// BIT 1, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(1, *u) } },
	// BIT 3, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(3, *u) } },
	// BIT 5, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(5, *u) } },
	// BIT 7, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.bit(7, *u) } },
	// RES 1, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(1, u) } },
	// RES 3, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(3, u) } },
	// RES 5, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(5, u) } },
	// RES 7, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.res(7, u) } },
	// SET 1, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(1, u) } },
	// SET 3, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(3, u) } },
	// SET 5, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(5, u) } },
	// SET 7, X
	func(u *uint8) func(*LR35902) { return func(c *LR35902) { c.set(7, u) } },
}

func generateCbInstructions() map[uint8]instruction {
	CBInstructions := map[uint8]instruction{}

	// Register mapping: indices 0-7 correspond to B, C, D, E, H, L, (HL), A.
	// nil indicates (HL) which is handled specially.
	regMapping := func(c *LR35902) []*uint8 {
		return []*uint8{
			&c.registers.b,
			&c.registers.c,
			&c.registers.d,
			&c.registers.e,
			&c.registers.h,
			&c.registers.l,
			nil, // (HL)
			&c.registers.a,
		}
	}

	for i := 0; i < 16; i++ {
		msb := uint8(i) << 4
		// Generate first half instructions
		for j, gen := range firstHalfCbInstructionsHelper {
			// Determine target register index (0-7)
			regIndex := j & 0x07
			CBOpcode := msb + uint8(j)
			if regIndex == 6 { // (HL) special handling
				CBInstructions[CBOpcode] = instruction{
					c: 16,
					op: func(c *LR35902) {
						addr := toRegisterPair(c.registers.h, c.registers.l)
						val := c.bus.Read(addr)
						c.rotate(&val, true, true)
						c.bus.Write(addr, val)
					},
				}
				continue
			}
			// Capture for closure
			genFun := gen
			// Use the register from mapping
			CBInstructions[CBOpcode] = instruction{
				c: 8,
				op: func(c *LR35902) {
					regs := regMapping(c)
					target := regs[regIndex]
					genFun(target)(c)
				},
			}
		}

		// Generate second half instructions
		for j, gen := range secondHalfCbInstructionsHelper {
			// For second half, lower 3 bits come from (0x08+j)
			regIndex := (0x08 + j) & 0x07
			CBOpcode := msb + 0x08 + uint8(j)
			if regIndex == 6 { // (HL) special handling
				CBInstructions[CBOpcode] = instruction{
					c: 16,
					op: func(c *LR35902) {
						addr := toRegisterPair(c.registers.h, c.registers.l)
						val := c.bus.Read(addr)
						c.shift(&val, false)
						c.bus.Write(addr, val)
					},
				}
				continue
			}
			// Capture for closure
			genFun := gen
			CBInstructions[CBOpcode] = instruction{
				c: 8,
				op: func(c *LR35902) {
					regs := regMapping(c)
					target := regs[regIndex]
					genFun(target)(c)
				},
			}
		}
	}

	return CBInstructions
}

var cbInstructions = generateCbInstructions()

// getCBMnemonic returns a string mnemonic for a given CB opcode.
func getCBMnemonic(op byte) string {
	// Get target register based on lower 3 bits
	reg := []string{"B", "C", "D", "E", "H", "L", "(HL)", "A"}[op&0x07]

	// Determine instruction based on upper bits
	// First half
	if op < 0x40 {
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
			return fmt.Sprintf("RES 2,%s", reg)
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
