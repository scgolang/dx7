package main

import (
	"encoding/json"
	"os"

	"github.com/scgolang/dx7/sysex"
	"github.com/scgolang/poly"
)

const (
	// polyphony is used to scale the gain of each synth voice.
	polyphony = 4

	// fmtAmtHi is the max value for op1amt.
	fmtAmtHi = float32(2000)

	// freqScaleLo is the min value for op2freqscale (as a power of 2).
	freqScaleLo = float32(-8)

	// freqScaleHi is the max value for op2freqscale (as a power of 2).
	freqScaleHi = float32(2)

	// decayLo is the min value for op2decay (in secs).
	decayLo = float32(0.0001)

	// decayHi is the max value for op2decay (in secs).
	decayHi = float32(10)
)

// DX7 encapsulates the synth architecture of the legendary Yamaha DX7.
type DX7 struct {
	*poly.Poly

	cfg           *config
	currentPreset *sysex.Sysex      // remove this
	presets       map[string]string // maps preset names to their synthdef names

	ctrls map[string]float32
}

// Run listens for events (either MIDI or OSC).
func (dx7 *DX7) Run() error {
	// Print a list of midi devices and exit.
	if dx7.cfg.listMidiDevices {
		poly.PrintMidiDevices(os.Stdout)
		return nil
	}

	// Load the current preset.
	if err := dx7.LoadPreset(dx7.cfg.preset); err != nil {
		return err
	}

	// Dump sysex data to stdout for the current preset and return.
	if dx7.cfg.dumpSysex {
		return json.NewEncoder(os.Stdout).Encode(dx7.currentPreset)
	}

	// Connect to scsynth.
	if err := dx7.Connect(dx7.cfg.localAddr, dx7.cfg.scsynthAddr); err != nil {
		return err
	}

	// Send all the synthdefs we need.
	if err := dx7.SendSynthdefs(); err != nil {
		return err
	}

	// Listen for MIDI events.
	return dx7.MidiListen(dx7.cfg.midiDeviceID)
}

// New returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func New(cfg *config) (*DX7, error) {
	dx7 := &DX7{
		cfg: cfg,
		ctrls: map[string]float32{
			"op1amt":       float32(defaultAmt),
			"op2freqscale": float32(1),
			"op2decay":     float32(defaultDecay),
			"op2sustain":   float32(defaultSustain),
		},
	}
	if p, err := poly.New(dx7); err != nil {
		return nil, err
	} else {
		dx7.Poly = p
	}
	return dx7, nil
}
