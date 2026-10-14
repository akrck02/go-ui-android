package material

import (
	"fmt"
	"math"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// Gestures a finger makes on lists and rows: pulling a list down to refresh it, swiping a row away,
// and sliding a row of cards. go-ui gives a finger to the first part it moves for (a list for up
// and down, these for sideways), so they sit inside each other freely.

// LongPressed tells whether the piece drawn just before it was long-pressed in this frame: a
// finger held still on it for half a second, or the mouse's right button. The phone gives its
// buzz then, as Android's own applications do. A list uses it to start choosing rows (Item.Selected);
// a touch that became a long press is no touch when it lifts.
func (look Look) LongPressed() bool {
	if !look.Context.LongPressed() {
		return false
	}
	window.FeelLongPress()
	return true
}

// refreshPull is how far a list is pulled before letting go refreshes it, and refreshSpace the room
// the indicator keeps over the list while it refreshes.
const (
	refreshPull  = 64
	refreshSpace = 64
)

// PullToRefresh is a scrolling list that refreshes when pulled down past its top and let go: a
// loading indicator follows the pull, and stays over the list while refreshing is true. It returns
// true in the frame the person asks for a refresh; the program then refreshes and sets refreshing
// until it is done. name tells it apart; content adds the list's rows.
func (look Look) PullToRefresh(name string, refreshing bool, layout ui.Layout, content func()) (refresh bool) {
	context := look.Context
	pulling, released := context.ScrollPull(name)
	refresh = released >= refreshPull && !refreshing
	// the indicator's room follows the finger, then eases to its place or away
	target, motion := float32(0), arrivingMotion
	switch {
	case pulling > 0:
		target, motion = pulling, 0
	case refreshing || refresh:
		target = refreshSpace
	}
	room := context.Animate(name+" pull", target, motion)
	layout.Direction = ui.Vertical
	context.Box(layout, func() {
		if room >= 1 {
			context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(room), AlignX: ui.Center, AlignY: ui.Center}, func() {
				look.ContainedLoadingIndicator(min(40, max(room-16, 8)))
			})
		}
		context.ScrollArea(name, ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, content)
	})
	return refresh
}

// SwipeToDismiss is a row a finger slides sideways to take it away, a bin showing in the gap it
// leaves: pushed past two fifths of its width, or flicked, it goes, and SwipeToDismiss returns
// true for the program to remove it; let go short of that, it slides back. name tells rows apart.
func (look Look) SwipeToDismiss(name string, content func()) (dismissed bool) {
	context, scheme := look.Context, look.Scheme
	context.SwipeBox(name, ui.Layout{Width: ui.Grow()}, ui.Horizontal, func(swipe ui.Swipe) {
		width := swipe.Area.Width
		if swipe.Released && width > 0 && (math.Abs(float64(swipe.Offset)) > float64(width)*0.4 || math.Abs(float64(swipe.Velocity)) > 1200) {
			dismissed = true
		}
		target, motion := float32(0), switchMotion
		if swipe.Held {
			target, motion = swipe.Offset, 0
		}
		offset := context.Animate(name+" offset", target, motion)
		gap := func() {
			if offset == 0 {
				return
			}
			layout := ui.Layout{Width: ui.Fixed(float32(math.Abs(float64(offset)))), Height: ui.Grow(), Padding: ui.PadXY(20, 0), AlignY: ui.Center}
			layout.AlignX = ui.Start
			if offset < 0 {
				layout.AlignX = ui.End
			}
			context.Frame(layout, ui.Look{Fill: scheme.ErrorContainer}, func() {
				context.IconSized("trash-2", 24, scheme.OnErrorContainer)
			})
		}
		row := ui.Layout{Width: ui.Grow()}
		if width > 0 {
			row.Width = ui.Fixed(width)
		}
		strip := ui.Layout{Width: ui.Grow()}
		if offset < 0 {
			strip.Padding.Left = offset
		}
		context.Box(strip, func() {
			if offset > 0 {
				gap()
			}
			context.Frame(row, ui.Look{Fill: scheme.Surface}, content)
			if offset < 0 {
				gap()
			}
		})
	})
	return dismissed
}

// Carousel is what a carousel remembers: where it rests, as the index of the first item shown.
type Carousel struct {
	rest float32
}

// Carousel is a row of items wider than the screen, slid sideways with a finger: let go, it glides
// on with the finger's speed and settles at an item's edge. Each item is itemWidth wide, with 8
// pixels between them; item draws one, given its index. name tells carousels apart.
func (look Look) Carousel(name string, carousel *Carousel, count int, itemWidth float32, item func(index int)) {
	context := look.Context
	const gap = 8
	step := itemWidth + gap
	context.SwipeBox(name, ui.Layout{Width: ui.Grow()}, ui.Horizontal, func(swipe ui.Swipe) {
		shown := float32(count)
		if swipe.Area.Width > 0 {
			shown = (swipe.Area.Width - 16) / step
		}
		last := max(float32(count)-shown, 0)
		position := carousel.rest - swipe.Offset/step
		target, motion := clampFloat(position, -0.15, last+0.15), time.Duration(0) // a little give at the ends
		if swipe.Released {
			// on with the finger's speed, a quarter of a second's worth, then to an item's edge
			carousel.rest = clampFloat(float32(math.Round(float64(position-swipe.Velocity*0.25/step))), 0, float32(math.Ceil(float64(last))))
			target, motion = carousel.rest, 3*arrivingMotion/2
		} else if !swipe.Held {
			target, motion = carousel.rest, 3*arrivingMotion/2
		}
		at := context.Animate(name+" position", target, motion)
		context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16 - at*step, Right: 16}, Gap: gap}, func() {
			for index := range count {
				context.Push(fmt.Sprint(name, " item ", index))
				context.Box(ui.Layout{Width: ui.Fixed(itemWidth)}, func() { item(index) })
				context.Pop()
			}
		})
	})
}

func clampFloat(value, minimum, maximum float32) float32 { return min(max(value, minimum), maximum) }

// CarouselItem is a carousel's usual item: a rounded tile in a colour, with a title over its foot.
func (look Look) CarouselItem(height float32, tint draw2d.Color, icon, title string) {
	tile := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(height), Padding: ui.Pad(16), AlignY: ui.End, Gap: 4}
	look.Context.Frame(tile, ui.Look{Fill: tint, Radius: 28}, func() {
		ink := readableOn(tint)
		if icon != "" {
			look.Context.IconSized(icon, 28, ink)
		}
		look.Context.Spacer()
		look.Text(title, 16, draw2d.Medium, ink)
	})
}
