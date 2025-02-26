package io

import (
	"testing"
)

func TestJoypad_Read(t *testing.T) {
	// Initialization, values should be 0xCF
	p := NewJoyPad(nil)
	if got := p.Read(0); got != 0xCF {
		t.Errorf("Initialization: got 0b%08b, want 0b11001111", got)
	}

	// Both modes selected (0x00): expecting lower nibble forced to 0xF, so 0x30|0xF = 0x3F.
	p.SetButton(Button_A, true)
	p.SetButton(Button_Up, true)
	p.Write(0, 0x30) // Both selected, lower nibble ignored.
	if got := p.Read(0); got != 0xFF {
		t.Errorf("Both selected: got 0b%08b, want 0b11111111", got)
	}

	// Buttons-only selected (Write 0x10): expecting 00011110 (button selected, A is pressed).
	p = NewJoyPad(nil)
	p.SetButton(Button_A, true)
	p.Write(0, 0x10)
	if got := p.Read(0); got != 0xDE {
		t.Errorf("Buttons only selected: got 0b%08b, want 0b11011110", got)
	}

	// D-pad-only selected (Write 0x20): expecting 00101101 (direction selected, Left is pressed).
	p = NewJoyPad(nil)
	p.SetButton(Button_Left, true)
	p.Write(0, 0x20)
	if got := p.Read(0); got != 0xED {
		t.Errorf("D-pad only selected: got 0b%08b, want 0b11101101", got)
	}
}
