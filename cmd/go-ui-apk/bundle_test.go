package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestTheBundleConfigurationNamesItsBundletool(t *testing.T) {
	// BundleConfig{bundletool: Bundletool{version: "1.17.2"}}
	want := []byte{0x0a, 0x08, 0x12, 0x06, '1', '.', '1', '7', '.', '2'}
	if got := bundleConfig("1.17.2"); !bytes.Equal(got, want) {
		t.Fatalf("BundleConfig.pb is % x; want % x", got, want)
	}
}

func TestLinkedFilesGoInTheBaseModule(t *testing.T) {
	for name, want := range map[string]string{
		"AndroidManifest.xml":                    "base/manifest/AndroidManifest.xml",
		"resources.pb":                           "base/resources.pb",
		"res/mipmap-xxxhdpi/icon_foreground.png": "base/res/mipmap-xxxhdpi/icon_foreground.png",
		"META-INF/notes.txt":                     "base/root/META-INF/notes.txt",
	} {
		if got := bundlePath(name); got != want {
			t.Errorf("%s goes to %s; want %s", name, got, want)
		}
	}
}

func TestABundleIsARelease(t *testing.T) {
	base := []string{"-name", "Tides", "-id", "io.example.tides", "-bundle"}
	if _, _, err := readArguments(base); err == nil {
		t.Fatal("a bundle was taken without -release")
	}
	keystore := filepath.Join(t.TempDir(), "release.keystore")
	if err := os.WriteFile(keystore, []byte("a store"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(storePasswordVariable, "secret")
	app, output, err := readArguments(append(base, "-release", "-keystore", keystore, "-key-alias", "tides"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Ext(output) != ".aab" || app.Debuggable {
		t.Fatalf("output %s, debuggable %v; want an .aab that is not debuggable", output, app.Debuggable)
	}
}
