// Command extract recovers an unrelocated executable for native asset analysis.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/olivierh59500/go-lockbuster3/internal/source"
)

func main() {
	input := flag.String("input", "", "packed TOS executable")
	output := flag.String("output", "", "unpacked executable destination")
	assets := flag.String("assets", "", "export native presentation assets to this directory")
	flag.Parse()
	if *input == "" || *output == "" && *assets == "" {
		log.Fatal("-input and either -output or -assets are required")
	}
	data, err := os.ReadFile(*input)
	if err != nil {
		log.Fatal(err)
	}
	unpacked, err := source.UnpackPRG(data)
	if err != nil {
		log.Fatal(err)
	}
	if *output != "" {
		if err := os.WriteFile(*output, unpacked, 0644); err != nil {
			log.Fatal(err)
		}
	}
	if *assets != "" {
		if err := source.ExportAssets(unpacked, *assets); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Recovered %d bytes", len(unpacked))
}
