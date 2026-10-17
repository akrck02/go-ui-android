# Reference

| Package or tool | What it is |
|---|---|
| [`material`](material.md) | Material Design 3 pieces for go-ui: `Look`, its pieces, `Scheme` and the state structs some pieces keep. Generated from the doc comments (`python3 docs/reference/generate.py`) |
| [The pieces](../pieces.md) | The same pieces by Material's catalogue, with pictures: the place to choose one |
| [`cmd/go-ui-apk`](../packaging.md) | Every flag of the packaging tool |

## The shape of every piece

| Kind of piece | Signature | Example |
|---|---|---|
| Touched | returns `true` in the frame it is touched | `FilledButton(icon, label) bool`, `ListItem(name, item) bool` |
| Edits a value | takes a pointer, returns `true` when it changed | `Switch(name, &on)`, `Slider(name, &value, minimum, maximum, step)`, `NavigationBar(destinations, &selected)` |
| Chooses one of a few | returns the index, or -1 | `Menu(icon, choices) int`, `FloatingToolbar(…) int`, `Banner(icon, text, actions...) int` |
| Asks | returns a `material.Answer`: `Open`, `Confirmed` or `Dismissed` | `Dialog`, `DatePicker`, `TimePicker` |
| Closes | returns `true` when the person closes it | `BottomSheet`, `SideSheet`, `NavigationDrawer` |
| Keeps state between frames | takes a struct the program keeps | `Snackbar`, `Navigation`, `Pager`, `Carousel`, `Dropdown`, `Search`, `Sheet`, `Stick`, `FloatingMenu`, `DatePicker` |
| Holds content | takes a `func()` | `Card(content)`, `Body(besideRail, content)`, `Screens(name, &navigation, screen)` |

Names: each piece takes the name go-ui tells it apart by (its label, its icon, or a name given);
repeated pieces get a scope each with `look.Context.Push`.

## Constants and helpers

| | |
|---|---|
| `material.NavigationBarHeight` | 80: what to lift floating things by over the navigation bar |
| `material.Indeterminate` | -1: a progress indicator with no known fraction |
| `material.UseRoboto()` | Draws every text in Roboto; before the window opens |
| `material.EnglishWords`, `EnglishDates`, `EnglishTimes` | The pieces' own words; replace them with `Look.Words`, `DateNames`, `TimeNames` in another language |
| `material.PalettesFromColour(colour)` | Tonal palettes from one colour |

---

[← Documentation](../README.md)
