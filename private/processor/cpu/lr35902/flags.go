package lr35902

type Flags struct {
	Zero,
	Subtract,
	HalfCarry,
	Carry bool
}

type FlagState int

const (
	Set FlagState = iota
	Reset
	Leave
)

func (c *LR35902) setFlags(zero, subtract, halfCarry, carry FlagState) {
	if zero == Set {
		c.Flags.Zero = true
	} else if zero == Reset {
		c.Flags.Zero = false
	}

	if subtract == Set {
		c.Flags.Subtract = true
	} else if subtract == Reset {
		c.Flags.Subtract = false
	}

	if halfCarry == Set {
		c.Flags.HalfCarry = true
	} else if halfCarry == Reset {
		c.Flags.HalfCarry = false
	}

	if carry == Set {
		c.Flags.Carry = true
	} else if carry == Reset {
		c.Flags.Carry = false
	}
}

func (f Flags) Read() uint8 {
	val := uint8(0)
	if f.Zero {
		val |= 1 << 7
	}
	if f.Subtract {
		val |= 1 << 6
	}
	if f.HalfCarry {
		val |= 1 << 5
	}
	if f.Carry {
		val |= 1 << 4
	}
	return val
}

func (f *Flags) Write(val uint8) {
	f.Zero = val&(1<<7) != 0
	f.Subtract = val&(1<<6) != 0
	f.HalfCarry = val&(1<<5) != 0
	f.Carry = val&(1<<4) != 0
}
