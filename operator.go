package main

import (
	"strconv"

	"github.com/scgolang/sc"
)

const (
	defaultFreq    = 440
	defaultGain    = 1
	defaultAmt     = 0
	defaultAttack  = 0.01
	defaultDecay   = 0.3
	defaultSustain = 0.5
	defaultRelease = 1
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

// defaults set default values for an operator.
func (op *Operator) defaults() {
	if op.Freq == nil {
		op.Freq = sc.C(defaultFreq)
	}
	if op.Gain == nil {
		op.Gain = sc.C(defaultGain)
	}
	if op.Amt == nil {
		op.Amt = sc.C(defaultAmt)
	}
	if op.FM == nil {
		op.FM = sc.C(0)
	}
	if op.A == nil {
		op.A = sc.C(defaultAttack)
	}
	if op.D == nil {
		op.D = sc.C(defaultDecay)
	}
	if op.S == nil {
		op.S = sc.C(defaultSustain)
	}
	if op.R == nil {
		op.R = sc.C(defaultRelease)
	}
	if op.Gate == nil {
		op.Gate = sc.C(1)
	}
}

// Rate creates a new ugen at a specific rate.
// If rate is an unsupported value this method will cause a runtime panic.
func (op Operator) Rate(rate int8) sc.Input {
	// Check the rate and set defaults.
	sc.CheckRate(rate)
	(&op).defaults()

	// Amp Envelope
	adsr := sc.EnvADSR{A: op.A, D: op.D, S: op.S, R: op.R}
	env := sc.EnvGen{
		Env:        adsr,
		Gate:       op.Gate,
		LevelScale: op.Gain,
		Done:       op.Done,
	}.Rate(sc.AR)

	// Modulate frequency with FM.
	freq := op.Freq.Add(op.FM.Mul(op.Amt))

	// Return the carrier.
	return sc.SinOsc{Freq: freq}.Rate(sc.AR).Mul(env)
}

// NewOperator creates an operator with a specific index
// and adds synth params to a synthdef.
func NewOperator(i int, p sc.Params, gate, fm sc.Input) sc.Input {
	var (
		name    = "op" + strconv.Itoa(i)
		freq    = p.Add(name+"freq", defaultFreq)
		gain    = p.Add(name+"gain", defaultGain)
		amt     = p.Add(name+"amt", defaultAmt)
		attack  = p.Add(name+"attack", defaultAttack)
		decay   = p.Add(name+"decay", defaultDecay)
		sustain = p.Add(name+"sustain", defaultSustain)
		release = p.Add(name+"release", defaultRelease)
	)

	return Operator{
		Gate: gate,
		Freq: freq,
		Gain: gain,
		FM:   fm,
		Amt:  amt,
		A:    attack,
		D:    decay,
		S:    sustain,
		R:    release,
		Done: sc.FreeEnclosing,
	}.Rate(sc.AR)
}
