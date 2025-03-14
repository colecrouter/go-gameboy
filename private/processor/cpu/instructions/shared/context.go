package shared

// Context represents intermediate state between micro-operations, within a single instruction.
type Context struct {
	Z uint8
	W uint8
}
