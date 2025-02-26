package io

type InterruptFlags struct {
	VBlank bool
	LCD    bool
	Timer  bool
	Serial bool
	Joypad bool
}

func (i *InterruptFlags) Read(addr uint16) uint8 {
	var val uint8
	if i.VBlank {
		val |= 1 << 0
	}
	if i.LCD {
		val |= 1 << 1
	}
	if i.Timer {
		val |= 1 << 2
	}
	if i.Serial {
		val |= 1 << 3
	}
	if i.Joypad {
		val |= 1 << 4
	}
	return val
}

func (i *InterruptFlags) Write(addr uint16, value uint8) {
	i.VBlank = value&(1<<0) > 0
	i.LCD = value&(1<<1) > 0
	i.Timer = value&(1<<2) > 0
	i.Serial = value&(1<<3) > 0
	i.Joypad = value&(1<<4) > 0
}
