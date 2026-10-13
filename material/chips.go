package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// AssistChip suggests an action related to what is on screen (add to the calendar, share), with
// an icon before its label. It returns true when touched.
func (look Look) AssistChip(icon, label string) bool {
	scheme := look.Scheme
	normal := ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}
	layout := ui.Layout{Height: ui.Fixed(32), Padding: ui.Padding{Left: 8, Right: 16}, Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("assist "+label, layout, pressedOnly(normal, stateLayer(scheme.OnSurface, 0.1)), func() {
		if icon != "" {
			look.Context.IconSized(icon, 18, scheme.Primary)
		}
		look.Text(label, 14, draw2d.Medium, scheme.OnSurface)
	})
}

// InputChip is something the person entered or chose (a recipient, a tag), with a cross to take it
// out. It tells whether the chip was touched and whether its cross was.
func (look Look) InputChip(label string) (touched, removed bool) {
	scheme := look.Scheme
	look.Context.Push("input " + label)
	defer look.Context.Pop()
	frame := ui.Layout{Height: ui.Fixed(32), Padding: ui.Padding{Left: 12, Right: 4}, Gap: 4, AlignY: ui.Center}
	look.Context.Frame(frame, ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}, func() {
		layout := ui.Layout{Height: ui.Grow(), AlignY: ui.Center}
		touched = look.Context.Clickable("label", layout, ui.Looks{}, func() {
			look.Text(label, 14, draw2d.Medium, scheme.OnSurfaceVariant)
		})
		cross := ui.Layout{Width: ui.Fixed(24), Height: ui.Fixed(24), AlignX: ui.Center, AlignY: ui.Center}
		removed = look.Context.Clickable("remove", cross, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnSurface, 0.12), Radius: 12}}, func() {
			look.Context.IconSized("x", 16, scheme.OnSurfaceVariant)
		})
	})
	return touched, removed
}

// Slider chooses a value between two ends by dragging, in steps (0 for any value): go-ui's slider
// in the scheme's colours and a finger's size. It returns true when the value changed.
func (look Look) Slider(name string, value *float32, minimum, maximum, step float32) bool {
	context := look.Context
	saved := context.Theme
	defer func() { context.Theme = saved }()
	context.Theme.ControlLine, context.Theme.Sunken = look.Scheme.SecondaryContainer, look.Scheme.SecondaryContainer
	if step <= 0 {
		return context.Slider(name, value, minimum, maximum)
	}
	return context.SliderStep(name, value, minimum, maximum, step)
}
