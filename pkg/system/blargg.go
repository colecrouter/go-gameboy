package system

// TestSerialDevice collects output from serial transfers.
type TestSerialDevice struct {
	output []uint8
}

// Transfer stores the received byte and echoes it back.
func (d *TestSerialDevice) Transfer(b byte) byte {
	d.output = append(d.output, b)
	return b
}

// SetupBlarggTestSystem builds a GameBoy system with the test serial device connected.
func SetupBlarggTestSystem() (*GameBoy, *TestSerialDevice) {
	gb := NewGameBoy()
	// Assume GameBoy has a Serial field of appropriate type.
	testDevice := &TestSerialDevice{}
	gb.ConnectSerialDevice(testDevice)
	return gb, testDevice
}
