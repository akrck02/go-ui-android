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

func TestADialogAnsweredFadesOutInsteadOfVanishing(t *testing.T) {
	asking := true
	screen := newScreen(t, func(look Look) {
		if asking && look.Dialog("light", Dialog{Title: "Light the north lamp?", Confirm: "Light"}) != Open {
			asking = false
		}
	})
	screen.clocked()
	for range 20 {
		screen.frame()
	}
	screen.tap("light", "button Light")
	if asking {
		t.Fatal("touching Light did not answer the dialog")
	}
	screen.frame()
	if !screen.drawn("Light the north lamp?") {
		t.Fatal("the dialog vanished at once")
	}
	for range 20 {
		screen.frame()
	}
	if screen.drawn("Light the north lamp?") {
		t.Fatal("the dialog is still there after it had time to fade")
	}
}
