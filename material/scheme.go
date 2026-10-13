package material

import (
	"github.com/akrck02/go-renderer/draw2d"
)

// Scheme holds the colour roles of Material Design 3: what each piece is painted with. Baseline
// gives Material's own; a program with colours of its own changes the roles it needs (WithPrimary
// for a brand's main colour).
type Scheme struct {
	Primary, OnPrimary, PrimaryContainer, OnPrimaryContainer draw2d.Color
	SecondaryContainer, OnSecondaryContainer                 draw2d.Color
	Error, OnError, ErrorContainer, OnErrorContainer         draw2d.Color
	// Surface is the screen's background; the containers are the surfaces raised over it, from the
	// lowest to the highest (a card, a navigation bar, a switch's track).
	Surface, OnSurface, OnSurfaceVariant                                                 draw2d.Color
	SurfaceContainerLow, SurfaceContainer, SurfaceContainerHigh, SurfaceContainerHighest draw2d.Color
	Outline, OutlineVariant                                                              draw2d.Color
	// The inverse roles paint what stands out against the screen: a snackbar.
	InverseSurface, InverseOnSurface, InversePrimary draw2d.Color
	// Scrim is the surface's colour, see-through: under a bar painted over a scene.
	Scrim draw2d.Color
}

// Baseline is Material Design 3's baseline scheme, light or dark.
func Baseline(dark bool) Scheme {
	if dark {
		return Scheme{
			Primary: draw2d.Hex("#d0bcff"), OnPrimary: draw2d.Hex("#381e72"),
			PrimaryContainer: draw2d.Hex("#4f378b"), OnPrimaryContainer: draw2d.Hex("#eaddff"),
			SecondaryContainer: draw2d.Hex("#4a4458"), OnSecondaryContainer: draw2d.Hex("#e8def8"),
			Error: draw2d.Hex("#f2b8b5"), OnError: draw2d.Hex("#601410"),
			ErrorContainer: draw2d.Hex("#8c1d18"), OnErrorContainer: draw2d.Hex("#f9dedc"),
			Surface: draw2d.Hex("#141218"), OnSurface: draw2d.Hex("#e6e0e9"), OnSurfaceVariant: draw2d.Hex("#cac4d0"),
			SurfaceContainerLow: draw2d.Hex("#1d1b20"), SurfaceContainer: draw2d.Hex("#211f26"),
			SurfaceContainerHigh: draw2d.Hex("#2b2930"), SurfaceContainerHighest: draw2d.Hex("#36343b"),
			Outline: draw2d.Hex("#938f99"), OutlineVariant: draw2d.Hex("#49454f"), Scrim: draw2d.Hex("#141218cc"),
			InverseSurface: draw2d.Hex("#e6e0e9"), InverseOnSurface: draw2d.Hex("#322f35"), InversePrimary: draw2d.Hex("#6750a4"),
		}
	}
	return Scheme{
		Primary: draw2d.Hex("#6750a4"), OnPrimary: draw2d.Hex("#ffffff"),
		PrimaryContainer: draw2d.Hex("#eaddff"), OnPrimaryContainer: draw2d.Hex("#21005d"),
		SecondaryContainer: draw2d.Hex("#e8def8"), OnSecondaryContainer: draw2d.Hex("#1d192b"),
		Error: draw2d.Hex("#b3261e"), OnError: draw2d.Hex("#ffffff"),
		ErrorContainer: draw2d.Hex("#f9dedc"), OnErrorContainer: draw2d.Hex("#410e0b"),
		Surface: draw2d.Hex("#fef7ff"), OnSurface: draw2d.Hex("#1d1b20"), OnSurfaceVariant: draw2d.Hex("#49454f"),
		SurfaceContainerLow: draw2d.Hex("#f7f2fa"), SurfaceContainer: draw2d.Hex("#f3edf7"),
		SurfaceContainerHigh: draw2d.Hex("#ece6f0"), SurfaceContainerHighest: draw2d.Hex("#e6e0e9"),
		Outline: draw2d.Hex("#79747e"), OutlineVariant: draw2d.Hex("#cac4d0"), Scrim: draw2d.Hex("#fef7ffcc"),
		InverseSurface: draw2d.Hex("#322f35"), InverseOnSurface: draw2d.Hex("#f5eff7"), InversePrimary: draw2d.Hex("#d0bcff"),
	}
}

// WithPrimary returns the scheme with a main colour of the program's own: the colour, the text on
// it, a soft container of it and the text on that container. A navigation bar's selected pill and
// a filter chip's fill take the container too.
func (scheme Scheme) WithPrimary(primary, onPrimary, container, onContainer draw2d.Color) Scheme {
	scheme.Primary, scheme.OnPrimary = primary, onPrimary
	scheme.PrimaryContainer, scheme.OnPrimaryContainer = container, onContainer
	scheme.SecondaryContainer, scheme.OnSecondaryContainer = container, onContainer
	return scheme
}
