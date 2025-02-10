package registers

type Timer struct {
	Divider uint8        // FF04
	Counter uint8        // FF05
	Modulo  uint8        // FF06
	Control TimerControl // FF07

	totalCycles uint64
}

func (t *Timer) Read(addr uint16) uint8 {
	switch addr {
	case 0x0:
		return t.Divider
	case 0x1:
		return t.Counter
	case 0x2:
		return t.Modulo
	case 0x3:
		return t.Control.Read(0)
	default:
		panic("Invalid address")
	}
}

func (t *Timer) Write(addr uint16, value uint8) {
	switch addr {
	case 0x0:
		t.Divider = 0
	case 0x1:
		t.Counter = value
	case 0x2:
		t.Modulo = value
	case 0x3:
		t.Control.Write(0, value)
	default:
		panic("Invalid address")
	}
}

func (t *Timer) Clock(interrupt func()) {
	t.totalCycles++
	if t.totalCycles%256 == 0 {
		t.Divider++
	}

	if !t.Control.Enabled {
		return
	}

	var interval uint64
	switch t.Control.Speed {
	case M256:
		interval = 256
	case M4:
		interval = 4
	case M16:
		interval = 16
	case M64:
		interval = 64
	}

	if t.totalCycles%interval == 0 {
		t.Counter++

		// Check for overflow
		if t.Counter == 0 {
			t.Counter = t.Modulo

			// Request interrupt
			interrupt()
		}
	}

}
