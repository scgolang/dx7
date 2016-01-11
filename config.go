package main

import (
	"flag"
	"os"
	"path"

	"github.com/rakyll/portmidi"
)

const (
	defaultMidiDeviceId = 0
	defaultLocalAddr    = "127.0.0.1:57110"
	defaultScsynthAddr  = "127.0.0.1:57120"
)

var (
	srcPath = path.Join(os.Getenv("GOPATH"), "src", "github.com", "scgolang", "dx7")
)

// config encapsulates info parsed from the CLI
type config struct {
	midiDeviceID    portmidi.DeviceId
	localAddr       string
	scsynthAddr     string
	eventsAddr      string
	assetsDir       string
	preset          string
	dumpSysex       bool
	listMidiDevices bool
	dumpOSC         bool
	algorithm       int
}

// getConfig gets a config from the CLI.
func getConfig() *config {
	var (
		cfg          = config{}
		fs           = flag.NewFlagSet("", flag.ExitOnError)
		midiDeviceId int
	)
	fs.IntVar(&midiDeviceId, "midi", defaultMidiDeviceId, "MIDI Device ID")
	fs.StringVar(&cfg.localAddr, "local", defaultLocalAddr, "local OSC address")
	fs.StringVar(&cfg.scsynthAddr, "scsynth", defaultScsynthAddr, "scsynth OSC address")
	fs.StringVar(&cfg.eventsAddr, "events", "", "events OSC address")
	fs.BoolVar(&cfg.listMidiDevices, "listmidi", false, "list MIDI devices")
	fs.BoolVar(&cfg.dumpOSC, "dumposc", false, "have scsynth dump OSC messages on stdout")
	fs.StringVar(&cfg.assetsDir, "assets-dir", path.Join(srcPath, "assets"), "path to assets directory")
	fs.StringVar(&cfg.preset, "preset", "organ1", "initial preset")
	fs.BoolVar(&cfg.dumpSysex, "dump-sysex", false, "print JSON-encoded presets to stdout ")
	fs.IntVar(&cfg.algorithm, "algorithm", -1, "DX7 algorithm")
	fs.Parse(os.Args[1:])
	cfg.midiDeviceID = portmidi.DeviceId(midiDeviceId)
	return &cfg
}
