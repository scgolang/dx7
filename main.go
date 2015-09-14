// dx7 is a SuperCollider-based FM synthesizer.
package main

import (
	"github.com/rakyll/portmidi"
	"github.com/scgolang/sc"
	. "github.com/scgolang/sc/types"
	. "github.com/scgolang/sc/ugens"
	"log"
)

var (
	defName = "dx7voice"

	// algorithm1
	//
	//        Op6
	//         |
	//        Op5
	//         |
	//  Op2   Op4
	//   |     |
	//  Op1   Op3
	//
	algorithm1 = sc.NewSynthdef(defName, func(p Params) Ugen {
		gate := p.Add("gate", 1)
		freq := p.Add("freq", 440)
		gain := p.Add("gain", 1)
		op1amt := p.Add("op1amt", 0)
		bus := C(0)

		// modulator
		op2 := Operator{
			Freq: freq,
			Gate: gate,
			Gain: gain,
			Done: FreeEnclosing,
		}.Rate(AR)

		// carrier
		op1 := Operator{
			Freq: freq,
			FM:   op2,
			Amt:  op1amt,
			Gate: gate,
			Gain: gain,
			Done: FreeEnclosing,
		}.Rate(AR)

		return Out{bus, op1}.Rate(AR)
	})
)

func main() {
	const (
		// polyphony is used to scale the gain of each synth voice
		polyphony      = 4
		midiDeviceId   = 3
		midiBufferSize = 1024
		localAddr      = "127.0.0.1:57110"
		scsynthAddr    = "127.0.0.1:57120"
	)
	client := sc.NewClient(localAddr)
	err := client.Connect(scsynthAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment this if you want scsynth to dump all the
	// midi messages it receives.
	// err = client.DumpOSC(sc.DumpAll)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dg, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}

	// send a synthdef
	err = client.SendDef(algorithm1)
	if err != nil {
		log.Fatal(err)
	}

	portmidi.Initialize()
	defer portmidi.Terminate()

	// Uncomment this section of code to figure out
	// which MIDI device to open.
	// midiDevices := portmidi.CountDevices()
	// fmt.Printf("%d midi devices:\n", midiDevices)
	// for i := 0; i < midiDevices; i++ {
	// 	did := portmidi.DeviceId(i)
	// 	fmt.Printf("%d\t%v\n", did, portmidi.GetDeviceInfo(did))
	// }

	midiDevice := portmidi.DeviceId(midiDeviceId)
	midiInput, err := portmidi.NewInputStream(midiDevice, midiBufferSize)
	if midiInput != nil {
		defer midiInput.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	// this slice keeps track of voice allocation
	synths := make([]*sc.Synth, 127)

	midiEvents := midiInput.Listen()
MidiLoop:
	for event := range midiEvents {
		switch event.Status {
		case NoteStatus:
			// note event
			note := NewNote(event)

			// note off
			if note.Velocity == 0 {
				if synths[note.Note] != nil {
					err = synths[note.Note].Set(map[string]float32{
						"gate": float32(0),
					})
					if err != nil {
						log.Fatal(err)
					}
				} else {
					// no synth node -- this should never happen
					log.Fatal("no synth node for note off event")
				}
				continue MidiLoop
			}

			// trigger the new note
			sid := client.NextSynthID()
			controls := map[string]float32{
				"gate":   float32(1),
				"freq":   sc.Midicps(note.Note),
				"gain":   float32(note.Velocity) / (127 * polyphony),
				"op1amt": float32(1000),
			}
			synth, err := dg.Synth(defName, sid, sc.AddToTail, controls)
			if err != nil {
				log.Fatal(err)
			}
			synths[note.Note] = synth
		}
	}
}
