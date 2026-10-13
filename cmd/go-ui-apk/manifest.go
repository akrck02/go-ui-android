package main

import (
	"bytes"
	"errors"
	"fmt"
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
	Debuggable  bool
	// Notifications names the permission Android 13 and later ask the person for before the
	// program's notifications show.
	Notifications bool
	Package       string // the Go package to build
}

// libraryName is the shared library's name in the package (libmain.so), which the manifest names.
const libraryName = "main"

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
	}
	return nil
}

// The manifest: Android's NativeActivity runs the program's library, with no Java (hasCode false).
// The activity keeps running through rotations and theme changes, which the program follows itself,
// and the on-screen keyboard shrinks the window rather than covering what is typed in.
var manifestTemplate = template.Must(template.New("manifest").Funcs(template.FuncMap{"xml": escapeXML}).Parse(
	`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="{{.ID}}" android:versionCode="{{.VersionCode}}" android:versionName="{{xml .VersionName}}">
  <uses-sdk android:minSdkVersion="{{.MinimumSDK}}" android:targetSdkVersion="{{.TargetSDK}}" />
  {{- if .Notifications}}
  <uses-permission android:name="android.permission.POST_NOTIFICATIONS" />
  {{- end}}
  <application android:label="{{xml .Name}}" android:hasCode="false" android:theme="@style/WindowTheme"
      {{- if .Icon}} android:icon="@mipmap/icon"{{end}}{{if .Debuggable}} android:debuggable="true"{{end}}>
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
func windowTheme(dark bool) []byte {
	parent, lightIcons, background := "Theme.DeviceDefault.Light.NoActionBar", "true", "white"
	if dark {
		parent, lightIcons, background = "Theme.DeviceDefault.NoActionBar", "false", "black"
	}
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
  <style name="WindowTheme" parent="@android:style/%s">
    <item name="android:statusBarColor">@android:color/transparent</item>
    <item name="android:navigationBarColor">@android:color/transparent</item>
    <item name="android:windowLightStatusBar">%s</item>
    <item name="android:windowLightNavigationBar">%s</item>
    <item name="android:windowBackground">@android:color/%s</item>
    <item name="android:windowLayoutInDisplayCutoutMode">shortEdges</item>
  </style>
</resources>
`, parent, lightIcons, lightIcons, background))
}
