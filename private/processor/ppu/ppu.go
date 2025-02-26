package ppu

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/layers"
	"github.com/colecrouter/gameboy-go/private/system"
)

type PPU struct {
	interrupt        *io.Interrupt
	vram             *vram.VRAM
	oam              *memory.OAM
	registers        *io.Registers
	lineCycleCounter uint16
	image            *image.Paletted
	clock            <-chan struct{}
}

const (
	// Row timings
	oamScanCycles       = 80
	pixelTransferCycles = 172
	hBlankCycles        = 204
	TotalCyclesPerLine  = oamScanCycles + pixelTransferCycles + hBlankCycles

	// Column timings
	visibleLines       = 144
	vBlankLines        = 10
	TotalLinesPerFrame = visibleLines + vBlankLines

	// Screen dimensions
	visibleColumns = 160
)

func NewPPU(broadcaster *system.Broadcaster, vram *vram.VRAM, oam *memory.OAM, registers *io.Registers, ie *io.Interrupt) *PPU {
	return &PPU{
		interrupt: ie,
		vram:      vram,
		oam:       oam,
		registers: registers,
		image:     image.NewPaletted(image.Rect(0, 0, visibleColumns, visibleLines), monochrome.Palette),
		clock:     broadcaster.SubscribeT(),
	}
}

/*
-------┌──────────┐-------
   oam │ transfer │ hblank
  80 c │   172 c  │ 204 c
       │ x 144 l  │
       │          │
-------└──────────┘-------
         vblank
         10 l
*/

// TClock emulates one PPU cycle.
func (p *PPU) TClock() {
	if p.registers.LY >= visibleLines {
		p.registers.LCDStatus.PPUMode = io.VBlank
		if p.registers.LY == visibleLines {
			p.interrupt.VBlank = true
		}
	} else {
		switch p.lineCycleCounter {
		case 0:
			p.registers.LCDStatus.PPUMode = io.OAMScan
			p.interrupt.LCD = true
		case oamScanCycles:
			p.registers.LCDStatus.PPUMode = io.Drawing
			p.interrupt.LCD = true
		case oamScanCycles + pixelTransferCycles:
			p.registers.LCDStatus.PPUMode = io.HBlank
			p.interrupt.LCD = true
		}
	}

	p.lineCycleCounter++
	if p.lineCycleCounter == TotalCyclesPerLine {
		p.registers.LY++
		p.lineCycleCounter = 0
	}
	if p.registers.LY == TotalLinesPerFrame {
		p.registers.LY = 0
	}

	<-p.clock
}

// compositeImage overlays src onto dst; it assumes pixel value 0 is transparent.
func compositeImage(dst, src *image.Paletted) {
	// Both images must have the same bounds.
	for i, pix := range src.Pix {
		if pix != 0 {
			dst.Pix[i] = pix
		}
	}
}

// DisplayClock updates p.image by compositing BG, window, and sprite layers.
func (p *PPU) DisplayClock() {
	// Create new layers with the screen bounds.
	screenRect := image.Rect(0, 0, visibleColumns, visibleLines)
	bgLayer := layers.NewBGLayer(p.vram, p.registers, screenRect)
	winLayer := layers.NewWindowLayer(p.vram, p.registers, screenRect)
	spriteLayer := layers.NewSpriteLayer(p.oam, p.vram, p.registers, screenRect)

	// Render layers.
	bgImg := bgLayer.Image()
	finalImg := image.NewPaletted(screenRect, monochrome.Palette)

	// Start with the background.
	copy(finalImg.Pix, bgImg.(*image.Paletted).Pix)

	// Composite window layer (if enabled).
	if p.registers.LCDControl.EnableWindow {
		winImg := winLayer.Image()
		compositeImage(finalImg, winImg.(*image.Paletted))
	}

	// Composite sprite layer.
	spriteImg := spriteLayer.Image()
	compositeImage(finalImg, spriteImg.(*image.Paletted))

	// Set the composite image as the PPU output.
	p.image = finalImg
}

func (p *PPU) Image() image.Image {
	return p.image
}

func (p *PPU) Run(close <-chan struct{}) {
	for {
		select {
		case <-close:
			return
		default:
			p.TClock()
		}
	}
}
