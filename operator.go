package main

import (
	. "github.com/scgolang/sc/types"
	. "github.com/scgolang/sc/ugens"
)

// Operator is a sine wave signal combined with an envelope generator.
type Operator struct {
	osc    SinOsc
	env    EnvGen
	Levels [4]Input
	Rates  [4]Input
}

// defaults
func (self *Operator) defaults() {
	if self.osc.Freq == nil {
		self.osc.Freq = C(440)
	}
	if self.osc.Phase == nil {
		self.osc.Phase = C(0)
	}
}

// Rate creates a new ugen at a specific rate.
// If rate is an unsupported value this method will cause a runtime panic.
func (self Operator) Rate(rate int8) Input {
	CheckRate(rate)
	return self.osc.Rate(AR).Mul(self.env.Rate(AR))
}

// NewOperator creates a new operator
func NewOperator(freq Input) Operator {
	return Operator{
		osc: SinOsc{Freq: freq},
		env: EnvGen{Env: EnvPairs{}},
	}
}
