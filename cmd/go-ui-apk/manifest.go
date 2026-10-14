package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

// application is what the package says about the program.
type application struct {
	Name        string
	ID          string
	VersionCode int
	VersionName string
	MinimumSDK  int
	TargetSDK   int
	Icon        string // a PNG file, or empty for Android's default icon
	// IconForeground, IconBackground and IconMonochrome make an adaptive icon (Android 8 and
	// later), which each launcher cuts to its own shape: a PNG for the foreground layer (108 dp,
	// its middle 66 dp always shown), a colour (#rrggbb) or a PNG behind it, and optionally a
	// one-colour PNG for themed icons (Android 13 and later).
	IconForeground, IconBackground, IconMonochrome string
	// SplashLight and SplashDark are the colours (#rrggbb) the window shows while the program
	// starts, in the light and the dark theme; Android 12 and later draw the icon over it.
	SplashLight, SplashDark string
	Debuggable              bool
	// Notifications names the permission Android 13 and later ask the person for before the
	// program's notifications show.
	Notifications bool
	// Internet names the permission to use the network, which Android grants when the package is
	// installed, with no question to the person.
	Internet bool
	Package  string // the Go package to build
	ABIs     []abi  // the processors it carries a library for
	// Release signs with the program's own key (Keystore, KeyAlias) instead of the debug key.
	Release            bool
	Keystore, KeyAlias string
}

// libraryName is the shared library's name in the package (libmain.so), which the manifest names.
const libraryName = "main"

// hexColour is a colour as Android's resources write it: #rrggbb.
var hexColour = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// An application id is two or more dotted parts, each a letter then letters, digits or _.
var applicationID = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+$`)

func (app application) check() error {
	switch {
	case app.Name == "":
		return errors.New("-name is required: the name under the icon")
	case !applicationID.MatchString(app.ID):
		return fmt.Errorf("-id %q is not an application id, like io.example.app", app.ID)
	case app.VersionCode < 1:
		return errors.New("-version-code must be 1 or more")
	case app.MinimumSDK < 26:
		return errors.New("-min-sdk must be 26 or more: the window needs Android 8")
	case app.TargetSDK < app.MinimumSDK:
		return errors.New("-target-sdk must not be under -min-sdk")
	case (app.IconBackground != "" || app.IconMonochrome != "") && app.IconForeground == "":
		return errors.New("-icon-background and -icon-monochrome go with -icon-foreground")
	case app.IconForeground != "" && app.IconBackground == "":
		return errors.New("-icon-foreground needs -icon-background: a colour (#rrggbb) or a PNG")
	}
	for flag, colour := range map[string]string{"-splash": app.SplashLight, "-splash-dark": app.SplashDark} {
		if colour != "" && !hexColour.MatchString(colour) {
			return fmt.Errorf("%s %q is not a colour like #1f5fc4", flag, colour)
		}
	}
	return app.checkRelease()
}

// The manifest: Android's NativeActivity runs the program's library, with no Java (hasCode false).
// The activity keeps running through rotations and theme changes, which the program follows itself,
// and the on-screen keyboard shrinks the window rather than covering what is typed in.
var manifestTemplate = template.Must(template.New("manifest").Funcs(template.FuncMap{"xml": escapeXML}).Parse(
	`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="{{.ID}}" android:versionCode="{{.VersionCode}}" android:versionName="{{xml .VersionName}}">
  <uses-sdk android:minSdkVersion="{{.MinimumSDK}}" android:targetSdkVersion="{{.TargetSDK}}" />
  {{- if .Internet}}
  <uses-permission android:name="android.permission.INTERNET" />
  {{- end}}
  {{- if .Notifications}}
  <uses-permission android:name="android.permission.POST_NOTIFICATIONS" />
  {{- end}}
  <application android:label="{{xml .Name}}" android:hasCode="false" android:theme="@style/WindowTheme"
      {{- if or .Icon .IconForeground}} android:icon="@mipmap/icon"{{end}}{{if .Debuggable}} android:debuggable="true"{{end}}>
    <activity android:name="android.app.NativeActivity" android:label="{{xml .Name}}"
        android:configChanges="orientation|screenSize|screenLayout|keyboardHidden|uiMode|density"
        android:windowSoftInputMode="adjustResize" android:exported="true">
      <meta-data android:name="android.app.lib_name" android:value="` + libraryName + `" />
      <intent-filter>
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
      </intent-filter>
    </activity>
  </application>
</manifest>
`))

func (app application) manifest() ([]byte, error) {
	var written bytes.Buffer
	if err := manifestTemplate.Execute(&written, app); err != nil {
		return nil, err
	}
	return written.Bytes(), nil
}

func escapeXML(text string) string {
	var escaped bytes.Buffer
	for _, character := range text {
		switch character {
		case '&':
			escaped.WriteString("&amp;")
		case '<':
			escaped.WriteString("&lt;")
		case '>':
			escaped.WriteString("&gt;")
		case '"':
			escaped.WriteString("&quot;")
		case '\'':
			escaped.WriteString("&apos;")
		default:
			escaped.WriteRune(character)
		}
	}
	return escaped.String()
}

// windowTheme is the window's theme, light or dark (for the system's night mode): it covers the
// whole screen under see-through bars and beside a camera's cutout (on its side when the phone lies
// on its side), whose backgrounds the program paints itself
// (ui.Context.Insets), with the bars' icons dark on a light interface and light on a dark one.
func windowTheme(dark bool, splash string) []byte {
	parent, lightIcons, background := "Theme.DeviceDefault.Light.NoActionBar", "true", "@android:color/white"
	if dark {
		parent, lightIcons, background = "Theme.DeviceDefault.NoActionBar", "false", "@android:color/black"
	}
	if splash != "" {
		background = splash // also behind the icon Android 12 and later show while the program starts
	}
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
  <style name="WindowTheme" parent="@android:style/%s">
    <item name="android:statusBarColor">@android:color/transparent</item>
    <item name="android:navigationBarColor">@android:color/transparent</item>
    <item name="android:windowLightStatusBar">%s</item>
    <item name="android:windowLightNavigationBar">%s</item>
    <item name="android:windowBackground">%s</item>
    <item name="android:windowLayoutInDisplayCutoutMode">shortEdges</item>
  </style>
</resources>
`, parent, lightIcons, lightIcons, background))
}

// resources are the files aapt2 links into the package: the manifest, the window's theme in the
// light and the dark, and the icon, read from the files the flags name.
func (app application) resources() (map[string][]byte, error) {
	manifest, err := app.manifest()
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{
		"AndroidManifest.xml":         manifest,
		"res/values/styles.xml":       windowTheme(false, app.SplashLight),
		"res/values-night/styles.xml": windowTheme(true, app.SplashDark),
	}
	pictures := map[string]string{
		// the legacy icon, and the adaptive icon's layers, drawn at the highest density so every
		// screen scales them down
		"res/mipmap/icon.png":                    app.Icon,
		"res/mipmap-xxxhdpi/icon_foreground.png": app.IconForeground,
		"res/mipmap-xxxhdpi/icon_monochrome.png": app.IconMonochrome,
	}
	if !hexColour.MatchString(app.IconBackground) {
		pictures["res/mipmap-xxxhdpi/icon_background.png"] = app.IconBackground
	}
	for name, path := range pictures {
		if path == "" {
			continue
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("the icon: %w", err)
		}
		files[name] = contents
	}
	if app.IconForeground != "" {
		files["res/mipmap-anydpi-v26/icon.xml"] = app.adaptiveIcon()
		if hexColour.MatchString(app.IconBackground) {
			files["res/values/icon_background.xml"] = []byte(`<?xml version="1.0" encoding="utf-8"?>
<resources><color name="icon_background">` + app.IconBackground + `</color></resources>
`)
		}
	}
	return files, nil
}

// adaptiveIcon names the icon's layers for launchers that cut icons to their own shape.
func (app application) adaptiveIcon() []byte {
	background := "@mipmap/icon_background"
	if hexColour.MatchString(app.IconBackground) {
		background = "@color/icon_background"
	}
	monochrome := ""
	if app.IconMonochrome != "" {
		monochrome = "\n  <monochrome android:drawable=\"@mipmap/icon_monochrome\" />"
	}
	return []byte(`<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
  <background android:drawable="` + background + `" />
  <foreground android:drawable="@mipmap/icon_foreground" />` + monochrome + `
</adaptive-icon>
`)
}
