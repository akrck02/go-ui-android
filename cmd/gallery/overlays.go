package main

import "github.com/akrck02/go-ui-android/material"

// overlaysState is what goes over every screen: a dialog, a sheet and a note at the bottom.
type overlaysState struct {
	dialogShown bool
	sheetShown  bool
	sideShown   bool
	onlyLit     bool
	onlyKept    bool
	snackbar    material.Snackbar
}

// draw shows what is over the screen, after the screen, so it is in front and answers the back
// first.
func (overlays *overlaysState) draw(program *gallery, look material.Look) {
	lights := program.lists.lighthouses
	if overlays.dialogShown {
		answer := look.Dialog("switch off", material.Dialog{Icon: "lightbulb-off", Title: "Switch every light off?",
			Text: "Ships will not see the coast until they are lit again.", Confirm: "Switch off", Dismiss: "Cancel"})
		if answer == material.Confirmed {
			for index := range lights {
				lights[index].lit = false
			}
			overlays.snackbar.Show("Every light is off", "Undo", after(6))
		}
		overlays.dialogShown = answer == material.Open
	}
	if overlays.sheetShown && look.BottomSheet("actions", func() {
		look.ListItem("share", material.Item{Icon: "share-2", Headline: "Share the chart"})
		look.ListItem("copy", material.Item{Icon: "copy", Headline: "Copy the position"})
		if look.ListItem("close", material.Item{Icon: "x", Headline: "Close"}) {
			overlays.sheetShown = false
		}
	}) {
		overlays.sheetShown = false
	}
	if overlays.sideShown && look.SideSheet("filters", "Filters", func() {
		look.Checkbox("Only lights that are lit", &overlays.onlyLit)
		look.Checkbox("Only lights with a keeper", &overlays.onlyKept)
		look.Context.Spacer()
		if look.FilledButton("", "Apply") {
			overlays.sideShown = false
		}
	}) {
		overlays.sideShown = false
	}
	overlays.snackbar.Above = aboveBars(look)
	if look.Snackbar(&overlays.snackbar) {
		for index := range lights {
			lights[index].lit = true
		}
	}
}
