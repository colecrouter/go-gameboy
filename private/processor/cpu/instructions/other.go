package instructions

import "github.com/colecrouter/gameboy-go/private/processor/cpu"

// Other
func decimalAdjust(c cpu.CPU) {
	a := c.Registers().A
	subtract := c.Flags().Subtract
	halfCarry := c.Flags().HalfCarry
	carry := c.Flags().Carry
	offset := uint8(0)

	if !subtract {
		if halfCarry || (a&0x0F) > 0x09 {
			offset |= 0x06
		}
		if carry || a > 0x99 {
			offset |= 0x60
		}
		a += offset
	} else {
		if halfCarry {
			offset |= 0x06
		}
		if carry {
			offset |= 0x60
		}
		a -= offset
	}

	c.Registers().A = a
	c.Flags().Zero = (a == 0)
	c.Flags().HalfCarry = false
	c.Flags().Carry = ((offset & 0x60) != 0)
}
