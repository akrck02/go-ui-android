package material

import (
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Dialog is what a dialog says and offers.
type Dialog struct {
	Icon    string // over the title, centred; none when empty
	Title   string
	Text    string // under the title, broken into lines
	Content func() // more under the text (fields, choices); none when nil
	Confirm string // the main action, at the end
	Dismiss string // the other action, before it; none when empty
}

// Answer is what the person did with a dialog in a frame.
type Answer int

const (
	Open      Answer = iota // nothing yet: keep showing it
	Confirmed               // the main action
	Dismissed               // the other action, a touch outside, Escape
)

// Dialog asks something over the screen, which is dimmed and does not react meanwhile: a card in
// the middle with a title, a text and one or two actions. The program shows it with an if while it
// wants an answer, and stops when the answer is not Open. name tells dialogs apart.
func (look Look) Dialog(name string, dialog Dialog) Answer {
	scheme := look.Scheme
	width, _ := look.Context.Size()
	answer := Open
	// it rises into place as the screen behind darkens
	arrived := look.Context.Appear("dialog "+name, arrivingMotion)
	padding := look.safePadding(24)
	padding.Top += (1 - arrived) * 48
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.Center, Padding: padding,
		Width: ui.Fixed(min(560, max(width-48, 200))), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion * 3 / 4}} // it fades out where it is
	card := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(24), Gap: 16}
	dismissed := look.Context.Modal(name, options, func() {
		look.Context.Frame(card, ui.Look{Fill: scheme.SurfaceContainerHigh, Radius: 28, Shadow: true}, func() {
			alignX := ui.Start
			if dialog.Icon != "" {
				alignX = ui.Center
				look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: ui.Center}, func() {
					look.Context.IconSized(dialog.Icon, 24, scheme.Primary)
				})
			}
			look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: alignX}, func() {
				look.Text(dialog.Title, 24, draw2d.Regular, scheme.OnSurface)
			})
			if dialog.Text != "" {
				look.Paragraph(dialog.Text, 14, draw2d.Regular, scheme.OnSurfaceVariant)
			}
			if dialog.Content != nil {
				dialog.Content()
			}
			look.Context.Box(ui.Layout{Width: ui.Grow(), Gap: 8, AlignX: ui.End}, func() {
				if dialog.Dismiss != "" && look.TextButton("", dialog.Dismiss) {
					answer = Dismissed
				}
				if dialog.Confirm != "" && look.TextButton("", dialog.Confirm) {
					answer = Confirmed
				}
			})
		})
	})
	if dismissed && answer == Open {
		answer = Dismissed
	}
	return answer
}

// BottomSheet shows content in a sheet rising from the bottom of the screen, which is dimmed
// meanwhile: choices for what was touched, a short form. It returns true when the person touches
// outside it (or presses Escape); the program then stops showing it. name tells sheets apart.
func (look Look) BottomSheet(name string, content func()) (dismissed bool) {
	scheme := look.Scheme
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	// it slides up from below the screen as the screen darkens
	arrived := look.Context.Appear("sheet "+name, arrivingMotion)
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: ui.Padding{Top: insets.Top + 72, Bottom: -(1 - arrived) * 400},
		Width: ui.Fixed(min(640, width)), Hidden: 1 - arrived,
		Leave: ui.Leave{Duration: leavingMotion, Y: 400}} // it slides back down
	sheet := ui.Layout{Direction: ui.Vertical, Width: ui.Grow(),
		Padding: ui.Padding{Left: 16 + insets.Left, Top: 16, Right: 16 + insets.Right, Bottom: 16 + insets.Bottom}, Gap: 8}
	return look.Context.Modal(name, options, func() {
		look.Context.Frame(sheet, ui.Look{Fill: scheme.SurfaceContainerLow, Radius: 28}, func() {
			handle := ui.Layout{Width: ui.Fixed(32), Height: ui.Fixed(4)}
			look.Context.Box(ui.Layout{Width: ui.Grow(), AlignX: ui.Center, Padding: ui.Padding{Bottom: 8}}, func() {
				look.Context.Frame(handle, ui.Look{Fill: stateLayer(scheme.OnSurfaceVariant, 0.4), Radius: 2}, nil)
			})
			content()
		})
	})
}

// Snackbar is a short note at the bottom of the screen for a few seconds, with an action that
// undoes or follows what was just done. A program keeps one, calls Show when something happens and
// draws it with Look.Snackbar in every frame.
type Snackbar struct {
	Text   string
	Action string // a button at its end; none when empty
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	until time.Time
}

// Show puts a note up for a while (Material's are 4 to 10 seconds).
func (snackbar *Snackbar) Show(text, action string, duration time.Duration) {
	snackbar.Text, snackbar.Action, snackbar.until = text, action, time.Now().Add(duration)
}

// Hide takes the note away before its time.
func (snackbar *Snackbar) Hide() { snackbar.until = time.Time{} }

// Shown tells whether the note is up.
func (snackbar *Snackbar) Shown() bool { return time.Now().Before(snackbar.until) }

// Snackbar draws the note while it is up, over everything, and asks for a frame when it is due to
// go. It returns true when its action is touched, and the note goes then.
func (look Look) Snackbar(snackbar *Snackbar) (actionTouched bool) {
	left := time.Until(snackbar.until)
	if left <= 0 {
		return false
	}
	look.Context.RequestFrame(left)
	scheme := look.Scheme
	insets := look.Context.Insets()
	width, _ := look.Context.Size()
	arrived := look.Context.Appear("snackbar", arrivingMotion) // it slides up from the bottom edge
	padding := ui.Padding{Left: 16 + insets.Left, Right: 16 + insets.Right,
		Bottom: 16 + insets.Bottom + snackbar.Above - (1-arrived)*(80+insets.Bottom+snackbar.Above)}
	options := ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: padding,
		Width: ui.Fixed(min(600, width-padding.Left-padding.Right)), Leave: ui.Leave{Duration: leavingMotion}}
	look.Context.Overlay("snackbar", options, func() {
		bar := ui.Layout{Width: ui.Grow(), Height: ui.Fixed(48), Padding: ui.Padding{Left: 16, Right: 8}, Gap: 8, AlignY: ui.Center}
		look.Context.Frame(bar, ui.Look{Fill: scheme.InverseSurface, Radius: 4, Shadow: true}, func() {
			look.Context.Box(ui.Layout{Width: ui.Grow()}, func() {
				look.Text(snackbar.Text, 14, draw2d.Regular, scheme.InverseOnSurface)
			})
			if snackbar.Action == "" {
				return
			}
			layout := ui.Layout{Height: ui.Fixed(40), Padding: ui.PadXY(12, 0), AlignY: ui.Center}
			normal := ui.Look{Radius: 20}
			actionTouched = look.Context.Clickable("snackbar action", layout, pressedOnly(normal, stateLayer(scheme.InversePrimary, 0.12)), func() {
				look.Text(snackbar.Action, 14, draw2d.Medium, scheme.InversePrimary)
			})
		})
	})
	if actionTouched {
		snackbar.Hide()
	}
	return actionTouched
}

// Menu is a button with an icon (three dots, usually) that opens a list of choices under it. It
// returns the index of the one chosen in this frame, or -1.
func (look Look) Menu(icon string, choices []string) int {
	chosen := -1
	look.Context.MenuButton(icon, "", func() {
		for index, choice := range choices {
			if look.Context.MenuItem("", choice, "") {
				chosen = index
			}
		}
	})
	return chosen
}
