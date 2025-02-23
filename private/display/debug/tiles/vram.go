package tiles

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/private/display"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
)

const (
	TILE_SIZE        = 8                  // each tile's size
	COLS             = 16                 // number of tile columns
	ROWS             = 24                 // number of tile rows
	CELL_SIZE        = TILE_SIZE + 1      // tile plus a 1px gap (gridline)
	FULL_GRID_WIDTH  = COLS*CELL_SIZE + 1 // extra 1 pixel at the right edge
	FULL_GRID_HEIGHT = ROWS*CELL_SIZE + 1 // extra 1 pixel at the bottom edge
)

type TileMenu struct {
	vram    *vram.VRAM
	img     *image.Paletted
	palette *color.Palette
	config  display.Config
}

func NewTileDebug(v *vram.VRAM, p *color.Palette) *TileMenu {
	return &TileMenu{
		vram:    v,
		palette: p,
		config:  display.Config{Title: "Tile Viewer"},
	}
}

func (t *TileMenu) Image() image.Image {
	return t.img
}

func (t *TileMenu) Clock() {
	// create a larger image that includes room for gridlines
	t.img = image.NewPaletted(image.Rect(0, 0, FULL_GRID_WIDTH, FULL_GRID_HEIGHT), *t.palette)

	// Construct a grid of tiles.
	// Draw each tile in its cell, offset by 1 pixel to allow for the grid border.
	for tileY := 0; tileY < ROWS; tileY++ {
		for tileX := 0; tileX < COLS; tileX++ {
			tileIndex := tileY*COLS + tileX
			tile := t.vram.GetTile(tileIndex)
			if tile == nil {
				continue
			}
			tileImage := tile.Pixels()
			if tileImage == nil {
				continue
			}

			// Calculate the top-left pixel where this tile should be drawn.
			offsetX := tileX*CELL_SIZE + 1
			offsetY := tileY*CELL_SIZE + 1

			// copy tile pixels into the image
			for y := 0; y < TILE_SIZE; y++ {
				for x := 0; x < TILE_SIZE; x++ {
					// get the color index for this pixel
					colorIndex := tileImage[y*TILE_SIZE+x]
					// get the actual color from the palette
					color := (*t.palette)[colorIndex]
					// set the pixel in the image
					t.img.Set(offsetX+x, offsetY+y, color)
				}
			}
		}
	}

	// Draw gridlines between cells
	gridColor := color.Black

	// Draw horizontal gridlines: at every CELL_SIZE boundary, including bottom edge
	for row := 0; row <= ROWS; row++ {
		y := row * CELL_SIZE
		for x := 0; x < FULL_GRID_WIDTH; x++ {
			t.img.Set(x, y, gridColor)
		}
	}

	// Draw vertical gridlines: at every CELL_SIZE boundary, including right edge
	for col := 0; col <= COLS; col++ {
		x := col * CELL_SIZE
		for y := 0; y < FULL_GRID_HEIGHT; y++ {
			t.img.Set(x, y, gridColor)
		}
	}
}

func (t *TileMenu) Config() *display.Config {
	return &t.config
}
