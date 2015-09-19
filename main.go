// dx7 is a SuperCollider-based FM synthesizer.
package main

import (
	"log"

	"github.com/rakyll/portmidi"
	"github.com/scgolang/sc"
)

var (
	defName = "dx7voice"
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

	if err := client.Connect(scsynthAddr); err != nil {
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
	if err := client.SendDef(algorithm1); err != nil {
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
					// set gate to 0
					ctls := map[string]float32{"gate": float32(0)}
					if err = synths[note.Note].Set(ctls); err != nil {
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
				"gate":    float32(1),
				"op1freq": sc.Midicps(note.Note),
				"op1gain": float32(note.Velocity) / (127 * polyphony),
				"op1amt":  float32(100),
			}
			synth, err := dg.Synth(defName, sid, sc.AddToTail, controls)
			if err != nil {
				log.Fatal(err)
			}
			synths[note.Note] = synth
		}
	}
}
