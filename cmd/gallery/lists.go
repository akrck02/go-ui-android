package main

import (
	"fmt"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
	"github.com/akrck02/go-ui-android/material"
)

// listsState is what the lists screen remembers.
type listsState struct {
	lighthouses []lighthouse
	added       int
	refreshing  int // frames left of a pretend refresh
	answers     [3]bool
	bannerShown bool
}

func newLists() listsState {
	return listsState{bannerShown: true, lighthouses: []lighthouse{{"North Cape", "Cliffs, north shore", true},
		{"Harbour Light", "End of the pier", true}, {"Seal Rock", "Two miles out", false},
		{"Old Tower", "The hill above the town", false}, {"East Point", "The lighthouse keeper's cottage", true},
		{"Sand Bar", "On the bar at low tide", false}}}
}

// darkCount is the navigation's badge: how many lighthouses are dark, none when all are lit.
func (lists *listsState) darkCount() string {
	dark := 0
	for _, light := range lists.lighthouses {
		if !light.lit {
			dark++
		}
	}
	if dark == 0 {
		return ""
	}
	return fmt.Sprint(dark)
}

// draw shows a list that refreshes when pulled, rows swiped away, a banner and rows that open.
func (lists *listsState) draw(program *gallery, look material.Look) {
	context, scheme := look.Context, look.Scheme
	notes := &program.overlays.snackbar
	program.topBar(look, "Lists", "Pull to refresh, swipe a row away")
	if lists.bannerShown {
		switch look.Banner("wifi-off", "The lights' positions were last brought up to date a week ago.", "Dismiss", "Refresh now") {
		case 0:
			lists.bannerShown = false
		case 1:
			lists.bannerShown, lists.refreshing = false, 90
		}
	}
	if lists.refreshing > 0 {
		lists.refreshing--
		context.RequestFrame(0)
		if lists.refreshing == 0 {
			notes.Show("Positions brought up to date", "", after(2))
		}
	}
	layout := ui.Layout{Width: ui.Grow(), Height: ui.Grow()}
	if look.PullToRefresh("lights", lists.refreshing > 0, layout, func() {
		for index := 0; index < len(lists.lighthouses); index++ {
			light := &lists.lighthouses[index]
			context.Push(light.name)
			if look.SwipeToDismiss("swipe", func() {
				look.ListItem("row", material.Item{Icon: "lightbulb", Headline: light.name, Supporting: light.place,
					Trailing: func() { look.Switch("lit", &light.lit) }})
			}) {
				notes.Show(light.name+" removed", "", after(3))
				lists.lighthouses = append(lists.lighthouses[:index], lists.lighthouses[index+1:]...)
				index--
			}
			context.Pop()
		}
		look.SectionTitle("People")
		for _, person := range []string{"Ana Lighthouse-Keeper", "Tom Boatman", "Lena Harbour"} {
			context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(56), Padding: ui.PadXY(16, 0), Gap: 16, AlignY: ui.Center}, func() {
				look.Avatar(person, "", draw2d.Color{}, 40)
				look.Text(person, 16, draw2d.Regular, scheme.OnSurface)
			})
		}
		look.SectionTitle("Questions")
		questions := []string{"When are the lights lit?", "Who keeps them?", "Can they be visited?"}
		answers := []string{"From sunset to sunrise, and in fog.", "Keepers live by three of them; the others run on their own.",
			"The Old Tower opens on summer weekends."}
		for index, question := range questions {
			look.ExpandableItem(fmt.Sprint("question ", index), material.Item{Headline: question}, &lists.answers[index], func() {
				look.Paragraph(answers[index], 14, draw2d.Regular, scheme.OnSurfaceVariant)
			})
		}
		context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Fixed(96)}, nil)
	}) {
		lists.refreshing = 90 // about a second and a half at sixty frames
	}
	context.Box(ui.Layout{Width: ui.Grow(), Padding: ui.Pad(16), AlignX: ui.End}, func() {
		if look.ExtendedFloatingButton("plus", "New light") {
			lists.added++
			lists.lighthouses = append(lists.lighthouses, lighthouse{name: fmt.Sprint("New light ", lists.added), place: "Not placed yet"})
		}
	})
}
