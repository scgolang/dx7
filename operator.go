package main

import (
	. "github.com/scgolang/sc/types"
	. "github.com/scgolang/sc/ugens"
)

// Operator is a sine wave signal combined with an envelope generator.
type Operator struct {
	// Freq is the oscillator frequency.
	Freq Input
	// FM is the frequency modulation input.
	FM Input
	// Amt controls the frequency modulation amount.
	Amt Input
	// Gain is the output gain.
	Gain Input
	// A is amp envelope attack (in seconds)
	A Input
	// D is amp envelope decay (in seconds)
	D Input
	// S is amp envelope sustain [0, 1]
	S Input
	// R is amp envelope release (in seconds)
	R Input
	// Gate trigger the envelope and holds it open while > 0
	Gate Input
	// Done is the ugen done action
	Done int
}

// defaults
func (self *Operator) defaults() {
	if self.Freq == nil {
		self.Freq = C(440)
	}
	if self.FM == nil {
		self.FM = C(1)
	}
	if self.Amt == nil {
		self.Amt = C(0)
	}
	if self.Gain == nil {
		self.Gain = C(1)
	}
	if self.A == nil {
		self.A = C(0.01)
	}
	if self.D == nil {
		self.D = C(0.3)
	}
	if self.S == nil {
		self.S = C(0.5)
	}
	if self.R == nil {
		self.R = C(1)
	}
	if self.Gate == nil {
		self.Gate = C(1)
	}
}

// Rate creates a new ugen at a specific rate.
// If rate is an unsupported value this method will cause a runtime panic.
func (self Operator) Rate(rate int8) Input {
	CheckRate(rate)
	(&self).defaults()

	adsr := EnvADSR{A: self.A, D: self.D, S: self.S, R: self.R}
	env := EnvGen{
		Env:        adsr,
		Gate:       self.Gate,
		LevelScale: self.Gain,
		Done:       self.Done,
	}.Rate(AR)

	freq := self.Freq.Add(self.FM.Mul(self.Amt))

	return SinOsc{Freq: freq}.Rate(AR).Mul(env)
}
