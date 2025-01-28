package display

type Display interface {
	Clock()
	DrawScanline(line uint8, pixels []uint8)
}
