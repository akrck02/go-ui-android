package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// abi is a kind of processor a package carries a library for.
type abi struct {
	name     string // the package's folder for it: lib/<name>/
	goArch   string // GOARCH
	goArm    string // GOARM, for 32-bit arm
	compiler string // the NDK's clang, before the API level
}

// abis are the processors a package can carry: phones (64-bit and old 32-bit ones) and the
// emulator on a PC.
var abis = map[string]abi{
	"arm64":  {name: "arm64-v8a", goArch: "arm64", compiler: "aarch64-linux-android"},
	"x86_64": {name: "x86_64", goArch: "amd64", compiler: "x86_64-linux-android"},
	"arm":    {name: "armeabi-v7a", goArch: "arm", goArm: "7", compiler: "armv7a-linux-androideabi"},
}

// readABIs turns -abi's list (arm64,x86_64) into processors, each once.
func readABIs(list string) ([]abi, error) {
	var chosen []abi
	seen := map[string]bool{}
	for _, name := range strings.Split(list, ",") {
		name = strings.TrimSpace(name)
		processor, known := abis[name]
		if !known {
			return nil, fmt.Errorf("-abi %q is not one of arm64, x86_64, arm", name)
		}
		if !seen[name] {
			seen[name] = true
			chosen = append(chosen, processor)
		}
	}
	return chosen, nil
}

// The passwords of a release key come from the environment, never from the command line, which
// other users of the computer can read in its list of processes.
const (
	storePasswordVariable = "GO_UI_APK_STORE_PASSWORD"
	keyPasswordVariable   = "GO_UI_APK_KEY_PASSWORD" // when the key's differs from the store's
)

// checkRelease makes sure a release names a key that is there, and its store's password.
func (app application) checkRelease() error {
	if !app.Release {
		return nil
	}
	if app.Keystore == "" || app.KeyAlias == "" {
		return errors.New("-release needs -keystore and -key-alias: the key the store knows the program by")
	}
	if _, err := os.Stat(app.Keystore); err != nil {
		return fmt.Errorf("-keystore: %w", err)
	}
	if os.Getenv(storePasswordVariable) == "" {
		return fmt.Errorf("-release reads the key store's password from %s", storePasswordVariable)
	}
	return nil
}

// releaseSigning is apksigner's arguments for the release key.
func (app application) releaseSigning() []string {
	arguments := []string{"--ks", app.Keystore, "--ks-key-alias", app.KeyAlias, "--ks-pass", "env:" + storePasswordVariable}
	if os.Getenv(keyPasswordVariable) != "" {
		arguments = append(arguments, "--key-pass", "env:"+keyPasswordVariable)
	}
	return arguments
}
