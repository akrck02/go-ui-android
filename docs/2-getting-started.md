# 🚀 Getting started

Two things for a [go-ui](https://github.com/akrck02/go-ui) program that should feel at home on a
phone: **`material`**, pieces drawn the way Android's own applications are (Material Design 3, with
Material 3 Expressive's newer pieces), and **`go-ui-apk`**, a tool that turns the program into an
Android package.

Neither needs Java, Gradle or Android Studio. The program stays an ordinary Go `main` package that
calls go-ui's `window.Run`: on the desktop it opens a window, built for WebAssembly it draws in a web
page, and on Android it is a shared library that Android's own `NativeActivity` loads (go-renderer's
`android` package).

## Using the pieces

A program makes a `material.Look` in every frame, from the go-ui context and a colour scheme, and
calls its pieces. Pieces that can be touched return `true` in the frame they are touched, as go-ui's
own controls do; the ones that keep something between frames (a pager, a carousel, a snackbar, a
picker) take a small struct of the program's.

```go
dark := window.SystemPrefersDark()
scheme := material.System(dark) // the wallpaper's colours on Android 12 and later
context.Theme = scheme.Theme(dark)
look := material.New(context, scheme)

look.TopAppBar(material.TopBar{Title: "Lighthouses"})
if look.FilledButton("plus", "Add one") { … }
look.NavigationBar(places, &screen)
```

Sizes are Material's: targets of at least 48 pixels, bars 64 and 80 pixels tall, 16-pixel text. A
finger has no hover, so pieces change only while pressed. Each piece takes the name go-ui tells it
apart by (its label, its icon or a name given); rows that repeat get a scope each with
`context.Push`, as in go-ui. `material.UseRoboto()`, before the window opens, draws every text in
Roboto, Android's typeface.

## Building against go-ui and go-renderer

go-ui-android builds against `go-ui` and `go-renderer` checked out beside it (`replace` in `go.mod`),
until they publish versions with touches, insets, the wallpaper's palettes and the Android window.

Then [package it for a phone](4-guides/5-packaging.md), and see [every piece](4-guides/1-pieces.md).

---

[← All the documentation](1-overview.md)
