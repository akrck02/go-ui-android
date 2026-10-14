// Command gallery shows every piece of the material package, on one screen for each part of
// Material's catalogue: actions, inputs, lists, feedback, and tabs over pages. The same code runs
// on a phone, on the desktop and in a web page:
//
//	go run ./cmd/gallery                                                   # a phone-sized window
//	go run ./cmd/go-ui-apk -name Gallery -id io.github.akrck02.gallery -notifications \
//	    -icon-foreground cmd/gallery/icon/foreground.png -icon-background '#6750a4' \
//	    -icon-monochrome cmd/gallery/icon/monochrome.png -splash '#fef7ff' -splash-dark '#141218' ./cmd/gallery
//	GOOS=js GOARCH=wasm go build -o gallery.wasm ./cmd/gallery             # a web page
//
// Each screen is a file of its own: actions.go, inputs.go, lists.go, feedback.go and tabs.go;
// overlays.go draws what goes over them.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui-android/material"
	"github.com/akrck02/go-ui/window"
)

// The screens, in the navigation bar's order.
const (
	actionsScreen = iota
	inputsScreen
	listsScreen
	feedbackScreen
	tabsScreen
)

// lighthouse is a row of the list: the program's own data.
type lighthouse struct {
	name, place string
	lit         bool
	chosen      bool // held, or touched while choosing, on the lists screen
}

// gallery is everything the program remembers between frames.
type gallery struct {
	screen     int
	dark       bool
	themeSet   bool // chosen in the inputs or the menu; until then it follows the system's
	colour     int  // the scheme's source: 0 the wallpaper, then one of seeds
	drawerOpen bool

	actions  actionsState
	inputs   inputsState
	lists    listsState
	feedback feedbackState
	tabs     tabsState
	overlays overlaysState
}

func main() {
	options := window.Options{Title: "Gallery", Width: 412, Height: 892}
	options.RegisterFlags(flag.CommandLine)
	screen := flag.Int("screen", actionsScreen, "the screen to open on, from 0 (Actions) to 4 (Tabs): for captures")
	flag.Parse()
	program := &gallery{dark: options.Dark, themeSet: options.Dark}
	program.inputs = newInputs()
	program.lists = newLists()
	program.tabs = newTabs()
	program.restore()
	if wasSet("screen") {
		program.screen = min(max(*screen, actionsScreen), tabsScreen)
	}
	if wasSet("dark") {
		program.dark, program.themeSet = options.Dark, true
	}
	options.Hidden = program.keep // a phone may end the gallery once it is out of sight
	material.UseRoboto()          // before the window makes its interface, which reads the fonts
	if err := window.Run(options, program.build); err != nil {
		log.Fatal(err)
	}
}

// wasSet tells whether a flag was given on the command line.
func wasSet(name string) bool {
	given := false
	flag.Visit(func(set *flag.Flag) { given = given || set.Name == name })
	return given
}

// keep keeps the screen, the theme and the colours chosen, for the next run.
func (program *gallery) keep() {
	theme := ""
	if program.themeSet {
		theme = map[bool]string{false: "light", true: "dark"}[program.dark]
	}
	for name, value := range map[string]string{"screen": strconv.Itoa(program.screen), "theme": theme, "colour": strconv.Itoa(program.colour)} {
		if err := window.KeepPreference(name, value); err != nil {
			log.Print(err)
			return
		}
	}
}

// restore takes up where the last run was left.
func (program *gallery) restore() {
	if screen, err := strconv.Atoi(preference("screen")); err == nil {
		program.screen = min(max(screen, actionsScreen), tabsScreen)
	}
	if colour, err := strconv.Atoi(preference("colour")); err == nil {
		program.colour = min(max(colour, 0), len(seeds))
	}
	switch preference("theme") {
	case "light":
		program.setDark(false)
	case "dark":
		program.setDark(true)
	}
}

func preference(name string) string {
	value, _ := window.Preference(name)
	return value
}

func (program *gallery) isDark() bool {
	if program.themeSet {
		return program.dark
	}
	return window.SystemPrefersDark()
}

func (program *gallery) setDark(dark bool) { program.dark, program.themeSet = dark, true }

// projectPage is where the gallery's pieces come from.
const projectPage = "https://github.com/akrck02/go-ui-android"

// seeds are the colours a scheme can be made from, besides the wallpaper.
var seeds = []draw2d.Color{draw2d.Hex("#00897b"), draw2d.Hex("#e8743b"), draw2d.Hex("#c2185b"), draw2d.Hex("#558b2f")}

// scheme is the wallpaper's colours (on Android 12 and later; Material's baseline elsewhere), or
// a whole scheme made from the colour chosen.
func (program *gallery) scheme(dark bool) material.Scheme {
	if program.colour == 0 {
		return material.System(dark)
	}
	return material.FromColour(seeds[program.colour-1], dark)
}

// places are the destinations, in the navigation bar, the rail and the drawer.
func (program *gallery) places() []material.Destination {
	return []material.Destination{{Icon: "mouse-pointer-click", Label: "Actions"}, {Icon: "text-cursor-input", Label: "Inputs"},
		{Icon: "list", Label: "Lists", Badge: program.lists.darkCount()}, {Icon: "loader", Label: "Feedback"},
		{Icon: "columns-3", Label: "Tabs"}}
}

// build describes the interface for this frame.
func (program *gallery) build(context *ui.Context) {
	dark := program.isDark()
	scheme := program.scheme(dark)
	context.Theme = scheme.Theme(dark)
	look := material.New(context, scheme)
	// last in the frame: what is over the screen answers the system's back before the screen does
	defer program.goBack(look)
	defer program.overlays.draw(program, look)
	defer program.drawer(look)
	// on its side, a phone has more width than height to spare: the destinations go down a rail
	if look.Wide() {
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.NavigationRail(program.places(), &program.screen)
			look.Body(true, func() { program.show(look) })
		})
		return
	}
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		look.Body(false, func() { program.show(look) })
		look.NavigationBar(program.places(), &program.screen)
	})
}

// show draws the screen chosen, rising into place as it arrives.
func (program *gallery) show(look material.Look) {
	look.Enter(fmt.Sprint("screen ", program.screen), func() {
		switch program.screen {
		case actionsScreen:
			program.actions.draw(program, look)
		case inputsScreen:
			program.inputs.draw(program, look)
		case listsScreen:
			program.lists.draw(program, look)
		case feedbackScreen:
			program.feedback.draw(program, look)
		default:
			program.tabs.draw(program, look)
		}
	})
}

// topBar is every screen's bar: its title, the drawer's button and a menu to choose the theme.
func (program *gallery) topBar(look material.Look, title, subtitle string) {
	look.TopAppBar(program.bar(look, title, subtitle))
}

// largeTopBar is the bar of a screen whose content is one page: large, it shrinks into a small one
// as the page scrolls up under it.
func (program *gallery) largeTopBar(look material.Look, title, subtitle, page string) {
	bar := program.bar(look, title, subtitle)
	bar.Large, bar.Scrolled = true, look.Context.ScrollOffset(page)
	look.TopAppBar(bar)
}

func (program *gallery) bar(look material.Look, title, subtitle string) material.TopBar {
	return material.TopBar{Title: title, Subtitle: subtitle, Actions: func() {
		if look.IconButton("menu") {
			program.drawerOpen = true
		}
		switch look.Menu("ellipsis-vertical", []string{"Light theme", "Dark theme", "About the gallery"}) {
		case 0:
			program.setDark(false)
		case 1:
			program.setDark(true)
		case 2:
			window.OpenLink(projectPage)
		}
	}}
}

// drawer shows every destination in a drawer from the left, opened from the top bar.
func (program *gallery) drawer(look material.Look) {
	if program.drawerOpen && look.NavigationDrawer("drawer", "Gallery", program.places(), &program.screen) {
		program.drawerOpen = false
	}
}

// goBack answers the system's back the overlays and the drawer left: from any destination back to
// the first; from the first, nobody answers, and the program goes to the background.
func (program *gallery) goBack(look material.Look) {
	if program.screen == listsScreen && look.GoBack(&program.lists.navigation) {
		return // from a lighthouse's details back to the list
	}
	if program.screen == listsScreen && program.lists.choosing() {
		if look.Back() {
			program.lists.chooseNone() // the back leaves choosing before it leaves the screen
		}
		return
	}
	if program.screen != actionsScreen && look.Back() {
		program.screen = actionsScreen
	}
}

// aboveBars is how far floating things sit over the bottom: the navigation bar's height, or none
// beside a rail.
func aboveBars(look material.Look) float32 {
	if look.Wide() {
		return 0
	}
	return material.NavigationBarHeight
}

// row lays pieces side by side with Material's 8-pixel gap, breaking onto more lines as needed.
func row(context *ui.Context, content func()) {
	context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, Wrap: true, AlignY: ui.Center}, content)
}

// page is a screen's scrolling content, inside Material's 16-pixel margins.
func page(context *ui.Context, name string, content func()) {
	context.ScrollArea(name, ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 8, Right: 16, Bottom: 160}, Gap: 12}, content)
	})
}

// after is a short wait's end, for notes that go by themselves.
func after(seconds float64) time.Duration { return time.Duration(seconds * float64(time.Second)) }
