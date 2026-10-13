package material

import (
	"testing"

	ui "github.com/akrck02/go-ui"
)

func TestTheNavigationBarGoesWhereItIsTouched(t *testing.T) {
	selected, changes := 0, 0
	places := []Destination{{Icon: "home", Label: "Home"}, {Icon: "map", Label: "Map"}, {Icon: "settings", Label: "Settings"}}
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
			look.NavigationBar([]Destination{{Icon: "home", Label: "Home"}}, &selected)
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

func TestOnItsSideTheScreenUsesARailAndKeepsClearOfTheEdges(t *testing.T) {
	selected := 0
	screen := newScreen(t, func(look Look) {
		if !look.Wide() {
			t.Fatal("a screen wider than tall is not wide")
		}
		look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
			look.NavigationRail([]Destination{{Icon: "home", Label: "Home"}, {Icon: "map", Label: "Map"}}, &selected)
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

func TestADrawerChoosesAndCloses(t *testing.T) {
	selected, open, closings := 0, true, 0
	places := []Destination{{Icon: "inbox", Label: "Inbox"}, {Icon: "send", Label: "Sent"}}
	screen := newScreen(t, func(look Look) {
		if open && look.NavigationDrawer("menu", "Mail", places, &selected) {
			open = false
			closings++
		}
	})
	screen.frame()
	screen.tap("menu", "drawer Sent")
	if selected != 1 || open || closings != 1 {
		t.Fatalf("selected %d, open %v: want Sent chosen and the drawer closed", selected, open)
	}
}

func TestTheBackGoesToTheScreenWhenNoDialogIsOpen(t *testing.T) {
	dialog, wentBack := true, false
	screen := newScreen(t, func(look Look) {
		if dialog {
			dialog = look.Dialog("ask", Dialog{Title: "Sure?", Confirm: "Yes"}) == Open
		}
		if look.Back() {
			wentBack = true
		}
	})
	screen.frame()
	screen.state.PressBack()
	screen.frame()
	if dialog || wentBack {
		t.Fatalf("dialog open %v, screen went back %v: the dialog closes first", dialog, wentBack)
	}
	screen.state.PressBack()
	screen.frame()
	if !wentBack {
		t.Fatal("the screen did not go back")
	}
}
