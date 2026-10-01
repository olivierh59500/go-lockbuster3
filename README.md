# Lockbuster 3 Go

A Go/Ebitengine port of Lockbuster's third Atari ST intro for DMA (1989), using Demo
Construction Kit. The screen combines turning DMA panels, three animated
letters, scrolling text over raster colors and three lower blocks.

The native executable supplies the artwork, text and animation data. DCK
provides reusable graphics and the YM playback facade. The presentation follows
the Atari ST's PAL clock.

## Run

```sh
go run ./cmd/lockbuster3
go run ./cmd/lockbuster3 -mute
```

Space or Escape closes the intro. The logical canvas is 320 × 240 pixels,
including the three volume blocks below the normal 200-line screen. Simulation
advances at 50 Hz independently of the monitor's refresh rate. The soundtrack is
Jochen Hippel's **Rollout tune 2**, opened through DCK's audio facade.

```sh
go run ./cmd/lockbuster3 -capture captures/preview -frame 450 -frames 300
```

This writes six seconds of native PNG frames without device audio. Runtime
rendering does not read pixels from the GPU.

## Rendering and native behavior

The three panel sheets preserve the original twenty-frame turns and their
forward/reverse selection program. Letter positions retain the native signed
word arithmetic and accelerating trajectories. Two retained screen buffers
restore the previous letter masks before drawing the next frame, preserving
overlap and layer order. DCK owns the retained sprite slots and the two
scrolling windows; glyph strips are sliced once, not recreated every frame.

The message switches between two original fonts and two 20-line scrolling
lanes, including the source's held-row handoffs and continuous restart. The
native raster palette and bitplane-zero text overlay use a compact Kage pass.
The three lower blocks read YM channel volumes on the fixed simulation clock.

Tests compare 43 states recorded from the original 68000 executable across
4,800 updates, including the complete first text loop. An optional GPU check
compares twelve complete 320 × 200 indexed frame buffers with the original:

```sh
LOCKBUSTER_NATIVE_BUFFERS=/path/to/native/buffers \
  go test -tags lockbuster_rendercheck ./internal/demo
```

The native buffer comparison is exact at those checkpoints. It excludes the
three YM-driven blocks below the main screen; no continuous hardware timing
equivalence is claimed. A separate lookup retains the original ST DAC levels.

## Native asset extraction

The extraction command recovers an unrelocated TOS executable from its packed
form:

```sh
go run ./cmd/extract -input path/to/freddy.prg -output /tmp/freddy-unpacked.prg
```

The decoder follows the original backward bitstream, validates bounds and
checks its word checksum before returning the recovered executable.
