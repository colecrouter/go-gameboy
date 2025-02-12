package registers

type Palette struct {
	colors [4]uint8
}

func (p *Palette) Write(addr uint16, data uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	p.colors[0] = data & 0x3
	p.colors[1] = (data >> 2) & 0x3
	p.colors[2] = (data >> 4) & 0x3
	p.colors[3] = (data >> 6) & 0x3
}

func (p *Palette) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	return p.colors[0] | (p.colors[1] << 2) | (p.colors[2] << 4) | (p.colors[3] << 6)
}

func (p *Palette) Reset() {
	p.colors = [4]uint8{0, 0, 0, 0}
}

// 4 colors, 2 bits each. Only 0-3 are valid (0 = white, 1 = light gray, 2 = dark gray, 3 = black)
func (p *Palette) Set(values [4]uint8) {
	for i, v := range values {
		p.colors[i] = v
		if p.colors[i] > 3 {
			panic("Invalid color value")
		}
	}
}

func (p *Palette) Match(val uint8) uint8 {
	if val > 3 {
		panic("Invalid color value")
	}

	return p.colors[val]
}
