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
	"log"

	"github.com/rakyll/portmidi"
)

func main() {
	// Initialize portmidi.
	portmidi.Initialize()
	defer portmidi.Terminate()

	// Initialize a new dx7.
	dx7, err := NewDX7(getConfig())
	if err != nil {
		log.Fatal(err)
	}

	if err := dx7.Run(); err != nil {
		log.Fatal(err)
	}
}
