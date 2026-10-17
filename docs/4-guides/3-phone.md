# 📱 Feeling at home on a phone

What a phone does that a desktop does not: fingers that swipe and hold, the back gesture, a
screen that turns, notifications, and the phone's own services. Most of it comes from go-ui and
go-renderer; the pieces use it the way Android's own applications do.

## Gestures

go-ui follows the one finger on a touch screen from where it comes down. Once it has moved a few
pixels, the way it went first decides who takes it: the innermost list for up and down, the
innermost of the sideways pieces (a carousel inside swiping pages) for sideways. The control it came
down on is then not touched, and a finger that hardly moves stays a touch. So the pieces sit inside
each other freely: the gallery has a carousel inside a list inside swiping pages, and rows swiped
away inside a list pulled to refresh.

- **Long press.** A finger held still on a piece for half a second (go-ui's `LongPressed`; with a
  mouse, its right button) long-presses it, and lifting it then is no touch. `look.LongPressed()`
  after a piece tells, and gives the phone's short buzz, as Android's own applications do. A list
  starts choosing with it: the row held is chosen (`Item.Selected`), touches choose more, and the
  top bar says how many, with a cross to stop (`TopBar.Close`) and actions for what is chosen.

  ```go
  touched := look.ListItem("row", material.Item{Headline: light.name, Selected: light.chosen})
  switch {
  case look.LongPressed():
  	light.chosen = true
  case touched && choosing:
  	light.chosen = !light.chosen
  }
  ```
- **Pull to refresh.** The indicator follows the pull (go-ui's `ScrollPull`); let go past 64 pixels,
  `PullToRefresh` returns true, and the program refreshes, keeping `refreshing` set until done.
- **Swiping pages.** A page let go past a third of the way, or flicked, turns to the next; the tab
  underline moves with the pages, between two tabs while a page is half way.
- **Rows swiped away.** Past two fifths of their width, or flicked, they go; short of that they slide
  back.
- **Carousels.** Let go, they glide on with the finger's speed and settle at a tile's edge.
- **The stick.** The finger that comes down on it holds it; `Forward` and `Sideways` go from -1 to 1.
  Other fingers stay free, so a program can turn a view with them (`context.Touches()`, leaving out
  `stick.Touching(touch)`). It can sit anywhere, in a `Push` scope, a list or swiping pages: the
  finger steering it is its own, so the pages and the list stay still under it.

## Back and motion

A phone's back gesture reaches the program. Dialogs, sheets, the drawer and the floating menu answer
it first and close; then `look.Back()` tells the screen, which stops choosing rows or goes back (to
the first destination, to the list a detail came from). A back nobody answers sends the program to
the background, as Android's own applications do at their first screen. Ask for it after drawing what goes over the
screen, so that gets it first:

```go
defer func() {
	if screen != firstScreen && look.Back() {
		screen = firstScreen
	}
}()
```

Pieces move as Material's do, and only while they move does the window draw: the switch's handle
slides and grows, the navigation bar's pill widens under the chosen icon, the tabs' underline slides,
dialogs rise in as the screen darkens, sheets and the drawer slide in, the snackbar slides up, the
floating menu's actions rise one after another, and `look.Enter(key, …)` lets a screen rise into
place when it is first shown.

They leave with motion too, though a program simply stops drawing them: a dialog answered fades
out, a sheet slides back down, the drawer slides back off the side, the snackbar fades when its time
is up, a side sheet slides back off its side, a floating menu's actions sink and fade as it
closes, and a row's more folds back under it. go-ui keeps what such a layer
painted last (`LayerOptions.Leave`) and lets that picture go; nothing in it can be touched on its
way out.

## Turning the phone

A phone on its side has more width than height to spare, and its camera's cutout at a side. Material
moves the destinations from a bar at the bottom to a rail down the left:

```go
if look.Wide() {
	context.Box(ui.Layout{Width: ui.Grow(), Height: ui.Grow()}, func() {
		look.NavigationRail(places, &screen)
		look.Body(true, func() { drawScreen(look) })
	})
	return
}
context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Height: ui.Grow()}, func() {
	look.Body(false, func() { drawScreen(look) })
	look.NavigationBar(places, &screen)
})
```

`Body` keeps the content clear of the cutout; go-ui-apk's window theme paints beside it
(`shortEdges`), so no black strip is left there.

## Notifications

A program's notifications go through go-ui's window: `window.Notify(id, title, text)` shows one in
the phone's shade, even while the program is in the background, and touching it opens the program.
Android 13 and later show them only once the person allows it: check `window.NotificationsAllowed()`
and ask with `window.AskForNotifications()` when they do something that needs them. Build the
package with `go-ui-apk -notifications`, which names the permission. In a browser they are the
page's notifications; on the desktop they are not shown yet.

## The phone's services

What a program asks of the phone goes through go-ui's window, and works in a browser too:

- **Links and sharing.** `window.OpenLink(link)` opens a web page in the browser, `mailto:` in the
  mail, `tel:` in the phone; `window.Share(title, text)` offers a text through the system's sheet.
  The gallery shares the lighthouses chosen in its list, and opens its project page from the menu.
- **The clipboard.** `look.Context.Copy(text)` puts a text on the phone's clipboard, and text fields
  copy and paste through it.
- **Preferences.** `window.Preference(name)` and `window.KeepPreference(name, value)` keep small
  values between runs. Android may end a program in the background without a word, so keep them
  in `Options.Hidden`, called when the program leaves the screen: the gallery keeps its screen, its
  theme and its colours there, and comes back as it was left.
- **Sound.** The feedback screen plays a note through [go-sound](https://github.com/akrck02/go-sound),
  which goes quiet when the gallery leaves the screen (`Options.Hidden`) and plays again when it
  comes back (`Options.Shown`). go-ui-apk passes the NDK's C++ compiler too, for the sound output's
  C++ (Oboe).
- **The person's text size.** Text grows with the phone's font size setting; icons and paddings
  keep theirs, as in Android's own applications.

---

[← All the documentation](../1-overview.md)
