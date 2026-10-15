module github.com/akrck02/go-ui-android

go 1.25.6

require (
	github.com/akrck02/go-renderer v0.0.0
	github.com/akrck02/go-sound v0.0.0
	github.com/akrck02/go-ui v0.0.0
)

require (
	github.com/ebitengine/oto/v3 v3.5.1 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/jfreymuth/pulse v0.1.3 // indirect
)

require (
	github.com/go-gl/gl v0.0.0-20260331235117-4566fea9a276 // indirect
	github.com/go-gl/glfw/v3.4/glfw v0.1.0-pre.1 // indirect
	golang.org/x/image v0.45.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)

// Until go-ui, go-renderer and go-sound publish versions with touches, insets, the Android window
// and the sound, build against local copies beside this one.
replace (
	github.com/akrck02/go-renderer => ../go-renderer
	github.com/akrck02/go-sound => ../go-sound
	github.com/akrck02/go-ui => ../go-ui
)
