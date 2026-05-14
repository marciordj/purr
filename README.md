# purr

A terminal UI for converting `.mov` files to `.mp4`. Wraps `ffmpeg` with a clean progress view.

## How it works

Files that use H.264/HEVC codecs are **remuxed** (container swap only — instant, lossless). Files with incompatible codecs (e.g. ProRes) fall back to **re-encode** with H.264/AAC, CRF 18.

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
