package io

import (
	"testing"
)

func TestJoypad_Read(t *testing.T) {
	// Initialization, values should be 0xCF
	p := NewJoyPad(nil)
	if got := p.Read(0); (got & 0x3F) != 0x0F { // expected: 0xCF & 0x3F = 0x0F
		t.Errorf("Initialization: got 0b%08b (lower 6 bits: 0b%06b), want 0b00001111", got, got&0x3F)
	}

	// Both modes selected (0x00): expecting lower nibble forced to 0xF, so 0x30|0xF = 0x3F.
	p.SetButton(Button_A, true)
	p.SetButton(Button_Up, true)
	p.Write(0, 0x30)                            // Both selected, lower nibble ignored.
	if got := p.Read(0); (got & 0x3F) != 0x3F { // expected: 0xFF & 0x3F = 0x3F
		t.Errorf("Both selected: got 0b%08b (lower 6 bits: 0b%06b), want 0b00111111", got, got&0x3F)
	}

	// Buttons-only selected (Write 0x10): expecting 00011110 (button selected, A is pressed).
	p = NewJoyPad(nil)
	p.SetButton(Button_A, true)
	p.Write(0, 0x10)
	if got := p.Read(0); (got & 0x3F) != 0x1E { // expected: 0xDE & 0x3F = 0x1E
		t.Errorf("Buttons only selected: got 0b%08b (lower 6 bits: 0b%06b), want 0b00011110", got, got&0x3F)
	}

	// D-pad-only selected (Write 0x20): expecting 00101101 (direction selected, Left is pressed).
	p = NewJoyPad(nil)
	p.SetButton(Button_Left, true)
	p.Write(0, 0x20)
	if got := p.Read(0); (got & 0x3F) != 0x2D { // expected: 0xED & 0x3F = 0x2D
		t.Errorf("D-pad only selected: got 0b%08b (lower 6 bits: 0b%06b), want 0b00101101", got, got&0x3F)
	}
}
