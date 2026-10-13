package material

import (
	"fmt"
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

// Pieces Material 3 Expressive added for actions: a floating button that opens into a few actions,
// a button split from its menu, buttons joined in a group, and a toolbar floating over the content.

// FloatingMenu is a floating action button that opens into a few related actions stacked above
// it, its icon turning into a cross. A program keeps one, and draws it with Look.FloatingMenu.
type FloatingMenu struct {
	Icon    string        // the button's icon while closed
	Actions []Destination // what it opens into: an icon and a label each
	// Above is how far over the bottom of the window it sits: a navigation bar's height
	// (NavigationBarHeight) when there is one.
	Above float32
	Open  bool
}

// FloatingMenu draws the button at the bottom end of the screen, over everything, and its actions
// while open. It returns the index of the action touched in this frame, or -1; touching one, the
// cross or the system's back closes it.
func (look Look) FloatingMenu(name string, menu *FloatingMenu) int {
	scheme := look.Scheme
	insets := look.Context.Insets()
	if menu.Open && look.Context.Back() {
		menu.Open = false
	}
	chosen := -1
	padding := ui.Padding{Right: 16 + insets.Right, Bottom: 16 + insets.Bottom + menu.Above}
	look.Context.Overlay(name, ui.LayerOptions{AlignX: ui.End, AlignY: ui.End, Padding: padding}, func() {
		look.Context.Box(ui.Layout{Direction: ui.Vertical, Gap: 8, AlignX: ui.End}, func() {
			if menu.Open {
				for index, action := range menu.Actions {
					// the actions rise in one after another
					arrived := look.Context.Appear(fmt.Sprint(name, " action ", index), time.Duration(150+50*index)*time.Millisecond)
					layout := ui.Layout{Height: ui.Fixed(56), Padding: ui.Padding{Left: 16, Right: 20, Top: (1 - arrived) * 16}, Gap: 12, AlignY: ui.Center}
					normal := ui.Look{Fill: scheme.PrimaryContainer, Radius: 28, Shadow: true}
					if look.Context.Clickable("action "+action.Label, layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
						look.Context.IconSized(action.Icon, 24, scheme.OnPrimaryContainer)
						look.Text(action.Label, 16, draw2d.Medium, scheme.OnPrimaryContainer)
					}) {
						chosen, menu.Open = index, false
					}
				}
			}
			// open, the button is a smaller round cross in the main colour
			icon, size, normal := menu.Icon, float32(56), ui.Look{Fill: scheme.PrimaryContainer, Radius: 16, Shadow: true}
			ink := scheme.OnPrimaryContainer
			if menu.Open {
				icon, size, normal, ink = "x", 56, ui.Look{Fill: scheme.Primary, Radius: 28, Shadow: true}, scheme.OnPrimary
			}
			layout := ui.Layout{Width: ui.Fixed(size), Height: ui.Fixed(size), AlignX: ui.Center, AlignY: ui.Center}
			if look.Context.Clickable("toggle", layout, pressedOnly(normal, scheme.SecondaryContainer), func() {
				look.Context.IconSized(icon, 24, ink)
			}) {
				menu.Open = !menu.Open
			}
		})
	})
	return chosen
}

// SplitButton is a main action joined to a small button that opens its other choices (send, and
// schedule): it tells which half was touched.
func (look Look) SplitButton(icon, label string) (action, more bool) {
	scheme := look.Scheme
	look.Context.Push("split " + label)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Height: ui.Fixed(40), Gap: 2}, func() {
		leading := ui.Look{Fill: scheme.Primary, Radius: 20}
		layout := ui.Layout{Height: ui.Grow(), Padding: ui.Padding{Left: 16, Right: 12}, Gap: 8, AlignY: ui.Center}
		action = look.Context.Clickable("action", layout, pressedOnly(leading, scheme.PrimaryContainer), func() {
			if icon != "" {
				look.Context.IconSized(icon, 18, scheme.OnPrimary)
			}
			look.Text(label, 14, draw2d.Medium, scheme.OnPrimary)
		})
		trailing := ui.Look{Fill: scheme.Primary, Radius: 20}
		layout = ui.Layout{Width: ui.Fixed(40), Height: ui.Grow(), AlignX: ui.Center, AlignY: ui.Center}
		more = look.Context.Clickable("more", layout, pressedOnly(trailing, scheme.PrimaryContainer), func() {
			look.Context.IconSized("chevron-down", 18, scheme.OnPrimary)
		})
	})
	return action, more
}

// ButtonGroup is a few choices joined side by side, the chosen one filled and fully rounded, the
// others squarer: Material 3 Expressive's connected group (a view, a size, an alignment). It
// changes selected and returns true when another is touched.
func (look Look) ButtonGroup(name string, labels []string, selected *int) bool {
	scheme := look.Scheme
	changed := false
	look.Context.Push(name)
	defer look.Context.Pop()
	look.Context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(40), Gap: 2}, func() {
		for index, label := range labels {
			chosen := index == *selected
			// the chosen one rounds out fully; its corners ease as the choice moves
			roundness := look.Context.Animate("round "+label, boolFloat(chosen), switchMotion)
			normal := ui.Look{Fill: scheme.SecondaryContainer, Radius: 8 + 12*roundness}
			ink := scheme.OnSecondaryContainer
			if chosen {
				normal.Fill, ink = scheme.Primary, scheme.OnPrimary
			}
			layout := ui.Layout{Width: ui.Share(1), Height: ui.Grow(), AlignX: ui.Center, AlignY: ui.Center}
			if look.Context.Clickable(label, layout, pressedOnly(normal, stateLayer(ink, 0.2)), func() {
				look.Text(label, 14, draw2d.Medium, ink)
			}) && !chosen {
				*selected, changed = index, true
			}
		}
	})
	return changed
}

// FloatingToolbar is a row of actions in a pill floating over the content at the bottom of the
// screen, for what is done to the content on show (edit, share, delete). above is how far over the
// bottom it sits. It returns the index of the icon touched in this frame, or -1.
func (look Look) FloatingToolbar(name string, icons []string, above float32) int {
	scheme := look.Scheme
	insets := look.Context.Insets()
	touched := -1
	padding := ui.Padding{Bottom: 16 + insets.Bottom + above}
	look.Context.Overlay(name, ui.LayerOptions{AlignX: ui.Center, AlignY: ui.End, Padding: padding}, func() {
		bar := ui.Layout{Height: ui.Fixed(64), Padding: ui.PadXY(8, 0), Gap: 4, AlignY: ui.Center}
		look.Context.Frame(bar, ui.Look{Fill: scheme.PrimaryContainer, Radius: 32, Shadow: true}, func() {
			for index, icon := range icons {
				layout := ui.Layout{Width: ui.Fixed(48), Height: ui.Fixed(48), AlignX: ui.Center, AlignY: ui.Center}
				if look.Context.Clickable("tool "+icon, layout, ui.Looks{Pressed: ui.Look{Fill: stateLayer(scheme.OnPrimaryContainer, 0.12), Radius: 24}}, func() {
					look.Context.IconSized(icon, 24, scheme.OnPrimaryContainer)
				}) {
					touched = index
				}
			}
		})
	})
	return touched
}
