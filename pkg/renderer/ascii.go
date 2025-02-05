package renderer

import (
	"fmt"
	"image"
	"image/color"
)

func RenderANSI(img image.Image) string {
	var ansi string

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			// Get the color of the pixel at (x, y)
			color := img.At(x, y)

			// Convert the color to an ANSI escape code
			ansi += colorToANSI(color)
		}
		// Print a newline
		ansi += "\n"
	}

	return ansi
}

func colorToANSI(c color.Color) string {
	// RGBA returns values in the range [0, 65535]. Convert these to 0-255.
	r, g, b, _ := c.RGBA()
	R := r >> 8
	G := g >> 8
	B := b >> 8

	// ANSI escape code for foreground color using 24-bit color.
	return fmt.Sprintf("\033[38;2;%d;%d;%dm██\033[0m", R, G, B)
}
