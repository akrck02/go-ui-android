# Components

Every Material piece of `material`, one by one: when to use it, concrete cases with their code,
how it behaves under a finger and with TalkBack. Each case is a real program's code, and the
pictures are captures of it in a phone-sized window. [The pieces](../pieces.md) shows them all at
a glance; this is the page to read when using one.

The cases assume the usual frame: `look` is the `material.Look` made in it
(`material.New(context, scheme)`), and variables such as `lit` or `place` belong to your program
and live between frames. `padded` stands for a column with the screen's margins:

```go
func padded(look material.Look, content func()) {
	look.Context.Box(ui.Layout{Direction: ui.Vertical, Width: ui.Grow(), Padding: ui.Pad(16), Gap: 12}, content)
}
```

| Page | Pieces |
|---|---|
| [Navigation](navigation.md) | `TopAppBar`, `NavigationBar`, `NavigationRail`, `NavigationDrawer`, `Tabs`, `Screens` |
| [Actions](actions.md) | `FilledButton`, `TonalButton`, `OutlinedButton`, `TextButton`, `IconButton`, `NamedIconButton`, `RoundButton`, `SplitButton`, `FloatingButton`, `ExtendedFloatingButton`, `FloatingMenu`, `ButtonGroup`, `FloatingToolbar` |
| [Selection](selection.md) | `Switch`, `Checkbox`, `RadioButton`, `SegmentedButtons`, `Slider`, `FilterChip`, `InputChip`, `AssistChip`, `Tag` |
| [Inputs](inputs.md) | `TextField`, `Dropdown`, `SearchBar`, `SearchView` |
| [Containment](containment.md) | `ListItem`, `ExpandableItem`, `Card`, `OutlinedCard`, `ElevatedCard`, `Avatar`, `Banner`, `SectionTitle`, `Divider` |
| [Overlays](overlays.md) | `Dialog`, `BottomSheet`, `StandardSheet`, `SideSheet`, `Snackbar`, `Menu` |
| [Pickers](pickers.md) | `DatePicker`, `TimePicker` |
| [Progress](progress.md) | `LoadingIndicator`, `CircularProgress`, `WavyProgress`, `LinearProgress`, `Skeleton` |
| [Gestures](gestures.md) | `Carousel`, `MultiBrowseCarousel`, `TabPager`, `PullToRefresh`, `SwipeToDismiss`, `LongPressed`, `Stick` |

## Rules every piece follows

- Pieces that can be touched **return `true` in the frame they are touched**; pieces that edit a
  value take a pointer to it. Overlays that ask return a `material.Answer`.
- A piece is told apart by its **name** (its label, its icon, or a name given). Rows that repeat
  get a scope each: `look.Context.Push(name)` … `look.Context.Pop()`.
- Targets are at least **48 pixels**, as Material asks. A finger has no hover: pieces change only
  while pressed.
- Colours come from the **scheme** (`look.Scheme.OnSurface`, `ErrorContainer`…), so the dark theme
  and the wallpaper's colours follow.
- Every piece is read by **TalkBack**; icon buttons need a name (`NamedIconButton`).

---

[← Documentation](../README.md)
