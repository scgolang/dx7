package main

import (
	"flag"
	"os"

	"github.com/rakyll/portmidi"
)

const (
	defaultMidiDeviceId = 0
	defaultLocalAddr    = "127.0.0.1:57110"
	defaultScsynthAddr  = "127.0.0.1:57120"
)

// config encapsulates info parsed from the CLI
type config struct {
	midiDeviceID    portmidi.DeviceId
	localAddr       string
	scsynthAddr     string
	eventsAddr      string
	listMidiDevices bool
	dumpOSC         bool
}

// parseConfig parses a config from the CLI.
func parseConfig() *config {
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
	fs.Parse(os.Args[1:])
	cfg.midiDeviceID = portmidi.DeviceId(midiDeviceId)
	return &cfg
}
