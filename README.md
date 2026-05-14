# purr

A terminal UI for converting `.mov` files to `.mp4`. Wraps `ffmpeg` with a clean progress view.

## How it works

By default, all files are **compressed** using HEVC (H.265) via macOS VideoToolbox hardware encoding at quality level 65 — typically ~70% smaller than the original with no perceptible quality loss. Audio is copied as-is if already AAC; otherwise re-encoded to AAC 192k.

Output files use the `hvc1` codec tag for compatibility with QuickTime, Safari, and iOS.

With `--fast`, compatible-codec files (H.264, HEVC, VP9…) are **remuxed** (container swap only — instant, lossless, original size preserved). Incompatible codecs (e.g. ProRes) still go through compression.

## Requirements

- Go 1.24+
- `ffmpeg` and `ffprobe` in PATH — install with `brew install ffmpeg`

## Install

```sh
git clone https://github.com/marciord/purr
cd purr
make install
```

This builds and copies the binary to `~/.local/bin/purr`. Ensure `~/.local/bin` is in your `PATH`.

To update after code changes: `make install` again. To uninstall: `make uninstall`.

## Usage

```sh
purr             # open file picker in current directory
purr file.mov    # convert directly (no picker)
purr *.mov       # convert all matched files
purr --fast *.mov  # remux instead of compress (instant, lossless, original size)
```

Output files are written next to the input (e.g. `video.mov` → `video.mp4`). Existing `.mp4` files are skipped.

## Keybindings

**File picker**

| Key          | Action                  |
|--------------|-------------------------|
| `↑`/`↓`, `k`/`j` | Move cursor        |
| `space`      | Toggle selection        |
| `enter`      | Start conversion        |
| `q` / `esc`  | Quit                    |

**Conversion**

| Key      | Action                                |
|----------|---------------------------------------|
| `esc`    | Cancel active job (confirm prompt)    |
| `ctrl+c` | Abort all jobs immediately            |
| `q`      | Quit (only when no job is running)    |
