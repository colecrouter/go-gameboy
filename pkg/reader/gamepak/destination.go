package gamepak

type DestinationCode uint8

const JAPANESE DestinationCode = 0x00
const NON_JAPANESE DestinationCode = 0x01

func (gp *GamePak) DestinationCode() DestinationCode {
	return DestinationCode(gp.buffer[0x014A])
}
