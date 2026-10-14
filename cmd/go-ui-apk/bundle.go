package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// An App Bundle is a zip of modules: here one, base, holding the manifest and the resources in
// protocol buffers (as aapt2 writes them with --proto-format), the libraries, and BundleConfig.pb,
// which names the bundletool version the bundle follows. Google Play makes the packages for each
// phone from it, and signs them with its own key; the bundle itself is signed with the upload key.

// bundletoolVersion is the version of bundletool whose layout the bundle follows.
const bundletoolVersion = "1.17.2"

// buildBundle links the resources in protocol buffers, lays them out as the base module with the
// libraries, and signs the bundle.
func buildBundle(app application, tools sdk, work string, libraries map[string]string, output string) error {
	linked := filepath.Join(work, "linked.zip")
	if err := packResources(app, tools, work, linked, true); err != nil {
		return err
	}
	if err := writeBundle(linked, libraries, output); err != nil {
		return err
	}
	return signBundle(app, output)
}

// writeBundle moves what aapt2 linked into the base module (the manifest under manifest/) and adds
// the libraries and the bundle's configuration.
func writeBundle(linked string, libraries map[string]string, output string) error {
	source, err := zip.OpenReader(linked)
	if err != nil {
		return err
	}
	defer source.Close()
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := zip.NewWriter(file)
	for _, entry := range source.File {
		header := entry.FileHeader
		header.Name = bundlePath(entry.Name)
		if err := copyEntry(writer, entry, &header); err != nil {
			return err
		}
	}
	names := make([]string, 0, len(libraries))
	for name := range libraries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := storeFile(writer, moduleFile(name), libraries[name]); err != nil {
			return err
		}
	}
	configuration, err := writer.Create("BundleConfig.pb")
	if err != nil {
		return err
	}
	if _, err := configuration.Write(bundleConfig(bundletoolVersion)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return file.Close()
}

// bundlePath is where a file aapt2 linked goes in the base module.
func bundlePath(name string) string {
	if name == "AndroidManifest.xml" {
		return "base/manifest/AndroidManifest.xml"
	}
	if name == "resources.pb" || strings.HasPrefix(name, "res/") {
		return "base/" + name
	}
	return "base/root/" + name // anything else travels as a file of the package's root
}

// moduleFile is where a file of a package goes in the base module: the classes under dex/, the
// libraries as they are.
func moduleFile(name string) string {
	if name == dexName {
		return "base/dex/" + name
	}
	return "base/" + name
}

// copyEntry copies a zip entry under another name, without compressing it again.
func copyEntry(writer *zip.Writer, entry *zip.File, header *zip.FileHeader) error {
	raw, err := entry.OpenRaw()
	if err != nil {
		return err
	}
	target, err := writer.CreateRaw(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(target, raw)
	return err
}

// bundleConfig is BundleConfig.pb: BundleConfig{bundletool: Bundletool{version}}, in protocol
// buffers' wire format (field 1, a message; inside it field 2, a string).
func bundleConfig(version string) []byte {
	inner := append([]byte{0x12, byte(len(version))}, version...)
	return append([]byte{0x0a, byte(len(inner))}, inner...)
}

// signBundle signs the bundle with the release key, as Google Play wants an upload signed:
// jarsigner reads the passwords from the same environment variables apksigner does.
func signBundle(app application, bundle string) error {
	arguments := []string{"-keystore", app.Keystore, "-storepass:env", storePasswordVariable}
	if os.Getenv(keyPasswordVariable) != "" {
		arguments = append(arguments, "-keypass:env", keyPasswordVariable)
	}
	return run("jarsigner", append(arguments, bundle, app.KeyAlias)...)
}
