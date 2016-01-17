package main

import (
	"encoding/json"
	"math"
	"os"

	"github.com/scgolang/dx7/sysex"
	"github.com/scgolang/poly"
	"github.com/scgolang/sc"
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

// SendSynthdefs sends all the synthdefs needed for the DX7.
func (dx7 *DX7) SendSynthdefs() error {
	for algo, f := range algorithms {
		if err := dx7.Client.SendDef(sc.NewSynthdef(algo, f)); err != nil {
			return err
		}
	}
	return nil
}

// FromNote provides params for a new synth voice.
func (dx7 *DX7) FromNote(note *poly.Note) map[string]float32 {
	return map[string]float32{
		"gate":         float32(1),
		"op1freq":      sc.Midicps(note.Note),
		"op1gain":      float32(note.Velocity) / (poly.MaxMIDI * polyphony),
		"op1amt":       dx7.ctrls["op1amt"],
		"op2freqscale": dx7.ctrls["op2freqscale"],
		"op2decay":     dx7.ctrls["op2decay"],
		"op2sustain":   dx7.ctrls["op2sustain"],
	}
}

// FromCtrl returns params for a MIDI CC event.
func (dx7 *DX7) FromCtrl(ctrl *poly.Ctrl) map[string]float32 {
	switch ctrl.Num {
	default:
		return nil
	case 106: // op1 FM Amt
		dx7.ctrls["op1amt"] = float32(ctrl.Value) * (fmtAmtHi / poly.MaxMIDI)
	case 107: // op2 Freq Scale
		dx7.ctrls["op2freqscale"] = getOp2FreqScale(ctrl.Value)
	case 108:
		dx7.ctrls["op2decay"] = linear(ctrl.Value, decayLo, decayHi)
	case 109:
		dx7.ctrls["op2sustain"] = float32(ctrl.Value) / poly.MaxMIDI
	}
	return dx7.ctrls
}

// getOp2FreqScale returns a frequency scaling value for op2.
func getOp2FreqScale(value int) float32 {
	exp := float64(linear(value, freqScaleLo, freqScaleHi))
	// return float32(math.Pow(2, float64((norm*(freqScaleHi-freqScaleLo))+freqScaleLo)))
	return float32(math.Pow(2, exp))
}

func linear(val int, min, max float32) float32 {
	norm := float32(val) / poly.MaxMIDI
	return (norm * (max - min)) + min
}

// Run listens for events (either MIDI or OSC).
func (dx7 *DX7) Run() error {
	// Print a list of midi devices and exit.
	if dx7.cfg.listMidiDevices {
		poly.PrintMidiDevices(os.Stdout)
		return nil
	}

	// Read all the sysex files.
	if dx7.cfg.preset != "" {
		if err := dx7.LoadPreset(dx7.cfg.preset); err != nil {
			return err
		}
	}

	// Dump sysex data to stdout.
	if dx7.cfg.dumpSysex {
		return json.NewEncoder(os.Stdout).Encode(dx7.currentPreset)
	}

	// Connect to scsynth.
	if err := dx7.Connect(dx7.cfg.localAddr, dx7.cfg.scsynthAddr); err != nil {
		return err
	}

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
