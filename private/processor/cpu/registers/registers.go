package registers

type Registers struct {
	// 8-bit registers
	A, B, C, D, E, H, L uint8

	// Stack Pointer
	SP uint16

	// Program Counter
	PC uint16
}
