// Command gallery shows every piece of the material package on three screens: the buttons and
// choices, a list, and a stick that sails a boat around a bay. The same code runs on a phone, on
// the desktop and in a web page:
//
//	go run ./cmd/gallery                                                   # a phone-sized window
//	go run ./cmd/go-ui-apk -name Gallery -id io.github.akrck02.gallery ./cmd/gallery
//	GOOS=js GOARCH=wasm go build -o gallery.wasm ./cmd/gallery             # a web page
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui-android/material"
	"github.com/akrck02/go-ui/window"
)

// The screens of the navigation bar.
const (
	piecesScreen = iota
	listScreen
	stickScreen
)

// lighthouse is a row of the list: the program's own data.
type lighthouse struct {
	name, place string
	lit         bool
}

// gallery is everything the program remembers between frames.
type gallery struct {
	screen      int
	dark        bool
	themeSet    bool // chosen with the switch; until then it follows the system's
	filters     [3]bool
	view        int
	lighthouses []lighthouse
	added       int

	stick    material.Stick
	boatX    float32 // the boat's place in the bay, from 0 to 1 across and down
	boatY    float32
	lastStep time.Time
}

func main() {
	options := window.Options{Title: "Gallery", Width: 412, Height: 892}
	options.RegisterFlags(flag.CommandLine)
	flag.Parse()
	program := &gallery{dark: options.Dark, themeSet: options.Dark, boatX: 0.5, boatY: 0.5, lastStep: time.Now(),
		lighthouses: []lighthouse{{"North Cape", "Cliffs, north shore", true}, {"Harbour Light", "End of the pier", true},
			{"Seal Rock", "Two miles out", false}, {"Old Tower", "The hill above the town", false}}}
	material.UseRoboto() // before the window makes its interface, which reads the fonts
	if err := window.Run(options, program.build); err != nil {
		log.Fatal(err)
	}
}

func (program *gallery) isDark() bool {
	if program.themeSet {
		return program.dark
	}
	return window.SystemPrefersDark()
}

// build describes the interface for this frame.
func (program *gallery) build(context *ui.Context) {
	dark := program.isDark()
	scheme := material.System(dark) // the wallpaper's colours on Android 12 and later
	context.Theme = scheme.Theme(dark)
	look := material.New(context, scheme)
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		switch program.screen {
		case piecesScreen:
			program.pieces(look)
		case listScreen:
			program.list(look)
		default:
			program.bay(look)
		}
		look.NavigationBar([]material.Destination{{Icon: "layout-grid", Label: "Pieces"}, {Icon: "list", Label: "List"},
			{Icon: "gamepad-2", Label: "Stick"}}, &program.screen)
	})
}

// pieces shows the buttons, the choices and a card.
func (program *gallery) pieces(look material.Look) {
	context, scheme := look.Context, look.Scheme
	look.TopAppBar(material.TopBar{Title: "Pieces", Subtitle: "Material Design 3, drawn by go-ui",
		Actions: func() { look.IconButton("search") }})
	context.ScrollArea("pieces", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		look.SectionTitle("Buttons")
		row(context, func() {
			look.FilledButton("plus", "Filled")
			look.TonalButton("", "Tonal")
		})
		row(context, func() {
			look.OutlinedButton("", "Outlined")
			look.TextButton("", "Text")
		})
		look.SectionTitle("Choices")
		row(context, func() {
			for index, label := range []string{"Lit", "Dark", "Far out"} {
				if look.FilterChip(label, program.filters[index]) {
					program.filters[index] = !program.filters[index]
				}
			}
		})
		context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 8)}, func() {
			look.SegmentedButtons("view", []string{"Day", "Week", "Month"}, &program.view)
		})
		look.ListItem("dark theme", material.Item{Icon: "moon", Headline: "Dark theme",
			Supporting: "Follows the system until chosen here", Trailing: func() {
				dark := program.isDark()
				if look.Switch("dark switch", &dark) {
					program.dark, program.themeSet = dark, true
				}
			}})
		look.SectionTitle("Cards and tags")
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.PadXY(16, 0), Gap: 12}, func() {
			look.Card(func() {
				look.Text("Tonight's watch", 18, draw2d.Medium, scheme.OnSurface)
				look.Paragraph("Two lighthouses lit, two dark. The wind turns west after midnight.", 14, draw2d.Regular, scheme.OnSurfaceVariant)
				context.Box(ui.Layout{Gap: 8}, func() {
					look.Tag("2 lit", scheme.OnPrimaryContainer, scheme.PrimaryContainer)
					look.Tag("2 dark", scheme.OnErrorContainer, scheme.ErrorContainer)
				})
			})
			look.OutlinedCard(func() {
				look.Text("An outlined card", 16, draw2d.Medium, scheme.OnSurface)
			})
		})
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(24)}, nil)
	})
}

// list shows the lighthouses, each with a switch, and a floating button that adds one.
func (program *gallery) list(look material.Look) {
	context := look.Context
	look.TopAppBar(material.TopBar{Title: "Lighthouses", Subtitle: fmt.Sprintf("%d on the island", len(program.lighthouses))})
	context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
		context.ScrollArea("lighthouses", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
			for index := range program.lighthouses {
				light := &program.lighthouses[index]
				context.Push(light.name)
				look.ListItem("row", material.Item{Icon: "lightbulb", Headline: light.name, Supporting: light.place,
					Trailing: func() { look.Switch("lit", &light.lit) }})
				context.Pop()
			}
		})
	})
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Pad(16), AlignX: ui.End}, func() {
		if look.FloatingButton("plus") {
			program.added++
			program.lighthouses = append(program.lighthouses, lighthouse{name: fmt.Sprintf("New light %d", program.added), place: "Not placed yet"})
		}
	})
}

// bay sails a boat with the stick: the boat keeps moving while the stick is held.
func (program *gallery) bay(look material.Look) {
	context, scheme := look.Context, look.Scheme
	now := time.Now()
	step := float32(now.Sub(program.lastStep).Seconds())
	program.lastStep = now
	program.boatX = clamp(program.boatX+program.stick.Sideways*step*0.4, 0, 1)
	program.boatY = clamp(program.boatY-program.stick.Forward*step*0.4, 0, 1)

	look.TopAppBar(material.TopBar{Title: "Stick", Subtitle: "Hold the circle and push"})
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Pad(16), Gap: 16}, func() {
		width, _ := context.Size()
		bay := width - 32
		boat := float32(28)
		padding := ui.Padding{Left: program.boatX * (bay - boat), Top: program.boatY * (bay - boat)}
		context.Frame(ui.Layout{Width: ui.Fixed(bay), Height: ui.Fixed(bay), Padding: padding},
			ui.Look{Fill: scheme.PrimaryContainer, Radius: 16}, func() {
				context.Frame(ui.Layout{Width: ui.Fixed(boat), Height: ui.Fixed(boat), AlignX: ui.Center, AlignY: ui.Center},
					ui.Look{Fill: scheme.Primary, Radius: boat / 2}, func() { context.IconSized("sailboat", 18, scheme.OnPrimary) })
			})
		context.Spacer()
	})
	// a stick is found by its name on the screen: keep it out of scopes and scroll areas
	look.Stick("stick", &program.stick, 128)
	context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(16)}, nil)
}

// row lays pieces side by side with Material's 8-pixel gap, inside the screen's margins.
func row(context *ui.Context, content func()) {
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 4), Gap: 8, Wrap: true, AlignY: ui.Center}, content)
}

func clamp(value, minimum, maximum float32) float32 {
	return min(max(value, minimum), maximum)
}
