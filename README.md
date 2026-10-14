# go-ui-android

> [AI disclosure](#ai-disclosure).

[go-ui](https://github.com/akrck02/go-ui) interfaces that look like Android's own applications
(Material Design 3), and a tool that packages any go-ui program as an Android application. The same
program still runs on the desktop and in a web page.

```go
dark := window.SystemPrefersDark()
scheme := material.System(dark) // the wallpaper's colours on Android 12 and later
context.Theme = scheme.Theme(dark)
look := material.New(context, scheme)

look.TopAppBar(material.TopBar{Title: "Lighthouses"})
if look.FilledButton("plus", "Add one") { … }
look.NavigationBar(places, &screen)
```

```bash
go run github.com/akrck02/go-ui-android/cmd/go-ui-apk -name Lighthouses -id io.example.lighthouses ./cmd/lighthouses
adb install lighthouses.apk
```

<p>
<img src="docs/images/actions.png" width="180" alt="Buttons, a split button, a button group and a floating toolbar">
<img src="docs/images/lists.png" width="180" alt="A banner over a list pulled to refresh">
<img src="docs/images/feedback.png" width="180" alt="Loading indicators, wavy progress and skeletons">
<img src="docs/images/tabs.png" width="180" alt="Tabs over swiping pages with a carousel">
</p>

Start with [the docs](docs/README.md): every piece by Material's catalogue (navigation, actions,
selection, inputs, containment, overlays, pickers, progress and gestures). `go run ./cmd/gallery`
shows them all.

## AI disclosure

I write most of my code by hand, so I want to be clear about where AI was used here.

- **Written by an AI assistant (Claude, by Anthropic) under my direction, in October 2026:** all of
  it so far: the material package, the go-ui-apk tool, the gallery and the docs. It is pending my
  review.

Commits made with AI carry a `Co-Authored-By: Claude` line.

The material package embeds the Roboto typeface (SIL Open Font License, licence in `material/fonts`).

## Licence

MIT, see [LICENSE](LICENSE). Roboto is under the SIL Open Font License (`material/fonts/OFL.txt`).
