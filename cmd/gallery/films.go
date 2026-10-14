package main

import (
	_ "embed"

	"github.com/akrck02/go-ui/window"
)

// coastFilm is six seconds of a made-up coast at dusk, a lighthouse's beam turning over the waves,
// painted and encoded for the gallery (H.264, which every phone decodes).
//
//go:embed coast.mp4
var coastFilm []byte

// film is the gallery's film, opened the first time it is shown; where films do not play yet (the
// desktop), it stays nil and the card says so.
type film struct {
	video  *window.Video
	tried  bool
	reason error
}

func (shown *film) open() *window.Video {
	if !shown.tried {
		shown.tried = true
		shown.video, shown.reason = window.OpenVideo(coastFilm)
		if shown.video != nil {
			shown.video.Loop = true
		}
	}
	return shown.video
}
