package main

import (
	"os"
	"path"

	"github.com/scgolang/dx7/sysex"
)

const syxExtension = ".syx"

// LoadPresets reads all the sysex files in a directory
// and returns a list of Sysex structs.
func (dx7 *DX7) LoadPreset(name string) error {
	// Read the sysex and load the appropriate synthdef.
	f, err := os.Open(path.Join(dx7.cfg.assetsDir, "syx", name+syxExtension))
	if err != nil {
		return err
	}
	syx, err := sysex.New(f)
	if err != nil {
		return err
	}
	dx7.currentPreset = syx
	return nil
}

// sendSynthdefs transforms a sysex preset into a synthdef.
// operators are wired up in the returned synthdef according
// to one of the 32 DX7 "algorithms".
// for a depiction of the dx7 algorithms, see
// http://www.polynominal.com/site/studio/gear/synth/yamaha-tx802/tx802-board2.jpg
func (dx7 *DX7) sendSynthdefs() error {
	return nil
}
