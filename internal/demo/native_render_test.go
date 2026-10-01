//go:build lockbuster_rendercheck

package demo

import (
	"encoding/binary"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/fidelity/ebiten/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) { os.Exit(testutil.RunGPU(m)) }

func TestNativeScreenIndices(t *testing.T) {
	g, err := NewGame(true)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	frames := []int{1, 19, 41, 309, 501, 801, 1144, 1600, 2400, 4000, 4600, 4800}
	root := os.Getenv("LOCKBUSTER_NATIVE_BUFFERS")
	if root == "" {
		t.Skip("native RAM capture directory not supplied")
	}
	indices, mono := make([]byte, Width*Height*4), make([]byte, Width*Height*4)
	for _, tick := range frames {
		for g.Tick() < tick {
			if err := g.Update(); err != nil {
				t.Fatal(err)
			}
		}
		raw, err := os.ReadFile(fmt.Sprintf("%s/frame-%04d.bin", root, tick))
		if err != nil {
			t.Fatal(err)
		}
		g.buffers[g.front].ReadPixels(indices)
		g.text.ReadPixels(mono)
		wrong := 0
		var first string
		for y := 0; y < 200; y++ {
			for x := 0; x < Width; x++ {
				want := byte(0)
				for p := 0; p < 4; p++ {
					word := binary.BigEndian.Uint16(raw[y*160+x/16*8+p*2:])
					want |= byte(word>>uint(15-x%16)&1) << p
				}
				at := (y*Width + x) * 4
				got := byte((int(indices[at]) + 8) / 17)
				if y >= 160 {
					got = got&^1 | mono[at]/255
				}
				if want != got {
					wrong++
					if first == "" {
						first = fmt.Sprintf("x%d y%d native%d Go%d", x, y, want, got)
					}
				}
			}
		}
		if wrong > 0 {
			t.Errorf("frame %d differs at %d native index pixels; first %s", tick, wrong, first)
		}
	}
	// The native VU artwork occupies three fixed 32x25 rectangles below the screen.
	dst := ebiten.NewImage(Width, Height)
	defer dst.Deallocate()
	g.Draw(dst)
}
