package io

import "testing"

// Test for overflow behavior without any intervening write.
func TestOverflowBehavior_NoWrite(t *testing.T) {
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	// Setup: enable timer with fastest period (M4) and overflow condition.
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Counter = 0xFF
	timer.Modulo = 0x23

	// Cycle A: Trigger overflow.
	for range 4 {
		timer.MRisingEdge()
		timer.MFallingEdge()
	}

	// Expect TIMA to be 0 (pending reload).
	if timer.Counter != 0 {
		t.Errorf("expected counter 0 immediately after overflow, got %d", timer.Counter)
	}

	// Cycle B: Reload from TMA and set the interrupt flag.
	timer.MRisingEdge()
	timer.MFallingEdge()
	if timer.Counter != 0x23 {
		t.Errorf("expected counter reloaded to 0x23, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Error("expected timer interrupt flag to be set")
	}
}

func TestWriteTIMA_CycleA(t *testing.T) {
	// Cycle A: Writing to TIMA cancels the overflow.
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Divider = 0x2C
	timer.Counter = 0xFF
	timer.Modulo = 80

	for range 2 {
		timer.MRisingEdge()
		timer.MFallingEdge()
	}

	// Write to TIMA before overflow can trigger.
	timer.Write(0x1, 55)

	for range 2 {
		timer.MRisingEdge()
		timer.MFallingEdge()
	}

	// Expected: no overflow reload; 55 is incremented normally.
	if timer.Counter != 56 {
		t.Errorf("expected counter 56, got %d", timer.Counter)
	}
	if intr.Timer {
		t.Errorf("expected timer interrupt flag not set")
	}
}

// TIMA should be reloaded from TMA despite the write.
func TestWriteTIMA_CycleB(t *testing.T) {
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Counter = 0xFF
	timer.Modulo = 0x23

	// Cycle A: Trigger overflow.
	for range 4 {
		timer.MRisingEdge()
		timer.MFallingEdge()
	}
	// Now the write occurs after we've moved into "cycle B."
	timer.Write(0x1, 99)

	// Complete cycle B, MFallingEdge.
	timer.MRisingEdge()

	if timer.Counter != 0x23 {
		t.Errorf("expected counter reloaded to 0x23, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Error("expected timer interrupt flag to be set")
	}
}

// The written value should remain, and the timer interrupt flag must not be set.
func TestWriteTIMA_DuringCycleA_BypassesOverflow(t *testing.T) {
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	// Setup: enable timer with fastest period (M4) and overflow condition.
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Counter = 0xFF
	timer.Modulo = 42 // arbitrary TMA value

	// Cycle A: trigger overflow.
	timer.MRisingEdge()
	timer.MFallingEdge()
	// Write to TIMA during cycle A.
	timer.Write(0x1, 77)
	// Cycle B: complete reload cycle.
	timer.MRisingEdge()
	timer.MFallingEdge()
	if timer.Counter != 77 {
		t.Errorf("expected counter to remain 77, got %d", timer.Counter)
	}
	if intr.Timer {
		t.Error("expected timer interrupt flag not set")
	}
}

func TestTimerIntervals(t *testing.T) {
	// For each mode, we expect that after (freq - (initialDivider mod freq)) cycles,
	// TIMA should increase by one.
	tests := []struct {
		name           string
		speed          Increment
		initialDivider uint16
	}{
		// For M256: freq = 256. With initialDivider = 128, it takes 256 - 128 = 128 cycles.
		{"M256", M256, 128},
		// For M4: freq = 4. With initialDivider = 2, it takes 4 - 2 = 2 cycles.
		{"M4", M4, 2},
		// For M16: freq = 16. With initialDivider = 8, it takes 16 - 8 = 8 cycles.
		{"M16", M16, 8},
		// For M64: freq = 64. With initialDivider = 32, it takes 64 - 32 = 32 cycles.
		{"M64", M64, 32},
	}

	for _, tc := range tests {
		intr := &Interrupt{}
		timer := NewTimer(nil, intr)
		timer.Control.Enabled = true
		timer.Control.Speed = tc.speed
		timer.Counter = 0
		timer.Divider = tc.initialDivider

		// Determine the frequency target.
		var freq uint16
		switch tc.speed {
		case M256:
			freq = 256
		case M4:
			freq = 4
		case M16:
			freq = 16
		case M64:
			freq = 64
		}

		// Calculate cycles needed to hit the next multiple of freq.
		cyclesNeeded := int(freq - (tc.initialDivider % freq))
		for i := 0; i < cyclesNeeded; i++ {
			timer.MRisingEdge()
			timer.MFallingEdge()
		}

		if timer.Counter != 1 {
			t.Errorf("%s: expected counter 1 after %d cycles, got %d", tc.name, cyclesNeeded, timer.Counter)
		}
	}
}
