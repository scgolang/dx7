package main

import (
	"flag"
	"os"

	"github.com/scgolang/sc"
)

// DX7 is a recreation of the legendary Yamaha DX7.
type DX7 struct {
	algorithm      int
	client         *sc.Client
	ctrls          map[string]float32
	flags          *flag.FlagSet
	midiDeviceName string
	scsynthAddr    string
}

// Connect connects to scsynth.
func (dx7 *DX7) Connect() error {
	return nil
}

// Listen listens for MIDI events.
func (dx7 *DX7) Listen() error {
	return nil
}

// run the dx7.
func (dx7 *DX7) run() error {
	// Connect to scsynth.
	if err := dx7.Connect(); err != nil {
		return err
	}
	// Send all the synthdefs we need.
	if err := dx7.SendSynthdefs(); err != nil {
		return err
	}
	// Listen for events.
	return dx7.Listen()
}

// New returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func New() (*DX7, error) {
	dx7 := &DX7{
		ctrls: map[string]float32{
			"op1amt":       float32(defaultAmt),
			"op2freqscale": float32(1),
			"op2decay":     float32(defaultDecay),
			"op2sustain":   float32(defaultSustain),
		},
		flags: flag.NewFlagSet("dx7", flag.ExitOnError),
	}
	dx7.flags.StringVar(&dx7.midiDeviceName, "d", "", "MIDI device name")
	dx7.flags.StringVar(&dx7.scsynthAddr, "scsynth", "127.0.0.1:57120", "scsynth UDP listening address")

	dx7.flags.Parse(os.Args[1:])

	return dx7, nil
}
