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
// It uses a maximum palette size of 256 colors (with a reserved transparency key).
// The algorithm uses median cut quantization (with an optional Floyd–Steinberg dithering)
// and then renders the image six rows at a time, with run‐length encoding.
func RenderSixel(img image.Image) string {
	// Use a maximum of 256 colors. (Note: color 0 is “reserved” for transparent pixels.)
	const maxColors = 256

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 || height == 0 {
		return ""
	}

	// Create a paletted copy of the image.
	var pal *image.Paletted
	if p, ok := img.(*image.Paletted); ok && len(p.Palette) <= maxColors {
		pal = p
	} else {
		// Use median cut quantizer (with maxColors-1 colors)
		quant := median.Quantizer(maxColors - 1)
		pal = quant.Paletted(img)
		// Apply Floyd–Steinberg dithering.
		draw.FloydSteinberg.Draw(pal, bounds, img, bounds.Min)
	}

	// Create an output bytes buffer with a reasonable preallocated capacity.
	outBuf := bytes.NewBuffer(make([]byte, 0, width*height/2))

	// Write the Sixel header:
	// DECSIXEL Introducer: ESC P 0;0;8q
	outBuf.Write([]byte{0x1b, 'P', '0', ';', '0', ';', '8', 'q'})

	// Write palette (graphics color definitions). Sixel expects definitions
	// like: "#<n>;2;<r>;<g>;<b>" with values scaled to a 0–100 range.
	for i, c := range pal.Palette {
		r32, g32, b32, _ := c.RGBA()
		r := int(r32 * 100 / 0xffff)
		g := int(g32 * 100 / 0xffff)
		b := int(b32 * 100 / 0xffff)
		writePaletteEntry(outBuf, i, r, g, b)
	}

	// Sixel image data is emitted six rows at a time.
	// For each six-row band we will, for each color in our palette (colors 1..N),
	// create a run-length encoded sequence.
	nColors := len(pal.Palette)
	// For each block of 6 rows (vertical band)
	yBands := (height + 5) / 6
	for z := 0; z < yBands; z++ {
		// “New line” indicator except at the top.
		if z > 0 {
			outBuf.WriteByte('-')
		}
		// Process each color sequentially.
		// (Sixel color indices are 1–nColors; we assume 0 is reserved for transparency.)
		for col := 1; col <= nColors; col++ {
			// Output a color change command.
			writeColorChange(outBuf, col)

			// For each column x in the image we form a 6‐bit value.
			var runCount int
			var currentVal byte // current sixel value (0..63)
			for x := 0; x < width; x++ {
				var sixel byte = 0
				// Process a vertical “column” of up to 6 pixels.
				for dy := 0; dy < 6; dy++ {
					y := z*6 + dy
					if y >= height {
						continue
					}
					// Only if the pixel isn’t fully transparent.
					px := img.At(x, y)
					_, _, _, a := px.RGBA()
					if a == 0 {
						// Transparent pixel; leave bit off.
						continue
					}
					// Check if the pixel’s palette index (shifted to 1..nColors) equals col.
					if int(pal.ColorIndexAt(x, y))+1 == col {
						sixel |= 1 << uint(dy)
					}
				}
				// On the first column, initialize the run.
				if x == 0 {
					currentVal = sixel
					runCount = 1
				} else if sixel == currentVal {
					runCount++
				} else {
					flushRun(outBuf, runCount, currentVal)
					currentVal = sixel
					runCount = 1
				}
			}
			flushRun(outBuf, runCount, currentVal)
		}
	}

	// Write the terminator: ESC \
	outBuf.Write([]byte{0x1b, '\\'})

	// Use an unsafe conversion to turn the buffer’s bytes into a string without a copy.
	return unsafeString(outBuf.Bytes())
}

// writePaletteEntry writes a Sixel palette definition for the given palette index
// (starting at 0, but the Sixel command uses index+1) and the red, green, blue values.
func writePaletteEntry(buf *bytes.Buffer, index, r, g, b int) {
	buf.WriteByte('#')
	// Palette numbering in Sixel starts at 1.
	buf.Write(strconv.AppendInt(nil, int64(index+1), 10))
	buf.WriteString(";2;")
	buf.Write(strconv.AppendInt(nil, int64(r), 10))
	buf.WriteByte(';')
	buf.Write(strconv.AppendInt(nil, int64(g), 10))
	buf.WriteByte(';')
	buf.Write(strconv.AppendInt(nil, int64(b), 10))
}

// writeColorChange writes a color change command for the given (1‐indexed) color.
func writeColorChange(buf *bytes.Buffer, col int) {
	buf.WriteByte('#')
	// Use decimal conversion for the color number.
	// (For numbers with more than one digit, strconv.AppendInt handles it correctly.)
	buf.Write(strconv.AppendInt(nil, int64(col), 10))
}

// flushRun writes run‑length encoded sixel data for a given run.
// In Sixel the sixel code is given by: 63 + <6‐bit value>.
// If count is 1, the sixel character (byte) is written directly;
// if count > 1 the sequence is written as: "!" followed by the count (in decimal) and then the sixel.
func flushRun(buf *bytes.Buffer, count int, val byte) {
	if count <= 0 {
		return
	}
	sixelChar := val + 63
	if count == 1 {
		buf.WriteByte(sixelChar)
	} else {
		buf.WriteByte('!')
		buf.Write(strconv.AppendInt(nil, int64(count), 10))
		buf.WriteByte(sixelChar)
	}
}

// unsafeString converts a byte slice to a string without copying the data.
// Use with care!
func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
