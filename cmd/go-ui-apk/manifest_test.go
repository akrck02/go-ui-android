package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheManifestNamesTheProgramAndItsLibrary(t *testing.T) {
	app := application{Name: `Tides & "Winds"`, ID: "io.example.tides", VersionCode: 3, VersionName: "1.2",
		MinimumSDK: 26, TargetSDK: 34, Icon: "icon.png", Debuggable: true, Notifications: true, Internet: true}
	manifest, err := app.manifest()
	if err != nil {
		t.Fatal(err)
	}
	if err := xml.Unmarshal(manifest, new(struct{})); err != nil {
		t.Fatalf("not XML: %v\n%s", err, manifest)
	}
	for _, wanted := range []string{`package="io.example.tides"`, `android:versionCode="3"`,
		`android:label="Tides &amp; &quot;Winds&quot;"`, `android:value="main"`, `android:icon="@mipmap/icon"`,
		`android:debuggable="true"`, `android:hasCode="true"`, `android.permission.POST_NOTIFICATIONS`, `android.permission.INTERNET`,
		`adjustResize`} {
		if !strings.Contains(string(manifest), wanted) {
			t.Fatalf("the manifest has no %s:\n%s", wanted, manifest)
		}
	}
}

func TestAPlainBuildLeavesOutTheIconAndDebugging(t *testing.T) {
	app := application{Name: "Tides", ID: "io.example.tides", VersionCode: 1, VersionName: "1.0", MinimumSDK: 26, TargetSDK: 34}
	manifest, _ := app.manifest()
	if strings.Contains(string(manifest), "android:icon") || strings.Contains(string(manifest), "debuggable") ||
		strings.Contains(string(manifest), "POST_NOTIFICATIONS") || strings.Contains(string(manifest), "INTERNET") {
		t.Fatalf("an icon or debugging that was not asked for:\n%s", manifest)
	}
}

func TestArgumentsAreChecked(t *testing.T) {
	for _, arguments := range [][]string{
		{"-id", "io.example.tides"},                    // no name
		{"-name", "Tides", "-id", "tides"},             // one part
		{"-name", "Tides", "-id", "io.example.9tides"}, // a part starting with a digit
		{"-name", "Tides", "-id", "io.example.tides", "-min-sdk", "21"},
		{"-name", "Tides", "-id", "io.example.tides", "one", "two"},
	} {
		if _, _, err := readArguments(arguments); err == nil {
			t.Fatalf("%v accepted", arguments)
		}
	}
	app, output, err := readArguments([]string{"-name", "Tides", "-id", "io.example.tides", "./cmd/Tides"})
	if err != nil || app.Package != "./cmd/Tides" || output != "tides.apk" {
		t.Fatalf("package %q, output %q, error %v", app.Package, output, err)
	}
}

func TestTheNewestToolsAreChosen(t *testing.T) {
	folder := t.TempDir()
	for _, name := range []string{"9.0.0", "36.0.0", "34.0.0-rc1", "android-34", "android-36"} {
		if err := os.Mkdir(filepath.Join(folder, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if found, _ := newest(folder, ""); filepath.Base(found) != "android-36" && filepath.Base(found) != "36.0.0" {
		t.Fatalf("newest build tools: %s", found)
	}
	if found, _ := newest(folder, "android-"); filepath.Base(found) != "android-36" {
		t.Fatalf("newest platform: %s, want android-36", found)
	}
	if !versionBefore("9.0.0", "36.0.0") || versionBefore("36.0.1", "36.0.0") {
		t.Fatal("versions compare as text, not numbers")
	}
}

func TestTheWindowThemeFollowsNightMode(t *testing.T) {
	if !strings.Contains(string(windowTheme(false, "")), "windowLightStatusBar\">true") ||
		!strings.Contains(string(windowTheme(true, "")), "windowLightStatusBar\">false") {
		t.Fatal("dark icons on a light interface and light ones on a dark interface, please")
	}
}

func TestProcessorsAreChosenByName(t *testing.T) {
	app, _, err := readArguments([]string{"-name", "Tides", "-id", "io.example.tides", "-abi", "arm64, x86_64,arm64"})
	if err != nil {
		t.Fatal(err)
	}
	if len(app.ABIs) != 2 || app.ABIs[0].name != "arm64-v8a" || app.ABIs[1].goArch != "amd64" {
		t.Fatalf("processors %+v; want arm64 and x86_64, once each", app.ABIs)
	}
	if _, _, err := readArguments([]string{"-name", "Tides", "-id", "io.example.tides", "-abi", "mips"}); err == nil {
		t.Fatal("an unknown processor was taken")
	}
}

func TestAReleaseNeedsItsKeyAndIsNotDebuggable(t *testing.T) {
	keystore := filepath.Join(t.TempDir(), "release.keystore")
	if err := os.WriteFile(keystore, []byte("a store"), 0o600); err != nil {
		t.Fatal(err)
	}
	release := []string{"-name", "Tides", "-id", "io.example.tides", "-release", "-keystore", keystore, "-key-alias", "tides"}
	t.Setenv(storePasswordVariable, "")
	if _, _, err := readArguments(release); err == nil {
		t.Fatal("a release was taken with no password for its store")
	}
	t.Setenv(storePasswordVariable, "secret")
	app, _, err := readArguments(release)
	if err != nil {
		t.Fatal(err)
	}
	if app.Debuggable {
		t.Fatal("a release is debuggable")
	}
	for _, argument := range app.releaseSigning() {
		if strings.Contains(argument, "secret") {
			t.Fatal("the password is on apksigner's command line")
		}
	}
	if _, _, err := readArguments([]string{"-name", "Tides", "-id", "io.example.tides", "-release"}); err == nil {
		t.Fatal("a release was taken with no key")
	}
	app, _, err = readArguments(append(release, "-debuggable"))
	if err != nil || !app.Debuggable {
		t.Fatal("asking for a debuggable release was not heard")
	}
}

func TestTheNetworkIsAllowedUnlessLeftOut(t *testing.T) {
	app, _, err := readArguments([]string{"-name", "Tides", "-id", "io.example.tides"})
	if err != nil || !app.Internet {
		t.Fatalf("internet %v (%v) by default; want the network allowed", app.Internet, err)
	}
	app, _, err = readArguments([]string{"-name", "Tides", "-id", "io.example.tides", "-internet=false"})
	if err != nil || app.Internet {
		t.Fatal("-internet=false left the network allowed")
	}
}

// pictureFile writes a stand-in PNG (only its bytes are copied) and gives its path.
func pictureFile(t *testing.T, name string) string {
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("\x89PNG "+name), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAnAdaptiveIconNamesItsLayers(t *testing.T) {
	app := application{Name: "Tides", ID: "io.example.tides", VersionCode: 1, VersionName: "1.0", MinimumSDK: 26, TargetSDK: 34,
		IconForeground: pictureFile(t, "foreground.png"), IconBackground: "#1f5fc4", IconMonochrome: pictureFile(t, "mono.png")}
	if err := app.check(); err != nil {
		t.Fatal(err)
	}
	files, err := app.resources()
	if err != nil {
		t.Fatal(err)
	}
	icon := string(files["res/mipmap-anydpi-v26/icon.xml"])
	for _, wanted := range []string{"@color/icon_background", "@mipmap/icon_foreground", "@mipmap/icon_monochrome"} {
		if !strings.Contains(icon, wanted) {
			t.Fatalf("the adaptive icon has no %s:\n%s", wanted, icon)
		}
	}
	if !strings.Contains(string(files["res/values/icon_background.xml"]), "#1f5fc4") ||
		!strings.HasPrefix(string(files["res/mipmap-xxxhdpi/icon_foreground.png"]), "\x89PNG") {
		t.Fatal("the background colour or the foreground picture is not in the package")
	}
	if !strings.Contains(string(files["AndroidManifest.xml"]), `android:icon="@mipmap/icon"`) {
		t.Fatal("the manifest does not name the adaptive icon")
	}
}

func TestAPictureCanBeTheIconsBackground(t *testing.T) {
	app := application{Name: "Tides", ID: "io.example.tides", VersionCode: 1, VersionName: "1.0", MinimumSDK: 26, TargetSDK: 34,
		IconForeground: pictureFile(t, "foreground.png"), IconBackground: pictureFile(t, "sea.png")}
	files, err := app.resources()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(files["res/mipmap-anydpi-v26/icon.xml"]), "@mipmap/icon_background") ||
		files["res/mipmap-xxxhdpi/icon_background.png"] == nil || files["res/values/icon_background.xml"] != nil {
		t.Fatal("a picture behind the icon is not its background layer")
	}
}

func TestIconAndSplashFlagsAreChecked(t *testing.T) {
	base := []string{"-name", "Tides", "-id", "io.example.tides"}
	for _, wrong := range [][]string{
		{"-icon-background", "#ffffff"}, // a background with no foreground
		{"-icon-foreground", "fg.png"},  // a foreground with no background
		{"-splash", "blue"},             // not a colour
		{"-splash-dark", "#12345"},      // too short
	} {
		if _, _, err := readArguments(append(append([]string{}, base...), wrong...)); err == nil {
			t.Errorf("%v was taken", wrong)
		}
	}
	app, _, err := readArguments(append(base, "-splash", "#eef1f5", "-splash-dark", "#0f141b"))
	if err != nil {
		t.Fatal(err)
	}
	files, _ := app.resources()
	if !strings.Contains(string(files["res/values/styles.xml"]), "#eef1f5") || !strings.Contains(string(files["res/values-night/styles.xml"]), "#0f141b") {
		t.Fatal("the splash colours are not the window's background")
	}
}
