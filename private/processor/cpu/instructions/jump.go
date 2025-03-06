package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
)

// Jump
func jump(c cpu.CPU, addr uint16, condition bool) {
	if !condition {
		return
	}

	c.Clock()
	c.Registers().PC = addr - 1
	c.Ack()

}

// JumpRelative now adds 2 to account for the two-byte instruction length.
func jumpRelative(c cpu.CPU, offset int8, condition bool) {
	if !condition {
		return
	}

	// Add 2 to account for the two-byte instruction length.
	// Then, subtract 1 for GetImmediate8's increment, and 1 again for the PC increment.
	// Woo, look we didn't have to add anything.

	c.Clock()
	c.Registers().PC = uint16(int32(c.Registers().PC) + int32(offset))
	c.Ack()
}

// Subroutines
func ret(c cpu.CPU, condition bool) {
	if !condition {
		return
	}
	c.Clock()
	c.Ack()

	c.Clock()
	low := c.Read(c.Registers().SP)
	c.Registers().SP++
	c.Ack()

	c.Clock()
	high := c.Read(c.Registers().SP)
	c.Registers().SP++
	c.Ack()

	c.Clock()
	addr := cpu.ToRegisterPair(high, low)
	c.Registers().PC = addr - 1 // Adjust for later PC increment
	c.Ack()
}

func call(c cpu.CPU, addr uint16, condition bool) {
	if !condition {
		return
	}

	// Theoretically, we want to set up the stack so that RET drops us exactly 1 PC before the next instruction.
	// However, the value pushed onto the stack matters. In the case of the Game Boy, the value pushed is the address
	// At this point in the CALL instruction (3 bytes long) we will be at initial PC + 2, so add 1 to get the correct return address.
	high, low := cpu.FromRegisterPair(c.Registers().PC + 1)

	c.Clock()
	c.Registers().SP--
	c.Write(c.Registers().SP, high)
	c.Ack()

	c.Clock()
	c.Registers().SP--
	c.Write(c.Registers().SP, low)
	c.Ack()

	c.Clock()
	c.Registers().PC = addr - 1
	c.Ack()
}

// func rst(c cpu.CPU, addr uint16) {
// 	// For RST, instruction size is 1 byte.
// 	// There is no immediate offset or anything like that, so we need to add 1 to return after the instruction.
// 	// retAddr := c.Registers().PC + 1
// 	high, low := cpu.FromRegisterPair(c.Registers().PC + 1)

// 	c.Clock()
// 	c.Registers().SP--
// 	c.Write(c.Registers().SP, low)
// 	c.Ack()

// 	c.Clock()
// 	c.Registers().SP--
// 	c.Write(c.Registers().SP, high)
// 	c.Ack()

// 	c.Registers().PC = addr - 1
// }
