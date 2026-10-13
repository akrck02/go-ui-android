package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// TopBar is what the top app bar shows.
type TopBar struct {
	Title    string
	Subtitle string // under the title, smaller; none when empty
	Back     bool   // an arrow before the title, to go back
	// Actions draws what goes at the bar's end, usually IconButtons; none when nil.
	Actions func()
	// OverScene paints the bar see-through (the scheme's scrim), over a scene behind the interface.
	OverScene bool
}

// TopAppBar is the bar at the top, painted up under the status bar, so the screen reaches the top
// edge as Android's own apps do. It returns true when the back arrow is touched.
func (look Look) TopAppBar(bar TopBar) (back bool) {
	insets := look.Context.Insets()
	scheme := look.Scheme
	fill := scheme.Surface
	if bar.OverScene {
		fill = scheme.Scrim
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(64 + insets.Top),
		Padding: ui.Padding{Left: 4, Top: insets.Top, Right: 4}, Gap: 4, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: fill}, func() {
		if bar.Back {
			back = look.IconButton("arrow-left")
		} else {
			look.Context.Box(ui.Layout{Width: ui.Fixed(12)}, nil)
		}
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 2}, func() {
			look.Text(bar.Title, 22, draw2d.Regular, scheme.OnSurface)
			if bar.Subtitle != "" {
				look.Text(bar.Subtitle, 14, draw2d.Regular, scheme.OnSurfaceVariant)
			}
		})
		if bar.Actions != nil {
			bar.Actions()
		}
	})
	return back
}

// Destination is a place the navigation bar leads to: an icon and a short label.
type Destination struct {
	Icon, Label string
}

// NavigationBar is the bar at the bottom, painted down under the gesture bar: three to five
// destinations, the selected one with its icon on a pill. It changes selected and returns true
// when another destination is touched.
func (look Look) NavigationBar(destinations []Destination, selected *int) bool {
	insets := look.Context.Insets()
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(80 + insets.Bottom),
		Padding: ui.Padding{Top: 12, Bottom: 16 + insets.Bottom}}
	changed := false
	look.Context.Frame(layout, ui.Look{Fill: look.Scheme.SurfaceContainer}, func() {
		for index, place := range destinations {
			if look.navigationItem(place, index == *selected) && index != *selected {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

func (look Look) navigationItem(place Destination, selected bool) bool {
	scheme := look.Scheme
	pill, ink, label, weight := ui.Look{Radius: 16}, scheme.OnSurfaceVariant, scheme.OnSurfaceVariant, draw2d.Medium
	if selected {
		pill.Fill, ink, label, weight = scheme.SecondaryContainer, scheme.OnSecondaryContainer, scheme.OnSurface, draw2d.Bold
	}
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Share(1), Height: ui.Grow(), Gap: 4, AlignX: ui.Center}
	return look.Context.Clickable("navigate "+place.Label, layout, ui.Looks{}, func() {
		indicator := ui.Layout{Width: ui.Fixed(64), Height: ui.Fixed(32), AlignX: ui.Center, AlignY: ui.Center}
		look.Context.Frame(indicator, pill, func() { look.Context.IconSized(place.Icon, 24, ink) })
		look.Text(place.Label, 12, weight, label)
	})
}

// StatusBarScrim paints the strip under the status bar alone, for a screen whose top is a scene
// rather than a bar: the system's icons stay readable over it.
func (look Look) StatusBarScrim() {
	top := look.Context.Insets().Top
	if top <= 0 {
		return
	}
	look.Context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(top)}, ui.Look{Fill: look.Scheme.Scrim}, nil)
}
