package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Field is what a text field says around what is typed.
type Field struct {
	Label       string // above the box
	Placeholder string // in the box, quiet, while it is empty
	Supporting  string // under the box: a hint, how many characters are left
	Error       string // under the box instead of the hint, in red, with the box red too
	MaxLength   int    // the most characters it takes; 0 = no limit
}

// TextField is Material's outlined text field: a label, the box typed in, and a hint or an error
// under it. On a phone the on-screen keyboard comes up while it is typed in. It returns true when
// the text changed. name tells fields apart.
func (look Look) TextField(name string, text *string, field Field) bool {
	scheme := look.Scheme
	changed := false
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 4}, func() {
		labelColour := scheme.OnSurfaceVariant
		if field.Error != "" {
			labelColour = scheme.Error
		}
		if field.Label != "" {
			look.Text(field.Label, 12, draw2d.Medium, labelColour)
		}
		look.withFieldTheme(field.Error != "", 4, scheme.Surface, func() {
			changed = look.Context.TextFieldWith(name, text, ui.TextFieldOptions{Placeholder: field.Placeholder, MaxLength: field.MaxLength})
		})
		switch {
		case field.Error != "":
			look.Paragraph(field.Error, 12, draw2d.Regular, scheme.Error)
		case field.Supporting != "":
			look.Paragraph(field.Supporting, 12, draw2d.Regular, scheme.OnSurfaceVariant)
		}
	})
	return changed
}

// SearchBar is a pill to search in, at the top of a screen. It returns true when the text changed.
func (look Look) SearchBar(name string, text *string, placeholder string) bool {
	scheme := look.Scheme
	changed := false
	bar := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 8}, Gap: 8, AlignY: ui.Center}
	look.Context.Frame(bar, ui.Look{Fill: scheme.SurfaceContainerHigh, Radius: 28}, func() {
		look.Context.IconSized("search", 24, scheme.OnSurfaceVariant)
		// the pill is the field's box: the field draws none of its own, nor a focus ring inside it
		changed = look.Context.TextFieldWith(name, text, ui.TextFieldOptions{Placeholder: placeholder, Plain: true})
		if *text != "" && look.IconButton("x") {
			*text, changed = "", true
		}
	})
	return changed
}

// withFieldTheme gives go-ui's text field, which reads the theme as it is laid out, Material's
// field look for the content's time: its corners, its fill, a red border for an error.
func (look Look) withFieldTheme(failed bool, radius float32, fill draw2d.Color, content func()) {
	context := look.Context
	saved := context.Theme
	defer func() { context.Theme = saved }()
	context.Theme.ControlRadius, context.Theme.ControlHeight = radius, 56
	context.Theme.Surface = fill
	context.Theme.FocusRing = draw2d.Transparent // Material shows the focus by the border alone
	if failed {
		context.Theme.ControlLine, context.Theme.Primary, context.Theme.InkQuiet = look.Scheme.Error, look.Scheme.Error, look.Scheme.Error
	}
	content()
}
