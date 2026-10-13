package material

import (
	"testing"
	"time"

	"github.com/akrck02/go-renderer/input"
)

func TestADialogSaysWhatWasChosen(t *testing.T) {
	answer := Open
	screen := newScreen(t, func(look Look) {
		if got := look.Dialog("delete", Dialog{Title: "Delete the light?", Text: "It cannot be undone.", Confirm: "Delete", Dismiss: "Cancel"}); got != Open {
			answer = got
		}
	})
	screen.frame()
	screen.tap("delete", "button Delete")
	if answer != Confirmed {
		t.Fatalf("answer %v after Delete; want Confirmed", answer)
	}
	answer = Open
	screen.tap("delete", "button Cancel")
	if answer != Dismissed {
		t.Fatalf("answer %v after Cancel; want Dismissed", answer)
	}
	answer = Open
	screen.state.MoveCursor(5, 5) // the dimmed screen around the card
	screen.state.SetButton(input.MouseLeft, true)
	screen.frame()
	screen.state.SetButton(input.MouseLeft, false)
	screen.frame()
	if answer != Dismissed {
		t.Fatalf("answer %v after a touch outside; want Dismissed", answer)
	}
}

func TestASnackbarShowsForItsTimeAndItsActionAnswers(t *testing.T) {
	var snackbar Snackbar
	undone := false
	screen := newScreen(t, func(look Look) {
		if look.Snackbar(&snackbar) {
			undone = true
		}
	})
	snackbar.Show("Light removed", "Undo", time.Minute)
	screen.frame()
	if _, found := screen.context.Area("snackbar"); !found || !snackbar.Shown() {
		t.Fatal("the note is not up")
	}
	screen.tap("snackbar", "snackbar action")
	if !undone || snackbar.Shown() {
		t.Fatalf("undone %v, still shown %v: the action answers and takes the note away", undone, snackbar.Shown())
	}
	snackbar.Show("Saved", "", time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	screen.frame()
	if snackbar.Shown() {
		t.Fatal("the note outlived its time")
	}
}

func TestChoicesAnswerTouches(t *testing.T) {
	ticked, chosen, tab := false, 0, 0
	screen := newScreen(t, func(look Look) {
		look.Checkbox("Lit at night", &ticked)
		for index, option := range []string{"North", "South"} {
			if look.RadioButton(option, index == chosen) {
				chosen = index
			}
		}
		look.Tabs("views", []string{"Map", "List"}, &tab)
	})
	screen.frame()
	screen.tap("check Lit at night")
	screen.tap("radio South")
	screen.tap("views", "List")
	if !ticked || chosen != 1 || tab != 1 {
		t.Fatalf("ticked %v, chosen %d, tab %d: want ticked, South, List", ticked, chosen, tab)
	}
}

func TestATextFieldTakesWhatIsTyped(t *testing.T) {
	name := ""
	screen := newScreen(t, func(look Look) {
		look.TextField("name", &name, Field{Label: "Name", Supporting: "As on the chart"})
	})
	screen.frame()
	screen.tap("name")
	if !screen.context.Typing() {
		t.Fatal("a touched field is not being typed in: the keyboard would not come up")
	}
	screen.state.AddCharacter('Ñ')
	screen.state.AddCharacter('o')
	screen.frame()
	if name != "Ño" {
		t.Fatalf("name %q; want what was typed", name)
	}
}
