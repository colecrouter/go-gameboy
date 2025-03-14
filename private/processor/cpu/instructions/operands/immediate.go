package operands

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/helpers"
)

type ImmediateOperand8 struct{}

func (i *ImmediateOperand8) Read(c cpu.CPU) uint8 {
	return c.Read(c.Registers().PC)
}

func (i *ImmediateOperand8) Write(c cpu.CPU, val uint8) {
	c.Write(c.Registers().PC, val)
}

type ImmediateOperand16 struct{}

func (i *ImmediateOperand16) Read(c cpu.CPU) uint16 {
	low := c.Read(c.Registers().PC)
	high := c.Read(c.Registers().PC + 1)
	return helpers.ToRegisterPair(high, low)
}

func (i *ImmediateOperand16) Write(c cpu.CPU, val uint16) {
	high, low := helpers.FromRegisterPair(val)
	c.Write(c.Registers().PC, low)
	c.Write(c.Registers().PC+1, high)
}

type ImmediateIndirectOperand struct{}

func (i *ImmediateIndirectOperand) Read(c cpu.CPU) uint8 {
	return c.Read(c.Registers().PC)
}

func (i *ImmediateIndirectOperand) Write(c cpu.CPU, val uint8) {
	c.Write(c.Registers().PC, val)
}
