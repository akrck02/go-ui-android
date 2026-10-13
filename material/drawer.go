package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// NavigationDrawer slides in from the left over the dimmed screen, with a title and the program's
// destinations, for more places than a navigation bar holds. It changes selected when one is
// touched, and returns true when it should close: a destination chosen, a touch outside, the
// system's back. name tells drawers apart.
func (look Look) NavigationDrawer(name, title string, destinations []Destination, selected *int) (closed bool) {
	scheme := look.Scheme
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	drawerWidth := min(360, width-56)
	arrived := look.Context.Appear("drawer "+name, arrivingMotion)
	options := ui.LayerOptions{AlignX: ui.Start, AlignY: ui.Start, Padding: ui.Padding{Left: -(1 - arrived) * drawerWidth},
		Width: ui.Fixed(drawerWidth), Height: ui.Grow(), Hidden: 1 - arrived}
	dismissed := look.Context.Modal(name, options, func() {
		sheet := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(),
			Padding: ui.Padding{Left: 12 + insets.Left, Top: 12 + insets.Top, Right: 12, Bottom: 12 + insets.Bottom}, Gap: 4}
		look.Context.Frame(sheet, ui.Look{Fill: scheme.SurfaceContainerLow, Radius: 16, Shadow: true}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 16, Bottom: 16}}, func() {
				look.Text(title, 14, draw2d.Medium, scheme.OnSurfaceVariant)
			})
			for index, place := range destinations {
				if look.drawerItem(place, index == *selected) {
					*selected, closed = index, true
				}
			}
		})
	})
	return closed || dismissed
}

func (look Look) drawerItem(place Destination, selected bool) bool {
	scheme := look.Scheme
	normal, ink := ui.Look{Radius: 28}, scheme.OnSurfaceVariant
	if selected {
		normal, ink = ui.Look{Fill: scheme.SecondaryContainer, Radius: 28}, scheme.OnSecondaryContainer
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 24}, Gap: 12, AlignY: ui.Center}
	return look.Context.Clickable("drawer "+place.Label, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
		look.Context.IconSized(place.Icon, 24, ink)
		look.Context.Box(ui.Layout{Width: ui.Grow()}, func() { look.Text(place.Label, 14, draw2d.Medium, ink) })
		if place.Badge != "" {
			look.Text(place.Badge, 14, draw2d.Medium, ink)
		}
	})
}
