package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akrck02/go-renderer/android/java"
)

// The package carries go-renderer's Java classes: ScreenReader, what TalkBack reads of the window,
// and TextInput, the view the on-screen keyboard types in, which Android only lets Java answer. They
// are compiled here from their source, against the platform's android.jar, and turned into
// Android's classes.dex by the build tools' d8.

// dexName is the classes' name in a package.
const dexName = "classes.dex"

// buildDex compiles the Java the program carries into classes.dex, in a folder of work, and tells
// where it is.
func buildDex(tools sdk, work string, minimumSDK int) (string, error) {
	javac, err := javaCompiler()
	if err != nil {
		return "", err
	}
	sources, err := writeSources(filepath.Join(work, "java"))
	if err != nil {
		return "", err
	}
	platform := filepath.Join(tools.platform, "android.jar")
	classes := filepath.Join(work, "classes")
	// Java 11's classes, which d8 reads
	if err := run(javac, append([]string{"--release", "11", "-classpath", platform, "-d", classes}, sources...)...); err != nil {
		return "", err
	}
	compiled, err := classFiles(classes)
	if err != nil {
		return "", err
	}
	dex := filepath.Join(work, "dex")
	if err := os.MkdirAll(dex, 0o755); err != nil {
		return "", err
	}
	arguments := append([]string{"--release", "--min-api", strconv.Itoa(minimumSDK), "--lib", platform, "--output", dex}, compiled...)
	if err := run(filepath.Join(tools.buildTools, "d8"), arguments...); err != nil {
		return "", err
	}
	return filepath.Join(dex, dexName), nil
}

// writeSources writes the Java go-renderer carries under a folder, for javac, and lists the files.
func writeSources(folder string) ([]string, error) {
	var written []string
	for _, source := range java.Sources {
		path := filepath.Join(folder, filepath.FromSlash(source.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(source.Text), 0o644); err != nil {
			return nil, err
		}
		written = append(written, path)
	}
	return written, nil
}

// javaCompiler finds javac: in JAVA_HOME, else on the path. A JDK is already needed for apksigner
// and keytool, which run on Java.
func javaCompiler() (string, error) {
	if home := os.Getenv("JAVA_HOME"); home != "" {
		javac := filepath.Join(home, "bin", "javac")
		if _, err := os.Stat(javac); err == nil {
			return javac, nil
		}
	}
	javac, err := exec.LookPath("javac")
	if err != nil {
		return "", errors.New("no javac: install a JDK (17 or later), or set JAVA_HOME to one")
	}
	return javac, nil
}

// classFiles lists the .class files javac wrote: a class and the ones inside it.
func classFiles(folder string) ([]string, error) {
	var found []string
	err := filepath.WalkDir(folder, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".class") {
			found = append(found, path)
		}
		return nil
	})
	if err == nil && len(found) == 0 {
		err = fmt.Errorf("javac wrote no classes in %s", folder)
	}
	return found, err
}
