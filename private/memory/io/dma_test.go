package io

import (
	"testing"
)

// FakeBus implements minimal Bus behavior for testing.
type FakeBus struct {
	mem map[uint16]uint8
}

func (fb *FakeBus) Read(addr uint16) uint8 {
	if v, ok := fb.mem[addr]; ok {
		return v
	}
	return 0
}

func (fb *FakeBus) Write(addr uint16, value uint8) {
	fb.mem[addr] = value
}

func TestDMA(t *testing.T) {
	// Prepopulate FakeBus memory for DMA source (0x8000-0x80A0).
	fakeBus := &FakeBus{mem: make(map[uint16]uint8)}
	const dmaSize = 0xA0
	const sourceBase = 0x8000
	const destBase = 0xFE00

	for i := 0; i < dmaSize; i++ {
		fakeBus.mem[sourceBase+uint16(i)] = uint8(i)
	}

	// Create Registers instance manually.
	regs := NewRegisters(nil, fakeBus, nil)

	// Trigger DMA transfer; passing 0x80 -> source = 0x80 << 8 = 0x8000.
	regs.Write(0x46, 0x80)

	// Verify that DMA copied dmaSize bytes from source to destination.
	for i := 0; i < dmaSize; i++ {
		expected := fakeBus.mem[sourceBase+uint16(i)]
		actual := fakeBus.mem[destBase+uint16(i)]
		if actual != expected {
			t.Errorf("DMA failed at offset %d: expected 0x%02X, got 0x%02X", i, expected, actual)
		}
	}
}
