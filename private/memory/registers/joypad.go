package registers

type JoyPad struct {
	buttons []bool

	// Used for reading
	initialized     bool
	selectButtons   bool
	selectDirection bool
	interrupt       *Interrupt
}

type Button uint8

const (
	Button_A Button = iota
	Button_B
	Button_Select
	Button_Start
	Button_Right
	Button_Left
	Button_Up
	Button_Down
)

// Button states are a bit weird
// 0 = pressed
// 1 = not pressed
// This also applies to the select bits

func (j *JoyPad) Read(addr uint16) uint8 {
	if !j.initialized {
		panic("JoyPad not initialized")
	}

	if addr != 0 {
		panic("Invalid address")
	}

	var val uint8 = 0b11111111 // Set bits 4 and 5, all buttons released

	// Correct upper nibble

	// Unset the bit if the mode is selected
	if j.selectButtons { // Unset bit 5
		val &= 0b11011111
	}
	if j.selectDirection { // Unset bit 4
		val &= 0b11101111
	}

	// Apply lower nibble

	// Not sure if this is correct behavior
	// Docs only mention when neither mode is selected, not when both are
	if j.selectButtons == j.selectDirection {
		return val
	}

	if j.selectButtons {
		if j.buttons[Button_A] {
			val &= 0b11111110
		}
		if j.buttons[Button_B] {
			val &= 0b11111101
		}
		if j.buttons[Button_Select] {
			val &= 0b11111011
		}
		if j.buttons[Button_Start] {
			val &= 0b11110111
		}
	} else if j.selectDirection {
		if j.buttons[Button_Right] {
			val &= 0b11111110
		}
		if j.buttons[Button_Left] {
			val &= 0b11111101
		}
		if j.buttons[Button_Up] {
			val &= 0b11111011
		}
		if j.buttons[Button_Down] {
			val &= 0b11110111
		}
	}

	return val
}

func (j *JoyPad) Write(addr uint16, value uint8) {
	if !j.initialized {
		panic("JoyPad not initialized")
	}

	if addr != 0 {
		panic("Invalid address")
	}

	j.selectButtons = value&0b00100000 == 0
	j.selectDirection = value&0b00010000 == 0
}

func (j *JoyPad) SetButton(button Button, pressed bool) {
	if !j.initialized {
		panic("JoyPad not initialized")
	}

	before := j.buttons[button]
	if before == pressed {
		return
	}

	j.buttons[button] = pressed

	if j.interrupt != nil {
		j.interrupt.Joypad = true
	}
}

func (j *JoyPad) GetButton(button Button) bool {
	if !j.initialized {
		panic("JoyPad not initialized")
	}

	return j.buttons[button]
}

func (j *JoyPad) ResetButtons() {
	for i := range j.buttons {
		j.buttons[i] = false
	}
}

func NewJoyPad(interrupt *Interrupt) *JoyPad {
	return &JoyPad{
		buttons:     make([]bool, 8),
		interrupt:   interrupt,
		initialized: true,
	}
}
