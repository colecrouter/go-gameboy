package registers

type Palette struct {
	Colors [4]uint8
}

func (p *Palette) Write(addr uint16, data uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	p.Colors[0] = data & 0x3
	p.Colors[1] = (data >> 2) & 0x3
	p.Colors[2] = (data >> 4) & 0x3
	p.Colors[3] = (data >> 6) & 0x3
}

func (p *Palette) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	return p.Colors[0] | (p.Colors[1] << 2) | (p.Colors[2] << 4) | (p.Colors[3] << 6)
}
