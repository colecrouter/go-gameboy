package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

// Jump
func jump(c cpu.CPU, addr uint16, condition bool) {
	if !condition {
		return
	}

	c.Registers().PC = addr - 1
}

// JumpRelative now adds 2 to account for the two-byte instruction length.
func jumpRelative(c cpu.CPU, offset int8, condition bool) {
	if !condition {
		return
	}

	// Add 2 for instruction length, then offset, then subtract 1 because MClock will increment PC.
	c.Registers().PC = uint16(int32(c.Registers().PC)+2+int32(offset)) - 1
}

// Subroutines
func ret(c cpu.CPU, condition bool) {
	if !condition {
		c.Clock()
		return
	}

	// Pop the return address in little-endian order.
	high, low := c.Read16(c.Registers().SP)
	c.Registers().SP += 2

	addr := cpu.ToRegisterPair(high, low)
	c.Registers().PC = addr - 1 // Adjust for later PC increment
}

func call(c cpu.CPU, addr uint16, condition bool) {
	if !condition {
		return
	}
	retAddr := c.Registers().PC + 3
	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, retAddr)
	// Update: subtract 1 to account for the later PC increment.
	c.Registers().PC = addr - 1
}

func rst(c cpu.CPU, addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.Registers().PC + 1

	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, retAddr)

	c.Registers().PC = addr - 1
}
