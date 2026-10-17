# 📦 Packaging: go-ui-apk

```bash
export ANDROID_HOME=…/android-sdk      # with build-tools and a platform
export ANDROID_NDK_HOME=…/ndk/<version>
# and a JDK, 17 or later: javac on the path or in JAVA_HOME
go run github.com/akrck02/go-ui-android/cmd/go-ui-apk -name Gallery -id io.example.gallery ./cmd/gallery
```

It builds the program as a shared library for each processor asked for (arm64 phones by default) with the NDK's compiler (aligned to 16 KB
pages, which Android 15 and later phones need), writes the manifest and a window theme that lets the
program paint under the system's bars (light or dark with the phone) and beside a camera's cutout,
compiles go-renderer's two Java classes (what a screen reader reads, and the view the on-screen
keyboard types in) with `javac` and the build tools' `d8`, and packs, aligns and signs the package with the SDK's own `aapt2`, `zipalign` and `apksigner`. The
window shrinks while the on-screen keyboard is up. It signs with the debug key in
`~/.android/debug.keystore`, making it when missing, for phones and the emulator; a release is
signed with the program's own key (below).

| Flag | What it sets |
|---|---|
| `-name` | The name under the icon (required) |
| `-id` | The application id, like `io.example.app` (required) |
| `-o` | The package to write; the program's folder name `.apk` by default |
| `-version-code`, `-version` | The version's number and the version people read |
| `-min-sdk`, `-target-sdk` | The oldest Android it runs on (26, Android 8, at least) and the one it is made for |
| `-icon` | A square PNG for the launcher |
| `-icon-foreground`, `-icon-background`, `-icon-monochrome` | An adaptive icon, which each launcher cuts to its own shape: a PNG of 108 dp whose middle 66 dp always shows, a colour (`#rrggbb`) or a PNG behind it, and a one-colour PNG for themed icons (Android 13+) |
| `-splash`, `-splash-dark` | The colour (`#rrggbb`) the window shows while the program starts, in the light and the dark theme; Android 12 and later draw the icon over it |
| `-debuggable` | Lets `adb` read its folder and debug it; on by default |
| `-notifications` | Names the permission Android 13 and later ask before showing notifications |
| `-internet` | Lets the program use the network: on by default, granted when the package is installed with no question; `-internet=false` leaves it out |
| `-abi` | The processors to carry a library for: `arm64` (phones, the default), `x86_64` (the emulator on a PC), `arm` (old 32-bit phones); `-abi arm64,x86_64` for several |
| `-release`, `-keystore`, `-key-alias` | Signs with the program's own key instead of the debug key, and leaves the package not debuggable: for a store |
| `-bundle` | With `-release`, writes an App Bundle (`.aab`), which Google Play asks for and makes each phone's packages from |

**A release** is signed with a key of the program's own, which a store knows it by from then on:
keep the key store safe and out of the repository. The store's password is read from
`GO_UI_APK_STORE_PASSWORD` (and the key's, when different, from `GO_UI_APK_KEY_PASSWORD`), never
from the command line, which other users of the computer can see:

```bash
keytool -genkeypair -keystore lighthouses.keystore -alias lighthouses -keyalg RSA -validity 10000
GO_UI_APK_STORE_PASSWORD=… go run ./cmd/go-ui-apk -name Lighthouses -id io.example.lighthouses \
    -release -keystore lighthouses.keystore -key-alias lighthouses -version-code 2 -version 1.1 ./cmd/lighthouses
```

**An App Bundle** (`-bundle`) is what Google Play asks for: the same program, its resources in
protocol buffers and its libraries in one `base` module, signed with the upload key; Play makes
each phone's packages from it. Google's `bundletool` validates it and makes those packages to try
on a phone (`bundletool build-apks --connected-device`, then `install-apks`); the gallery's bundle
was tried that way on the emulator.

What the program writes to its standard error (a panic, `log.Fatal`) is in Android's log:
`adb logcat -s Go`.

---

[← All the documentation](../1-overview.md)
