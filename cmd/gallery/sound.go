package main

import (
	"time"

	sound "github.com/akrck02/go-sound"
	"github.com/akrck02/go-ui-android/material"
)

// soundState plays a note through go-sound, the gallery's way to try the phone's speaker. The
// sound output opens the first time a note is played: a gallery left silent opens none.
type soundState struct {
	context   *sound.Context
	note      *sound.Player
	volume    float32
	suspended bool
}

// play plays the note from its start, opening the sound output the first time; it says what went
// wrong in a note at the bottom.
func (sounds *soundState) play(notes *material.Snackbar) {
	if sounds.context == nil {
		context, err := sound.NewContext(sound.Options{Name: "Gallery"})
		if err != nil {
			notes.Show("No sound: "+err.Error(), "", after(4))
			return
		}
		sounds.context = context
	}
	if sounds.note == nil {
		player, err := sounds.context.NewPlayer(sound.Tone(660, 600*time.Millisecond, 0.6))
		if err != nil {
			notes.Show("No sound: "+err.Error(), "", after(4))
			return
		}
		sounds.note = player
	}
	sounds.note.SetVolume(float64(sounds.volume))
	sounds.note.Rewind()
	sounds.note.Play()
}

// setVolume follows the slider while the note plays.
func (sounds *soundState) setVolume() {
	if sounds.note != nil {
		sounds.note.SetVolume(float64(sounds.volume))
	}
}

// suspend stops the sound when the gallery leaves the screen, as a phone's applications do.
func (sounds *soundState) suspend() {
	if sounds.context != nil && sounds.context.Suspend() == nil {
		sounds.suspended = true
	}
}

// resume plays again once the gallery is back on the screen (window.Options.Shown).
func (sounds *soundState) resume() {
	if sounds.suspended && sounds.context.Resume() == nil {
		sounds.suspended = false
	}
}
