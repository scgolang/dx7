package main

import (
	"os"
	"path"

	"github.com/scgolang/dx7/sysex"
)

const syxExt = ".syx"

// LoadPresets reads all the sysex files in a directory
// and returns a list of Sysex structs.
func (dx7 *DX7) LoadPreset(name string) error {
	logger.Printf("loading preset %s\n", name)

	// Read the sysex and load the appropriate synthdef.
	f, err := os.Open(path.Join(dx7.cfg.assetsDir, "syx", name+syxExt))
	if err != nil {
		return err
	}
	syx, err := sysex.New(f)
	if err != nil {
		return err
	}
	dx7.currentPreset = syx
	dx7.Poly.Def = getDefName(syx.Data.Algorithm)
	logger.Printf("set current synthef to %s\n", dx7.Poly.Def)
	return nil
}
