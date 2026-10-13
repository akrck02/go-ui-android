package material

import (
	"fmt"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// DateNames are the words a date picker shows, for a program in another language than English.
type DateNames struct {
	Title    string     // over the date chosen
	Months   [12]string // January first
	Weekdays [7]string  // one or two letters each, Monday first
	Confirm  string     // the action that keeps the date
	Dismiss  string     // the one that leaves it
}

// EnglishDates are the words a date picker shows by default.
var EnglishDates = DateNames{Title: "Select date",
	Months:   [12]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	Weekdays: [7]string{"M", "T", "W", "T", "F", "S", "S"}, Confirm: "OK", Dismiss: "Cancel"}

// DatePicker is what a date picker remembers while shown: the month on show and the day chosen.
type DatePicker struct {
	Names DateNames // the zero value shows EnglishDates
	month time.Time // the first of the month on show
	day   time.Time // the day chosen, not yet kept
	open  bool
}

// DatePicker asks for a day in a dialog: a month at a time, with arrows to the months around it.
// It starts at the date given; on Confirmed it writes the day chosen there. The program shows it
// with an if while it wants an answer, as Dialog.
func (look Look) DatePicker(name string, picker *DatePicker, date *time.Time) Answer {
	names := picker.Names
	if names.Confirm == "" {
		names = EnglishDates
	}
	if !picker.open {
		picker.open, picker.day = true, dayOf(*date)
		picker.month = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	}
	scheme := look.Scheme
	answer := look.Dialog(name, Dialog{Title: names.Title, Confirm: names.Confirm, Dismiss: names.Dismiss, Content: func() {
		look.Text(fmt.Sprintf("%d %s %d", picker.day.Day(), names.Months[picker.day.Month()-1], picker.day.Year()), 28, draw2d.Regular, scheme.OnSurface)
		look.Divider()
		look.monthBar(picker, names)
		look.monthGrid(picker, names)
	}})
	if answer == Confirmed {
		*date = time.Date(picker.day.Year(), picker.day.Month(), picker.day.Day(), date.Hour(), date.Minute(), date.Second(), 0, date.Location())
	}
	if answer != Open {
		picker.open = false
	}
	return answer
}

// monthBar names the month on show, between arrows to the one before and the one after.
func (look Look) monthBar(picker *DatePicker, names DateNames) {
	look.Context.Box(ui.Layout{Width: ui.Grow(), AlignY: ui.Center}, func() {
		look.Context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Padding{Left: 12}}, func() {
			look.Text(fmt.Sprintf("%s %d", names.Months[picker.month.Month()-1], picker.month.Year()), 14, draw2d.Medium, look.Scheme.OnSurfaceVariant)
		})
		if look.IconButton("chevron-left") {
			picker.month = picker.month.AddDate(0, -1, 0)
		}
		if look.IconButton("chevron-right") {
			picker.month = picker.month.AddDate(0, 1, 0)
		}
	})
}

// monthGrid lays the month's days under the weekdays, Monday first, the chosen one filled and
// today ringed.
func (look Look) monthGrid(picker *DatePicker, names DateNames) {
	scheme := look.Scheme
	today := dayOf(time.Now())
	week := func(content func()) { look.Context.Box(ui.Layout{Width: ui.Grow()}, content) }
	cell := ui.Layout{Width: ui.Share(1), Height: ui.Fixed(44), AlignX: ui.Center, AlignY: ui.Center}
	week(func() {
		for _, weekday := range names.Weekdays {
			look.Context.Box(cell, func() { look.Text(weekday, 14, draw2d.Regular, scheme.OnSurface) })
		}
	})
	blanks := (int(picker.month.Weekday()) + 6) % 7 // Monday first
	days := picker.month.AddDate(0, 1, -1).Day()
	for start := 1 - blanks; start <= days; start += 7 {
		week(func() {
			for day := start; day < start+7; day++ {
				if day < 1 || day > days {
					look.Context.Box(cell, nil)
					continue
				}
				date := picker.month.AddDate(0, 0, day-1)
				if look.dayCell(date, date.Equal(picker.day), date.Equal(today), cell) {
					picker.day = date
				}
			}
		})
	}
}

func (look Look) dayCell(date time.Time, chosen, today bool, cell ui.Layout) bool {
	scheme := look.Scheme
	circle, ink := ui.Look{Radius: 20}, scheme.OnSurface
	switch {
	case chosen:
		circle, ink = ui.Look{Fill: scheme.Primary, Radius: 20}, scheme.OnPrimary
	case today:
		circle, ink = ui.Look{Border: scheme.Primary, BorderWidth: 1, Radius: 20}, scheme.Primary
	}
	touched := false
	look.Context.Box(cell, func() {
		layout := ui.Layout{Width: ui.Fixed(40), Height: ui.Fixed(40), AlignX: ui.Center, AlignY: ui.Center}
		touched = look.Context.Clickable(date.Format("day 2006-01-02"), layout, pressedOnly(circle, stateLayer(scheme.OnSurface, 0.12)), func() {
			look.Text(fmt.Sprint(date.Day()), 14, draw2d.Regular, ink)
		})
	})
	return touched
}

// dayOf is a date's day, at midnight, in its own zone.
func dayOf(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
