package instructions

import (
	"github.com/colecrouter/gameboy-go/private/processor/cpu"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/flags"
)

// Register manipulation
func inc8(c cpu.CPU, r *uint8) {
	old := *r
	*r++
	zero := flags.Reset
	hc := flags.Reset
	// Half-carry set when lower nibble overflows (0xF -> 0x0)
	if old&0xF == 0xF {
		hc = flags.Set
	}
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Reset, hc, flags.Leave)
}
func inc16(c cpu.CPU, high, low *uint8) {
	c.Clock()
	combined := cpu.ToRegisterPair(*high, *low)
	combined++
	c.Ack()
	*high, *low = cpu.FromRegisterPair(combined)
}
func dec8(c cpu.CPU, r *uint8) {
	old := *r
	*r--
	zero := flags.Reset
	hc := flags.Reset
	// Half-carry set when lower nibble underflows (0x0 -> 0xF)
	if old&0xF == 0x0 {
		hc = flags.Set
	}
	if *r == 0 {
		zero = flags.Set
	}
	c.Flags().Set(zero, flags.Set, hc, flags.Leave)
}
func dec16(c cpu.CPU, high, low *uint8) {
	c.Clock()
	combined := cpu.ToRegisterPair(*high, *low)
	combined--
	c.Ack()
	*high, *low = cpu.FromRegisterPair(combined)
}
