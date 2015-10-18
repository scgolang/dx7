package main

import (
	"fmt"
	"log"

	"github.com/scgolang/sc"
)

const (
	// polyphony is used to scale the gain of each synth voice
	polyphony = 4

	// numVoices is the number of possible voices
	numVoices = 127

	// defName is the name of the synthdef
	defaultDefName = "defaultDX7voice"
)

// DX7 encapsulates the synth architecture of the legendary Yamaha DX7.
type DX7 struct {
	cfg      *config
	client   *sc.Client
	group    *sc.Group
	voices   [numVoices]*sc.Synth
	curVoice string
}

// Play plays a note. This can either turn a voice on or
// off depending on if velocity is > 0.
func (dx7 *DX7) Play(note Note) error {
	if note.Velocity == 0 {
		if dx7.voices[note.Note] != nil {
			// set gate to 0
			ctls := map[string]float32{"gate": float32(0)}
			if err := dx7.voices[note.Note].Set(ctls); err != nil {
				return err
			}
		} else {
			// no synth node -- this should never happen
			return fmt.Errorf("no synth node for note off event")
		}
		return nil
	}
	sid := dx7.client.NextSynthID()
	controls := map[string]float32{
		"gate":    float32(1),
		"op1freq": sc.Midicps(note.Note),
		"op1gain": float32(note.Velocity) / (127 * polyphony),
		"op1amt":  float32(100),
	}
	synth, err := dx7.group.Synth(dx7.curVoice, sid, sc.AddToTail, controls)
	if err != nil {
		return err
	}
	dx7.voices[note.Note] = synth
	return nil
}

// Control provides control over the DX7 using control messages.
func (dx7 *DX7) Control(ctrl Ctrl) error {
	// TODO: implement
	return nil
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
	}, nil
}
