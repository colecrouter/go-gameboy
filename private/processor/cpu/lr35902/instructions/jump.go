package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

// Jump
func jump(c cpu.CPU, addr uint16, condition bool) {
	if condition {
		c.Registers().PC = addr - 1
		return
	}

	c.Registers().PC += 3 // 3-byte instruction
}
func jumpRelative(c cpu.CPU, offset int8, condition bool) {
	if condition {
		c.Registers().PC = uint16(int32(c.Registers().PC)+2+int32(offset)) - 1
		return
	}

	// Updated: For a 2-byte jumpRelative, false branch adds 1 (MClock will add 1 later)
	c.Registers().PC += 1
}

// Subroutines
func ret(c cpu.CPU, condition bool) {
	if !condition {
		c.Clock()
		// Reverse the PC increment caused by Clock() so that PC remains unchanged.
		c.Registers().PC--
		return
	}

	// Pop the return address in little-endian order.
	high, low := c.Read16(c.Registers().SP)
	c.Registers().SP += 2

	addr := cpu.ToRegisterPair(high, low)
	c.Registers().PC = addr - 1 // Adjust for later PC increment
}

func call(c cpu.CPU, addr uint16, condition bool) {
	// // Use the instruction's starting PC (c.lastPC) to compute return address
	// if !condition {
	// 	c.Registers().PC += 3
	// 	return
	// }

	// // True branch: push return address = c.lastPC+3
	// retAddr := c.lastPC + 3
	// c.Registers().SP -= 2
	// c.Write16(c.Registers().SP, retAddr)

	// c.Registers().PC = addr - 1 // Adjust for MClock increment.
	panic("Not implemented")
}

func rst(c cpu.CPU, addr uint16) {
	// For RST, instruction size is 1 byte.
	retAddr := c.Registers().PC + 1

	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, retAddr)

	c.Registers().PC = addr - 1
}
