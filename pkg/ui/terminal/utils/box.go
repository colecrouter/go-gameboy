package utils

import (
	"strings"

	"github.com/colecrouter/gameboy-go/pkg/display"
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

func DrawBox(d display.Display, options *BoxOptions) []string {
	border := borderOptions[options.Border]
	var content []string
	var width int
	var isImage bool
	switch v := d.(type) {
	case display.ImageDisplay:
		content = append(content, renderer.RenderSixel(v.Image()))

		imageSize := v.Image().Bounds().Dx() / 11 // Wooo magic number
		boxMinSize := len(d.Config().Title) + 4   // Title + 2 spaces + 2 borders

		width = max(imageSize, boxMinSize)
		isImage = true
	case display.TextDisplay:
		content = append(content, v.Text()...)
		width = max(d.Config().Width/CHAR_WIDTH, len(d.Config().Title)+4)
	}

	horizontalCount := width

	title := d.Config().Title

	topLine := border.TopLeft + " " + title + " " + repeat(border.Horizontal, horizontalCount-len(title)-2) + border.TopRight
	bottomLine := border.BottomLeft + repeat(border.Horizontal, horizontalCount) + border.BottomRight

	if isImage {
		// If it's an image, the content will be a single line.
		return []string{topLine, content[0], bottomLine}
	}

	height := max(len(content), d.Config().Height/CHAR_HEIGHT)
	var lines []string
	for i := 0; i < height; i++ {
		var line string
		if i < len(content) {
			line = content[i]
		}
		lines = append(lines, border.Vertical+line+repeat(" ", horizontalCount-len(line))+border.Vertical)
	}

	return append([]string{topLine}, append(lines, bottomLine)...)
}

func repeat(s string, n int) string {
	if n < 0 {
		n = 0
	}
	return strings.Repeat(s, n)
}
