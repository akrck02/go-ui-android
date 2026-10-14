package material

import (
	"fmt"
	"testing"
	"time"

	"github.com/akrck02/go-renderer/input"
	ui "github.com/akrck02/go-ui"
)

// clocked gives the screen a clock that moves a frame at a time, as fingers' speeds need.
func (screen *screen) clocked() {
	now := time.Unix(1000, 0)
	screen.state.Now = func() time.Time { return now }
	frame := screen.build
	screen.build = func(look Look) {
		now = now.Add(16 * time.Millisecond)
		frame(look)
	}
}

// drag puts a finger down, moves it in ten frames, and lifts it when told.
func (screen *screen) drag(fromX, fromY, byX, byY float64, lift bool) {
	screen.state.SetTouch(1, fromX, fromY)
	screen.state.MoveCursor(fromX, fromY)
	screen.state.SetButton(input.MouseLeft, true)
	screen.frame()
	for step := 1; step <= 10; step++ {
		x, y := fromX+byX*float64(step)/10, fromY+byY*float64(step)/10
		screen.state.SetTouch(1, x, y)
		screen.state.MoveCursor(x, y)
		screen.frame()
	}
	if lift {
		screen.state.SetButton(input.MouseLeft, false)
		screen.state.LiftTouch(1)
		screen.frame()
	}
}

func rows(look Look, count int) {
	for index := range count {
		look.ListItem(fmt.Sprint("row ", index), Item{Headline: fmt.Sprint("Light ", index)})
	}
}

func TestPullingAListDownAndLettingGoRefreshesIt(t *testing.T) {
	refreshes := 0
	screen := newScreen(t, func(look Look) {
		if look.PullToRefresh("lights", false, ui.Layout{Width: ui.Grow(), Height: ui.Fixed(600)}, func() { rows(look, 30) }) {
			refreshes++
		}
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(200, 100, 0, 60, true) // too short a pull
	if refreshes != 0 {
		t.Fatal("a short pull refreshed the list")
	}
	screen.drag(200, 100, 0, 300, true)
	screen.frame()
	if refreshes != 1 {
		t.Fatalf("%d refreshes after a long pull; want one", refreshes)
	}
}

func TestARowSwipedFarGoesAndOneSwipedLittleStays(t *testing.T) {
	dismissed := 0
	screen := newScreen(t, func(look Look) {
		if look.SwipeToDismiss("row", func() { look.ListItem("light", Item{Headline: "Harbour Light"}) }) {
			dismissed++
		}
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(200, 30, 60, 2, true)
	for range 20 {
		screen.frame()
	}
	if dismissed != 0 {
		t.Fatal("a short swipe took the row away")
	}
	screen.drag(200, 30, -260, 2, true)
	if dismissed != 1 {
		t.Fatalf("%d dismissals after a long swipe; want one", dismissed)
	}
}

func TestPagesTurnWithAFingerAndWithTheirTabs(t *testing.T) {
	var pager Pager
	screen := newScreen(t, func(look Look) {
		look.TabPager("views", []string{"Map", "List", "Notes"}, &pager, func(index int) { rows(look, 3) })
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(300, 300, -250, 5, true)
	if pager.Selected != 1 {
		t.Fatalf("page %d after swiping left; want the next one", pager.Selected)
	}
	screen.tap("views", "Notes")
	if pager.Selected != 2 {
		t.Fatalf("page %d after touching Notes; want it", pager.Selected)
	}
	for range 40 {
		screen.frame()
	}
	if pager.shown != 2 {
		t.Fatalf("pages shown at %v; want them settled on Notes", pager.shown)
	}
}

func TestACarouselSettlesOnAnItem(t *testing.T) {
	var carousel Carousel
	screen := newScreen(t, func(look Look) {
		look.Carousel("places", &carousel, 8, 160, func(index int) {
			look.CarouselItem(200, look.Scheme.PrimaryContainer, "", fmt.Sprint("Place ", index))
		})
	})
	screen.clocked()
	screen.frame()
	screen.frame()
	screen.drag(300, 100, -250, 3, true)
	if carousel.rest < 1 || carousel.rest != float32(int(carousel.rest)) {
		t.Fatalf("rests at %v after a swipe; want a later item's edge", carousel.rest)
	}
}
