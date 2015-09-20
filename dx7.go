package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

const (
	// polyphony is used to scale the gain of each synth voice
	polyphony = 4

	// numVoices is the number of possible voices
	numVoices = 127

	// defName is the name of the synthdef
	defName = "dx7voice"
)

var (
	// defaultAlgorithm
	//
	//    Op2
	//     |
	//    Op1
	//
	defaultAlgorithm = sc.NewSynthdef(defName, func(p sc.Params) sc.Ugen {
		var (
			gate = p.Add("gate", 1)

			op1freq = p.Add("op1freq", 440)
			op1gain = p.Add("op1gain", 1)
			op1amt  = p.Add("op1amt", 0)

			op2freq = p.Add("op2freq", 440)
			op2gain = p.Add("op2gain", 1)
			op2amt  = p.Add("op2amt", 0)

			bus = sc.C(0)
		)

		// modulator
		op2 := Operator{
			Freq: op2freq,
			Amt:  op2amt,
			Gate: gate,
			Gain: op2gain,
			Done: sc.FreeEnclosing,
		}.Rate(sc.AR)

		// carrier
		op1 := Operator{
			Freq: op1freq,
			FM:   op2,
			Amt:  op1amt,
			Gate: gate,
			Gain: op1gain,
			Done: sc.FreeEnclosing,
		}.Rate(sc.AR)

		return sc.Out{bus, op1}.Rate(sc.AR)
	})
)

// DX7 encapsulates the synth architecture of the legendary Yamaha DX7.
type DX7 struct {
	client *sc.Client
	group  *sc.Group
	voices [numVoices]*sc.Synth
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
	synth, err := dx7.group.Synth(defName, sid, sc.AddToTail, controls)
	if err != nil {
		return err
	}
	dx7.voices[note.Note] = synth
	return nil
}

// Control provides control over the DX7 using MIDI CC.
func (dx7 *DX7) Control(cc CC) error {
	// TODO: implement
	return nil
}

// NewDX7 returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func NewDX7(client *sc.Client, group *sc.Group) MIDIHandler {
	return &DX7{
		client: client,
		group:  group,
	}
}

// Wire represents a connection between two operators.
type Wire struct {
	// Input is the input side of the wire.
	Input int `xml:"input"`
	// Output is the output side of the wire.
	Output int `xml:"output"`
}
