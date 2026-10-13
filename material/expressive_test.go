package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

func fixed(size float32) ui.Sizing { return ui.Fixed(size) }

func TestLoadingIndicatorsPaintAndKeepMoving(t *testing.T) {
	for name, draw := range map[string]func(look Look){
		"loading":   func(look Look) { look.LoadingIndicator(48) },
		"contained": func(look Look) { look.ContainedLoadingIndicator(48) },
		"circular":  func(look Look) { look.CircularProgress(Indeterminate, 48, true) },
		"wavy":      func(look Look) { look.WavyProgress(0.4) },
		"skeleton":  func(look Look) { look.Skeleton(fixed(200), fixed(16), 8) },
	} {
		screen := newScreen(t, draw)
		screen.frame()
		painted := 0
		for _, command := range screen.context.Drawing().Commands {
			if command.Kind == draw2d.RectangleKind {
				painted++
			}
		}
		// a skeleton is one block, and its light only while it crosses; shapes are many rectangles
		if wanted := map[bool]int{true: 1, false: 5}[name == "skeleton"]; painted < wanted {
			t.Fatalf("%s painted %d rectangles; want a shape made of many", name, painted)
		}
		if _, wanted := screen.state.TakeFrameRequest(); !wanted {
			t.Fatalf("%s does not ask for its next frame", name)
		}
	}
}

func TestAKnownProgressStandsStill(t *testing.T) {
	screen := newScreen(t, func(look Look) { look.CircularProgress(0.5, 48, false) })
	screen.frame()
	if _, wanted := screen.state.TakeFrameRequest(); wanted {
		t.Fatal("a still ring asks for frames: a screen showing it never sleeps")
	}
}

func TestTheFloatingMenuOpensChoosesAndCloses(t *testing.T) {
	menu := &FloatingMenu{Icon: "plus", Actions: []Destination{{Icon: "file", Label: "Note"}, {Icon: "image", Label: "Picture"}}}
	chosen := -1
	screen := newScreen(t, func(look Look) {
		if index := look.FloatingMenu("create", menu); index >= 0 {
			chosen = index
		}
	})
	screen.frame()
	screen.tap("create", "toggle")
	if !menu.Open {
		t.Fatal("touching the button does not open the menu")
	}
	screen.frame() // the actions rise in
	screen.tap("create", "action Picture")
	if chosen != 1 || menu.Open {
		t.Fatalf("chosen %d, open %v: want Picture, and the menu closed", chosen, menu.Open)
	}
	screen.tap("create", "toggle")
	screen.state.PressBack()
	screen.frame()
	if menu.Open {
		t.Fatal("the back does not close the menu")
	}
}

func TestGroupedAndSplitButtonsAnswerTheirParts(t *testing.T) {
	selected, action, more, tool := 0, false, false, -1
	screen := newScreen(t, func(look Look) {
		look.ButtonGroup("size", []string{"S", "M", "L"}, &selected)
		if main, menu := look.SplitButton("send", "Send"); main || menu {
			action, more = action || main, more || menu
		}
		if index := look.FloatingToolbar("tools", []string{"pencil", "share-2", "trash-2"}, 0); index >= 0 {
			tool = index
		}
	})
	screen.frame()
	screen.tap("size", "L")
	screen.tap("split Send", "more")
	screen.tap("tools", "tool trash-2")
	if selected != 2 || action || !more || tool != 2 {
		t.Fatalf("selected %d, action %v, more %v, tool %d: want L, the menu half alone, and the bin", selected, action, more, tool)
	}
}
