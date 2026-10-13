package material

import (
	"embed"

	"github.com/akrck02/go-renderer/opengl/overlay"
	ui "github.com/akrck02/go-ui"
)

// Theme gives go-ui's own controls (its slider, text fields, menus) the scheme's colours and a size
// for fingers: controls 48 pixels tall, 16-pixel text, pill-shaped corners. Set it as the
// context's theme in every frame, before the pieces.
func (scheme Scheme) Theme(dark bool) ui.Theme {
	theme := ui.LightTheme()
	if dark {
		theme = ui.DarkTheme()
	}
	theme.Background, theme.Surface, theme.Sunken = scheme.Surface, scheme.SurfaceContainerLow, scheme.SurfaceContainerHighest
	theme.Ink, theme.InkSoft, theme.InkQuiet = scheme.OnSurface, scheme.OnSurfaceVariant, scheme.Outline
	theme.Line, theme.ControlLine = scheme.OutlineVariant, scheme.Outline
	theme.Primary, theme.OnPrimary, theme.PrimaryTint = scheme.Primary, scheme.OnPrimary, scheme.SecondaryContainer
	theme.FocusRing, theme.Selection = scheme.Primary, scheme.PrimaryContainer
	theme.TextSize, theme.SmallTextSize, theme.HeadingSize, theme.IconSize = 16, 14, 22, 24
	theme.ControlHeight, theme.Gap, theme.Padding = 48, 12, 16
	theme.ControlRadius, theme.PanelRadius = 20, 16
	return theme
}

//go:embed fonts/Roboto-Regular.ttf fonts/Roboto-Medium.ttf fonts/Roboto-Bold.ttf
var roboto embed.FS

// UseRoboto draws every text in Roboto, Android's own typeface (SIL Open Font License, in fonts/).
// Call it before the window makes its interface, which reads the fonts. A program with a typeface
// of its own gives it to go-renderer's overlay.UseFonts instead.
func UseRoboto() {
	regular, _ := roboto.ReadFile("fonts/Roboto-Regular.ttf") // embedded: they are always there
	medium, _ := roboto.ReadFile("fonts/Roboto-Medium.ttf")
	bold, _ := roboto.ReadFile("fonts/Roboto-Bold.ttf")
	overlay.UseFonts(overlay.Fonts{Regular: regular, Medium: medium, Bold: bold})
}
