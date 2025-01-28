package gamepak

type CartridgeType uint8

const ROM_ONLY CartridgeType = 0x00
const MBC1 CartridgeType = 0x01
const MBC1_RAM CartridgeType = 0x02
const MBC1_RAM_BATTERY CartridgeType = 0x03
const MBC2 CartridgeType = 0x05
const MBC2_BATTERY CartridgeType = 0x06
const ROM_RAM CartridgeType = 0x08
const ROM_RAM_BATTERY CartridgeType = 0x09
const MMM01 CartridgeType = 0x0B
const MMM01_RAM CartridgeType = 0x0C
const MMM01_RAM_BATTERY CartridgeType = 0x0D
const MBC3_TIMER_BATTERY CartridgeType = 0x0F
const MBC3_TIMER_RAM_BATTERY CartridgeType = 0x10
const MBC3 CartridgeType = 0x11
const MBC3_RAM CartridgeType = 0x12
const MBC3_RAM_BATTERY CartridgeType = 0x13
const MBC5 CartridgeType = 0x19
const MBC5_RAM CartridgeType = 0x1A
const MBC5_RAM_BATTERY CartridgeType = 0x1B
const MBC5_RUMBLE CartridgeType = 0x1C
const MBC5_RUMBLE_RAM CartridgeType = 0x1D
const MBC5_RUMBLE_RAM_BATTERY CartridgeType = 0x1E
const MBC6 CartridgeType = 0x20
const MBC7_SENSOR_RUMBLE_RAM_BATTERY CartridgeType = 0x22
const POCKET_CAMERA CartridgeType = 0xFC
const BANDAI_TAMA5 CartridgeType = 0xFD
const HUC3 CartridgeType = 0xFE
const HUC1_RAM_BATTERY CartridgeType = 0xFF

func (gp *GamePak) CartridgeType() CartridgeType {
	return CartridgeType(gp.buffer[0x147])
}
