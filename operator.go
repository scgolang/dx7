package main

import (
	"github.com/scgolang/sc"
)

// Operator is a sine wave signal combined with an envelope generator.
type Operator struct {
	// Freq is the oscillator frequency.
	Freq sc.Input

	// FM is the frequency modulation input.
	FM sc.Input

	// Amt controls the frequency modulation amount.
	Amt sc.Input

	// Gain is the output gain.
	Gain sc.Input

	// A is amp envelope attack (in seconds)
	A sc.Input

	// D is amp envelope decay (in seconds)
	D sc.Input

	// S is amp envelope sustain [0, 1]
	S sc.Input

	// R is amp envelope release (in seconds)
	R sc.Input

	// Gate trigger the envelope and holds it open while > 0
	Gate sc.Input

	// Done is the ugen done action
	Done int
}

// defaults
func (op *Operator) defaults() {
	if op.Freq == nil {
		op.Freq = sc.C(440)
	}
	if op.FM == nil {
		op.FM = sc.C(1)
	}
	if op.Amt == nil {
		op.Amt = sc.C(0)
	}
	if op.Gain == nil {
		op.Gain = sc.C(1)
	}
	if op.A == nil {
		op.A = sc.C(0.01)
	}
	if op.D == nil {
		op.D = sc.C(0.3)
	}
	if op.S == nil {
		op.S = sc.C(0.5)
	}
	if op.R == nil {
		op.R = sc.C(1)
	}
	if op.Gate == nil {
		op.Gate = sc.C(1)
	}
}

// Rate creates a new ugen at a specific rate.
// If rate is an unsupported value this method will cause a runtime panic.
func (op Operator) Rate(rate int8) sc.Input {
	sc.CheckRate(rate)
	(&op).defaults()

	adsr := sc.EnvADSR{A: op.A, D: op.D, S: op.S, R: op.R}
	env := sc.EnvGen{
		Env:        adsr,
		Gate:       op.Gate,
		LevelScale: op.Gain,
		Done:       op.Done,
	}.Rate(sc.AR)

	freq := op.Freq.Add(op.FM.Mul(op.Amt))

	return sc.SinOsc{Freq: freq}.Rate(sc.AR).Mul(env)
}
