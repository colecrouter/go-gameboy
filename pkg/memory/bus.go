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

// PrintMemory prints the contents of the memory bus within the specified address range.
func (b *Bus) PrintMemory(start, end uint16) {
	if start > end {
		panic("Start address must be less than or equal to end address")
	}

	// Add column headers
	fmt.Printf("\n        ") // Adjusted spacing from 7 to 8 spaces for alignment
	for i := 0; i < 16; i++ {
		fmt.Printf("%02X ", i)
	}
	fmt.Println()

	for addr := start; addr <= end; addr++ {
		data := b.Read(addr)
		if addr%16 == 0 {
			fmt.Printf("\n0x%04X: ", addr)
		}
		if data == 0x00 {
			fmt.Printf("\033[90m%02X \033[0m", data) // Gray color for 0x00
		} else {
			fmt.Printf("%02X ", data)
		}
	}

	fmt.Println()
}
