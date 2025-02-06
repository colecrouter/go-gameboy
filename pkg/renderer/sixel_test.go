package renderer

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestRenderSixel(t *testing.T) {
	// Create a simple 1x1 red image.
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Call RenderSixel.
	result := RenderSixel(img)

	// Validate that the result is non-empty.
	if result == "" {
		t.Error("RenderSixel returned an empty string")
	}

	// Validate the Sixel header.
	expectedHeader := "\x1bP0;0;8q"
	if !strings.HasPrefix(result, expectedHeader) {
		t.Errorf("expected header %q, got %q", expectedHeader, result[:len(expectedHeader)])
	}

	// Validate the Sixel terminator.
	expectedTerminator := "\x1b\\"
	if !strings.HasSuffix(result, expectedTerminator) {
		t.Errorf("expected terminator %q, got %q", expectedTerminator, result[len(result)-len(expectedTerminator):])
	}
}

// ...additional tests if needed...
