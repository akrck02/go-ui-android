package material

import (
	"testing"
	"unicode/utf8"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui/window"
)

// measurer stands in for the fonts: no display needed.
type measurer struct{}

func (measurer) MeasureText(text string, size float32, _ draw2d.Weight, _ float32) draw2d.TextMetrics {
	return draw2d.TextMetrics{Width: float32(utf8.RuneCountInString(text)) * size / 2, Ascent: size * 0.8, Descent: size * 0.2}
}

func (measurer) HasIcon(string) bool { return true }

// screen drives a phone-sized interface frame by frame.
type screen struct {
	t       *testing.T
	context *ui.Context
	state   *input.State
	build   func(look Look)
}

func newScreen(t *testing.T, build func(look Look)) *screen {
	state := input.NewState()
	state.WindowWidth, state.WindowHeight = 400, 800
	state.FramebufferWidth, state.FramebufferHeight = 400, 800
	return &screen{t: t, context: ui.NewWithMeasurer(measurer{}), state: state, build: build}
}

func (screen *screen) frame() {
	screen.context.Begin(screen.state)
	screen.build(New(screen.context, Baseline(false)))
	screen.context.End()
	screen.state.EndFrame()
}

// tap presses and lifts the pointer in the middle of a piece laid out in the last frame.
func (screen *screen) tap(path ...string) {
	screen.t.Helper()
	area, found := screen.context.Area(path...)
	if !found {
		screen.t.Fatalf("nothing called %v was laid out", path)
	}
	screen.state.MoveCursor(float64(area.X+area.Width/2), float64(area.Y+area.Height/2))
	screen.state.SetButton(input.MouseLeft, true)
	screen.frame()
	screen.state.SetButton(input.MouseLeft, false)
	screen.frame()
}

func TestASwitchFlipsWhenTouched(t *testing.T) {
	on, flips := false, 0
	screen := newScreen(t, func(look Look) {
		if look.Switch("sound", &on) {
			flips++
		}
	})
	screen.frame()
	screen.tap("sound")
	if !on || flips != 1 {
		t.Fatalf("on %v after %d flips; want on after one", on, flips)
	}
}

func TestTheNavigationBarGoesWhereItIsTouched(t *testing.T) {
	selected, changes := 0, 0
	places := []Destination{{"home", "Home"}, {"map", "Map"}, {"settings", "Settings"}}
	screen := newScreen(t, func(look Look) {
		if look.NavigationBar(places, &selected) {
			changes++
		}
	})
	screen.frame()
	screen.tap("navigate Map")
	screen.tap("navigate Map")
	if selected != 1 || changes != 1 {
		t.Fatalf("selected %d after %d changes; want Map, changed once", selected, changes)
	}
}

func TestBarsReachUnderTheSystemsBars(t *testing.T) {
	var topHeight, bottomHeight float32
	screen := newScreen(t, func(look Look) {
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.TopAppBar(TopBar{Title: "Inbox"})
			look.Context.Spacer()
			selected := 0
			look.NavigationBar([]Destination{{"home", "Home"}}, &selected)
		})
	})
	screen.context.SetInsets(ui.Insets{Top: 24, Bottom: 48})
	screen.frame()
	if area, found := screen.context.Area("navigate Home"); found {
		bottomHeight = 800 - area.Y
	}
	topHeight = 64 + screen.context.Insets().Top
	if topHeight != 88 || bottomHeight < 48 {
		t.Fatalf("top bar %v tall, navigation reaching %v from the bottom; want 88 and below the gesture bar", topHeight, bottomHeight)
	}
}

func TestTheBackArrowSaysItWasTouched(t *testing.T) {
	wentBack := false
	screen := newScreen(t, func(look Look) {
		if look.TopAppBar(TopBar{Title: "Details", Back: true}) {
			wentBack = true
		}
	})
	screen.frame()
	screen.tap("icon arrow-left")
	if !wentBack {
		t.Fatal("touching the back arrow did not say so")
	}
}

func TestSegmentedButtonsChooseOne(t *testing.T) {
	selected := 0
	screen := newScreen(t, func(look Look) {
		look.SegmentedButtons("speed", []string{"Slow", "Fast"}, &selected)
	})
	screen.frame()
	screen.tap("speed", "Fast")
	if selected != 1 {
		t.Fatalf("selected %d; want Fast", selected)
	}
}

func TestTheStickFollowsItsFingerAndNoOther(t *testing.T) {
	var stick Stick
	screen := newScreen(t, func(look Look) { look.Stick("walk", &stick, 120) })
	screen.frame()
	area, _ := screen.context.Area("walk")
	middleX, middleY := float64(area.X+area.Width/2), float64(area.Y+area.Height/2)

	screen.state.SetTouch(7, 300, 300) // a finger elsewhere does not take it
	screen.frame()
	if stick.Held {
		t.Fatal("a finger away from the stick took it")
	}
	screen.state.SetTouch(8, middleX, middleY) // on it
	screen.frame()
	if !stick.Held || stick.Finger != 8 || stick.Forward != 0 {
		t.Fatalf("held %v by %d, forward %v; want held by the finger on it, in the middle", stick.Held, stick.Finger, stick.Forward)
	}
	screen.state.SetTouch(8, middleX, middleY-200) // dragged far up
	screen.frame()
	if stick.Forward < 0.99 || stick.Sideways != 0 {
		t.Fatalf("forward %v, sideways %v; want all the way ahead", stick.Forward, stick.Sideways)
	}
	if !stick.Touching(input.Touch{Finger: 8}) || stick.Touching(input.Touch{Finger: 7}) {
		t.Fatal("the stick does not tell its finger from the other one")
	}
	screen.state.LiftTouch(8)
	screen.frame()
	if stick.Held || stick.Forward != 0 {
		t.Fatal("letting go does not bring the stick back")
	}
}

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

// palettesOf makes palettes whose shades are greys from white (0) to black (1000), with each accent
// tinted so the roles can be told apart.
func palettesOf() window.Palettes {
	var palettes window.Palettes
	for index := range window.PaletteTones {
		grey := 1 - float32(index)/12
		palettes.Accent1[index] = draw2d.Color{grey, 0, 0, 1}
		palettes.Accent2[index] = draw2d.Color{0, grey, 0, 1}
		palettes.Accent3[index] = draw2d.Color{0, 0, grey, 1}
		palettes.Neutral1[index] = draw2d.Color{grey, grey, grey, 1}
		palettes.Neutral2[index] = draw2d.Color{grey, grey, grey * 0.9, 1}
	}
	return palettes
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

func distance(first, second draw2d.Color) float32 {
	return abs(first[0]-second[0]) + abs(first[1]-second[1]) + abs(first[2]-second[2])
}

func abs(value float32) float32 { return max(value, -value) }

func TestAParagraphBreaksIntoLines(t *testing.T) {
	screen := newScreen(t, func(look Look) {
		look.Context.Box(ui.Layout{Width: ui.Fixed(200)}, func() {
			look.Context.Push("text")
			look.Paragraph("one two three four five six seven eight nine ten eleven twelve", 14, draw2d.Regular, look.Scheme.OnSurface)
			look.Context.Pop()
		})
	})
	screen.frame()
	lines := map[float32]bool{}
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind != draw2d.TextKind {
			continue
		}
		lines[command.Baseline] = true
		if command.Area.X+command.Area.Width > 200.5 {
			t.Fatalf("%q ends at %v, past the 200-pixel box", command.Text, command.Area.X+command.Area.Width)
		}
	}
	if len(lines) < 2 {
		t.Fatalf("%d lines; want the words broken into several in 200 pixels", len(lines))
	}
}

func TestOnItsSideTheScreenUsesARailAndKeepsClearOfTheEdges(t *testing.T) {
	selected := 0
	screen := newScreen(t, func(look Look) {
		if !look.Wide() {
			t.Fatal("a screen wider than tall is not wide")
		}
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.NavigationRail([]Destination{{"home", "Home"}, {"map", "Map"}}, &selected)
			look.Body(true, func() { look.FilledButton("", "Inside") })
		})
	})
	screen.state.WindowWidth, screen.state.WindowHeight = 800, 400
	screen.context.SetInsets(ui.Insets{Top: 24, Bottom: 20, Left: 30, Right: 40})
	screen.frame()
	screen.tap("navigate Map")
	if selected != 1 {
		t.Fatal("the rail did not go to Map")
	}
	rail, _ := screen.context.Area("navigate Home")
	button, _ := screen.context.Area("button Inside")
	if rail.X < 30 || button.X < 30+80 {
		t.Fatalf("rail item at %v, content at %v: want them clear of the 30-pixel cutout and the rail", rail.X, button.X)
	}
}
