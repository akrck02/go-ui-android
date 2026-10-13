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
	formsScreen
	listScreen
	loadingScreen
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

	// the forms screen
	search, name, reach string
	formTab, kind       int
	nightly             bool
	progress            float32
	dialogShown         bool
	sheetShown          bool
	snackbar            material.Snackbar
	date                time.Time
	datePicker          material.DatePicker
	choosingDate        bool
	drawerOpen          bool
	tags                []string
	volume              float32
	loaded              bool
	size                int
	create              material.FloatingMenu

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
		date: time.Now(), tags: []string{"North shore", "Cliffs", "Harbour"}, volume: 0.4,
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
	places := []material.Destination{{Icon: "layout-grid", Label: "Pieces"}, {Icon: "text-cursor-input", Label: "Forms"},
		{Icon: "list", Label: "List", Badge: program.darkCount()}, {Icon: "loader", Label: "Loading"},
		{Icon: "gamepad-2", Label: "Stick"}}
	// last in the frame: what is over the screen answers the system's back before the screen does
	defer program.goBack(look)
	defer program.layers(look, look.Wide())
	defer program.drawer(look, places)
	// on its side, a phone has more width than height to spare: the destinations go down a rail
	if look.Wide() {
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.NavigationRail(places, &program.screen)
			look.Body(true, func() { program.show(look) })
		})
		return
	}
	context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		look.Body(false, func() { program.show(look) })
		look.NavigationBar(places, &program.screen)
	})
}

// show draws the screen chosen in the navigation, rising into place as it arrives.
func (program *gallery) show(look material.Look) {
	look.Enter(fmt.Sprint("screen ", program.screen), func() { program.showScreen(look) })
}

func (program *gallery) showScreen(look material.Look) {
	switch program.screen {
	case piecesScreen:
		program.pieces(look)
	case formsScreen:
		program.forms(look)
	case listScreen:
		program.list(look)
	case loadingScreen:
		program.loading(look)
	default:
		program.bay(look)
	}
}

// pieces shows the buttons, the choices and a card.
func (program *gallery) pieces(look material.Look) {
	context, scheme := look.Context, look.Scheme
	look.TopAppBar(material.TopBar{Title: "Pieces", Subtitle: "Material Design 3, drawn by go-ui",
		Actions: func() {
			if look.IconButton("menu") {
				program.drawerOpen = true
			}
			switch look.Menu("ellipsis-vertical", []string{"Light theme", "Dark theme"}) {
			case 0:
				program.dark, program.themeSet = false, true
			case 1:
				program.dark, program.themeSet = true, true
			}
		}})
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
		look.SectionTitle("Chips and sliders")
		row(context, func() {
			look.AssistChip("calendar-plus", "Add to calendar")
			look.AssistChip("share-2", "Share")
		})
		row(context, func() {
			for index := 0; index < len(program.tags); index++ {
				if _, removed := look.InputChip(program.tags[index]); removed {
					program.tags = append(program.tags[:index], program.tags[index+1:]...)
					index--
				}
			}
		})
		context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 0), Gap: 12, AlignY: ui.Center}, func() {
			context.IconSized("volume-2", 24, scheme.OnSurfaceVariant)
			look.Slider("volume", &program.volume, 0, 1, 0)
		})
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
		if look.ExtendedFloatingButton("plus", "New light") {
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
	wide := look.Wide()
	// tall, the stick goes under the bay; on its side, beside it
	layout := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow(), Padding: ui.Pad(16), Gap: 16}
	if wide {
		layout.Direction, layout.AlignY = ui.Horizontal, ui.End
	}
	context.Box(layout, func() {
		size := program.baySize(look, wide)
		boat := float32(28)
		padding := ui.Padding{Left: program.boatX * (size - boat), Top: program.boatY * (size - boat)}
		context.Frame(ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), Padding: padding},
			ui.Look{Fill: scheme.PrimaryContainer, Radius: 16}, func() {
				context.Frame(ui.Layout{Width: ui.Fixed(boat), Height: ui.Fixed(boat), AlignX: ui.Center, AlignY: ui.Center},
					ui.Look{Fill: scheme.Primary, Radius: boat / 2}, func() { context.IconSized("sailboat", 18, scheme.OnPrimary) })
			})
		context.Spacer()
		// a stick is found by its name on the screen: keep it out of scopes and scroll areas
		look.Stick("stick", &program.stick, 128)
	})
}

// baySize is the square bay's side: the screen's width when tall, the height left under the top
// bar when on its side.
func (program *gallery) baySize(look material.Look, wide bool) float32 {
	width, height := look.Context.Size()
	insets := look.Context.Insets()
	if wide {
		return max(height-insets.Top-insets.Bottom-64-32, 64)
	}
	return max(width-insets.Left-insets.Right-32, 64)
}

// row lays pieces side by side with Material's 8-pixel gap, inside the screen's margins.
func row(context *ui.Context, content func()) {
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.PadXY(16, 4), Gap: 8, Wrap: true, AlignY: ui.Center}, content)
}

func clamp(value, minimum, maximum float32) float32 {
	return min(max(value, minimum), maximum)
}

// darkCount is the list's badge: how many lighthouses are dark, none when all are lit.
func (program *gallery) darkCount() string {
	dark := 0
	for _, light := range program.lighthouses {
		if !light.lit {
			dark++
		}
	}
	if dark == 0 {
		return ""
	}
	return fmt.Sprint(dark)
}

// forms shows the fields and choices, and opens a dialog, a sheet, a note and a notification.
func (program *gallery) forms(look material.Look) {
	context := look.Context
	look.TopAppBar(material.TopBar{Title: "Forms", Subtitle: "Fields, choices and what goes over the screen", Large: true})
	context.ScrollArea("forms", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.PadXY(16, 8), Gap: 16}, func() {
			look.SearchBar("search", &program.search, "Search the lighthouses")
			look.Tabs("form tabs", []string{"Details", "Watch"}, &program.formTab)
			if program.formTab == 0 {
				program.details(look)
			} else {
				program.watch(look)
			}
			look.SectionTitle("Over the screen")
			context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, Wrap: true}, func() {
				if look.TonalButton("message-square", "Dialog") {
					program.dialogShown = true
				}
				if look.TonalButton("panel-bottom", "Sheet") {
					program.sheetShown = true
				}
				if look.TonalButton("info", "Snackbar") {
					program.snackbar.Show("Harbour Light switched off", "Undo", 6*time.Second)
				}
				if look.TonalButton("bell", "Notification") {
					program.notify()
				}
			})
			context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(24)}, nil)
		})
	})
}

func (program *gallery) details(look material.Look) {
	reachError := ""
	for _, character := range program.reach {
		if character < '0' || character > '9' {
			reachError = "Nautical miles, in digits"
			break
		}
	}
	look.TextField("name", &program.name, material.Field{Label: "Name", Placeholder: "Harbour Light", Supporting: "As the chart writes it"})
	look.TextField("reach", &program.reach, material.Field{Label: "Reach", Placeholder: "18", Supporting: "In nautical miles", Error: reachError, MaxLength: 3})
	look.Checkbox("Lit every night", &program.nightly)
	if look.OutlinedButton("calendar", "Last inspection: "+program.date.Format("2 Jan 2006")) {
		program.choosingDate = true
	}
	for index, kind := range []string{"Fixed light", "Flashing", "Occulting"} {
		if look.RadioButton(kind, index == program.kind) {
			program.kind = index
		}
	}
}

func (program *gallery) watch(look material.Look) {
	scheme := look.Scheme
	look.Text("Keeper's round", 16, draw2d.Medium, scheme.OnSurface)
	look.LinearProgress(program.progress)
	look.Context.SliderStep("round", &program.progress, 0, 1, 0.05)
	look.Text("Waiting for the weather", 16, draw2d.Medium, scheme.OnSurface)
	look.LinearProgress(-1)
}

// notify sends a notification, asking first where the phone wants it allowed.
func (program *gallery) notify() {
	if !window.NotificationsAllowed() {
		window.AskForNotifications()
		program.snackbar.Show("Allow notifications, then touch again", "", 4*time.Second)
		return
	}
	if window.Notify(1, "Harbour Light", "Lit for the night, reaching 18 miles") {
		program.snackbar.Show("Notification sent: look at the shade", "", 4*time.Second)
		return
	}
	program.snackbar.Show("Notifications are not shown here", "", 4*time.Second)
}

// layers draws what goes over the screen: the dialog, the sheet and the note.
func (program *gallery) layers(look material.Look, wide bool) {
	if program.dialogShown {
		answer := look.Dialog("switch off", material.Dialog{Icon: "lightbulb-off", Title: "Switch every light off?",
			Text: "Ships will not see the coast until they are lit again.", Confirm: "Switch off", Dismiss: "Cancel"})
		if answer == material.Confirmed {
			for index := range program.lighthouses {
				program.lighthouses[index].lit = false
			}
			program.snackbar.Show("Every light is off", "Undo", 6*time.Second)
		}
		program.dialogShown = answer == material.Open
	}
	if program.choosingDate {
		program.choosingDate = look.DatePicker("inspection", &program.datePicker, &program.date) == material.Open
	}
	if program.sheetShown && look.BottomSheet("actions", func() {
		look.ListItem("share", material.Item{Icon: "share-2", Headline: "Share the chart"})
		look.ListItem("copy", material.Item{Icon: "copy", Headline: "Copy the position"})
		if look.ListItem("close", material.Item{Icon: "x", Headline: "Close"}) {
			program.sheetShown = false
		}
	}) {
		program.sheetShown = false
	}
	program.snackbar.Above = material.NavigationBarHeight
	if wide {
		program.snackbar.Above = 0
	}
	if look.Snackbar(&program.snackbar) {
		for index := range program.lighthouses {
			program.lighthouses[index].lit = true
		}
	}
}

// drawer shows every destination in a drawer from the left, opened from the Pieces screen's menu.
func (program *gallery) drawer(look material.Look, places []material.Destination) {
	if program.drawerOpen && look.NavigationDrawer("drawer", "Gallery", places, &program.screen) {
		program.drawerOpen = false
	}
}

// goBack answers the system's back the dialogs, sheets and drawer left: from any destination back
// to the first; from the first, nobody answers, and the program goes to the background.
func (program *gallery) goBack(look material.Look) {
	if program.screen != piecesScreen && look.Back() {
		program.screen = piecesScreen
	}
}

// loading shows Material 3 Expressive's indicators, skeletons standing in for a list, and the
// pieces for actions: a button group, a split button, a floating toolbar and a floating menu.
func (program *gallery) loading(look material.Look) {
	context, scheme := look.Context, look.Scheme
	look.TopAppBar(material.TopBar{Title: "Loading", Subtitle: "Indicators, skeletons and actions"})
	context.ScrollArea("loading", ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
		context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Padding{Left: 16, Top: 8, Right: 16, Bottom: 160}, Gap: 16}, func() {
			look.SectionTitle("Indicators")
			context.Box(ui.Layout{Width: ui.Grow(), Gap: 16, AlignY: ui.Center, Wrap: true}, func() {
				look.LoadingIndicator(48)
				look.ContainedLoadingIndicator(48)
				look.CircularProgress(material.Indeterminate, 48, false)
				look.CircularProgress(material.Indeterminate, 48, true)
				look.CircularProgress(program.progress, 48, true)
			})
			look.WavyProgress(material.Indeterminate)
			look.WavyProgress(program.progress)
			look.Slider("progress", &program.progress, 0, 1, 0)
			look.SectionTitle("Skeletons")
			look.ListItem("loaded", material.Item{Headline: "Loaded", Trailing: func() { look.Switch("loaded switch", &program.loaded) }})
			for index, light := range program.lighthouses {
				if program.loaded {
					look.ListItem(fmt.Sprint("light ", index), material.Item{Icon: "lightbulb", Headline: light.name, Supporting: light.place})
					continue
				}
				context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(72), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}, func() {
					look.Skeleton(ui.Fixed(40), ui.Fixed(40), 20)
					context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Gap: 8}, func() {
						look.Skeleton(ui.Percent(0.6), ui.Fixed(14), 7)
						look.Skeleton(ui.Percent(0.4), ui.Fixed(12), 6)
					})
				})
			}
			look.SectionTitle("Actions")
			look.ButtonGroup("size", []string{"Small", "Medium", "Large"}, &program.size)
			if action, more := look.SplitButton("send", "Send the report"); action {
				program.snackbar.Show("Report sent", "", 3*time.Second)
			} else if more {
				program.snackbar.Show("Later: schedule it, or send it to someone else", "", 3*time.Second)
			}
			look.Text("A toolbar and a menu float over the screen.", 14, draw2d.Regular, scheme.OnSurfaceVariant)
		})
	})
	above := float32(material.NavigationBarHeight)
	if look.Wide() {
		above = 0
	}
	switch look.FloatingToolbar("tools", []string{"pencil", "share-2", "trash-2"}, above) {
	case 0:
		program.snackbar.Show("Edit", "", 2*time.Second)
	case 1:
		program.snackbar.Show("Share", "", 2*time.Second)
	case 2:
		program.snackbar.Show("Deleted", "Undo", 4*time.Second)
	}
	program.create.Icon, program.create.Above = "plus", above
	program.create.Actions = []material.Destination{{Icon: "file", Label: "Note"}, {Icon: "image", Label: "Picture"}, {Icon: "clock", Label: "Reminder"}}
	if chosen := look.FloatingMenu("create", &program.create); chosen >= 0 {
		program.snackbar.Show("New "+program.create.Actions[chosen].Label, "", 2*time.Second)
	}
}
