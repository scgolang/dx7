// dx7 is a SuperCollider-based FM synthesizer.
package main

import (
	"flag"
	"fmt"
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
		listDevices  = fs.Bool("list", false, "list MIDI devices")
		dumposc      = fs.Bool("dumposc", false, "have scsynth dump OSC messages on stdout")
	)

	// parse cli
	fs.Parse(os.Args[1:])

	// initialize MIDI
	portmidi.Initialize()
	defer portmidi.Terminate()

	if *listDevices {
		midiDevices := portmidi.CountDevices()
		fmt.Printf("%d midi devices:\n", midiDevices)
		for i := 0; i < midiDevices; i++ {
			did := portmidi.DeviceId(i)
			fmt.Printf("%d\t%v\n", did, portmidi.GetDeviceInfo(did))
		}
		return
	}

	client := sc.NewClient(localAddr)

	if err := client.Connect(*scsynthAddr); err != nil {
		log.Fatal(err)
	}

	// Tell scsynth to dump all the midi messages it receives.
	if *dumposc {
		if err := client.DumpOSC(sc.DumpAll); err != nil {
			log.Fatal(err)
		}
	}

	// add the default group
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}

	// initialize a new dx7
	dx7, err := NewDX7(client, defaultGroup)
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment this section of code to figure out
	// which MIDI device to open.

	midiDevice := portmidi.DeviceId(*midiDeviceId)
	midiInput, err := portmidi.NewInputStream(midiDevice, midiBufferSize)
	if midiInput != nil {
		defer midiInput.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	for event := range midiInput.Listen() {
		switch event.Status {
		case NoteStatus:
			if err := dx7.Play(NewNote(event)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
