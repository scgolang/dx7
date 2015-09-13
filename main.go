// dx7 is a SuperCollider-based FM synthesizer.
package main

import (
	"github.com/scgolang/sc"
	. "github.com/scgolang/sc/types"
	. "github.com/scgolang/sc/ugens"
	"log"
)

func main() {
	const (
		localAddr   = "127.0.0.1:57110"
		scsynthAddr = "127.0.0.1:57120"
		defName     = "dx7voice"
	)
	client := sc.NewClient(localAddr)
	err := client.Connect(scsynthAddr)
	if err != nil {
		log.Fatal(err)
	}
	g, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}

	// send a synthdef
	def := sc.NewSynthdef(defName, func(p Params) Ugen {
		bus, freq := C(0), C(440)
		sig := SinOsc{Freq: freq}.Rate(AR)
		return Out{bus, sig}.Rate(AR)
	})
	err = client.SendDef(def)
	if err != nil {
		log.Fatal(err)
	}

	// create a new synth node
	id := client.NextSynthID()
	_, err = g.Synth(defName, id, sc.AddToTail, nil)
	if err != nil {
		log.Fatal(err)
	}
}
