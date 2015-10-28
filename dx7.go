package main

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/scgolang/sc"
)

const (
	// polyphony is used to scale the gain of each synth voice.
	polyphony = 4

	// numVoices is the number of possible voices.
	numVoices = 127

	// defName is the name of the synthdef.
	defaultDefName = "defaultDX7voice"

	// fmtAmtHi is the max value for op1amt.
	fmtAmtHi = float32(2000)

	// freqScaleLo is the min value for op2freqscale (as a power of 2).
	freqScaleLo = float32(-8)

	// freqScaleHi is the max value for op2freqscale (as a power of 2).
	freqScaleHi = float32(4)
)

// Common errors.
var (
	ErrNilNote       = errors.New("nil note")
	ErrNilCtrl       = errors.New("nil ctrl")
	ErrMissingNoteOn = errors.New("received a note off for a voice that is not on")
)

// DX7 encapsulates the synth architecture of the legendary Yamaha DX7.
type DX7 struct {
	cfg      *config
	client   *sc.Client
	group    *sc.Group
	voices   [numVoices]*sc.Synth
	curVoice string
	ctrls    map[string]float32 // synth param values
}

// Play plays a note. This can either turn a voice on or
// off depending on if velocity is > 0.
func (dx7 *DX7) Play(note *Note) error {
	if note == nil {
		return ErrNilNote
	}

	if note.Velocity == 0 {
		if dx7.voices[note.Note] != nil {
			// set gate to 0
			ctls := map[string]float32{"gate": float32(0)}
			if err := dx7.voices[note.Note].Set(ctls); err != nil {
				return err
			}
		} else {
			// received note off for a voice that is not on
			return ErrMissingNoteOn
		}
		// set voice to nil so we don't send any more messages to it
		dx7.voices[note.Note] = nil
		return nil
	}
	sid := dx7.client.NextSynthID()
	controls := map[string]float32{
		"gate":         float32(1),
		"op1freq":      sc.Midicps(note.Note),
		"op1gain":      float32(note.Velocity) / (maxMIDI * polyphony),
		"op1amt":       dx7.ctrls["op1amt"],
		"op2freqscale": dx7.ctrls["op2freqscale"],
	}
	synth, err := dx7.group.Synth(dx7.curVoice, sid, sc.AddToTail, controls)
	if err != nil {
		return err
	}
	dx7.voices[note.Note] = synth
	return nil
}

// Control provides control over the DX7 using control messages.
func (dx7 *DX7) Control(ctrl *Ctrl) error {
	if ctrl == nil {
		return ErrNilCtrl
	}

	if changed := dx7.setCtrls(ctrl); !changed {
		return nil
	}

	// TODO: get rid of data race
	for _, voice := range dx7.voices {
		if voice != nil {
			if err := voice.Set(dx7.ctrls); err != nil {
				return err
			}
		}
	}

	return nil
}

// setCtrls sets controller values.
func (dx7 *DX7) setCtrls(ctrl *Ctrl) bool {
	// TODO: allow configurable controller mappings
	switch ctrl.Num {
	default:
		return false
	case 106: // op1 FM Amt
		dx7.ctrls["op1amt"] = float32(ctrl.Value) * (fmtAmtHi / maxMIDI)
		return true
	case 107: // op2 Freq Scale
		freqscale := getOp2FreqScale(ctrl.Value)
		fmt.Printf("freqscale %f\n", freqscale)
		dx7.ctrls["op2freqscale"] = freqscale
		return true
	}
}

func getOp2FreqScale(value int) float32 {
	norm := float32(value) / maxMIDI
	return float32(math.Pow(2, float64((norm*(freqScaleHi-freqScaleLo))+freqScaleLo)))
}

// Listen listens for events (either MIDI or OSC).
func (dx7 *DX7) Listen() error {
	// Listen for MIDI or OSC events, depending on
	// whether an events address was specified.
	if dx7.cfg.eventsAddr == "" {
		if err := MidiListen(dx7.cfg.midiDeviceID, dx7); err != nil {
			return err
		}
	} else {
		// Listen for OSC events.
	}
	return nil
}

// NewDX7 returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func NewDX7(cfg *config) (*DX7, error) {
	// Initialize a new supercollider client.
	client := sc.NewClient(cfg.localAddr)
	if err := client.Connect(cfg.scsynthAddr); err != nil {
		log.Fatal(err)
	}

	// Tell scsynth to dump all the midi messages it receives.
	if cfg.dumpOSC {
		if err := client.DumpOSC(sc.DumpAll); err != nil {
			log.Fatal(err)
		}
	}

	// Add the default group.
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		log.Fatal(err)
	}

	// Send a synthdef.
	if err := client.SendDef(defaultAlgorithm); err != nil {
		return nil, err
	}

	return &DX7{
		cfg:      cfg,
		client:   client,
		group:    defaultGroup,
		curVoice: defaultDefName,
		ctrls: map[string]float32{
			"op1amt":       float32(defaultAmt),
			"op2freqscale": float32(1),
		},
	}, nil
}
