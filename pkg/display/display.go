package display

const WIDTH = 160
const HEIGHT = 144

type Display interface {
	Clock()
	DrawScanline(line uint8, pixels []uint8)
}
