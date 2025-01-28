package memory

type memoryMapping struct {
	Start  uint16
	End    uint16
	Device *Device
}

type Bus struct {
	mapping []memoryMapping
}

func (b *Bus) AddDevice(start uint16, end uint16, device Device) {
	b.mapping = append(b.mapping, memoryMapping{Start: start, End: end, Device: &device})
}

func (b *Bus) Read(addr uint16) byte {
	for _, mapping := range b.mapping {
		println("Checking", mapping.Start, mapping.End)
		if addr >= mapping.Start && addr <= mapping.End {
			println("Reading from device", mapping.Device, "at address", addr)
			// Get adjusted address
			addr -= mapping.Start

			return (*mapping.Device).Read(addr)
		}
	}
	panic("No device found for address")
}

func (b *Bus) Write(addr uint16, data byte) {
	for _, mapping := range b.mapping {
		if addr >= mapping.Start && addr <= mapping.End {
			// Get adjusted address
			addr -= mapping.Start

			(*mapping.Device).Write(addr, data)
			return
		}
	}
}
