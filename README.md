# Lockbuster Intro Go

A Go/Ebitengine port of Lockbuster's 1989 Atari ST intro for DMA, using Demo
Construction Kit. The screen combines turning DMA panels, three animated
letters, scrolling text over raster colors and three lower blocks.

The native executable supplies the artwork, text and animation data. DCK
provides reusable graphics and the YM playback facade. The presentation follows
the Atari ST's PAL clock.

## Native asset extraction

The extraction command recovers an unrelocated TOS executable from its packed
form:

```sh
go run ./cmd/extract -input path/to/freddy.prg -output /tmp/freddy-unpacked.prg
```

The decoder follows the original backward bitstream, validates bounds and
checks its word checksum before returning the recovered executable.
