package main

import (
	"time"

	"github.com/akrck02/go-renderer/draw2d"
	"github.com/akrck02/go-ui-android/material"
)

// inputsState is what the inputs screen remembers.
type inputsState struct {
	search, name, reach string
	kind, view          int
	nightly             bool
	filters             [3]bool
	tags                []string
	volume              float32
	date                time.Time
	datePicker          material.DatePicker
	timePicker          material.TimePicker
	choosingDate        bool
	choosingTime        bool
}

func newInputs() inputsState {
	return inputsState{date: time.Now(), tags: []string{"North shore", "Cliffs", "Harbour"}, volume: 0.4}
}

// draw shows the fields, the choices and the pickers.
func (inputs *inputsState) draw(program *gallery, look material.Look) {
	context, scheme := look.Context, look.Scheme
	program.topBar(look, "Inputs", "Fields, choices and pickers")
	page(context, "inputs", func() {
		look.SearchBar("search", &inputs.search, "Search the lighthouses")
		look.Heading("Text fields")
		reachError := ""
		for _, character := range inputs.reach {
			if character < '0' || character > '9' {
				reachError = "Nautical miles, in digits"
				break
			}
		}
		look.TextField("name", &inputs.name, material.Field{Label: "Name", Placeholder: "Harbour Light", Supporting: "As the chart writes it"})
		look.TextField("reach", &inputs.reach, material.Field{Label: "Reach", Placeholder: "18", Supporting: "In nautical miles", Error: reachError, MaxLength: 3})
		look.Heading("Choices")
		look.Checkbox("Lit every night", &inputs.nightly)
		for index, kind := range []string{"Fixed light", "Flashing", "Occulting"} {
			if look.RadioButton(kind, index == inputs.kind) {
				inputs.kind = index
			}
		}
		dark := program.isDark()
		look.ListItem("dark theme", material.Item{Icon: "moon", Headline: "Dark theme", Supporting: "Follows the system until chosen here",
			Trailing: func() {
				if look.Switch("dark switch", &dark) {
					program.setDark(dark)
				}
			}})
		look.SegmentedButtons("view", []string{"Day", "Week", "Month"}, &inputs.view)
		row(context, func() {
			for index, label := range []string{"Lit", "Dark", "Far out"} {
				if look.FilterChip(label, inputs.filters[index]) {
					inputs.filters[index] = !inputs.filters[index]
				}
			}
		})
		row(context, func() {
			for index := 0; index < len(inputs.tags); index++ {
				if _, removed := look.InputChip(inputs.tags[index]); removed {
					inputs.tags = append(inputs.tags[:index], inputs.tags[index+1:]...)
					index--
				}
			}
		})
		look.Context.Box(uiRow(12), func() {
			context.IconSized("volume-2", 24, scheme.OnSurfaceVariant)
			look.Slider("volume", &inputs.volume, 0, 1, 0)
		})
		look.Heading("Pickers")
		look.Text("Last inspection", 14, draw2d.Medium, scheme.OnSurfaceVariant)
		row(context, func() {
			if look.OutlinedButton("calendar", inputs.date.Format("2 Jan 2006")) {
				inputs.choosingDate = true
			}
			if look.OutlinedButton("clock", inputs.date.Format("15:04")) {
				inputs.choosingTime = true
			}
		})
	})
	if inputs.choosingDate {
		inputs.choosingDate = look.DatePicker("inspection day", &inputs.datePicker, &inputs.date) == material.Open
	}
	if inputs.choosingTime {
		inputs.choosingTime = look.TimePicker("inspection time", &inputs.timePicker, &inputs.date) == material.Open
	}
}
