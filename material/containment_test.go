package material

import "testing"

func TestAnExpandableRowOpensAndCloses(t *testing.T) {
	open, shownMore := false, false
	screen := newScreen(t, func(look Look) {
		shownMore = false
		look.ExpandableItem("faq", Item{Headline: "When is it lit?"}, &open, func() { shownMore = true })
	})
	screen.frame()
	screen.tap("faq")
	if !open || !shownMore {
		t.Fatal("touching the row did not open it")
	}
	screen.tap("faq")
	if open || shownMore {
		t.Fatal("touching it again did not close it")
	}
}

func TestInitialsAreTheFirstLettersOfTwoWords(t *testing.T) {
	for name, want := range map[string]string{"ana maría garcía": "AM", "Ñandú": "Ñ", "": ""} {
		if got := initialsOf(name); got != want {
			t.Fatalf("initials of %q = %q; want %q", name, got, want)
		}
	}
}
