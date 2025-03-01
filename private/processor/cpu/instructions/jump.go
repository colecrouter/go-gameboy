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

	// Add 2 to account for the two-byte instruction length.
	// Then, subtract 1 for GetImmediate8's increment, and 1 again for the PC increment.
	// Woo, look we didn't have to add anything.
	c.Registers().PC = uint16(int32(c.Registers().PC) + int32(offset))
}

// Subroutines
func ret(c cpu.CPU, condition bool) {
	if !condition {
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

	// Theoretically, we want to set up the stack so that RET drops us exactly 1 PC before the next instruction.
	// However, the value pushed onto the stack matters. In the case of the Game Boy, the value pushed is the address
	// At this point in the CALL instruction (3 bytes long) we will be at initial PC + 2, so add 1 to get the correct return address.
	retAddr := c.Registers().PC + 1

	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, retAddr)
	// Subtract 1 to account for the later PC increment.
	c.Registers().PC = addr - 1
}

func rst(c cpu.CPU, addr uint16) {
	// For RST, instruction size is 1 byte.
	// There is no immediate offset or anything like that, so we need to add 1 to return after the instruction.
	retAddr := c.Registers().PC + 1

	c.Registers().SP -= 2
	c.Write16(c.Registers().SP, retAddr)

	c.Registers().PC = addr - 1
}
