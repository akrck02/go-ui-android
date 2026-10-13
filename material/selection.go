package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Switch is Material's switch: a 52 × 32 track whose handle grows and moves to the end when on. It
// flips the value and returns true when touched.
func (look Look) Switch(name string, on *bool) bool {
	scheme := look.Scheme
	track := ui.Look{Fill: scheme.SurfaceContainerHighest, Border: scheme.Outline, BorderWidth: 2, Radius: 16}
	handle, handleSize, align := scheme.Outline, float32(16), ui.Start
	if *on {
		track = ui.Look{Fill: scheme.Primary, Radius: 16}
		handle, handleSize, align = scheme.OnPrimary, 24, ui.End
	}
	layout := ui.Layout{Width: ui.Fixed(52), Height: ui.Fixed(32), Padding: ui.PadXY(4, 0), AlignX: align, AlignY: ui.Center}
	touched := look.Context.Clickable(name, layout, ui.Looks{Normal: track, Hovered: track, Pressed: track}, func() {
		size := ui.Fixed(handleSize)
		look.Context.Frame(ui.Layout{Width: size, Height: size}, ui.Look{Fill: handle, Radius: handleSize / 2}, nil)
	})
	if touched {
		*on = !*on
	}
	return touched
}

// FilterChip is a choice that is on or off, among others in a row: outlined, or filled with a tick
// when selected. It returns true when touched; the program flips its own value.
func (look Look) FilterChip(label string, selected bool) bool {
	scheme := look.Scheme
	normal, ink := ui.Look{Border: scheme.OutlineVariant, BorderWidth: 1, Radius: 8}, scheme.OnSurfaceVariant
	if selected {
		normal, ink = ui.Look{Fill: scheme.SecondaryContainer, Radius: 8}, scheme.OnSecondaryContainer
	}
	layout := ui.Layout{Height: ui.Fixed(32), Padding: ui.PadXY(12, 0), Gap: 8, AlignY: ui.Center}
	return look.Context.Clickable("chip "+label, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
		if selected {
			look.Context.IconSized("check", 18, ink)
		}
		look.Text(label, 14, draw2d.Medium, ink)
	})
}

// Tag is a label that cannot be touched, in a pill of its own colours: a state ("3 running"), a
// count, a kind. ink is the text's colour and tint the pill's.
func (look Look) Tag(text string, ink, tint draw2d.Color) {
	look.Context.Frame(ui.Layout{Height: ui.Fixed(26), Padding: ui.PadXY(10, 0), AlignY: ui.Center},
		ui.Look{Fill: tint, Radius: 13}, func() { look.Text(text, 12, draw2d.Bold, ink) })
}

// SegmentedButtons are two to five choices side by side sharing the row's width, one of them
// chosen (a view, a speed, a period). It changes selected and returns true when another is
// touched. name tells the group apart from another one on the same screen.
func (look Look) SegmentedButtons(name string, options []string, selected *int) bool {
	scheme := look.Scheme
	changed := false
	look.Context.Push(name)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: 8}, func() {
		for index, option := range options {
			chosen := index == *selected
			normal, ink := ui.Look{Border: scheme.Outline, BorderWidth: 1, Radius: 20}, scheme.OnSurface
			if chosen {
				normal, ink = ui.Look{Fill: scheme.SecondaryContainer, Border: scheme.Outline, BorderWidth: 1, Radius: 20}, scheme.OnSecondaryContainer
			}
			layout := ui.Layout{Width: ui.Share(1), Height: ui.Fixed(40), Gap: 6, AlignX: ui.Center, AlignY: ui.Center}
			touched := look.Context.Clickable(option, layout, pressedOnly(normal, stateLayer(ink, 0.12)), func() {
				if chosen {
					look.Context.IconSized("check", 18, ink)
				}
				look.Text(option, 14, draw2d.Medium, ink)
			})
			if touched && !chosen {
				*selected, changed = index, true
			}
		}
	})
	return changed
}
