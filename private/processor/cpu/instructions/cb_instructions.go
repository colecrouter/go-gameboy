package instructions

import (
	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/enums"
	. "github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/generators"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/shared"
)

// All of the CB instructions are all the same per first bit
// The only difference is the register that the instruction operates on
// E.g. 0x00-0x07 are all RLC instructions, B, C, D, E, H, L, (HL), A

// Here we'll define two arrays of helper functions to help us construct the final instruction map

type instructionGenerator func(operands.Operand[uint8]) []shared.MicroOp

var firstHalfCbInstructionsHelper = [16]instructionGenerator{
	// RLC
	func(u operands.Operand[uint8]) []shared.MicroOp { return Rotate(u, true, false, true) },
	// RL
	func(u operands.Operand[uint8]) []shared.MicroOp { return Rotate(u, true, true, true) },
	// SLA
	func(u operands.Operand[uint8]) []shared.MicroOp { return Shift(u, true, false) },
	// SWAP
	func(u operands.Operand[uint8]) []shared.MicroOp { return Swap(u) },
	// BIT 0, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 0) },
	// BIT 2, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 2) },
	// BIT 4, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 4) },
	// BIT 6, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 6) },
	// RES 0, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 0) },
	// RES 2, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 2) },
	// RES 4, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 4) },
	// RES 6, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 6) },
	// SET 0, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 0) },
	// SET 2, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 2) },
	// SET 4, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 4) },
	// SET 6, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 6) },
}

var secondHalfCbInstructionsHelper = [16]instructionGenerator{
	// RRC
	func(u operands.Operand[uint8]) []shared.MicroOp { return Rotate(u, false, false, true) },
	// RR
	func(u operands.Operand[uint8]) []shared.MicroOp { return Rotate(u, false, true, true) },
	// SRA
	func(u operands.Operand[uint8]) []shared.MicroOp { return Shift(u, false, true) },
	// SRL
	func(u operands.Operand[uint8]) []shared.MicroOp { return Shift(u, false, false) },
	// BIT 1, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 1) },
	// BIT 3, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 3) },
	// BIT 5, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 5) },
	// BIT 7, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ReadBit(u, 7) },
	// RES 1, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 1) },
	// RES 3, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 3) },
	// RES 5, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 5) },
	// RES 7, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return ResetBit(u, 7) },
	// SET 1, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 1) },
	// SET 3, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 3) },
	// SET 5, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 5) },
	// SET 7, X
	func(u operands.Operand[uint8]) []shared.MicroOp { return SetBit(u, 7) },
}

func generateCbInstructions() [0x100]shared.Instruction {
	CBInstructions := [0x100]shared.Instruction{}

	// Register mapping: indices 0-7 correspond to B, C, D, E, H, L, (HL), A.
	regMapping := [8]operands.Operand[uint8]{
		B,
		C,
		D,
		E,
		H,
		L,
		HL_,
		A,
	}

	for y := range 16 {
		msb := uint8(y) << 4
		for x := range 16 {
			lsb := uint8(x)
			opcode := msb + lsb

			if x < 8 {
				// First half: use y as the op code selector.
				gen := firstHalfCbInstructionsHelper[y]
				CBInstructions[opcode] = gen(regMapping[x])

			} else {
				// Second half: use y as the op code selector.
				gen := secondHalfCbInstructionsHelper[y]
				CBInstructions[opcode] = gen(regMapping[x-8])
			}

		}
	}

	return CBInstructions
}

var CBInstructions = generateCbInstructions()
