package gamepak

// RomSize returns the size of the ROM in bytes
func (gp *GamePak) RomSize() uint {
	val := uint(gp.buffer[0x148])
	return (32 * 1024) * (1 << val)
}

// BankCount returns the number of ROM banks
func (gp *GamePak) BankCount() uint {
	val := uint(gp.buffer[0x148])
	return 1 << val
}

// RamSize returns the size of the RAM in bytes
func (gp *GamePak) RamSize() uint {
	val := uint(gp.buffer[0x149])
	switch val {
	case 1:
		return 2 * 1024
	case 2:
		return 8 * 1024
	case 3:
		return 32 * 1024
	case 4:
		return 128 * 1024
	case 5:
		return 64 * 1024
	default:
		return 0
	}
}

// RamBankCount returns the number of RAM banks
func (gp *GamePak) RamBankCount() uint {
	val := uint(gp.buffer[0x149])
	switch val {
	case 2:
		return 1
	case 3:
		return 4
	case 4:
		return 16
	case 5:
		return 8
	default:
		return 0
	}
}
