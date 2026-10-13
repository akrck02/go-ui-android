package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
)

func TestABrandsColourTakesThePrimaryRoles(t *testing.T) {
	brand := draw2d.Hex("#00668b")
	scheme := Baseline(false).WithPrimary(brand, draw2d.Hex("#ffffff"), draw2d.Hex("#c4e7ff"), draw2d.Hex("#001e2c"))
	if scheme.Primary != brand || scheme.SecondaryContainer != draw2d.Hex("#c4e7ff") {
		t.Fatalf("scheme %+v: the brand's colour and container are not used", scheme)
	}
	if theme := scheme.Theme(false); theme.ControlHeight != 48 || theme.Primary != brand {
		t.Fatalf("theme: controls %v tall in %v; want 48, in the brand's colour", theme.ControlHeight, theme.Primary)
	}
}

func TestTextOnAMarkIsReadable(t *testing.T) {
	if readableOn(draw2d.Hex("#ffeb3b"))[0] != 0 || readableOn(draw2d.Hex("#1a237e"))[0] != 1 {
		t.Fatal("want dark text on yellow and light text on navy")
	}
}

func TestTheWallpapersColoursMakeTheScheme(t *testing.T) {
	palettes := palettesOf()
	light, dark := FromPalettes(palettes, false), FromPalettes(palettes, true)
	if light.Primary != palettes.Accent1[shade600] || light.PrimaryContainer != palettes.Accent1[shade100] {
		t.Fatalf("light primary %v, container %v: want the first accent's 600 and 100", light.Primary, light.PrimaryContainer)
	}
	if dark.Primary != palettes.Accent1[shade200] || dark.SecondaryContainer != palettes.Accent2[shade700] {
		t.Fatalf("dark primary %v, secondary container %v: want accent 1's 200 and accent 2's 700", dark.Primary, dark.SecondaryContainer)
	}
	// surfaces rise towards the text's colour: each container a little further from the surface
	for _, scheme := range []Scheme{light, dark} {
		steps := []draw2d.Color{scheme.Surface, scheme.SurfaceContainerLow, scheme.SurfaceContainer, scheme.SurfaceContainerHigh, scheme.SurfaceContainerHighest}
		for index := 1; index < len(steps); index++ {
			if distance(steps[index], scheme.Surface) <= distance(steps[index-1], scheme.Surface) {
				t.Fatalf("surface containers do not rise in order: %v", steps)
			}
		}
	}
	if light.Error != Baseline(false).Error {
		t.Fatal("the wallpaper changed the error colour")
	}
}
