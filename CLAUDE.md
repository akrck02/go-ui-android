# go-ui-android

Material Design 3 pieces for go-ui programs (`material`), a tool that packages a go-ui program as an
Android application (`cmd/go-ui-apk`), and a gallery of every piece (`cmd/gallery`). Start with
`docs/1-overview.md`.

## Commands

```bash
go vet ./...
go test ./...                       # no display needed: a fake measurer, fake pointer and fingers
go run ./cmd/gallery                # needs a display
go run ./cmd/gallery -capture frame.png
ANDROID_HOME=… ANDROID_NDK_HOME=… go run ./cmd/go-ui-apk -name Gallery -id io.github.akrck02.gallery ./cmd/gallery
```

It builds against `../go-ui` and `../go-renderer` (`replace` in `go.mod`).

## Rules

- **Generic only.** This project knows nothing about any particular world, game, company or story.
  Never put names, places, dates, numbers or settings from a user's own project in code, comments,
  tests, examples, docs or commit messages. Examples use made-up places (lighthouses on an island).
- **Readable code.** Full-word names (no acronyms, no one-letter variables), small functions that do
  one thing, early returns instead of nested branches, comments that say *why*. Match the code
  around you.
- **Tests with every change.** `go vet ./...` and `go test ./...` must pass before committing.
- **Docs live in `docs/`**, written for people: explanatory, friendly, with images where they help.
  Keep the README short and point to `docs/`. Update the docs in the same commit as the behaviour.
- **AI disclosure.** Commits written with an AI assistant end with a `Co-Authored-By: Claude …` line,
  and the *AI disclosure* section of the README stays true. The code is pending human review.
- **Git.** Commit only when asked; never push, force-push, rewrite history or delete branches without
  the author's explicit permission. Do not change the global git config.
- **What belongs elsewhere.** Window, input or Android changes go to go-renderer or go-ui, generic and
  on a `feature/` branch there; this repository only uses their public packages.

## Human and AI code

- `main` holds hand-written code only. AI work goes on feature branches (`feature/…`) and is merged by
  the author after review.
