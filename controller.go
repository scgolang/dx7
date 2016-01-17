package main

import (
	"math"

	"github.com/scgolang/poly"
	"github.com/scgolang/sc"
)

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
