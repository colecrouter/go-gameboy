package instructions

import "github.com/colecrouter/gameboy-go/private/processor/cpu"

// All of the CB instructions are all the same per first bit
// The only difference is the register that the instruction operates on
// E.g. 0x00-0x07 are all RLC instructions, B, C, D, E, H, L, (HL), A

// Here we'll define two arrays of helper functions to help us construct the final instruction map

type instructionGenerator func(*uint8) func(cpu.CPU)

var firstHalfCbInstructionsHelper = [16]instructionGenerator{
	// RLC
	func(u *uint8) func(cpu.CPU) {
		return func(c cpu.CPU) { rotate(c, u, true, false, true) }
	},
	// RL
	func(u *uint8) func(cpu.CPU) {
		return func(c cpu.CPU) { rotate(c, u, true, true, true) }
	},
	// SLA
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { shift(c, u, true, false) } },
	// SWAP
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { swap(c, u) } },
	// BIT 0, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 0, *u) } },
	// BIT 2, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 2, *u) } },
	// BIT 4, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 4, *u) } },
	// BIT 6, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 6, *u) } },
	// RES 0, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 0, u) } },
	// RES 2, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 2, u) } },
	// RES 4, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 4, u) } },
	// RES 6, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 6, u) } },
	// SET 0, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 0, u) } },
	// SET 2, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 2, u) } },
	// SET 4, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 4, u) } },
	// SET 6, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 6, u) } },
}

var secondHalfCbInstructionsHelper = [16]instructionGenerator{
	// RRC
	func(u *uint8) func(cpu.CPU) {
		return func(c cpu.CPU) { rotate(c, u, false, false, true) }
	},
	// RR
	func(u *uint8) func(cpu.CPU) {
		return func(c cpu.CPU) { rotate(c, u, false, true, true) }
	},
	// SRA
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { shift(c, u, false, true) } },
	// SRL
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { shift(c, u, false, false) } },
	// BIT 1, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 1, *u) } },
	// BIT 3, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 3, *u) } },
	// BIT 5, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 5, *u) } },
	// BIT 7, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { bit(c, 7, *u) } },
	// RES 1, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 1, u) } },
	// RES 3, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 3, u) } },
	// RES 5, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 5, u) } },
	// RES 7, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { res(c, 7, u) } },
	// SET 1, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 1, u) } },
	// SET 3, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 3, u) } },
	// SET 5, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 5, u) } },
	// SET 7, X
	func(u *uint8) func(cpu.CPU) { return func(c cpu.CPU) { set(c, 7, u) } },
}

func generateCbInstructions() [0x100]Instruction {
	CBInstructions := [0x100]Instruction{}

	// Register mapping: indices 0-7 correspond to B, C, D, E, H, L, (HL), A.
	regMapping := func(c cpu.CPU) [8]*uint8 {
		return [8]*uint8{
			&c.Registers().B,
			&c.Registers().C,
			&c.Registers().D,
			&c.Registers().E,
			&c.Registers().H,
			&c.Registers().L,
			nil, // (HL) handled specially below
			&c.Registers().A,
		}
	}

	for y := range 16 {
		msb := uint8(y) << 4
		for x := range 16 {
			lsb := uint8(x)
			opcode := msb + lsb

			if x < 8 {
				// First half: use y as the op code selector.
				gen := firstHalfCbInstructionsHelper[y]
				if x == 6 {
					CBInstructions[opcode] = func(c cpu.CPU) {
						c.Clock()
						addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
						val := c.Read(addr)
						c.Ack()

						gen(&val)(c)

						// If the instruction is BIT, don't write back to memory.
						if opcode>>6 != 0x1 {
							c.Write(addr, val)
						}

						c.Clock()
						// TODO
						c.Ack()

					}
				} else {
					CBInstructions[opcode] = func(c cpu.CPU) {
						regs := regMapping(c)
						target := regs[x]
						gen(target)(c)
					}
				}
			} else {
				// Second half: use y as the op code selector.
				gen := secondHalfCbInstructionsHelper[y]
				if x == 14 {
					CBInstructions[opcode] = func(c cpu.CPU) {
						c.Clock()
						addr := cpu.ToRegisterPair(c.Registers().H, c.Registers().L)
						val := c.Read(addr)
						c.Ack()

						gen(&val)(c)

						// If the instruction is BIT, don't write back to memory.
						if opcode>>6 != 0x1 {
							c.Write(addr, val)
						}
						c.Clock()
						// TODO
						c.Ack()
					}
				} else {
					CBInstructions[opcode] = func(c cpu.CPU) {
						regs := regMapping(c)
						target := regs[x-8]
						gen(target)(c)
					}
				}
			}

		}
	}

	return CBInstructions
}

var CBInstructions = generateCbInstructions()
