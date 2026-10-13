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
// edge as Android's own apps do. Inside Body, it keeps clear of a cutout at the sides too. It returns true when the back arrow is touched.
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
		Padding: ui.Padding{Left: insets.Left, Top: 12, Right: insets.Right, Bottom: 16 + insets.Bottom}}
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

// NavigationRail is the navigation bar's form on a wide screen (a phone on its side, a tablet): the
// destinations down the left edge, painted under the status bar and beside a camera's cutout. It
// changes selected and returns true when another destination is touched. Wide tells when to use it.
func (look Look) NavigationRail(destinations []Destination, selected *int) bool {
	insets := look.Context.Insets()
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Fixed(80 + insets.Left), Height: ui.Grow(),
		Padding: ui.Padding{Left: insets.Left, Top: insets.Top + 12, Bottom: insets.Bottom + 12}, Gap: 12}
	changed := false
	look.Context.Frame(layout, ui.Look{Fill: look.Scheme.SurfaceContainer}, func() {
		for index, place := range destinations {
			if look.railItem(place, index == *selected) && index != *selected {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

func (look Look) railItem(place Destination, selected bool) bool {
	scheme := look.Scheme
	pill, ink, label, weight := ui.Look{Radius: 16}, scheme.OnSurfaceVariant, scheme.OnSurfaceVariant, draw2d.Medium
	if selected {
		pill.Fill, ink, label, weight = scheme.SecondaryContainer, scheme.OnSecondaryContainer, scheme.OnSurface, draw2d.Bold
	}
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Fixed(56), Gap: 4, AlignX: ui.Center}
	return look.Context.Clickable("navigate "+place.Label, layout, ui.Looks{}, func() {
		indicator := ui.Layout{Width: ui.Fixed(56), Height: ui.Fixed(32), AlignX: ui.Center, AlignY: ui.Center}
		look.Context.Frame(indicator, pill, func() { look.Context.IconSized(place.Icon, 24, ink) })
		look.Text(place.Label, 12, weight, label)
	})
}

// Wide tells whether the screen is wide enough for a navigation rail instead of a bar: wider than
// tall, or 600 pixels across (Material's medium width, a small tablet).
func (look Look) Wide() bool {
	width, height := look.Context.Size()
	return width > height || width >= 600
}

// Body holds a screen's content between the bars: it grows to fill them and keeps its content
// clear of what the system draws at the sides (a camera's cutout, a navigation bar on the side).
// Beside a navigation rail (besideRail), the rail keeps the left side clear, and with no
// navigation bar under it, Body keeps clear of the gesture bar itself.
func (look Look) Body(besideRail bool, content func()) {
	insets := look.Context.Insets()
	padding := ui.Padding{Left: insets.Left, Right: insets.Right}
	if besideRail {
		padding.Left, padding.Bottom = 0, insets.Bottom
	}
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: padding}, content)
}
