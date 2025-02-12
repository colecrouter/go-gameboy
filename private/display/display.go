package display

import "image"

type Config struct {
	Title  string // Title to be displayed
	Width  int    // Width in pixels. Only applies to text displays
	Height int    // Height in pixels. Only applies to text displays
}

type Display interface {
	Clock()
	Config() *Config
}

type ImageDisplay interface {
	Image() image.Image
	Clock()
	Config() *Config
}

type TextDisplay interface {
	Text() []string
	Clock()
	Config() *Config
}
