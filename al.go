package main

import "github.com/scgolang/sc"

// defaultAlgorithm
//    Op2 ---> Op1
var defaultAlgorithm = sc.NewSynthdef(defaultDefName, func(p sc.Params) sc.Ugen {
	gate, bus := p.Add("gate", 1), sc.C(0)

	// modulator
	op2 := NewOperator(2, p, gate, nil)

	// carrier
	op1 := NewOperator(1, p, gate, op2)

	// output signal
	sig := sc.Multi(op1, op1)

	return sc.Out{bus, sig}.Rate(sc.AR)
})
