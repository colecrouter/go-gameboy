package renderer

import (
	"bytes"
	"image"
	"strconv"
	"unsafe"

	"github.com/soniakeys/quant/median"
)

// RenderSixel converts an image.Image into a Sixel-encoded string.
// It uses a maximum palette size of 256 colors and enforces a 1:1 aspect ratio.
func RenderSixel(img image.Image) string {
	// Use a maximum of 256 colors (0 is reserved for transparency).
	const maxColors = 256

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return ""
	}

	// Convert image to paletted form using a median cut quantizer.
	var paletted *image.Paletted
	if p, ok := img.(*image.Paletted); ok && len(p.Palette) <= maxColors {
		paletted = p
	} else {
		quant := median.Quantizer(maxColors - 1)
		paletted = quant.Paletted(img)
		// Optionally: apply dithering with Floyd–Steinberg.
	}

	// Preallocate buffer for Sixel data.
	outBuf := bytes.NewBuffer(make([]byte, 0, width*height/2))

	// Write Sixel Introducer and raster attributes (enforcing 1:1 pixel aspect ratio).
	outBuf.Write([]byte{0x1b, 'P', '0', ';', '0', ';', '8', 'q'})

	// Set pixel aspect ratio to 1:1.
	outBuf.WriteByte('"')
	outBuf.WriteString("1;1")
	// Removed unused width/height write.

	// Define palette: reserve index 0 for transparency.
	outBuf.WriteByte('#')
	outBuf.WriteString("0;2;0;0;0")

	for i, c := range paletted.Palette {
		r32, g32, b32, _ := c.RGBA()
		r := int(r32 * 100 / 0xffff)
		g := int(g32 * 100 / 0xffff)
		b := int(b32 * 100 / 0xffff)
		writePaletteEntry(outBuf, i, r, g, b)
	}

	// Encode image data in 6-row bands.
	nColors := len(paletted.Palette)
	yBands := (height + 5) / 6

	for z := 0; z < yBands; z++ {
		if z > 0 {
			outBuf.WriteByte('-') // Newline command for subsequent bands.
		}

		// Determine which colors are used in the current 6-row band.
		usedColors := make([]bool, nColors+1)
		for x := 0; x < width; x++ {
			for dy := 0; dy < 6; dy++ {
				y := z*6 + dy
				if y >= height {
					continue
				}
				_, _, _, alpha := img.At(x, y).RGBA()
				if alpha != 0 {
					colIdx := int(paletted.ColorIndexAt(x, y)) + 1
					usedColors[colIdx] = true
				}
			}
		}

		// Process each used color in the band.
		for col := 1; col <= nColors; col++ {
			if !usedColors[col] {
				continue
			}
			outBuf.WriteByte('$') // Start new color band
			writeColorChange(outBuf, col)

			runCount := 0
			currentVal := byte(0)
			// Iterate over columns to compute sixel value from six rows.
			for x := 0; x < width; x++ {
				var sixel byte
				for dy := 0; dy < 6; dy++ {
					y := z*6 + dy
					if y >= height {
						continue
					}
					_, _, _, alpha := img.At(x, y).RGBA()
					if alpha != 0 && int(paletted.ColorIndexAt(x, y))+1 == col {
						sixel |= 1 << uint(dy)
					}
				}
				// Run-length encode repeated sixel values.
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
			flushRun(outBuf, runCount, currentVal)
		}
	}

	// End Sixel data.
	outBuf.Write([]byte{0x1b, '\\'})

	// Convert buffer to string without copying.
	return unsafeString(outBuf.Bytes())
}

// writePaletteEntry writes a Sixel palette definition using a 1-indexed color.
func writePaletteEntry(buf *bytes.Buffer, index, r, g, b int) {
	buf.WriteByte('#')
	buf.WriteString(strconv.Itoa(index + 1))
	buf.WriteString(";2;")
	buf.WriteString(strconv.Itoa(r))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(g))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(b))
}

// writeColorChange writes a command to change the current drawing color.
func writeColorChange(buf *bytes.Buffer, col int) {
	buf.WriteByte('#')
	buf.WriteString(strconv.Itoa(col))
}

// flushRun writes run‑length encoded sixel data for a run of identical values.
func flushRun(buf *bytes.Buffer, count int, val byte) {
	if count <= 0 {
		return
	}
	sixelChar := val + 63
	for count > 255 {
		buf.WriteByte('!')
		buf.WriteString(strconv.Itoa(255))
		buf.WriteByte(sixelChar)
		count -= 255
	}
	if count == 1 {
		buf.WriteByte(sixelChar)
	} else if count > 1 {
		buf.WriteByte('!')
		buf.WriteString(strconv.Itoa(count))
		buf.WriteByte(sixelChar)
	}
}

// unsafeString converts a []byte to a string without extra memory allocation.
func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
