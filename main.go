// dx7 is a SuperCollider-based FM synthesizer.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/rakyll/portmidi"
	"github.com/scgolang/sc"
)

func main() {
	const (
		defaultMidiDeviceId = 0
		midiBufferSize      = 1024
		localAddr           = "127.0.0.1:57110"
		defaultScsynthAddr  = "127.0.0.1:57120"
	)

	var (
		fs           = flag.NewFlagSet("", flag.ExitOnError)
		midiDeviceId = fs.Int("device", defaultMidiDeviceId, "MIDI Device ID")
		scsynthAddr  = fs.String("scsynth", defaultScsynthAddr, "scsynth address")
	)

	fs.Parse(os.Args[1:])

	client := sc.NewClient(localAddr)

	if err := client.Connect(*scsynthAddr); err != nil {
		log.Fatal(err)
	}

	// Uncomment this if you want scsynth to dump all the
	// midi messages it receives.
	// err = client.DumpOSC(sc.DumpAll)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}

	dx7 := NewDX7(client, defaultGroup)

	// TODO: have the dx7 send whatever synthdef it is using
	if err := client.SendDef(defaultAlgorithm); err != nil {
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

	midiDevice := portmidi.DeviceId(*midiDeviceId)
	midiInput, err := portmidi.NewInputStream(midiDevice, midiBufferSize)
	if midiInput != nil {
		defer midiInput.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	midiEvents := midiInput.Listen()

	for event := range midiEvents {
		switch event.Status {
		case NoteStatus:
			if err := dx7.Play(NewNote(event)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
