package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Card is a filled container for one subject: a few lines, a picture, its actions.
func (look Look) Card(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Scheme.SurfaceContainerHighest, Radius: 12}, content)
}

// ElevatedCard is a card raised over what is beneath it: over a map or a scene.
func (look Look) ElevatedCard(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Scheme.SurfaceContainerLow, Radius: 12, Shadow: true}, content)
}

// OutlinedCard is a card drawn by its border, among others on the surface.
func (look Look) OutlinedCard(content func()) {
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}
	look.Context.Frame(layout, ui.Look{Fill: look.Scheme.Surface, Border: look.Scheme.OutlineVariant, BorderWidth: 1, Radius: 12}, content)
}

// Item is what a list item shows.
type Item struct {
	Icon       string       // in a round avatar at the start; none when empty
	Mark       draw2d.Color // the avatar's colour; the primary container when zero
	Headline   string
	Supporting string // the line under the headline; none when empty
	// Trailing draws what goes at the end (a Switch, an IconButton, a Tag); none when nil.
	Trailing func()
}

// ListItem is a row of a list, 72 pixels tall with a supporting line and 56 without. It returns
// true when the row is touched. name tells it apart from the other rows.
func (look Look) ListItem(name string, item Item) bool {
	scheme := look.Scheme
	height := float32(56)
	if item.Supporting != "" {
		height = 72
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(height), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}
	looks := ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnSurface, 0.1)}}
	return look.Context.Clickable(name, layout, looks, func() {
		if item.Icon != "" {
			look.avatar(item)
		}
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 2}, func() {
			look.Text(item.Headline, 16, draw2d.Regular, scheme.OnSurface)
			if item.Supporting != "" {
				look.Text(item.Supporting, 14, draw2d.Regular, scheme.OnSurfaceVariant)
			}
		})
		if item.Trailing != nil {
			item.Trailing()
		}
	})
}

func (look Look) avatar(item Item) {
	fill, ink := look.Scheme.PrimaryContainer, look.Scheme.OnPrimaryContainer
	if item.Mark != (draw2d.Color{}) {
		fill, ink = item.Mark, readableOn(item.Mark)
	}
	layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
	look.Context.Frame(layout, ui.Look{Fill: fill, Radius: 20}, func() { look.Context.IconSized(item.Icon, 20, ink) })
}

// readableOn is black or white, whichever reads better on a colour.
func readableOn(colour draw2d.Color) draw2d.Color {
	lightness := 0.299*colour[0] + 0.587*colour[1] + 0.114*colour[2]
	if lightness > 0.6 {
		return draw2d.Color{0, 0, 0, 0.87}
	}
	return draw2d.Color{1, 1, 1, 1}
}

// SectionTitle is a small heading over a group of rows, in the primary colour.
func (look Look) SectionTitle(text string) {
	look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 16, Right: 16, Bottom: 8}}, func() {
		look.Text(text, 14, draw2d.Medium, look.Scheme.Primary)
	})
}

// Divider is a thin line across, between groups.
func (look Look) Divider() {
	look.Context.Frame(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(1)}, ui.Look{Fill: look.Scheme.OutlineVariant}, nil)
}
