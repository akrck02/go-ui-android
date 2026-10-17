<p align="center"><img alt="" src="images/logo.svg" width="120"></p>

# Welcome to go-ui-android 👋

**go-ui-android makes Android apps in Go that look and feel like Android's own**: Material Design pieces (bars, lists, buttons, dialogs) and a tool that turns the program into an app you install on a phone. No Java, no Android Studio.

<p align="center"><img alt="Three phone screens of the gallery: buttons, a list of lighthouses with photos, tabs over pictures of a coast" src="images/hero-light.png" width="720"></p>

You write the app once in Go with go-ui; the `material` pieces give it Android's looks (the wallpaper's colours, the back gesture, the on-screen keyboard), and `go-ui-apk` packs it into an `.apk`. The same program also opens on a computer, to try it as you go.

## Is it for me?

- ✅ You write Go and want an **Android app**.
- ✅ You want the app to feel **at home on a phone**: gestures, dark theme, TalkBack.
- ✅ You already know a little of [go-ui](https://github.com/akrck02/go-ui).

## Where to start

Pick the line that sounds like you.

| If you… | Read | Time |
|---|---|---|
| just want to **see** every piece | [Try it without writing code](#try-it-without-writing-code) | 2 min |
| are **new** and want the first lines | [Getting started](2-getting-started.md) | 10 min |
| want to **learn** step by step, up to a phone | 🏴‍☠️ [The quest](3-the-quest/0-the-quest.md): a lighthouse app, in 5 levels | 90 min |
| want to know about colours, gestures, screen readers or packaging | [The guides](4-guides/1-pieces.md) | — |
| want to use **one piece** (a top bar, a dialog, a carousel) | [Components](5-components/0-components.md): each one with real examples | — |
| are **looking up** a function | [The reference](6-reference/0-reference.md) | — |

## Try it without writing code

With Go installed and the projects copied (see [Getting started](2-getting-started.md)), the gallery shows every piece in a phone-sized window:

```sh
go run ./cmd/gallery
```

## All the pages, in reading order

The menu on the left (or the files in `docs/`) follows this order.

**Start**

- [Overview](1-overview.md): this page
- [Getting started](2-getting-started.md)

**Learn**

- [The quest](3-the-quest/0-the-quest.md): five levels, from a first screen to a phone

**Guides**

- [The pieces](4-guides/1-pieces.md): the whole catalogue at a glance
- [Colours](4-guides/2-colours.md)
- [Feeling at home on a phone](4-guides/3-phone.md)
- [Screen readers](4-guides/4-screen-readers.md)
- [Packaging: go-ui-apk](4-guides/5-packaging.md)
- [The gallery](4-guides/6-gallery.md)

**Use**

- [Components](5-components/0-components.md): every piece, one by one

**Reference**

- [The reference](6-reference/0-reference.md)

## 🚧 Not done yet

- TalkBack itself, on a phone (see [screen readers](4-guides/4-screen-readers.md)).
- Keyboards for Chinese and Japanese (they compose as Android's own keyboard does, which is tried),
  right-to-left layout and Arabic letters joined, and emoji joined into one (a
  family, a skin tone): other scripts and colour emoji are drawn with the phone's own fonts.
- Sound in films: a film plays its pictures only (go-sound plays other sounds).

## Words you will meet

| Word | What it means |
|---|---|
| Go | A programming language. To try these projects you install it once, from go.dev, and type `go run …` in a terminal. |
| Terminal | The window where you type commands: *Terminal* on macOS and Linux, *PowerShell* on Windows. |
| Repository | A project's folder of files, kept on GitHub. `git clone <address>` copies it to your computer. |
| `.apk` file | An Android app in one file, ready to install on a phone. |
| Immediate mode | A way of making interfaces: every moment the program describes the whole screen again from its data, so the screen can never show something out of date. |
| Quest | The tutorial of these docs: short levels, each ending with something that works and an achievement. |

## The family

go-ui-android is one of a family of small projects that work together:

| Project | What it does |
|---|---|
| [go-renderer](https://github.com/akrck02/go-renderer) | draws 3D scenes and flat interfaces |
| [go-ui](https://github.com/akrck02/go-ui) | buttons, panels and whole applications |
| **go-ui-android** (this one) | Android apps that look like Android's own |
| [go-sound](https://github.com/akrck02/go-sound) | sound and music |
| [go-scene-builder](https://github.com/akrck02/go-scene-builder) | 3D scenes from a flat map |
| [go-simulator](https://github.com/akrck02/go-simulator) | simulations that change a scene |
| [go-simulator-terrain](https://github.com/akrck02/go-simulator-terrain) | ground shaped by rivers |
| [go-simulator-climate](https://github.com/akrck02/go-simulator-climate) | weather and biomes |
| [go-simulator-flora](https://github.com/akrck02/go-simulator-flora) | forests that grow |
| [go-simulator-urban](https://github.com/akrck02/go-simulator-urban) | towns that grow |
| [go-simulator-motion](https://github.com/akrck02/go-simulator-motion) | a scene that moves |
| [go-3ds](https://github.com/akrck02/go-3ds) | Go on the Nintendo 3DS |

> ⚠️ All of them are **educational experiments**, written with an AI assistant and pending review by their author: good for learning and playing, not for production.
