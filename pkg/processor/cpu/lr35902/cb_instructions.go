package lr35902

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

	for i := 0; i < 16; i++ {
		// We're iterating over rows e.g. 0x10, 0x20, 0x30, etc.
		// Each entry will be the LSB e.g. 0x01, 0x02, 0x03, etc.
		msb := uint8(i) << 4

		// Generate the first half of the instructions
		for j, gen := range firstHalfCbInstructionsHelper {
			if j == 6 {
				CBInstructions[msb+0x06] = instruction{
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

			CBInstructions[msb+uint8(j)] = instruction{
				c: 8,
				op: func(c *LR35902) {
					gen(&c.registers.b)(c)
				},
			}
		}

		// Generate the second half of the instructions
		for j, gen := range secondHalfCbInstructionsHelper {
			if j == 6 {
				CBInstructions[msb+0x16] = instruction{
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

			CBInstructions[msb+0x08+uint8(j)] = instruction{
				c: 8,
				op: func(c *LR35902) {
					gen(&c.registers.b)(c)
				},
			}
		}
	}

	return CBInstructions
}

var cbInstructions = generateCbInstructions()
