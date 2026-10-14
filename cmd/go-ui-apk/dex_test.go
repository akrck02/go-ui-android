package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestTheClassesGoInTheBaseModulesDexFolder(t *testing.T) {
	for name, want := range map[string]string{
		"classes.dex":                "base/dex/classes.dex",
		"lib/arm64-v8a/libmain.so":   "base/lib/arm64-v8a/libmain.so",
		"lib/armeabi-v7a/libmain.so": "base/lib/armeabi-v7a/libmain.so",
	} {
		if got := moduleFile(name); got != want {
			t.Errorf("%s goes to %s; want %s", name, got, want)
		}
	}
}

func TestJavacIsFoundInJavaHome(t *testing.T) {
	home := t.TempDir()
	javac := filepath.Join(home, "bin", "javac")
	if err := os.MkdirAll(filepath.Dir(javac), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(javac, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("JAVA_HOME", home)
	if found, err := javaCompiler(); err != nil || found != javac {
		t.Fatalf("found %q (%v); want JAVA_HOME's %q", found, err, javac)
	}
}

// The class compiles for real when the computer has a JDK and the Android SDK.
func TestTheScreenReaderCompilesIntoClassesDex(t *testing.T) {
	if _, err := javaCompiler(); err != nil {
		t.Skip(err)
	}
	if os.Getenv("ANDROID_HOME") == "" || os.Getenv("ANDROID_NDK_HOME") == "" {
		t.Skip("no Android SDK")
	}
	tools, err := findSDK()
	if err != nil {
		t.Skip(err)
	}
	dex, err := buildDex(tools, t.TempDir(), 26)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(dex)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(contents, []byte("dex\n")) || !bytes.Contains(contents, []byte("Lio/github/akrck02/gorenderer/ScreenReader;")) {
		t.Fatal("classes.dex holds the screen reader's class")
	}
}
