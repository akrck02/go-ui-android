package material

import (
	"testing"
	"time"
)

func TestADatePickerKeepsTheDayChosenOnlyWhenConfirmed(t *testing.T) {
	date := time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)
	var picker DatePicker
	answer := Open
	shown := true
	screen := newScreen(t, func(look Look) {
		if shown {
			answer = look.DatePicker("when", &picker, &date)
			shown = answer == Open
		}
	})
	screen.frame()
	screen.tap("when", "icon chevron-right") // to April
	screen.tap("when", "day 2026-04-21")
	if !date.Equal(time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)) {
		t.Fatal("the date changed before it was confirmed")
	}
	screen.tap("when", "button OK")
	if answer != Confirmed || !date.Equal(time.Date(2026, time.April, 21, 9, 30, 0, 0, time.UTC)) {
		t.Fatalf("answer %v, date %v: want Confirmed and 21 April, keeping the time of day", answer, date)
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

func TestAnInputChipTellsItsCrossApart(t *testing.T) {
	removed := false
	screen := newScreen(t, func(look Look) {
		if _, gone := look.InputChip("harbour"); gone {
			removed = true
		}
	})
	screen.frame()
	screen.tap("input harbour", "remove")
	if !removed {
		t.Fatal("touching the cross did not remove the chip")
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
