package main

import "github.com/scgolang/sc"

// defaultAlgorithm
//    Op2 ---> Op1
var defaultAlgorithm = sc.NewSynthdef(defaultDefName, func(p sc.Params) sc.Ugen {
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
