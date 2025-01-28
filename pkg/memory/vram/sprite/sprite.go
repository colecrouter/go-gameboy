package sprite

type Sprite struct {
	buffer [4]byte
	// Y          byte
	// X          byte
	// Tile       byte
	// Priority   bool
	// FlipY      bool
	// FlipX      bool
	// DmgPalette bool
	// Bank       bool
	// CGBPalette byte
}

func (s *Sprite) Y() uint8 {
	return s.buffer[0]
}

func (s *Sprite) X() uint8 {
	return s.buffer[1]
}

func (s *Sprite) Tile() uint8 {
	return s.buffer[2]
}

func (s *Sprite) Priority() bool {
	return s.buffer[3]&0b1000_0000 != 0
}

func (s *Sprite) FlipY() bool {
	return s.buffer[3]&0b0100_0000 != 0
}

func (s *Sprite) FlipX() bool {
	return s.buffer[3]&0b0010_0000 != 0
}

func (s *Sprite) DmgPalette() bool {
	return s.buffer[3]&0b0001_0000 != 0
}

func (s *Sprite) Bank() bool {
	return s.buffer[3]&0b0000_1000 != 0
}

func (s *Sprite) CGBPalette() uint8 {
	return s.buffer[3] & 0b0000_0111
}

func (s *Sprite) ReadPixel(y, x uint8) uint8 {
	if s.FlipX() {
		x = 7 - x
	}
	// TODO
	if s.FlipY() {
		y = 7 - y
	}

	return s.Tile() & (1 << (7 - x))
}

func NewSprite(buffer [4]byte) *Sprite {
	return &Sprite{buffer}
}
