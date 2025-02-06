package utils

import (
	"image"
	"strings"

	"github.com/colecrouter/gameboy-go/pkg/renderer"
)

type Border uint

const (
	BorderNone Border = iota
	BorderSingle
	BorderDouble
)

type BoxOptions struct {
	Border
	Title string
}

type borderOption struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

var borderOptions = map[Border]borderOption{
	BorderSingle: {
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
	},
	BorderDouble: {
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		Vertical:    "║",
	},
	BorderNone: {
		TopLeft:     " ",
		TopRight:    " ",
		BottomLeft:  " ",
		BottomRight: " ",
		Horizontal:  " ",
		Vertical:    " ",
	},
}

func DrawBox(img image.Image, options *BoxOptions) []string {
	border := borderOptions[options.Border]
	content := renderer.RenderSixel(img)

	bounds := img.Bounds()
	width := bounds.Dx()/8 + 2
	horizontalCount := (width - 2) * 2

	topLine := border.TopLeft + " " + options.Title + " " + repeat(border.Horizontal, horizontalCount-len(options.Title)-2) + border.TopRight
	middleLine := content
	bottomLine := border.BottomLeft + repeat(border.Horizontal, horizontalCount) + border.BottomRight

	return []string{topLine, middleLine, bottomLine}
}

func repeat(s string, n int) string {
	if n < 0 {
		n = 0
	}
	return strings.Repeat(s, n)
}
