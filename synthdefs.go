package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

var synthdefs = map[string]sc.UgenFunc{
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
	"dx7_algo23": func(p sc.Params) sc.Ugen {
		gate, bus := p.Add("gate", 1), sc.C(0)
		op6 := NewOperator(6, p, gate, nil)
		op5 := NewOperator(5, p, gate, op6)
		op4 := NewOperator(4, p, gate, op6)
		op3 := NewOperator(3, p, gate, nil)
		op2 := NewOperator(2, p, gate, op3)
		op1 := NewOperator(1, p, gate, nil)
		sig := sc.Mix(sc.AR, []sc.Input{op1, op2, op4, op5})
		sig = sc.Multi(sig, sig)
		return sc.Out{Bus: bus, Channels: sig}.Rate(sc.AR)
	},
}

// We don't have feedback yet, so alias algorithms with
// feedback to their no-feedback versions.
var synthdefAliases = map[string]string{
	"dx7_algo2": "dx7_algo1",
	"dx7_algo4": "dx7_algo3",
	"dx7_algo6": "dx7_algo5",
}

// getDefName gets a synthdef name from an algorithm number.
func getDefName(algo int8) string {
	def := fmt.Sprintf("dx7_algo%d", algo)
	if _, ok := synthdefs[def]; ok {
		return def
	}
	if alias, ok := synthdefAliases[def]; ok {
		return alias
	}
	panic(fmt.Sprintf("No synthdef for algorithm %d", algo))
}

// SendSynthdefs sends all the synthdefs needed for the DX7.
func (dx7 *DX7) SendSynthdefs() error {
	logger.Println("sending synthdefs")
	for def, f := range synthdefs {
		if err := dx7.Client.SendDef(sc.NewSynthdef(def, f)); err != nil {
			return err
		}
		logger.Printf("sent synthdef %s\n", def)
	}
	return nil
}
