package tilemap

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/private/display"
	font "github.com/colecrouter/gameboy-go/private/font/4x4" // import the 4x4 font package
	"github.com/colecrouter/gameboy-go/private/memory/vram"
)

const (
	TILE_SIZE        = 8
	COLS             = 32
	ROWS             = 32
	CELL_SIZE        = TILE_SIZE + 1
	FULL_GRID_WIDTH  = COLS*CELL_SIZE + 1
	FULL_GRID_HEIGHT = ROWS*CELL_SIZE + 1
)

type TilemapMenu struct {
	vram    *vram.VRAM
	img     *image.Paletted
	palette *color.Palette
	config  display.Config
}

func NewTilemapDebug(v *vram.VRAM, p *color.Palette) *TilemapMenu {
	return &TilemapMenu{
		vram:    v,
		palette: p,
		config:  display.Config{Title: "Tilemap Viewer"},
	}
}

func (t *TilemapMenu) Image() image.Image {
	return t.img
}

func (t *TilemapMenu) Clock() {
	t.img = image.NewPaletted(image.Rect(0, 0, FULL_GRID_WIDTH, FULL_GRID_HEIGHT), *t.palette)

	// Render tilemap cells by drawing 2 hex digits (high nibble and low nibble) using a 4x4 font.
	// Each cell is 8x8; we center the 4x4 sprites vertically (with a 2-pixel vertical offset),
	// and draw the hi digit at the left half and the lo digit at the right half.
	for tileY := 0; tileY < ROWS; tileY++ {
		for tileX := 0; tileX < COLS; tileX++ {
			tileIndex := tileY*COLS + tileX
			tileNumber := t.vram.GetTileMapValue(vram.Map0, tileIndex)
			hiVal := (tileNumber >> 4) & 0xF
			loVal := tileNumber & 0xF

			var hiDigit, loDigit rune
			if hiVal < 10 {
				hiDigit = '0' + rune(hiVal)
			} else {
				hiDigit = 'A' + rune(hiVal-10)
			}
			if loVal < 10 {
				loDigit = '0' + rune(loVal)
			} else {
				loDigit = 'A' + rune(loVal-10)
			}

			hiSprite := font.CharMap[hiDigit]
			loSprite := font.CharMap[loDigit]

			offsetX := tileX*CELL_SIZE + 1
			offsetY := tileY*CELL_SIZE + 1
			// Vertical centering: 8 - 4 = 4, split evenly = 2
			yOffset := 2

			// Draw hiSprite at left (x offset 0) and loSprite at right (x offset 4)
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					var pixelColor color.Color
					if hiSprite.Pixels[i][j] == 1 {
						pixelColor = color.Black
					} else {
						pixelColor = color.White
					}
					t.img.Set(offsetX+j, offsetY+yOffset+i, pixelColor)
				}
			}
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					var pixelColor color.Color
					if loSprite.Pixels[i][j] == 1 {
						pixelColor = color.Black
					} else {
						pixelColor = color.White
					}
					t.img.Set(offsetX+4+j, offsetY+yOffset+i, pixelColor)
				}
			}
		}
	}

	// Draw gridlines.
	gridColor := color.Black
	for row := 0; row <= ROWS; row++ {
		y := row * CELL_SIZE
		for x := 0; x < FULL_GRID_WIDTH; x++ {
			t.img.Set(x, y, gridColor)
		}
	}

	for col := 0; col <= COLS; col++ {
		x := col * CELL_SIZE
		for y := 0; y < FULL_GRID_HEIGHT; y++ {
			t.img.Set(x, y, gridColor)
		}
	}
}

func (t *TilemapMenu) Config() *display.Config {
	return &t.config
}
