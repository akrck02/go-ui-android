package main

import (
	"testing"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

type measurer struct{}

func (measurer) MeasureText(text string, size float32, _ draw2d.Weight, _ float32) draw2d.TextMetrics {
	return draw2d.TextMetrics{Width: float32(utf8.RuneCountInString(text)) * size / 2, Ascent: size * 0.8, Descent: size * 0.2}
}

func (measurer) HasIcon(string) bool { return true }

// Every destination can be reached from every other, by its place in the navigation bar.
func TestEveryScreenLeadsToEveryOther(t *testing.T) {
	program := &gallery{inputs: newInputs(), lists: newLists(), tabs: newTabs()}
	context, state := ui.NewWithMeasurer(measurer{}), input.NewState()
	state.WindowWidth, state.WindowHeight = 412, 892
	state.FramebufferWidth, state.FramebufferHeight = 412, 892
	var finger input.TouchPointer
	frame := func() {
		context.Begin(state)
		program.build(context)
		context.End()
		state.EndFrame()
	}
	tap := func(label string) {
		area, found := context.Area("navigate " + label)
		if !found {
			t.Fatalf("no destination %s on screen %d", label, program.screen)
		}
		// as a phone's finger does: down, a frame, up, a frame
		x, y := area.X+area.Width/2, area.Y+area.Height/2
		finger.Begin(state, 1, x, y, 1)
		frame()
		finger.End(state, 1, x, y, 1)
		frame()
		frame()
	}
	labels := []string{"Actions", "Inputs", "Lists", "Feedback", "Tabs"}
	frame()
	for from := range labels {
		for to := range labels {
			tap(labels[from])
			tap(labels[to])
			if program.screen != to {
				t.Fatalf("from %s, touching %s left the gallery on screen %d", labels[from], labels[to], program.screen)
			}
		}
	}
}
