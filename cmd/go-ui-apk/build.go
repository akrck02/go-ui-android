package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// sdk is where the Android tools are.
type sdk struct {
	buildTools string // aapt2, zipalign, apksigner
	platform   string // android.jar
	compilers  string // the NDK's folder of clang compilers
}

// compiler is the NDK's clang for a processor and the oldest Android the package runs on.
func (tools sdk) compiler(processor abi, minimumSDK int) string {
	return filepath.Join(tools.compilers, fmt.Sprintf("%s%d-clang", processor.compiler, minimumSDK))
}

// build makes the package: the library, the resources and manifest, then aligned and signed.
func build(app application, output string) error {
	tools, err := findSDK()
	if err != nil {
		return err
	}
	work, err := os.MkdirTemp("", "go-ui-apk-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	libraries := map[string]string{} // the package's name for each library (and the classes), and where it was built
	for _, processor := range app.ABIs {
		name := "lib/" + processor.name + "/lib" + libraryName + ".so"
		library := filepath.Join(work, filepath.FromSlash(name))
		if err := buildLibrary(app.Package, library, processor, tools.compiler(processor, app.MinimumSDK)); err != nil {
			return err
		}
		libraries[name] = library
	}
	dex, err := buildDex(tools, work, app.MinimumSDK)
	if err != nil {
		return err
	}
	libraries[dexName] = dex
	if app.Bundle {
		return buildBundle(app, tools, work, libraries, output)
	}
	unaligned := filepath.Join(work, "unaligned.apk")
	if err := packResources(app, tools, work, unaligned, false); err != nil {
		return err
	}
	if err := addLibraries(unaligned, libraries); err != nil {
		return err
	}
	aligned := filepath.Join(work, "aligned.apk")
	if err := run(filepath.Join(tools.buildTools, "zipalign"), "-f", "-p", "16", unaligned, aligned); err != nil {
		return err
	}
	return sign(app, tools, aligned, output)
}

// buildLibrary builds the program as a shared library for a processor, aligned to 16 KB pages,
// which phones with such pages (Android 15 and later) need.
func buildLibrary(pkg, library string, processor abi, compiler string) error {
	if _, err := os.Stat(compiler); err != nil {
		return fmt.Errorf("the NDK has no compiler for %s: %w", processor.name, err)
	}
	command := exec.Command("go", "build", "-buildmode=c-shared", "-trimpath",
		"-ldflags=-extldflags=-Wl,-z,max-page-size=16384", "-o", library, pkg)
	// C++ too, for the libraries that carry some (go-sound's audio output, through Oboe)
	command.Env = append(os.Environ(), "CGO_ENABLED=1", "GOOS=android", "GOARCH="+processor.goArch, "CC="+compiler, "CXX="+compiler+"++")
	if processor.goArm != "" {
		command.Env = append(command.Env, "GOARM="+processor.goArm)
	}
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("building %s for Android: %w", pkg, err)
	}
	// cgo's header for the library's exports, when it writes one: the package needs none
	if err := os.Remove(strings.TrimSuffix(library, ".so") + ".h"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// packResources writes the manifest, the window theme and the icon, and links them into a package:
// in the binary form a phone reads, or in the protocol buffers an App Bundle holds (proto).
func packResources(app application, tools sdk, work, unaligned string, proto bool) error {
	files, err := app.resources()
	if err != nil {
		return err
	}
	for name, contents := range files {
		path := filepath.Join(work, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, contents, 0o644); err != nil {
			return err
		}
	}
	aapt := filepath.Join(tools.buildTools, "aapt2")
	compiled := filepath.Join(work, "resources.zip")
	if err := run(aapt, "compile", "--dir", filepath.Join(work, "res"), "-o", compiled); err != nil {
		return err
	}
	arguments := []string{"link", "-o", unaligned, "-I", filepath.Join(tools.platform, "android.jar"),
		"--manifest", filepath.Join(work, "AndroidManifest.xml"), "-R", compiled, "--auto-add-overlay"}
	if proto {
		arguments = append(arguments, "--proto-format")
	}
	return run(aapt, arguments...)
}

// addLibraries adds the libraries and the classes to the package, stored rather than compressed:
// Android maps the libraries straight from the package.
func addLibraries(packagePath string, libraries map[string]string) error {
	original, err := zip.OpenReader(packagePath)
	if err != nil {
		return err
	}
	defer original.Close()
	rewritten := packagePath + ".new"
	file, err := os.Create(rewritten)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	for _, entry := range original.File {
		if err := writer.Copy(entry); err != nil {
			return err
		}
	}
	names := make([]string, 0, len(libraries))
	for name := range libraries {
		names = append(names, name)
	}
	sort.Strings(names) // the same package from the same program
	for _, name := range names {
		if err := storeFile(writer, name, libraries[name]); err != nil {
			return err
		}
	}
	if err := writer.Close(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(rewritten, packagePath)
}

func storeFile(writer *zip.Writer, name, path string) error {
	source, err := os.Open(path)
	if err != nil {
		return err
	}
	defer source.Close()
	target, err := writer.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Store})
	if err != nil {
		return err
	}
	_, err = io.Copy(target, source)
	return err
}

// sign signs a release with the program's own key, and anything else with the debug key.
func sign(app application, tools sdk, aligned, output string) error {
	signing, err := debugSigning()
	if app.Release {
		signing, err = app.releaseSigning(), nil
	}
	if err != nil {
		return err
	}
	arguments := append([]string{"sign"}, signing...)
	arguments = append(arguments, "--out", output, aligned)
	if err := run(filepath.Join(tools.buildTools, "apksigner"), arguments...); err != nil {
		return err
	}
	_ = os.Remove(output + ".idsig") // apksigner's v4 signature, for streaming installs only
	return nil
}

// debugSigning is apksigner's arguments for the debug key, made the first time as Android's own
// tools do.
func debugSigning() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	keystore := filepath.Join(home, ".android", "debug.keystore")
	if _, err := os.Stat(keystore); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(keystore), 0o755); err != nil {
			return nil, err
		}
		if err := run("keytool", "-genkeypair", "-keystore", keystore, "-storepass", "android",
			"-alias", "androiddebugkey", "-keypass", "android", "-keyalg", "RSA", "-validity", "10000",
			"-dname", "CN=Android Debug,O=Android,C=US"); err != nil {
			return nil, err
		}
	}
	// the debug key's passwords are Android's well-known ones
	return []string{"--ks", keystore, "--ks-pass", "pass:android", "--ks-key-alias", "androiddebugkey", "--key-pass", "pass:android"}, nil
}

func run(program string, arguments ...string) error {
	command := exec.Command(program, arguments...)
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(program), err)
	}
	return nil
}

// findSDK finds the newest build tools and platform in ANDROID_HOME, and the NDK's compilers.
func findSDK() (sdk, error) {
	home, ndk := os.Getenv("ANDROID_HOME"), os.Getenv("ANDROID_NDK_HOME")
	if home == "" || ndk == "" {
		return sdk{}, errors.New("set ANDROID_HOME to the Android SDK and ANDROID_NDK_HOME to the NDK")
	}
	buildTools, err := newest(filepath.Join(home, "build-tools"), "")
	if err != nil {
		return sdk{}, fmt.Errorf("no build tools in the SDK: %w", err)
	}
	platform, err := newest(filepath.Join(home, "platforms"), "android-")
	if err != nil {
		return sdk{}, fmt.Errorf("no platform in the SDK: %w", err)
	}
	compilers := filepath.Join(ndk, "toolchains", "llvm", "prebuilt", runtime.GOOS+"-x86_64", "bin")
	if _, err := os.Stat(compilers); err != nil {
		return sdk{}, fmt.Errorf("the NDK has no compilers for this computer: %w", err)
	}
	return sdk{buildTools: buildTools, platform: platform, compilers: compilers}, nil
}

// newest is the folder whose name, after a prefix, is the highest version (36.0.0 above 9.0.0).
func newest(folder, prefix string) (string, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return "", err
	}
	var versions []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			versions = append(versions, strings.TrimPrefix(entry.Name(), prefix))
		}
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("nothing in %s", folder)
	}
	sort.Slice(versions, func(first, second int) bool { return versionBefore(versions[first], versions[second]) })
	return filepath.Join(folder, prefix+versions[len(versions)-1]), nil
}

// versionBefore compares dotted versions part by part as numbers; a part that is not a number
// (a preview's "rc1") counts as zero.
func versionBefore(first, second string) bool {
	firstParts, secondParts := strings.Split(first, "."), strings.Split(second, ".")
	for index := 0; index < max(len(firstParts), len(secondParts)); index++ {
		firstNumber, secondNumber := part(firstParts, index), part(secondParts, index)
		if firstNumber != secondNumber {
			return firstNumber < secondNumber
		}
	}
	return false
}

func part(parts []string, index int) int {
	if index >= len(parts) {
		return 0
	}
	number, _ := strconv.Atoi(parts[index])
	return number
}
