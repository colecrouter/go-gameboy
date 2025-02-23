package tile

const TILE_SIZE = 8

type Tile struct {
	initialized bool
	pixels      [TILE_SIZE * TILE_SIZE]uint8
}

// NewTile converts 16 bytes into a Tile.
func NewTile(bytes [16]uint8) *Tile {
	t := &Tile{initialized: true}
	for i := uint8(0); i < TILE_SIZE; i++ {
		for j := uint8(0); j < TILE_SIZE; j++ {
			msb := (bytes[i*2] >> (7 - j)) & 1
			lsb := (bytes[i*2+1] >> (7 - j)) & 1
			t.pixels[i*TILE_SIZE+j] = msb | (lsb << 1)
		}
	}
	return t
}

// Pixels returns the pixel data for the tile.
func (t *Tile) Pixels() []uint8 {
	return t.pixels[:]
}
