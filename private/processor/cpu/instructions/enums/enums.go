package enums

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/conditions"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/instructions/operands"
)

var (
	A = &operands.RegisterOperand{Register: operands.A}
	B = &operands.RegisterOperand{Register: operands.B}
	C = &operands.RegisterOperand{Register: operands.C}
	D = &operands.RegisterOperand{Register: operands.D}
	E = &operands.RegisterOperand{Register: operands.E}
	H = &operands.RegisterOperand{Register: operands.H}
	L = &operands.RegisterOperand{Register: operands.L}
	F = &operands.RegisterOperand{Register: operands.F}

	AF = &operands.RegisterPairOperand{RegisterPair: operands.AF}
	BC = &operands.RegisterPairOperand{RegisterPair: operands.BC}
	DE = &operands.RegisterPairOperand{RegisterPair: operands.DE}
	HL = &operands.RegisterPairOperand{RegisterPair: operands.HL}
	SP = &operands.RegisterPairOperand{RegisterPair: operands.SP}

	AF_ = &operands.IndirectOperand{Indirectable: operands.AF_}
	BC_ = &operands.IndirectOperand{Indirectable: operands.BC_}
	DE_ = &operands.IndirectOperand{Indirectable: operands.DE_}
	HL_ = &operands.IndirectOperand{Indirectable: operands.HL_}
	SP_ = &operands.IndirectOperand{Indirectable: operands.SP_}
	A_  = &operands.IndirectOperand{Indirectable: operands.A_}
	B_  = &operands.IndirectOperand{Indirectable: operands.B_}
	C_  = &operands.IndirectOperand{Indirectable: operands.C_}
	D_  = &operands.IndirectOperand{Indirectable: operands.D_}
	E_  = &operands.IndirectOperand{Indirectable: operands.E_}
	H_  = &operands.IndirectOperand{Indirectable: operands.H_}
	L_  = &operands.IndirectOperand{Indirectable: operands.L_}
	F_  = &operands.IndirectOperand{Indirectable: operands.F_}

	D8   = &operands.ImmediateOperand8{}        // 8-bit immediate
	D16  = &operands.ImmediateOperand16{}       // 16-bit immediate
	A8_  = &operands.ImmediateIndirectOperand{} // 8-bit indirect
	A16_ = &operands.ImmediateIndirectOperand{} // 16-bit indirect

	Always   = conditions.Always // Always
	Zero     = conditions.Z      // If the Zero flag is set
	NotZero  = conditions.NZ     // If the Zero flag is not set
	Carry    = conditions.C      // If the Carry flag is set
	NotCarry = conditions.NC     // If the Carry flag is not set
)
