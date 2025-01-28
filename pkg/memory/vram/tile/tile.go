package tile

type Tile struct {
	Bytes [16]uint8
}

func (t *Tile) ReadPixel(row, col uint8) uint8 {
	line := []uint8{t.Bytes[row*2], t.Bytes[row*2+1]}

	msb := (line[0]) >> (7 - col) & 1
	lsb := (line[1]) >> (7 - col) & 1

	// Combine them to get the color id (0-3)
	return msb | (lsb << 1)
}

func FromBytes(bytes [16]uint8) Tile {
	return Tile{Bytes: bytes}
}
