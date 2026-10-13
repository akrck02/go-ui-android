package material

import (
	"testing"

	"github.com/akrck02/go-renderer/draw2d"
	ui "github.com/akrck02/go-ui"
)

func TestAParagraphBreaksIntoLines(t *testing.T) {
	screen := newScreen(t, func(look Look) {
		look.Context.Box(ui.Layout{Width: ui.Fixed(200)}, func() {
			look.Context.Push("text")
			look.Paragraph("one two three four five six seven eight nine ten eleven twelve", 14, draw2d.Regular, look.Scheme.OnSurface)
			look.Context.Pop()
		})
	})
	screen.frame()
	lines := map[float32]bool{}
	for _, command := range screen.context.Drawing().Commands {
		if command.Kind != draw2d.TextKind {
			continue
		}
		lines[command.Baseline] = true
		if command.Area.X+command.Area.Width > 200.5 {
			t.Fatalf("%q ends at %v, past the 200-pixel box", command.Text, command.Area.X+command.Area.Width)
		}
	}
	if len(lines) < 2 {
		t.Fatalf("%d lines; want the words broken into several in 200 pixels", len(lines))
	}
}

func TestATextFieldTakesWhatIsTyped(t *testing.T) {
	name := ""
	screen := newScreen(t, func(look Look) {
		look.TextField("name", &name, Field{Label: "Name", Supporting: "As on the chart"})
	})
	screen.frame()
	screen.tap("name")
	if !screen.context.Typing() {
		t.Fatal("a touched field is not being typed in: the keyboard would not come up")
	}
	screen.state.AddCharacter('Ñ')
	screen.state.AddCharacter('o')
	screen.frame()
	if name != "Ño" {
		t.Fatalf("name %q; want what was typed", name)
	}
}
