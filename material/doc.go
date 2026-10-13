// Package material draws go-ui interfaces the way Android's own applications look (Material
// Design 3, with Material 3 Expressive's newer pieces). The pieces are made of go-ui boxes, texts
// and clickables, as any program can make controls of its own, so they draw the same on a phone, on
// a desktop and in a browser. A program makes a Look for each frame and calls its pieces:
//
//	scheme := material.System(dark) // the wallpaper's colours on Android 12 and later
//	context.Theme = scheme.Theme(dark)
//	look := material.New(context, scheme)
//	look.TopAppBar(material.TopBar{Title: "Inbox"})
//	if look.FilledButton("plus", "New message") { … }
//
// The files follow Material's catalogue:
//
//	look.go         the Look, text and paragraphs
//	scheme.go       colour schemes: Material's baseline, a brand's colour
//	dynamic.go      the phone's wallpaper colours (Material You)
//	theme.go        go-ui's own controls in the scheme, and Roboto
//	navigation.go   top bars, the navigation bar, rail and drawer, tabs and swiping pages, the back
//	actions.go      buttons, floating buttons and their menu, split buttons, button groups, toolbars
//	selection.go    switches, check boxes, radio buttons, segmented buttons, sliders
//	chips.go        filter, assist and input chips, and tags
//	inputs.go       text fields and the search bar
//	containment.go  cards, list items, avatars, banners, section titles, dividers
//	overlays.go     dialogs, bottom sheets, snackbars, menus
//	pickers.go      date and time pickers
//	progress.go     loading indicators, circular, linear and wavy progress, skeletons
//	gestures.go     pull to refresh, rows swiped away, carousels
//	stick.go        an on-screen stick
package material
