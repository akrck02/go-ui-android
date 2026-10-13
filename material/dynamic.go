package material

import (
	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-ui/window"
)

// System is the phone's own scheme: on Android 12 and later, the colours it made from the
// wallpaper (Material You), so the program matches the system's applications and changes with the
// wallpaper; elsewhere, Material's baseline. Call it in every frame.
func System(dark bool) Scheme {
	palettes, found := window.SystemPalettes()
	if !found {
		return Baseline(dark)
	}
	return FromPalettes(palettes, dark)
}

// The shades of a system palette, by Android's numbers (window.PaletteTones): 0 is white, 1000
// black. Material's tone T is Android's shade (100 - T) × 10.
const (
	shade0 = iota
	shade10
	shade50
	shade100
	shade200
	shade300
	shade400
	shade500
	shade600
	shade700
	shade800
	shade900
	shade1000
)

// FromPalettes builds the scheme from a phone's palettes, as Android's own dynamic colours do: the
// main colour from the first accent palette, the containers from the second, the surfaces from the
// neutral ones. The error colours stay Material's, which a wallpaper does not change.
func FromPalettes(palettes window.Palettes, dark bool) Scheme {
	accent, second, neutral, variant := palettes.Accent1, palettes.Accent2, palettes.Neutral1, palettes.Neutral2
	scheme := Baseline(dark)
	if dark {
		scheme.Primary, scheme.OnPrimary = accent[shade200], accent[shade800]
		scheme.PrimaryContainer, scheme.OnPrimaryContainer = accent[shade700], accent[shade100]
		scheme.SecondaryContainer, scheme.OnSecondaryContainer = second[shade700], second[shade100]
		scheme.Surface = mix(neutral[shade900], neutral[shade1000], 0.4) // tone 6
		scheme.OnSurface, scheme.OnSurfaceVariant = neutral[shade100], variant[shade200]
		scheme.SurfaceContainerLow = neutral[shade900]                               // tone 10
		scheme.SurfaceContainer = mix(neutral[shade900], neutral[shade800], 0.2)     // tone 12
		scheme.SurfaceContainerHigh = mix(neutral[shade900], neutral[shade800], 0.7) // tone 17
		scheme.SurfaceContainerHighest = mix(neutral[shade800], neutral[shade700], 0.2)
		scheme.Outline, scheme.OutlineVariant = variant[shade400], variant[shade700]
		scheme.InverseSurface, scheme.InverseOnSurface, scheme.InversePrimary = neutral[shade100], neutral[shade800], accent[shade600]
	} else {
		scheme.Primary, scheme.OnPrimary = accent[shade600], accent[shade0]
		scheme.PrimaryContainer, scheme.OnPrimaryContainer = accent[shade100], accent[shade900]
		scheme.SecondaryContainer, scheme.OnSecondaryContainer = second[shade100], second[shade900]
		scheme.Surface = mix(neutral[shade10], neutral[shade50], 0.25) // tone 98
		scheme.OnSurface, scheme.OnSurfaceVariant = neutral[shade900], variant[shade700]
		scheme.SurfaceContainerLow = mix(neutral[shade10], neutral[shade50], 0.75)  // tone 96
		scheme.SurfaceContainer = mix(neutral[shade50], neutral[shade100], 0.2)     // tone 94
		scheme.SurfaceContainerHigh = mix(neutral[shade50], neutral[shade100], 0.6) // tone 92
		scheme.SurfaceContainerHighest = neutral[shade100]                          // tone 90
		scheme.Outline, scheme.OutlineVariant = variant[shade500], variant[shade200]
		scheme.InverseSurface, scheme.InverseOnSurface, scheme.InversePrimary = neutral[shade800], neutral[shade50], accent[shade200]
	}
	scheme.Scrim = draw2d.Color{scheme.Surface[0], scheme.Surface[1], scheme.Surface[2], 0.8}
	return scheme
}

// mix is a colour a share of the way from one to another.
func mix(from, to draw2d.Color, share float32) draw2d.Color {
	var mixed draw2d.Color
	for channel := range mixed {
		mixed[channel] = from[channel] + (to[channel]-from[channel])*share
	}
	return mixed
}
