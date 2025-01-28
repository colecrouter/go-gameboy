package registers

type Registers struct {
	buffer [128]uint8
}

type InterruptFlags struct {
	VBlank bool
	LCD    bool
	Timer  bool
	Serial bool
	Joypad bool
}

func (r *Registers) Read(addr uint16) uint8 {
	return r.buffer[addr]
}

func (r *Registers) Write(addr uint16, value uint8) {
	r.buffer[addr] = value
}

// TODO
func (r *Registers) ReadP1() uint8 {
	return r.buffer[0x00]
}

func (r *Registers) WriteP1(value uint8) {
	r.buffer[0x00] = value
}

func (r *Registers) ReadDivider() uint8 {
	return r.buffer[0x04]
}

func (r *Registers) WriteDivider(value uint8) {
	r.buffer[0x04] = value
}

func (r *Registers) ReadTimerCounter() uint8 {
	return r.buffer[0x05]
}

func (r *Registers) WriteTimerCounter(value uint8) {
	r.buffer[0x05] = value
}

func (r *Registers) ReadTimerModulo() uint8 {
	return r.buffer[0x06]
}

func (r *Registers) WriteTimerModulo(value uint8) {
	r.buffer[0x06] = value
}

func (r *Registers) ReadTimerControl() uint8 {
	return r.buffer[0x07]
}

func (r *Registers) WriteTimerControl(value uint8) {
	r.buffer[0x07] = value
}

func (r *Registers) ReadInterrupts() InterruptFlags {
	return InterruptFlags{
		VBlank: r.buffer[0x0F]&0x01 == 0x01,
		LCD:    r.buffer[0x0F]&0x02 == 0x02,
		Timer:  r.buffer[0x0F]&0x04 == 0x04,
		Serial: r.buffer[0x0F]&0x08 == 0x08,
		Joypad: r.buffer[0x0F]&0x10 == 0x10,
	}
}

func (r *Registers) WriteInterrupts(flags InterruptFlags) {
	var value uint8
	if flags.VBlank {
		value |= 0x01
	}
	if flags.LCD {
		value |= 0x02
	}
	if flags.Timer {
		value |= 0x04
	}
	if flags.Serial {
		value |= 0x08
	}
	if flags.Joypad {
		value |= 0x10
	}
	r.buffer[0x0F] = value
}
