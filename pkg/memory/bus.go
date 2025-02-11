package memory

import "fmt"

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
		if addr >= mapping.Start && addr <= mapping.End {
			// fmt.Printf("Reading from device %v at address 0x%X\n", mapping.Device, addr)
			// Get adjusted address
			addr -= mapping.Start

			return (*mapping.Device).Read(addr)
		}
	}
	panic("No device found for address")
}

func (b *Bus) Read16(addr uint16) (high, low uint8) {
	low = b.Read(addr)
	high = b.Read(addr + 1)

	return high, low
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

func (b *Bus) Write16(addr uint16, data uint16) {
	// Extract low and high bytes
	low := uint8(data & 0xFF)
	high := uint8((data >> 8) & 0xFF)

	b.Write(addr, low)
	b.Write(addr+1, high)
}

// Print prints the contents of the memory bus within the specified address range.
func Print(d Device, start, end uint16) {
	if start > end {
		panic("Start address must be less than or equal to end address")
	}

	// Add column headers
	fmt.Printf("\nREL    FIXED   ") // Header for relative and fixed offsets
	for i := 0; i < 16; i++ {
		fmt.Printf("%02X ", i)
	}
	fmt.Println()

	for addr := start; addr <= end; addr++ {
		var data uint8
		defer func() {
			if r := recover(); r != nil {
				data = 0x00
			}
		}()
		data = d.Read(addr)

		// Print new row with relative and fixed offsets
		if (addr-start)%16 == 0 {
			rel := addr - start
			fmt.Printf("\033[90m0d%04d\033[0m 0x%04X: ", rel, addr)
		}

		// Apply color for 0x00, normal print otherwise.
		if data == 0x00 {
			fmt.Printf("\033[90m%02X \033[0m", data)
		} else {
			fmt.Printf("%02X ", data)
		}

		// Print new line after 16 bytes
		if (addr-start+1)%16 == 0 {
			fmt.Println()
		}
	}

	fmt.Println()
}
