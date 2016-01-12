package main

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"sync"

	"github.com/scgolang/dx7/sysex"
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

// Common errors.
var (
	ErrNilNote       = errors.New("nil note")
	ErrNilCtrl       = errors.New("nil ctrl")
	ErrMissingNoteOn = errors.New("received a note off for a voice that is not on")
)

// DX7 encapsulates the synth architecture of the legendary Yamaha DX7.
type DX7 struct {
	cfg           *config
	client        *sc.Client
	group         *sc.Group
	voices        [maxMIDI]*sc.Synth
	voicesMutex   *sync.Mutex
	currentDef    string
	currentPreset *sysex.Sysex       // remove this
	presets       map[string]string  // maps preset names to their synthdef names
	ctrls         map[string]float32 // synth param values
}

// Connect to scsynth and load synthdefs.
func (dx7 *DX7) Connect() error {
	// Initialize a new supercollider client.
	client, err := sc.NewClient("udp", dx7.cfg.localAddr, dx7.cfg.scsynthAddr)
	if err != nil {
		return err
	}
	dx7.client = client

	// Tell scsynth to dump all the midi messages it receives.
	if dx7.cfg.dumpOSC {
		if err := client.DumpOSC(sc.DumpAll); err != nil {
			return err
		}
	}

	// Add the default group.
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		return err
	}
	dx7.group = defaultGroup
	return nil
}

// Play plays a note. This can either turn a voice on or
// off depending on if velocity is > 0.
func (dx7 *DX7) Play(note *Note) error {
	if note == nil {
		return ErrNilNote
	}

	dx7.voicesMutex.Lock()
	defer dx7.voicesMutex.Unlock()
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
		"op2decay":     dx7.ctrls["op2decay"],
		"op2sustain":   dx7.ctrls["op2sustain"],
	}
	synth, err := dx7.group.Synth(dx7.currentDef, sid, sc.AddToTail, controls)
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

	dx7.voicesMutex.Lock()
	defer dx7.voicesMutex.Unlock()

	for _, voice := range dx7.voices {
		if voice != nil {
			if err := voice.Set(dx7.ctrls); err != nil {
				return err
			}
		}
	}

	return nil
}

// setCtrls sets controller values and returns a bool
// indicating whether any were changed.
func (dx7 *DX7) setCtrls(ctrl *Ctrl) bool {
	// TODO: allow configurable controller mappings
	switch ctrl.Num {
	default:
		return false
	case 106: // op1 FM Amt
		dx7.ctrls["op1amt"] = float32(ctrl.Value) * (fmtAmtHi / maxMIDI)
	case 107: // op2 Freq Scale
		dx7.ctrls["op2freqscale"] = getOp2FreqScale(ctrl.Value)
	case 108:
		dx7.ctrls["op2decay"] = linear(ctrl.Value, decayLo, decayHi)
	case 109:
		dx7.ctrls["op2sustain"] = float32(ctrl.Value) / maxMIDI
	}
	return true
}

// getOp2FreqScale returns a frequency scaling value for op2.
func getOp2FreqScale(value int) float32 {
	exp := float64(linear(value, freqScaleLo, freqScaleHi))
	// return float32(math.Pow(2, float64((norm*(freqScaleHi-freqScaleLo))+freqScaleLo)))
	return float32(math.Pow(2, exp))
}

func linear(val int, min, max float32) float32 {
	norm := float32(val) / maxMIDI
	return (norm * (max - min)) + min
}

// Run listens for events (either MIDI or OSC).
func (dx7 *DX7) Run() error {
	// Print a list of midi devices and exit.
	if dx7.cfg.listMidiDevices {
		PrintMidiDevices(os.Stdout)
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

// New returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func New(cfg *config) (*DX7, error) {
	return &DX7{
		cfg: cfg,
		ctrls: map[string]float32{
			"op1amt":       float32(defaultAmt),
			"op2freqscale": float32(1),
			"op2decay":     float32(defaultDecay),
			"op2sustain":   float32(defaultSustain),
		},
	}, nil
}
