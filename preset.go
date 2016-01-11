package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/scgolang/dx7/sysex"
)

// LoadPresets reads all the sysex files in a directory
// and returns a list of Sysex structs.
func (dx7 *DX7) LoadPresets(cfg *config) error {
	loadPreset := func(path string, info os.FileInfo, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		if _, err := sysex.New(f); err != nil {
			return err
		}
		log.Printf("path %s\n", path)
		return nil
	}
	return filepath.Walk(path.Join(cfg.assetsDir, "syx"), loadPreset)
}

// sendSynthdefs transforms a sysex preset into a synthdef.
// operators are wired up in the returned synthdef according
// to one of the 32 DX7 "algorithms".
// for a depiction of the dx7 algorithms, see
// http://www.polynominal.com/site/studio/gear/synth/yamaha-tx802/tx802-board2.jpg
func (dx7 *DX7) sendSynthdefs() error {
	return nil
}
