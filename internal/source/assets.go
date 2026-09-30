package source

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

const nativeImageSize = 125440

// ExportAssets writes only presentation data, not the original executable.
func ExportAssets(prg []byte, directory string) error {
	if len(prg) < 28+nativeImageSize || binary.BigEndian.Uint16(prg) != 0x601a {
		return fmt.Errorf("source: incomplete native executable")
	}
	b := prg[28 : 28+nativeImageSize]
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	write := func(name string, data []byte) error {
		return os.WriteFile(filepath.Join(directory, name), data, 0644)
	}
	for _, part := range []struct {
		name       string
		start, end int
	}{
		{"motion.bin", 0x1058, 0x1088},
		{"panel-order.bin", 0x15bc, 0x16bc},
		{"palette.bin", 0x1502, 0x1522},
		{"rasters.bin", 0x1522, 0x15bc},
		{"message.bin", 0x1088, 0x1088 + 1144},
	} {
		if err := write(part.name, b[part.start:part.end]); err != nil {
			return err
		}
	}
	images := map[string]*image.NRGBA{"artwork.png": decodePlanar(b[0x16bc:], 320, 125)}
	for bank := 0; bank < 3; bank++ {
		images[fmt.Sprintf("panels-%d.png", bank)] = decodePlanar(b[0xbdfe+bank*25600:], 320, 160)
	}
	for letter, x := range map[string]int{"a": 0, "d": 48, "m": 96} {
		img := image.NewNRGBA(image.Rect(0, 0, 48, 35))
		for y := 0; y < 35; y++ {
			for dx := 0; dx < 48; dx++ {
				paint := images["artwork.png"].NRGBAAt(x+dx, y)
				if paint.R == 0 {
					paint.A = 0 // The original sprite mask is NOT(plane0|...|plane3).
				}
				img.SetNRGBA(dx, y, paint)
			}
		}
		images["letter-"+letter+".png"] = img
	}
	// The font uses one 32-bit row per glyph, with 20 rows. Control bytes
	// retain their native wrapped glyph indices rather than inventing artwork.
	for bank, base := range []int{0x93fa, 0xa8fc} {
		img := image.NewNRGBA(image.Rect(0, 0, 16*32, 16*20))
		for glyph := 0; glyph < 256; glyph++ {
			for y := 0; y < 20; y++ {
				word := binary.BigEndian.Uint32(b[base+glyph*80+y*4:])
				for x := 0; x < 32; x++ {
					v := byte(word >> uint(31-x) & 1)
					img.SetNRGBA(glyph%16*32+x, glyph/16*20+y, color.NRGBA{R: v * 255, A: 255})
				}
			}
		}
		images[fmt.Sprintf("font-%d.png", bank)] = img
	}
	for name, img := range images {
		file, err := os.Create(filepath.Join(directory, name))
		if err != nil {
			return err
		}
		err = png.Encode(file, img)
		closeErr := file.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	hash := sha256.Sum256(prg)
	manifest := map[string]any{"production": "Lockbuster DMA intro", "year": 1989, "native_unpacked_sha256": hex.EncodeToString(hash[:]), "logical_width": 320, "logical_height": 240, "ticks_per_second": 50, "indexed_encoding": "red = index * 17; alpha = native sprite mask", "scroll_font_width": 32, "scroll_font_height": 20, "panel_width": 64, "panel_height": 40, "panel_frames_per_bank": 20}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return write("manifest.json", append(data, '\n'))
}

func decodePlanar(data []byte, width, height int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	stride := width / 16 * 8
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var index byte
			for plane := 0; plane < 4; plane++ {
				word := binary.BigEndian.Uint16(data[y*stride+x/16*8+plane*2:])
				index |= byte(word>>uint(15-x%16)&1) << plane
			}
			img.SetNRGBA(x, y, color.NRGBA{R: index * 17, A: 255})
		}
	}
	return img
}
