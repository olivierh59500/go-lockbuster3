package source

import (
	"encoding/csv"
	"github.com/olivierh59500/go-lockbuster3/assets"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestClocksMatchOriginal68000Checkpoints(t *testing.T) {
	read := func(name string) []byte {
		b, e := assets.Files.ReadFile("original/" + name)
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	p := NewPanelClock(read("panel-order.bin"))
	s := NewScrollClock(read("message.bin"))
	letters := ReadLetters(read("motion.bin"))
	for range 2 {
		for i := range letters {
			letters[i].Step()
		}
	}
	file, err := os.Open("testdata/native-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	tick := 0
	for _, row := range rows[1:] {
		want := make([]int, len(row))
		for i, x := range row {
			v, e := strconv.Atoi(x)
			if e != nil {
				t.Fatal(e)
			}
			want[i] = v
		}
		for tick < want[0] {
			p.Step()
			s.Step()
			for i := range letters {
				letters[i].Step()
			}
			tick++
		}
		got := []int{tick, p.StepN, p.Offset, p.Cursor, s.Remaining, s.Glyph, s.Column}
		for _, letter := range letters {
			for _, word := range letter.Words {
				got = append(got, int(uint16(word)))
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("native checkpoint %d: got %v, want %v", tick, got, want)
		}
	}
}
