package material

import (
	"testing"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

func fixed(size float32) ui.Sizing { return ui.Fixed(size) }

// measurer stands in for the fonts: no display needed.
type measurer struct{}

func (measurer) MeasureText(text string, size float32, _ draw2d.Weight, _ float32) draw2d.TextMetrics {
	return draw2d.TextMetrics{Width: float32(utf8.RuneCountInString(text)) * size / 2, Ascent: size * 0.8, Descent: size * 0.2}
}

func (measurer) HasIcon(string) bool { return true }

// screen drives a phone-sized interface frame by frame.
type screen struct {
	t       *testing.T
	context *ui.Context
	state   *input.State
	build   func(look Look)
}

func newScreen(t *testing.T, build func(look Look)) *screen {
	state := input.NewState()
	state.WindowWidth, state.WindowHeight = 400, 800
	state.FramebufferWidth, state.FramebufferHeight = 400, 800
	return &screen{t: t, context: ui.NewWithMeasurer(measurer{}), state: state, build: build}
}

func (screen *screen) frame() {
	screen.context.Begin(screen.state)
	screen.build(New(screen.context, Baseline(false)))
	screen.context.End()
	screen.state.EndFrame()
}

// tap presses and lifts the pointer in the middle of a piece laid out in the last frame.
func (screen *screen) tap(path ...string) {
	screen.t.Helper()
	area, found := screen.context.Area(path...)
	if !found {
		screen.t.Fatalf("nothing called %v was laid out", path)
	}
	screen.state.MoveCursor(float64(area.X+area.Width/2), float64(area.Y+area.Height/2))
	screen.state.SetButton(input.MouseLeft, true)
	screen.frame()
	screen.state.SetButton(input.MouseLeft, false)
	screen.frame()
}

// palettesOf makes palettes whose shades are greys from white (0) to black (1000), with each accent
// tinted so the roles can be told apart.
func palettesOf() window.Palettes {
	var palettes window.Palettes
	for index := range window.PaletteTones {
		grey := 1 - float32(index)/12
		palettes.Accent1[index] = draw2d.Color{grey, 0, 0, 1}
		palettes.Accent2[index] = draw2d.Color{0, grey, 0, 1}
		palettes.Accent3[index] = draw2d.Color{0, 0, grey, 1}
		palettes.Neutral1[index] = draw2d.Color{grey, grey, grey, 1}
		palettes.Neutral2[index] = draw2d.Color{grey, grey, grey * 0.9, 1}
	}
	return palettes
}

func distance(first, second draw2d.Color) float32 {
	return abs(first[0]-second[0]) + abs(first[1]-second[1]) + abs(first[2]-second[2])
}

func abs(value float32) float32 { return max(value, -value) }
