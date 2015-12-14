package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

// loadAlgo transforms a sysex preset into a synthdef.
// operators are wired up in the returned synthdef according
// to one of the 32 DX7 "algorithms".
// for a depiction of the dx7 algorithms, see
// http://www.polynominal.com/site/studio/gear/synth/yamaha-tx802/tx802-board2.jpg
func loadAlgo(algo int) *sc.Synthdef {
	switch algo {
	default:
		panic(fmt.Sprintf("algorithm not in range 1-32: %d", algo))
	case 1, 2:
		return sc.NewSynthdef("dx7_algorithm1", func(p sc.Params) sc.Ugen {
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
		})
	case 3, 4:
		return sc.NewSynthdef("dx7_algorithm2", func(p sc.Params) sc.Ugen {
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
		})
	case 5, 6:
		return sc.NewSynthdef("dx7_algorithm3", func(p sc.Params) sc.Ugen {
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
		})
	}
}
