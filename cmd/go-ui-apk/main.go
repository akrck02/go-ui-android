// Command go-ui-apk builds a go-ui program as an Android package (.apk), ready to install on a
// phone or the emulator.
//
//	go-ui-apk -name "Lighthouses" -id io.example.lighthouses -o lighthouses.apk ./cmd/lighthouses
//
// The program is the same main package that opens a window on the desktop (go-ui's window.Run).
// On Android it is a shared library that Android's own NativeActivity loads: no Java, no Gradle.
// The tool builds the library with the NDK's compiler, writes the manifest and a window theme that
// lets the program paint under the system's bars, and packs, aligns and signs the package with the
// Android SDK's own tools. It needs ANDROID_HOME (the SDK, with build-tools and a platform) and
// ANDROID_NDK_HOME (the NDK).
//
// The package is signed with the debug key in ~/.android/debug.keystore, made when missing: for
// trying on a device, not for a store.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app, output, err := readArguments(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "go-ui-apk:", err)
		os.Exit(2)
	}
	if err := build(app, output); err != nil {
		fmt.Fprintln(os.Stderr, "go-ui-apk:", err)
		os.Exit(1)
	}
	fmt.Println("built", output)
}

// readArguments reads the flags and the program's package (the current folder when none).
func readArguments(arguments []string) (application, string, error) {
	flags := flag.NewFlagSet("go-ui-apk", flag.ContinueOnError)
	var app application
	var output string
	flags.StringVar(&app.Name, "name", "", "the name under the icon (required)")
	flags.StringVar(&app.ID, "id", "", "the application id, like io.example.app (required)")
	flags.StringVar(&output, "o", "", "the package to write (default: the program's folder name .apk)")
	flags.IntVar(&app.VersionCode, "version-code", 1, "the version's number, higher for each release")
	flags.StringVar(&app.VersionName, "version", "1.0", "the version as people read it")
	flags.IntVar(&app.MinimumSDK, "min-sdk", 26, "the oldest Android API level it runs on (26 is Android 8)")
	flags.IntVar(&app.TargetSDK, "target-sdk", 34, "the Android API level it is made for")
	flags.StringVar(&app.Icon, "icon", "", "a square PNG for the launcher (default: Android's own)")
	flags.BoolVar(&app.Debuggable, "debuggable", true, "let adb read its folder (run-as) and debug it")
	if err := flags.Parse(arguments); err != nil {
		return application{}, "", err
	}
	app.Package = "."
	if flags.NArg() > 1 {
		return application{}, "", fmt.Errorf("one program at a time, got %v", flags.Args())
	}
	if flags.NArg() == 1 {
		app.Package = flags.Arg(0)
	}
	if err := app.check(); err != nil {
		return application{}, "", err
	}
	if output == "" {
		folder, err := filepath.Abs(app.Package)
		if err != nil {
			return application{}, "", err
		}
		output = strings.ToLower(filepath.Base(folder)) + ".apk"
	}
	return app, output, nil
}
