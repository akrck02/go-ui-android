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
