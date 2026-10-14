package material

import (
	"math"
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

// tone is a colour's Material tone, its L*.
func tone(colour draw2d.Color) float64 {
	return cielabLightness(decodeGamma(float64(colour[0])), decodeGamma(float64(colour[1])), decodeGamma(float64(colour[2])))
}

func TestEveryShadeHasItsTone(t *testing.T) {
	palette := tonalPalette(2.2, 0.115) // a blue-green
	for index, shade := range []int{0, 10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000} {
		want := 100 - float64(shade)/10
		if got := tone(palette[index]); math.Abs(got-want) > 0.6 {
			t.Errorf("shade %d has tone %.1f; want %.0f", shade, got, want)
		}
	}
}

func TestAColoursSchemeKeepsItsHue(t *testing.T) {
	seed := draw2d.Hex("#e8743b") // an orange
	_, _, seedHue := oklch(seed)
	for _, dark := range []bool{false, true} {
		scheme := FromColour(seed, dark)
		_, chroma, hue := oklch(scheme.Primary)
		if difference := math.Abs(math.Remainder(hue-seedHue, 2*math.Pi)); difference > 0.15 || chroma < 0.05 {
			t.Errorf("dark %v: the primary's hue is %.2f, %.2f from the seed's, chroma %.3f; want the seed's hue, coloured", dark, hue, difference, chroma)
		}
	}
}

func TestAColoursSchemeReadsWell(t *testing.T) {
	for _, seed := range []draw2d.Color{draw2d.Hex("#1a73e8"), draw2d.Hex("#2e7d32"), draw2d.Hex("#c2185b"), draw2d.Hex("#ffeb3b")} {
		for _, dark := range []bool{false, true} {
			scheme := FromColour(seed, dark)
			for _, pair := range [][2]draw2d.Color{
				{scheme.Primary, scheme.OnPrimary},
				{scheme.PrimaryContainer, scheme.OnPrimaryContainer},
				{scheme.Surface, scheme.OnSurface},
			} {
				if contrast := math.Abs(tone(pair[0]) - tone(pair[1])); contrast < 40 {
					t.Errorf("seed %v, dark %v: tones %.0f and %.0f too close to read", seed, dark, tone(pair[0]), tone(pair[1]))
				}
			}
		}
	}
}
