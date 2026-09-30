// Command extract recovers an unrelocated executable for native asset analysis.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/olivierh59500/go-lockbusterintro/internal/source"
)

func main() {
	input := flag.String("input", "", "packed TOS executable")
	output := flag.String("output", "", "unpacked executable destination")
	flag.Parse()
	if *input == "" || *output == "" {
		log.Fatal("-input and -output are required")
	}
	data, err := os.ReadFile(*input)
	if err != nil {
		log.Fatal(err)
	}
	unpacked, err := source.UnpackPRG(data)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*output, unpacked, 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("Recovered %d bytes", len(unpacked))
}
