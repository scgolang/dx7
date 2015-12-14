// `dx7` is an FM synthesizer built with Go and SuperCollider
// that is inspired by the legendary Yamaha DX7.
// The synth architecture is built on the same principles as the DX7,
// but is much more flexible and extensible.
// Operators consist of a sine wave and envelope generator, and
// are combined into different algorithms to modulate one another.
// Instead of being limited to 6 operators and 32 algorithms though,
// `dx7` allows you to use any number of operators and combine them
// into algorithms of your own design.
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/rakyll/portmidi"
	"github.com/scgolang/dx7/sysex"
)

const (
	ExitSuccess = 0
)

func main() {
	// Initialize portmidi.
	portmidi.Initialize()
	defer portmidi.Terminate()

	cfg := parseConfig()

	// Print a list of midi devices and exit.
	if cfg.listMidiDevices {
		PrintMidiDevices(os.Stdout)
		os.Exit(ExitSuccess)
	}

	// Dump sysex data to stdout.
	if cfg.dumpSysex != "" {
		if err := dumpSysex(cfg.dumpSysex); err != nil {
			log.Fatal(err)
		}
		os.Exit(ExitSuccess)
	}

	// Initialize a new dx7.
	dx7, err := NewDX7(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := dx7.Listen(); err != nil {
		log.Fatal(err)
	}
}

// dumpSysex prints a JSON-encoded sysex structure to stdout.
func dumpSysex(sysexPath string) error {
	sysexFile, err := os.Open(sysexPath)
	if err != nil {
		return err
	}

	syx, err := sysex.New(sysexFile)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(os.Stdout).Encode(syx); err != nil {
		return err
	}

	return nil
}
