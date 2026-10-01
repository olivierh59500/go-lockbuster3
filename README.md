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

## Android

```sh
./scripts/run-android.sh
./scripts/run-android.sh --build-only
```

The script builds the ARM64 Ebitengine library and debug APK, then installs and
launches **Lockbuster 3** on an authorized USB device. It uses Ebitengine 2.9.11,
Java 17, SDK 36, Android Gradle Plugin 8.10.1 and the included Gradle 8.11.1
wrapper. Set `ANDROID_HOME` and `JAVA_HOME` when tools are outside their usual
locations; `ANDROID_SERIAL` selects a device when more than one is connected.

The APK is `android/app/build/outputs/apk/debug/app-debug.apk`; its application
ID is `com.olivierh.lockbuster3`. The scene runs in landscape, centered with
its 4:3 aspect ratio preserved. Android Back closes the activity; swiping from
the screen edge temporarily reveals the system controls. The window stays awake
during playback, and Ebitengine suspends/resumes with the activity lifecycle.
Audio and graphics initialize in the first mobile Update after the Android
application context is available. All artwork and music are embedded.

An optional muted checkpoint records runtime counters without changing the
normal launch:

```sh
adb shell am start -S -W -n com.olivierh.lockbuster3/.MainActivity \
  --ei lockbuster3_verify_tick 450
```

A positive checkpoint can be displayed over the lock screen for verification;
it does not dismiss the keyguard. Normal launches start at tick zero with music.
Generated AARs, APKs, build caches and machine settings remain local.

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
