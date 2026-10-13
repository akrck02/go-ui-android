package material

import (
	"testing"

	"github.com/akrck02/go-renderer/input"
)

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
