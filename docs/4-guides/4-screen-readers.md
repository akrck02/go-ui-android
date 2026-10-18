# 🗣️ Screen readers

Every piece reads under TalkBack as Android's own do: a finger moved over the screen reads what is
under it ("Lamp, switch, on"), a swipe goes to the next thing, a double tap presses it, and the
thing read is ringed. Go UI works the list out from the boxes ([its docs](https://github.com/akrck02/go-ui/blob/main/docs/4-guide.md#screen-readers)):
a control is read by the texts inside it, texts side by side are read together, a list holds its
rows, and what is behind a dialog or a sheet is left out. The pieces add what the boxes cannot
tell:

| Piece | Read as |
|---|---|
| `Switch`, `Checkbox`, `RadioButton` | A switch, a checkbox, a radio button, and whether it is on |
| `SegmentedButtons`, `ButtonGroup`, a dropdown's choices, the date picker's days | The chosen one selected |
| `NavigationBar`, `NavigationRail`, `Tabs` | Tabs, the one shown selected; a badge after the label |
| `TopAppBar`, `Dialog` | The title as a heading; the arrow as "Back" (or "Close") |
| `TextField` | The label over the box names the field, read with what is typed (not a secret one's); the screen reader's own typing, pasting and dictation fill it |
| `Slider` | Its range; the volume keys or a swipe up or down move it a step |
| `FloatingMenu`, `FloatingToolbar` | The menu's `Label` ("Create"), "Close" while open; each tool's label |
| Icons with no label (a field's cross, the month's arrows, a `Menu`'s button, the split button's arrow) | What they do ("More options"): `Look.Words`, English by default (`EnglishWords`) |

An icon button of the program's own reads its icon's name ("share 2") unless it has one:
`NamedIconButton("share-2", "Share")`, or `Describe` after any piece. The gallery's test goes to
every screen with double taps alone, and checks that every control has a name.

The package carries one small Java class for it (`go-ui-apk` compiles it in, see [packaging](5-packaging.md)). It was tried
on the Android 14 emulator with an accessibility service acting as TalkBack does: the gallery's
tree reaches the system with every name, role and state, and a click goes to another destination,
the focus rings a checkbox, a second click ticks it, a scroll moves the list on and text set on a
field ("North light 🌊") lands in it whole. TalkBack itself, on a phone, is not tried yet.

---

[← All the documentation](../1-overview.md)
