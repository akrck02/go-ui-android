package material

import (
	"math"
	"time"

	ui "github.com/akrck02/go-ui"
)

// LinearProgress is a bar across that fills as work gets done: fraction from 0 to 1. A negative
// fraction is work of unknown length: a piece of the bar runs along it, and it asks for frames
// while shown.
func (look Look) LinearProgress(fraction float32) {
	scheme := look.Scheme
	track := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(4)}
	if fraction >= 0 {
		look.Context.Frame(track, ui.Look{Fill: scheme.SecondaryContainer, Radius: 2}, func() {
			look.Context.Frame(ui.Layout{Width: ui.Share(max(min(fraction, 1), 0.001)), Height: ui.Grow()}, ui.Look{Fill: scheme.Primary, Radius: 2}, nil)
			look.Context.Box(ui.Layout{Width: ui.Share(max(1-fraction, 0.001))}, nil)
		})
		return
	}
	// a third of the bar sweeps across it every 1.5 seconds
	phase := float32(math.Mod(float64(time.Now().UnixMilli())/1500, 1))
	before := max(phase*1.3-0.3, 0)
	length := min(phase*1.3, 1) - before
	look.Context.Frame(track, ui.Look{Fill: scheme.SecondaryContainer, Radius: 2}, func() {
		look.Context.Box(ui.Layout{Width: ui.Share(before + 0.001)}, nil)
		look.Context.Frame(ui.Layout{Width: ui.Share(length + 0.001), Height: ui.Grow()}, ui.Look{Fill: scheme.Primary, Radius: 2}, nil)
		look.Context.Box(ui.Layout{Width: ui.Share(max(1-before-length, 0) + 0.001)}, nil)
	})
	look.Context.RequestFrame(16 * time.Millisecond)
}
