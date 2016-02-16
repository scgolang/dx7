package main

import (
	"encoding/json"
	"os"
	"path"

	"github.com/scgolang/dx7/sysex"
	"github.com/scgolang/poly"
)

var (
	srcPath       = path.Join(os.Getenv("GOPATH"), "src", "github.com", "scgolang", "dx7")
	defaultPreset = "organ1"
)

// DX7 is a recreation of the legendary Yamaha DX7.
type DX7 struct {
	*poly.Poly

	// config
	assetsDir string
	preset    string
	dumpSysex bool
	algorithm int

	currentPreset *sysex.Sysex

	ctrls map[string]float32
}

// run the dx7.
func (dx7 *DX7) run() error {
	// Load the current preset.
	if err := dx7.LoadPreset(dx7.preset); err != nil {
		return err
	}

	// Dump sysex data to stdout for the current preset and return.
	if dx7.dumpSysex {
		return json.NewEncoder(os.Stdout).Encode(dx7.currentPreset)
	}

	// Connect to scsynth.
	if err := dx7.Connect(); err != nil {
		return err
	}

	// Send all the synthdefs we need.
	if err := dx7.SendSynthdefs(); err != nil {
		return err
	}

	// Listen for events.
	return poly.Listen(dx7.Poly)
}

// New returns a DX7 using the defaultAlgorithm.
// client will be used to create synth nodes, and all the synth
// nodes will be added to the provided group.
func New() (*DX7, error) {
	dx7 := &DX7{
		ctrls: map[string]float32{
			"op1amt":       float32(defaultAmt),
			"op2freqscale": float32(1),
			"op2decay":     float32(defaultDecay),
			"op2sustain":   float32(defaultSustain),
		},
	}
	p, err := poly.New(dx7)
	if err != nil {
		return nil, err
	}

	p.FlagSet.StringVar(&dx7.assetsDir, "assets", path.Join(srcPath, "assets"), "path to assets directory")
	p.FlagSet.StringVar(&dx7.preset, "preset", defaultPreset, "initial preset")
	p.FlagSet.BoolVar(&dx7.dumpSysex, "dump-sysex", false, "print JSON-encoded presets to stdout ")
	p.FlagSet.IntVar(&dx7.algorithm, "algorithm", -1, "DX7 algorithm")

	p.FlagSet.Parse(os.Args[1:])

	dx7.Poly = p

	return dx7, nil
}
