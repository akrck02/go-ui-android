# Lesson 3: fields, dialogs and notes

In this lesson a lighthouse's screen edits it, a dialog asks before removing it, and a note at the
bottom offers to undo.

| Mission | Pieces | Time |
|---|---|---|
| Editing, removing and undoing | `TextField`, `Slider`, `Dialog`, `Snackbar` | 15 min |

#### The idea

A dialog is a **question that needs an answer** before anything else happens; a snackbar is a **sticky note** on the fridge: it says what just happened and offers to take it back, and falls off by itself after a few seconds.

## Editing

Lesson 2's `details` already has the pieces:

| Piece | Edits |
|---|---|
| `TextField(name, &text, material.Field{Label, Placeholder, Supporting, Error, MaxLength, Counter, Filled, Prefix, Suffix, Secret, Reveal})` | A string. On a phone the on-screen keyboard comes up, suggests words and underlines the one it composes; a held word is selected, with handles and a bar to cut, copy and paste |
| `Slider(name, &value, minimum, maximum, step)` | A number |
| `Switch(name, &on)` | A `bool` |
| `Dropdown`, `SegmentedButtons`, `Checkbox`, `RadioButton`, `DatePicker`, `TimePicker` | See [the pieces](../4-guides/1-pieces.md) |

## A dialog

`Dialog(name, material.Dialog{…})` is a question over the dimmed screen. Show it while your program
wants (here, while `removing` names a lighthouse); it returns `material.Open` while it waits, then
`Confirmed` or `Dismissed` (the other button, a touch outside, or the back gesture) once.

## A note with Undo

A `material.Snackbar` is kept by the program: `Show(text, action, duration)` when something
happens, and `look.Snackbar(&snackbar)` at every frame, which returns `true` when its action is
touched. `Above` lifts it over the navigation bar.

```go
// overlays draws what goes over the screens: the dialog that asks before removing, and the note.
func (program *app) overlays(look material.Look) {
	if program.removing != "" {
		answer := look.Dialog("remove", material.Dialog{Icon: "trash-2", Title: "Remove " + program.removing + "?",
			Text: "It leaves the list.", Confirm: "Remove", Dismiss: "Cancel"})
		if answer == material.Confirmed {
			program.remove(program.removing)
			program.navigation.Back()
			program.notes.Show(program.removed.name+" removed", "Undo", 5*time.Second)
		}
		if answer != material.Open {
			program.removing = ""
		}
	}
	program.notes.Above = material.NavigationBarHeight
	if look.Snackbar(&program.notes) && program.removed != nil {
		program.lighthouses = append(program.lighthouses, *program.removed) // Undo
		program.removed = nil
	}
}

func (program *app) remove(name string) {
	for index, light := range program.lighthouses {
		if light.name == name {
			program.removed = &light
			program.lighthouses = append(program.lighthouses[:index], program.lighthouses[index+1:]...)
			return
		}
	}
}
```

![The dialog before East Cape is removed](../images/tutorial/phone-dialog.jpg)

`overlays` is deferred in `build`, so it is drawn over the screens and answers the back before
them.

#### Checkpoint

The bin on a lighthouse's screen asks *Remove East Cape?*; *Remove* takes it out, goes back to the list and shows *East Cape removed* with *Undo*, which puts it back.

![Achievement unlocked: Lamplighter](../images/tutorial/quest/achievement-lamplighter.svg)

## Bonus (optional)

Give the field an error, "A place is needed", while it is empty.

<details>
<summary>🔓 Show the solution</summary>

```go
field := material.Field{Label: "Where it stands", MaxLength: 60, Counter: true}
if light.place == "" {
	field.Error = "A place is needed"
}
look.TextField("place", &light.place, field)
```

</details>

---

[← Lesson 2](2-screens-and-back.md) · [Lesson 4: colours and settings →](4-colours-and-settings.md)
