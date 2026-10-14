package material

import (
	"math"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-ui/window"
)

// FromColour makes a whole scheme from one colour (a brand's, a place's), as Material's own tools
// do: palettes of its hue at every tone, the way Android 12 and later make them from a wallpaper
// (FromPalettes), then the roles from them.
func FromColour(seed draw2d.Color, dark bool) Scheme {
	return FromPalettes(PalettesFromColour(seed), dark)
}

// PalettesFromColour makes the five palettes Android makes from a wallpaper, from one colour:
// the main one of its hue, a quieter one for the containers, a third of the hue turned a sixth of
// the way round, and two near-greys of its hue for the surfaces. The chromas are Material's
// "tonal spot" ones, in OKLCH's units.
func PalettesFromColour(seed draw2d.Color) window.Palettes {
	_, _, hue := oklch(seed)
	return window.Palettes{
		Accent1:  tonalPalette(hue, 0.115),
		Accent2:  tonalPalette(hue, 0.05),
		Accent3:  tonalPalette(hue+math.Pi/3, 0.077),
		Neutral1: tonalPalette(hue, 0.016),
		Neutral2: tonalPalette(hue, 0.026),
	}
}

// tonalPalette is a hue at a chroma, at each of Android's shades: shade s is Material's tone
// 100 - s/10, the colour's lightness on CIELAB's scale (L*).
func tonalPalette(hue, chroma float64) [13]draw2d.Color {
	var palette [13]draw2d.Color
	for index, shade := range window.PaletteTones {
		palette[index] = colourAtTone(hue, chroma, 100-float64(shade)/10)
	}
	return palette
}

// colourAtTone finds the colour of a hue and chroma whose lightness is a tone: the OKLab lightness
// giving that L*, with the chroma lowered as far as needed to stay a colour a screen shows.
func colourAtTone(hue, chroma, tone float64) draw2d.Color {
	if tone >= 100 {
		return draw2d.Color{1, 1, 1, 1}
	}
	if tone <= 0 {
		return draw2d.Color{0, 0, 0, 1}
	}
	low, high := 0.0, 1.0
	var red, green, blue float64
	for range 30 {
		lightness := (low + high) / 2
		red, green, blue = inGamut(lightness, chroma, hue)
		if cielabLightness(red, green, blue) < tone {
			low = lightness
		} else {
			high = lightness
		}
	}
	return draw2d.Color{float32(encodeGamma(red)), float32(encodeGamma(green)), float32(encodeGamma(blue)), 1}
}

// inGamut is the linear colour of an OKLCH one, its chroma lowered until it fits sRGB.
func inGamut(lightness, chroma, hue float64) (red, green, blue float64) {
	red, green, blue = fromOKLab(lightness, chroma*math.Cos(hue), chroma*math.Sin(hue))
	if fits(red, green, blue) {
		return red, green, blue
	}
	low, high := 0.0, chroma
	for range 20 {
		middle := (low + high) / 2
		if candidateRed, candidateGreen, candidateBlue := fromOKLab(lightness, middle*math.Cos(hue), middle*math.Sin(hue)); fits(candidateRed, candidateGreen, candidateBlue) {
			low = middle
		} else {
			high = middle
		}
	}
	return fromOKLab(lightness, low*math.Cos(hue), low*math.Sin(hue))
}

func fits(red, green, blue float64) bool {
	const margin = 1e-4
	return red >= -margin && red <= 1+margin && green >= -margin && green <= 1+margin && blue >= -margin && blue <= 1+margin
}

// oklch is a colour's OKLCH lightness, chroma and hue (in radians).
func oklch(colour draw2d.Color) (lightness, chroma, hue float64) {
	lightness, a, b := toOKLab(decodeGamma(float64(colour[0])), decodeGamma(float64(colour[1])), decodeGamma(float64(colour[2])))
	return lightness, math.Hypot(a, b), math.Atan2(b, a)
}

// toOKLab and fromOKLab turn linear sRGB into OKLab and back (Björn Ottosson's matrices).
func toOKLab(red, green, blue float64) (lightness, a, b float64) {
	long := math.Cbrt(0.4122214708*red + 0.5363325363*green + 0.0514459929*blue)
	medium := math.Cbrt(0.2119034982*red + 0.6806995451*green + 0.1073969566*blue)
	short := math.Cbrt(0.0883024619*red + 0.2817188376*green + 0.6299787005*blue)
	return 0.2104542553*long + 0.7936177850*medium - 0.0040720468*short,
		1.9779984951*long - 2.4285922050*medium + 0.4505937099*short,
		0.0259040371*long + 0.7827717662*medium - 0.8086757660*short
}

func fromOKLab(lightness, a, b float64) (red, green, blue float64) {
	long := cube(lightness + 0.3963377774*a + 0.2158037573*b)
	medium := cube(lightness - 0.1055613458*a - 0.0638541728*b)
	short := cube(lightness - 0.0894841775*a - 1.2914855480*b)
	return 4.0767416621*long - 3.3077115913*medium + 0.2309699292*short,
		-1.2684380046*long + 2.6097574011*medium - 0.3413193965*short,
		-0.0041960863*long - 0.7034186147*medium + 1.7076147010*short
}

func cube(value float64) float64 { return value * value * value }

// cielabLightness is a linear colour's L*: Material's tone.
func cielabLightness(red, green, blue float64) float64 {
	luminance := 0.2126*red + 0.7152*green + 0.0722*blue
	if luminance <= 216.0/24389 {
		return luminance * 24389.0 / 27
	}
	return 116*math.Cbrt(luminance) - 16
}

func decodeGamma(channel float64) float64 {
	if channel <= 0.04045 {
		return channel / 12.92
	}
	return math.Pow((channel+0.055)/1.055, 2.4)
}

func encodeGamma(channel float64) float64 {
	channel = min(max(channel, 0), 1)
	if channel <= 0.0031308 {
		return channel * 12.92
	}
	return 1.055*math.Pow(channel, 1/2.4) - 0.055
}
