package io

import "github.com/colecrouter/gameboy-go/private/system"

type Timer struct {
	Divider uint16       // FF04 - DIV
	Counter uint8        // FF05 - TIMA
	Modulo  uint8        // FF06 - TMA
	Control TimerControl // FF07 - TAC

	interruptFlags *Interrupt
	initialized    bool

	// Synchronization state.
	prevBit       bool
	overflow      bool
	pendingReload bool

	clock <-chan struct{}
}

func (t *Timer) Read(addr uint16) uint8 {
	if !t.initialized {
		panic("Timer not initialized")
	}

	switch addr {
	case 0x0:
		return uint8(t.Divider >> 8)
	case 0x1:
		return t.Counter
	case 0x2:
		return t.Modulo
	case 0x3:
		return t.Control.Read(0)
	default:
		panic("Invalid address")
	}
}

func (t *Timer) Write(addr uint16, value uint8) {
	if !t.initialized {
		panic("Timer not initialized")
	}

	switch addr {
	case 0x0:
		// DIV is reset to 0 when written to.
		t.Divider = 0
	case 0x1:
		// If an overflow is pending, ignore the write so that the reload occurs.
		if t.overflow {
			return
		}
		t.Counter = value
	case 0x2:
		t.Modulo = value
	case 0x3:
		// If TAC speed is changed, reset the falling edge detection.
		if t.Control.Speed != Increment(value&0x3) {
			t.prevBit = false
		}

		t.Control.Write(0, value)
	default:
		panic("Invalid address")
	}
}

func (t *Timer) MClock() {
	// If a reload is pending, do it and then immediately return.
	if t.pendingReload {
		t.Counter = t.Modulo
		t.interruptFlags.Timer = true
		t.pendingReload = false
		return
	}

	// Always increment DIV
	t.Divider++

	if !t.Control.Enabled {
		return
	}

	// Determine the bit offset (using updated offsets when counting per M-cycle)
	var offset uint16
	switch t.Control.Speed {
	case M256:
		offset = 7
	case M4:
		offset = 1
	case M16:
		offset = 3
	case M64:
		offset = 5
	}

	// Falling edge detection
	old := t.Counter
	bit := (t.Divider >> offset) & 1
	if t.prevBit && bit == 0 {
		t.Counter++
	}
	t.prevBit = bit != 0

	// Check for overflow: if TIMA wraps from 0xFF to 0
	if t.Counter == 0 && old == 0xFF {
		// Mark that an overflow occurred. Do not apply TMA yet.
		t.pendingReload = true
	}
}

func (t *Timer) Run(close <-chan struct{}) {
	if !t.initialized {
		panic("CPU not initialized")
	}

	for {
		select {
		case <-close:
			return
		case <-t.clock:
			t.MClock()
		}
	}
}

func NewTimer(broadcaster *system.Broadcaster, interrupt *Interrupt) *Timer {
	timer := &Timer{initialized: true}
	timer.interruptFlags = interrupt

	if broadcaster != nil {
		timer.clock = broadcaster.SubscribeM()
	}

	return timer
}

type Increment uint8

const (
	M256 Increment = 0 // 256 M-cycles, 1024 T-cycles
	M4   Increment = 1 // 4 M-cycles, 16 T-cycles
	M16  Increment = 2 // 16 M-cycles, 64 T-cycles
	M64  Increment = 3 // 64 M-cycles, 256 T-cycles
)

type TimerControl struct {
	Speed   Increment
	Enabled bool
}

func (t *TimerControl) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	val := uint8(0)
	if t.Enabled {
		val |= 1 << 2
	}
	val |= uint8(t.Speed)
	return val
}

func (t *TimerControl) Write(addr uint16, value uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	// Only the lower 3 bits are writable.
	t.Enabled = value&(1<<2) != 0
	t.Speed = Increment(value & 0x3)
}
