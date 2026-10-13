package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Buttons are 40 pixels tall pills with a 48-pixel target around them, as Material's are; the
// icon, when given, goes before the label. An empty icon draws none.

// FilledButton is the screen's main action, in the primary colour.
func (look Look) FilledButton(icon, label string) bool {
	scheme := look.Scheme
	return look.button(icon, label, ui.Look{Fill: scheme.Primary, Radius: 20}, scheme.PrimaryContainer, scheme.OnPrimary, 24)
}

// TonalButton is a main action that should not stand out as much, on the secondary container.
func (look Look) TonalButton(icon, label string) bool {
	scheme := look.Scheme
	normal := ui.Look{Fill: scheme.SecondaryContainer, Radius: 20}
	return look.button(icon, label, normal, stateLayer(scheme.OnSecondaryContainer, 0.16), scheme.OnSecondaryContainer, 24)
}

// OutlinedButton is a second action beside the main one.
func (look Look) OutlinedButton(icon, label string) bool {
	scheme := look.Scheme
	normal := ui.Look{Border: scheme.Outline, BorderWidth: 1, Radius: 20}
	return look.button(icon, label, normal, stateLayer(scheme.Primary, 0.12), scheme.Primary, 24)
}

// TextButton is the quietest action: a dialog's "Cancel", a card's "More".
func (look Look) TextButton(icon, label string) bool {
	scheme := look.Scheme
	return look.button(icon, label, ui.Look{Radius: 20}, stateLayer(scheme.Primary, 0.12), scheme.Primary, 12)
}

func (look Look) button(icon, label string, normal ui.Look, pressedFill, ink draw2d.Color, padding float32) bool {
	layout := ui.Layout{Height: ui.Fixed(40), Padding: ui.PadXY(padding, 0), Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("button "+label, layout, pressedOnly(normal, pressedFill), func() {
		if icon != "" {
			look.Context.IconSized(icon, 18, ink)
		}
		look.Text(label, 14, draw2d.Medium, ink)
	})
}

// IconButton is a round 48-pixel target with a 24-pixel icon, in the bars and beside list items.
func (look Look) IconButton(icon string) bool {
	ink := look.Scheme.OnSurfaceVariant
	layout := ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}
	looks := ui.Looks{Pressed: ui.Look{Fill: stateLayer(ink, 0.12), Radius: 24}}
	return look.Context.Clickable("icon "+icon, layout, looks, func() { look.Context.IconSized(icon, 24, ink) })
}

// FloatingButton is the floating action button: 56 pixels, rounded square corners, raised over
// what scrolls beneath it. Place it at the end of the screen, above the navigation bar.
func (look Look) FloatingButton(icon string) bool {
	scheme := look.Scheme
	normal := ui.Look{Fill: scheme.PrimaryContainer, Radius: 16, Shadow: true}
	layout := ui.Layout{Width: ui.Fixed(56), Height: ui.Fixed(56), AlignX: ui.Center, AlignY: ui.Center}
	return look.Context.Clickable("floating "+icon, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(icon, 24, scheme.OnPrimaryContainer)
	})
}

// RoundButton is a round floating button of any size, over a map or a scene (zoom in, zoom out,
// change the view). primary fills it with the main colour; otherwise it is a raised surface.
func (look Look) RoundButton(icon string, size float32, primary bool) bool {
	scheme := look.Scheme
	normal := ui.Look{Fill: scheme.SurfaceContainerLow, Border: scheme.OutlineVariant, BorderWidth: 1, Radius: size / 2, Shadow: true}
	ink := scheme.OnSurface
	if primary {
		normal, ink = ui.Look{Fill: scheme.Primary, Radius: size / 2, Shadow: true}, scheme.OnPrimary
	}
	layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), AlignX: ui.Center, AlignY: ui.Center}
	return look.Context.Clickable("round "+icon, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
		look.Context.IconSized(icon, size*0.45, ink)
	})
}
