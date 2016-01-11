package main

import "github.com/scgolang/sc"

var algorithms = map[string]sc.UgenFunc{
	"dx7_algo1": func(p sc.Params) sc.Ugen {
		gate, bus := p.Add("gate", 1), sc.C(0)
		op6 := NewOperator(6, p, gate, nil)
		op5 := NewOperator(5, p, gate, op6)
		op4 := NewOperator(4, p, gate, op5)
		op3 := NewOperator(3, p, gate, op4)
		op2 := NewOperator(2, p, gate, nil)
		op1 := NewOperator(1, p, gate, op2)
		sig := op1.Add(op3)
		sig = sc.Multi(sig, sig)
		return sc.Out{Bus: bus, Channels: sig}.Rate(sc.AR)
	},
	"dx7_algo3": func(p sc.Params) sc.Ugen {
		gate, bus := p.Add("gate", 1), sc.C(0)
		op6 := NewOperator(6, p, gate, nil)
		op5 := NewOperator(5, p, gate, op6)
		op4 := NewOperator(4, p, gate, op5)
		op3 := NewOperator(3, p, gate, nil)
		op2 := NewOperator(2, p, gate, op3)
		op1 := NewOperator(1, p, gate, op2)
		sig := op1.Add(op4)
		sig = sc.Multi(sig, sig)
		return sc.Out{Bus: bus, Channels: sig}.Rate(sc.AR)
	},
	"dx7_algo5": func(p sc.Params) sc.Ugen {
		gate, bus := p.Add("gate", 1), sc.C(0)
		op6 := NewOperator(6, p, gate, nil)
		op5 := NewOperator(5, p, gate, op6)
		op4 := NewOperator(4, p, gate, nil)
		op3 := NewOperator(3, p, gate, op4)
		op2 := NewOperator(2, p, gate, nil)
		op1 := NewOperator(1, p, gate, op2)
		sig := sc.Mix(sc.AR, []sc.Input{op1, op3, op5})
		sig = sc.Multi(sig, sig)
		return sc.Out{Bus: bus, Channels: sig}.Rate(sc.AR)
	},
}

func init() {
	// We don't have feedback yet, so alias algorithms with
	// feedback to their no-feedback versions.
	algorithms["dx7_algo2"] = algorithms["dx7_algo1"]
	algorithms["dx7_algo4"] = algorithms["dx7_algo3"]
	algorithms["dx7_algo6"] = algorithms["dx7_algo5"]
}
