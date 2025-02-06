package renderer

import (
	"bytes"
	"image"
	"image/draw"
	"strconv"
	"unsafe"

	"github.com/soniakeys/quant/median"
)

// RenderSixel converts an image.Image into a string containing its Sixel‐encoded data.
// It uses a maximum palette size of 256 colors and a 1:1 aspect ratio so the image
// will not appear 2× as tall as it really is.
func RenderSixel(img image.Image) string {
	// Use a maximum of 256 colors (0 is a reserved transparency key).
	const maxColors = 256

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return ""
	}

	// Convert to paletted. (Optional dithering with Floyd–Steinberg.)
	var paletted *image.Paletted
	if p, ok := img.(*image.Paletted); ok && len(p.Palette) <= maxColors {
		paletted = p
	} else {
		quant := median.Quantizer(maxColors - 1)
		paletted = quant.Paletted(img)
		draw.FloydSteinberg.Draw(paletted, bounds, img, bounds.Min)
	}

	// Preallocate a buffer.
	outBuf := bytes.NewBuffer(make([]byte, 0, width*height/2))

	//----------------------------------------------------------------------
	// 1) Write the Sixel Introducer and Raster Attributes (aspect ratio)
	//----------------------------------------------------------------------
	// ESC P 0;0;8q   (introduces Sixel mode)
	outBuf.Write([]byte{0x1b, 'P', '0', ';', '0', ';', '8', 'q'})

	// "1;1;<width>;<height>   => sets the pixel aspect ratio to 1:1
	outBuf.WriteByte('"')
	outBuf.WriteString("1;1;")
	outBuf.WriteString(strconv.Itoa(width))
	outBuf.WriteByte(';')
	outBuf.WriteString(strconv.Itoa(height))

	// Optionally you could add more parameters if you want to scale or define
	// an offset, but the above is sufficient to force a 1:1 aspect.

	//----------------------------------------------------------------------
	// 2) Write palette definitions
	//----------------------------------------------------------------------
	for i, c := range paletted.Palette {
		r32, g32, b32, _ := c.RGBA()
		r := int(r32 * 100 / 0xffff)
		g := int(g32 * 100 / 0xffff)
		b := int(b32 * 100 / 0xffff)
		writePaletteEntry(outBuf, i, r, g, b)
	}

	//----------------------------------------------------------------------
	// 3) Encode the image data row by row (in 6-row bands)
	//----------------------------------------------------------------------
	nColors := len(paletted.Palette)
	yBands := (height + 5) / 6

	for z := 0; z < yBands; z++ {
		if z > 0 {
			// New line in sixel
			outBuf.WriteByte('-')
		}
		// For each color in our palette, generate sixels.
		for col := 1; col <= nColors; col++ {
			writeColorChange(outBuf, col)

			var runCount int
			var currentVal byte // will hold the 6 bits
			for x := 0; x < width; x++ {
				var sixel byte
				for dy := 0; dy < 6; dy++ {
					y := z*6 + dy
					if y >= height {
						continue
					}
					_, _, _, alpha := img.At(x, y).RGBA()
					if alpha == 0 {
						// Transparent
						continue
					}
					if int(paletted.ColorIndexAt(x, y))+1 == col {
						sixel |= 1 << uint(dy)
					}
				}

				// Run-length encode
				if x == 0 {
					currentVal = sixel
					runCount = 1
				} else {
					if sixel == currentVal {
						runCount++
					} else {
						flushRun(outBuf, runCount, currentVal)
						currentVal = sixel
						runCount = 1
					}
				}
			}
			// Flush out the final run in this color band
			flushRun(outBuf, runCount, currentVal)
		}
	}

	// End the sixel data: ESC \
	outBuf.Write([]byte{0x1b, '\\'})

	// Convert []byte to string with an unsafe pointer to avoid copying.
	return unsafeString(outBuf.Bytes())
}

// writePaletteEntry writes a Sixel palette definition for the given palette index
// (starting at 0, but the Sixel command uses index+1) and the red, green, blue values.
func writePaletteEntry(buf *bytes.Buffer, index, r, g, b int) {
	buf.WriteByte('#')
	buf.WriteString(strconv.Itoa(index + 1)) // Sixel color indices start from 1
	buf.WriteString(";2;")
	buf.WriteString(strconv.Itoa(r))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(g))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(b))
}

// writeColorChange writes a color change command for the given (1-indexed) color.
func writeColorChange(buf *bytes.Buffer, col int) {
	buf.WriteByte('#')
	buf.WriteString(strconv.Itoa(col))
}

// flushRun writes run‑length encoded sixel data for a given run.
func flushRun(buf *bytes.Buffer, count int, val byte) {
	if count <= 0 {
		return
	}
	// Each pixel pattern is offset by 63 in Sixel
	sixelChar := val + 63
	// If count>1, we do "!<count><sixel_char>" for run-length encoding.
	if count == 1 {
		buf.WriteByte(sixelChar)
	} else {
		buf.WriteByte('!')
		buf.WriteString(strconv.Itoa(count))
		buf.WriteByte(sixelChar)
	}
}

// unsafeString allows us to convert a []byte to string without copying the data.
func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
