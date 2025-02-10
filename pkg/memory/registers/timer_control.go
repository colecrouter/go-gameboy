package registers

type Increment uint8

const (
	M256 Increment = 0
	M4   Increment = 1
	M16  Increment = 2
	M64  Increment = 3
)

type TimerControl struct {
	Speed   Increment
	Enabled bool
}

func (t *TimerControl) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	val := uint8(0)
	if t.Enabled {
		val |= 1 << 2
	}
	val |= uint8(t.Speed)
	return val
}

func (t *TimerControl) Write(addr uint16, value uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	t.Enabled = value&(1<<2) > 0
	t.Speed = Increment(value & 0x3)
}
