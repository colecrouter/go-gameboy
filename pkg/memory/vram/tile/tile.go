package tile

const TILE_SIZE = 8

type Tile struct {
	initialized bool
	Bytes       [16]uint8
	Pixels      [TILE_SIZE][TILE_SIZE]uint8
}

func (t *Tile) readPixel(row, col uint8) uint8 {
	msb := (t.Bytes[row*2] >> (7 - col)) & 1
	lsb := (t.Bytes[row*2+1] >> (7 - col)) & 1

	return msb | (lsb << 1)
}

func FromBytes(bytes [16]uint8) *Tile {
	t := &Tile{Bytes: bytes, initialized: true}

	for row := uint8(0); row < TILE_SIZE; row++ {
		for col := uint8(0); col < TILE_SIZE; col++ {
			t.Pixels[row][col] = t.readPixel(row, col)
		}
	}

	return t
}
