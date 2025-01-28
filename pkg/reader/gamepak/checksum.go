package gamepak

func (gp *GamePak) headerChecksum() uint8 {
	return gp.buffer[0x014D]
}

func (gp *GamePak) globalChecksum() uint16 {
	var val uint16
	val = uint16(gp.buffer[0x014E]) << 8
	val |= uint16(gp.buffer[0x014F])

	return val
}

func (gp *GamePak) ComputeHeaderChecksum() bool {
	var sum uint8
	for i := 0x0134; i <= 0x014C; i++ {
		sum = sum - gp.buffer[i] - 1
	}

	return sum == gp.headerChecksum()
}

func (gp *GamePak) ComputeGlobalChecksum() bool {
	var sum uint16
	for i := 0; i < len(gp.buffer); i++ {
		if i == 0x014E || i == 0x014F {
			continue
		}
		sum += uint16(gp.buffer[i])
	}

	return sum == gp.globalChecksum()
}
