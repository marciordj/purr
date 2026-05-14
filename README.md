# purr

A terminal UI for converting `.mov` files to `.mp4`. Wraps `ffmpeg` with a clean progress view.

## How it works

Files that use H.264/HEVC codecs are **remuxed** (container swap only — instant, lossless). Files with incompatible codecs (e.g. ProRes) fall back to **re-encode** with H.264/AAC, CRF 18.

## Requirements

- Go 1.24+
- `ffmpeg` and `ffprobe` in PATH — install with `brew install ffmpeg`

## Install

```sh
go install github.com/marciord/purr/cmd/purr@latest
```

Or build from source:

```sh
git clone https://github.com/marciord/purr
cd purr
go build -o purr ./cmd/purr
```

## Usage

```sh
purr file.mov
purr *.mov
purr recording1.mov recording2.mov
```

Output files are written next to the input (e.g. `video.mov` → `video.mp4`). Existing `.mp4` files are skipped.

## Keybindings

| Key      | Action                                |
|----------|---------------------------------------|
| `esc`    | Cancel active job (confirm prompt)    |
| `ctrl+c` | Abort all jobs immediately            |
| `q`      | Quit (only when no job is running)    |
