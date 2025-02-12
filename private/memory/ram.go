package memory

// Generic memory device interface
type Memory struct {
	Buffer []byte
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.Buffer[addr]
}

func (m *Memory) Write(addr uint16, data uint8) {
	m.Buffer[addr] = data
}
