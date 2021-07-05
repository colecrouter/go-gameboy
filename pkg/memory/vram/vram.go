package vram

type VRAM struct {
	data [0x2000]uint8 // 8KB
}

func (v *VRAM) Read(addr uint16) uint8 {
	return v.data[addr]
}

func (v *VRAM) Write(addr uint16, data uint8) {
	v.data[addr] = data
}
