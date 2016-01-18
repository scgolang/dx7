package main

import (
	"fmt"
	"math"

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

var (
	ops = []int{1, 2, 3, 4, 5, 6}
)

func ctrlName(op int, name string) string {
	return fmt.Sprintf("op%d%s", op, name)
}

// FromNote implements poly.Controller.
func (dx7 *DX7) FromNote(note *poly.Note) map[string]float32 {
	var (
		ctrls = map[string]float32{"gate": float32(1)}
		freq  = sc.Midicps(note.Note)
		gain  = float32(note.Velocity) / (poly.MaxMIDI * polyphony)
	)

	for op := range ops {
		ctrls[ctrlName(op, "freq")] = freq
		ctrls[ctrlName(op, "gain")] = gain
		ctrls[ctrlName(op, "amt")] = dx7.ctrls["op1amt"]
		ctrls[ctrlName(op, "freqscale")] = dx7.ctrls["op2freqscale"]
		ctrls[ctrlName(op, "decay")] = dx7.ctrls["op2decay"]
		ctrls[ctrlName(op, "sustain")] = dx7.ctrls["op2sustain"]
	}

	return ctrls
}

// FromCtrl implements poly.Controller.
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
