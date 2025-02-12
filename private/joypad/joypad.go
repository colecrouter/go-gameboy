package joypad

// type JoyPadState struct {
// 	buffer uint8
// }

// func (j *JoyPadState) Read() uint8 {
// 	return j.buffer
// }

// func (j *JoyPadState) Write(value uint8) {
// 	j.buffer = value
// }

// type Button uint

// const (
// 	A Button = iota
// 	B
// 	Up
// 	Down
// 	Left
// 	Right
// 	Select
// 	Start
// )

// func (j *JoyPadState) UpdateButton(button Button, pressed bool) {
// 	var val uint8
// 	if !pressed {
// 		val = 1
// 	}

// 	switch button {
// 	case A:
// 		fallthrough
// 	case Right:
// 		j.buffer |= val
// 	case B:
// 		fallthrough
// 	case Left:
// 		j.buffer |= val << 1
// 	case Select:
// 		fallthrough
// 	case Up:
// 		j.buffer |= val << 2
// 	case Start:
// 		fallthrough
// 	case Down:
// 		j.buffer |= val << 3
// 	}

// 	switch button {
// 	case A:
// 		fallthrough
// 	case B:
// 		fallthrough
// 	case Select:
// 		fallthrough
// 	case Start:
// 		j.buffer |= 0b00100000
// 	case Up:
// 		fallthrough
// 	case Down:
// 		fallthrough
// 	case Left:
// 		fallthrough
// 	case Right:
// 		j.buffer |= 0b00010000
// 	}
