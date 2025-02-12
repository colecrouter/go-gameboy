package registers

type Interrupt struct {
	Joypad bool
	Serial bool
	Timer  bool
	LCD    bool
	VBlank bool
}

func (i *Interrupt) Read() uint8 {
	var interrupt uint8
	if i.Joypad {
		interrupt |= 1 << 4
	}
	if i.Serial {
		interrupt |= 1 << 3
	}
	if i.Timer {
		interrupt |= 1 << 2
	}
	if i.LCD {
		interrupt |= 1 << 1
	}
	if i.VBlank {
		interrupt |= 1 << 0
	}
	return interrupt
}

func (i *Interrupt) Write(data uint8) {
	i.Joypad = data&(1<<4) != 0
	i.Serial = data&(1<<3) != 0
	i.Timer = data&(1<<2) != 0
	i.LCD = data&(1<<1) != 0
	i.VBlank = data&(1<<0) != 0
}
