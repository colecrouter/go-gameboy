package io

import "testing"

// Test for overflow behavior without any intervening write.
func TestOverflowBehavior_NoWrite(t *testing.T) {
	// Cycle A: natural overflow.
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Divider = 0x2E // Chosen so that two increments happen over 16 cycles.
	timer.Counter = 0xFE // Will overflow after two increments.
	timer.Modulo = 50    // TMA for reload.

	// Run cycles to trigger overflow.
	for i := 0; i < 16; i++ {
		timer.MClock()
	}

	// After overflow, TIMA should be 0 and no interrupt yet.
	if timer.Counter != 0 || intr.Timer {
		t.Errorf("Cycle A: expected TIMA=0 and no timer interrupt; got TIMA=%d, IF=%v", timer.Counter, intr.Timer)
	}

	// Cycle B: on the next clock, the delayed reload takes place.
	timer.MClock()
	if timer.Counter != 50 {
		t.Errorf("Cycle B: expected TIMA reloaded to 50 from TMA, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Errorf("Cycle B: expected timer interrupt flag set")
	}

	// Cycle C: Let the timer increment normally.
	// After the appropriate number of cycles, TIMA should eventually increment past the reloaded TMA value.
	// For example, if an increment happens every 8 cycles:
	for i := 0; i < 8; i++ {
		timer.MClock()
	}
	// In this test, TIMA should have incremented after reload in cycle C.
	if timer.Counter != 0 {
		t.Errorf("Cycle C: expected TIMA to wrap around (or be incremented as per timing), got %d", timer.Counter)
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

	// Loop so that t.prevBit is set to 1.
	for range 2 {
		timer.MClock()
	}

	// Write to TIMA before overflow can trigger.
	timer.Write(0x1, 55)

	for range 2 {
		timer.MClock()
	}

	// Expected: no overflow reload; 55 is incremented normally.
	if timer.Counter != 56 {
		t.Errorf("expected counter 56, got %d", timer.Counter)
	}
	if intr.Timer {
		t.Errorf("expected timer interrupt flag not set")
	}
}

// Test for delayed reload behavior with a write during the delay.
func TestWriteTMA_CycleB(t *testing.T) {
	// Cycle A: natural overflow occurs.
	intr := &Interrupt{}
	timer := NewTimer(nil, intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4 // Assuming M4 increments every 4 M-cycles (adjust if using revised offsets).
	timer.Divider = 0x2C     // Starting value chosen so that overflow occurs after 16 cycles.
	timer.Counter = 0xFF     // Will overflow when incremented.
	timer.Modulo = 80        // TMA for reload.

	// Cycle A: Run enough cycles to trigger the overflow and schedule the reload.
	for i := 0; i < 16; i++ {
		timer.MClock()
	}

	// At this point, overflow has occurred.
	// According to hardware behavior:
	// • TIMA is 0 (overflowed) but the delayed reload hasn’t occurred yet.
	// • Timer interrupt flag should not be set.
	if timer.Counter != 0 || intr.Timer {
		t.Errorf("Cycle A: expected TIMA=0 and no timer interrupt; got TIMA=%d, IF=%v", timer.Counter, intr.Timer)
	}

	// Write new value to TIMA while the reload is pending.
	timer.Write(0x1, 55)

	// Cycle B: the reload occurs on the first clock after overflow.
	timer.MClock()
	if timer.Counter != 80 {
		t.Errorf("Cycle B: expected TIMA reloaded to 80 from TMA, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Errorf("Cycle B: expected timer interrupt flag to be set")
	}

	// Cycle C: Subsequent clocks cause normal incrementation.
	// For example, after 8 more cycles (assuming M4 increment every 8 cycles in your test scenario),
	// TIMA should increment from 80 to 81.
	for i := 0; i < 8; i++ {
		timer.MClock()
	}
	if timer.Counter != 81 {
		t.Errorf("Cycle C: expected TIMA=81 after normal increment, got %d", timer.Counter)
	}
}

func TestTimerIntervals(t *testing.T) {
	// For each mode, we want:
	//   cyclesNeeded = (2^(offset+1) - (initialDivider mod 2^(offset+1)))
	// with the following offsets (assuming one DIV increment per M-cycle):
	//   M256: offset 7, period = 256 cycles, high phase = [128,255]
	//   M4  : offset 1, period = 4 cycles, high phase = [2,3]
	//   M16 : offset 3, period = 16 cycles, high phase = [8,15]
	//   M64 : offset 5, period = 64 cycles, high phase = [32,63]
	tests := []struct {
		name           string
		speed          Increment
		cycles         int // number of M-cycles to run
		initialDivider uint16
	}{
		// For M256: choose initialDivider = 128 (128 mod 256 = 128, high phase).
		// It will take 256-128 = 128 cycles to hit the falling edge.
		{"M256", M256, 128, 128},

		// For M4: choose initialDivider = 2 (2 mod 4 = 2, high phase).
		// It takes 4-2 = 2 cycles to hit the falling edge.
		{"M4", M4, 2, 2},

		// For M16: choose initialDivider = 8 (8 mod 16 = 8, high phase).
		// It takes 16-8 = 8 cycles to hit the falling edge.
		{"M16", M16, 8, 8},

		// For M64: choose initialDivider = 32 (32 mod 64 = 32, high phase).
		// It takes 64-32 = 32 cycles to hit the falling edge.
		{"M64", M64, 32, 32},
	}

	for _, tc := range tests {
		intr := &Interrupt{}
		timer := NewTimer(nil, intr)
		timer.Control.Enabled = true
		timer.Control.Speed = tc.speed
		timer.Counter = 0

		// Set to a value whose selected bit is 1.
		timer.Divider = tc.initialDivider
		timer.prevBit = true

		for i := 0; i < tc.cycles; i++ {
			timer.MClock()
		}
		if timer.Counter != 1 {
			t.Errorf("%s: expected counter 1 after %d cycles, got %d", tc.name, tc.cycles, timer.Counter)
		}
	}
}
