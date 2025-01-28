package lr35902

type Flags struct{ Zero, Subtract, HalfCarry, Carry bool }

type FlagState int

const (
	Set FlagState = iota
	Reset
	Leave
)

func (c *LR35902) setFlags(zero, subtract, halfCarry, carry FlagState) {
	if zero == Set {
		c.flags.Zero = true
	} else if zero == Reset {
		c.flags.Zero = false
	}

	if subtract == Set {
		c.flags.Subtract = true
	} else if subtract == Reset {
		c.flags.Subtract = false
	}

	if halfCarry == Set {
		c.flags.HalfCarry = true
	} else if halfCarry == Reset {
		c.flags.HalfCarry = false
	}

	if carry == Set {
		c.flags.Carry = true
	} else if carry == Reset {
		c.flags.Carry = false
	}
}
