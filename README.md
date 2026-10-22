# purr

purr converts QuickTime `.mov` and `.qt` files to `.mp4` from the terminal. It wraps `ffmpeg` with an interactive file picker and a progress view, and re-encodes recordings to HEVC using macOS VideoToolbox hardware encoding.

Português (Brasil): [README.pt-BR.md](README.pt-BR.md)

## Requirements

- **macOS** for the default compression path (Apple's VideoToolbox encoder, already included in Homebrew's `ffmpeg`). The `--fast` remux path works wherever `ffmpeg`/`ffprobe` are available.
- `ffmpeg` and `ffprobe` in your `PATH` — `brew install ffmpeg` on macOS. purr exits right away if either is missing.
- Go 1.26.3 or newer (the version declared in `go.mod`), only to build the binary.

## Quickstart

```sh
git clone https://github.com/marciordj/purr
cd purr
make build          # produces ./purr
./purr              # open the picker in the current directory
```

The picker lists only `.mov` and `.qt` files in the current directory. To convert a specific video directly, pass it as an argument:

```sh
./purr ~/Desktop/recording.mov
```

If you don't have a `.mov` at hand, create a screen recording (QuickTime Player → File → New Screen Recording) and pass the saved file's path to `./purr` while still in the project directory.

## Optional install

```sh
make install        # builds and copies the binary to ~/.local/bin/purr
make uninstall      # removes it
```

Make sure `~/.local/bin` is in your `PATH`; if it isn't, add to your `~/.zshrc` (or `~/.bashrc`):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

After code changes run `make install` again to update the installed binary.

## Usage

```sh
purr                    # file picker (current directory only, no subdirectories)
purr video.mov          # convert one file directly, no picker
purr video.qt           # QuickTime .qt files work the same way
purr clip1.mov clip2.qt  # convert several files, in order (mixing .mov and .qt is fine)
purr --fast video.mov   # remux without re-encoding (-f is the same)
```

- Arguments must be existing files whose extension is `.mov` or `.qt` (case-insensitive); anything else exits with an error.
- The shell expands `purr *.mov` (or `purr *.qt`) before launch, so the command receives the `.mov`/`.qt` files the glob matches on disk.
- Output goes next to the input with the same base name: `video.mov` → `video.mp4`, `video.qt` → `video.mp4`.
- If the output file already exists, the job is skipped (shown as "output exists, skipped"); purr never overwrites. Note that `.mov` and `.qt` inputs share output names — `video.mov` and `video.qt` in the same directory both target `video.mp4` — so converting both, in one run or across runs, produces a single output file and the second job is skipped.
- Jobs run one at a time, in list order.
- purr probes each file with `ffprobe` before converting with `ffmpeg`.

### Default: compress

Without flags every file is re-encoded to **HEVC (H.265)** so the result is usually much smaller than the source; the savings depend on the footage and are not guaranteed.

- Video: `hevc_videotoolbox` (VideoToolbox, macOS) at quality 65, tagged as `hvc1` so QuickTime Player, Safari and iOS play it; software encoding is allowed when hardware is unavailable.
- Audio: copied as-is when AAC, re-encoded to AAC 192k otherwise.
- Because it relies on the VideoToolbox encoder, the default path works on macOS only; there is no cross-platform fallback encoder configured.

### `--fast`: remux

`--fast` (or `-f`) probes the file and, when its codecs are considered MP4-compatible, copies it into an MP4 container without re-encoding. This is usually faster than compression, preserves the encoded quality and produces roughly the original size.

- Compatible video: H.264, HEVC, VP9, AV1, MPEG-4.
- Compatible audio: AAC, MP3, Opus, AC3, E-AC3, FLAC.
- With incompatible codecs (e.g. ProRes), purr falls back to full compression, which requires macOS. Conversion may still fail if FFmpeg cannot process the input or mux its codecs into MP4.

## Keys

**File picker**

| Key                 | Action                                        |
|---------------------|-----------------------------------------------|
| `↑`/`↓`, `k`/`j`    | Move cursor                                   |
| `space`             | Toggle selection                              |
| `enter`             | Start conversion (requires at least one pick) |
| `q`, `esc`, `ctrl+c`| Quit                                          |

**Conversion**

| Key      | Action                                                                       |
|----------|------------------------------------------------------------------------------|
| `esc`    | Prompt "Abort remaining jobs? [y/N]": `y` cancels the running job and skips the rest; any other key resumes |
| `ctrl+c` | Kill the running conversion immediately and exit                              |
| `q`      | Quit — only available while no job is running                                 |

Each job shows ○ queued / ▸ running (with progress bar, speed and ETA) / ✓ done / – skipped / ✗ failed (ffmpeg's error output below the line) / ◌ cancelled. A cancelled or aborted job removes its partial output file.

## Known limits

- Direct arguments and the picker accept `.mov` and `.qt` only; there is no recursive directory scan. Accepting either extension does not mean every codec inside will convert — the codec notes above still apply.
- Only one conversion at a time; the queue is sequential.
- Errors surface ffmpeg's stderr under the failed entry.
- No configuration file, no custom encoder selection, no output directory override.
