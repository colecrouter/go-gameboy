package io

import "github.com/colecrouter/gameboy-go/private/system"

// Added new fields for overflow management.
type Timer struct {
	Divider uint16       // FF04 - DIV
	Counter uint8        // FF05 - TIMA
	Modulo  uint8        // FF06 - TMA
	Control TimerControl // FF07 - TAC

	interruptFlags *Interrupt
	initialized    bool

	// Two clock channels from the broadcaster:
	clockRising  <-chan struct{}
	clockFalling <-chan struct{}

	// State for overflow management.
	pendingOverflow   bool
	enableWriteCancel bool
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
		// // Writing to TIMA.
		// // If a pending overflow is active, only cancel it if the cancellation window is open.
		// if t.pendingOverflow {
		// 	if t.enableWriteCancel {
		// 		// In cycle A: allow cancellation and update TIMA normally.
		// 		t.pendingOverflow = false
		// 	} else {
		// 		// In cycle B: ignore the write so that the reload from TMA occurs.
		// 		return
		// 	}
		// }
		t.Counter = value
	case 0x2:
		t.Modulo = value
	case 0x3:
		t.Control.Write(0, value)
	default:
		panic("Invalid address")
	}
}

// MRisingEdge should be called on the m-cycle rising edge.
// It closes the cancellation window and reloads TIMA if overflow is still pending.
func (t *Timer) MRisingEdge() {
	// Close cancellation window.
	t.enableWriteCancel = false
	if t.pendingOverflow {
		t.Counter = t.Modulo
		t.interruptFlags.Timer = true
		t.pendingOverflow = false
	}
}

// MFallingEdge should be called on the m-cycle falling edge.
// It increments the divider and updates TIMA, potentially triggering an overflow.
func (t *Timer) MFallingEdge() {
	oldDiv := t.Divider
	t.Divider++

	// If the timer is disabled, nothing more to do.
	if !t.Control.Enabled {
		return
	}

	var freq uint16
	switch t.Control.Speed {
	case M256:
		freq = 256
	case M4:
		freq = 4
	case M16:
		freq = 16
	case M64:
		freq = 64
	}

	increase := (t.Divider / freq) - (oldDiv / freq)
	if increase > 0 {
		newVal := uint16(t.Counter) + increase
		if newVal > 0xFF {
			// Overflow: set TIMA to 0 and open cancellation window.
			t.Counter = 0
			t.pendingOverflow = true
			t.enableWriteCancel = true
		} else {
			t.Counter = uint8(newVal)
		}
	}
}

// Run listens for rising and falling edge ticks and calls the proper functions.
func (t *Timer) Run(close <-chan struct{}) {
	for {
		select {
		case <-close:
			return
		case <-t.clockRising:
			t.MRisingEdge()
		case <-t.clockFalling:
			t.MFallingEdge()
		}
	}
}

func NewTimer(broadcaster *system.Broadcaster, interrupt *Interrupt) *Timer {
	timer := &Timer{initialized: true}
	timer.interruptFlags = interrupt

	if broadcaster != nil {
		// Subscribe separately to the m-cycle rising and falling edges.
		timer.clockRising = broadcaster.Subscribe(system.MRisingEdge)
		timer.clockFalling = broadcaster.Subscribe(system.MFallingEdge)
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
