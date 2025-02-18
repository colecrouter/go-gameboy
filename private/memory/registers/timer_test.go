package registers

import "testing"

func TestOverflowBehavior_NoWrite(t *testing.T) {
	// Cycle A: natural overflow.
	intr := &Interrupt{}
	timer := NewTimer(intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Divider = 0x2E
	timer.Counter = 0xFE
	timer.Modulo = 50

	// Loop so that t.prevBit is set to 1.
	for range 6 {
		timer.MClock()
	}

	// Expected: overflow triggers, TIMA becomes 0, and interrupt flag is not set.
	if timer.Counter != 0 || intr.Timer {
		t.Errorf("cycle A: expected counter 0 and no interrupt, got counter %d, IF %v", timer.Counter, intr.Timer)
	}

	// Cycle B: overflow should trigger again, TIMA should reload and increment.
	timer.MClock()

	// Expected: TIMA reloads and does not increment.
	if timer.Counter != 50 {
		t.Errorf("expected counter 50, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Errorf("expected timer interrupt flag set")
	}

	// Cycle C: overflow should trigger again, TIMA should reload and increment.
	timer.MClock()
}

func TestWriteTIMA_CycleA(t *testing.T) {
	// Cycle A: Writing to TIMA cancels the overflow.
	intr := &Interrupt{}
	timer := NewTimer(intr)
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

func TestWriteTMA_CycleB(t *testing.T) {
	// Cycle A: natural overflow occurs.
	intr := &Interrupt{}
	timer := NewTimer(intr)
	timer.Control.Enabled = true
	timer.Control.Speed = M4
	timer.Divider = 0x2C
	timer.Counter = 0xFF
	timer.Modulo = 80

	// Loop so that t.prevBit is set to 1.
	for range 4 {
		timer.MClock()
	}

	// Cycle A: expect overflow to trigger delayed reload:
	// TIMA becomes 0 and interrupt flag is not yet set.
	if timer.Counter != 0 || intr.Timer {
		t.Errorf("cycle A: expected counter 0 and no interrupt, got counter %d, IF %v", timer.Counter, intr.Timer)
	}

	// Write new TIMA
	timer.Write(0x1, 55)

	// Cycle B: first clock under delay - TIMA should be reloaded from TMA.
	timer.MClock()
	if timer.Counter != 80 {
		t.Errorf("cycle B: expected counter 80, got %d", timer.Counter)
	}
	if !intr.Timer {
		t.Errorf("cycle B: expected timer interrupt flag set")
	}

	// Cycle C: second clock under delay - TIMA should increment normally.
	for range 4 {
		timer.MClock()
	}
	if timer.Counter != 81 {
		t.Errorf("cycle C: expected counter 81, got %d", timer.Counter)
	}
}
