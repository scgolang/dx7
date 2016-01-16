package main

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/scgolang/portmidi"
)

const (
	// NoteStatus is the status value of a MIDI Note message.
	NoteStatus = int64(144)

	// CCStatus is the status value of a MIDI CC message.
	CCStatus = int64(176)

	// maxMIDI is the max value of a MIDI CC or Note.
	maxMIDI = float32(127)

	// midiBufferSize hardcoded buffer size for MIDI data.
	midiBufferSize = 1024

	// readMax is the maximum number of MIDI messages to read
	// at a single time.
	readMax = 1024
)

// Common errors.
var (
	ErrNotNote = errors.New("portmidi event is not a note event")
	ErrNotCtrl = errors.New("e is not a control change event")
)

// InitializeMIDI initializes the MIDI system.
func InitializeMIDI() error {
	portmidi.Initialize()
	return nil
}

// TerminateMIDI terminates the MIDI system.
func TerminateMIDI() error {
	portmidi.Terminate()
	return nil
}

// MIDINote creates a new MIDI note event.
// It panics if the event status does not indicate a MIDI note.
func MIDINote(e portmidi.Event) (*Note, error) {
	if e.Status != NoteStatus {
		return nil, ErrNotNote
	}
	return &Note{int(e.Data1), int(e.Data2)}, nil
}

// MIDICtrl create a new MIDI control change event.
// It panics if the event status does not indicate a control change.
func MIDICtrl(e portmidi.Event) (*Ctrl, error) {
	if e.Status != CCStatus {
		return nil, ErrNotCtrl
	}
	return &Ctrl{int(e.Data1), int(e.Data2)}, nil
}

// PrintMidiDevices prints a list of portmidi devices on stdout.
func PrintMidiDevices(w io.Writer) {
	midiDevices := portmidi.CountDevices()

	fmt.Fprintln(w, "| ID | Interface |         Name         | Input | Output |")
	fmt.Fprintln(w, "|----|-----------|----------------------|-------|--------|")

	row := "| %-2d | %-9s | %-20s | %-5t | %-6t |\n"
	for i := 0; i < midiDevices; i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceId(i))
		fmt.Printf(row, i, info.Interface, info.Name, info.IsInputAvailable, info.IsOutputAvailable)
	}
}

// MidiListen listens for MIDI events and handles them with
// a MIDIHandler. This function blocks.
func MidiListen(midiDeviceID int, handler EventHandler) chan error {
	var (
		errch = make(chan error)
		did   = portmidi.DeviceId(midiDeviceID)
	)
	go func() {
		midiInput, err := portmidi.NewInputStream(did, midiBufferSize)
		if midiInput != nil {
			defer midiInput.Close()
		}
		if err != nil {
			errch <- err
			return
		}

		log.Printf("listening for midi events on device %d\n", did)

		for events, err := midiInput.Read(1024); err == nil; events, err = midiInput.Read(1024) {
			for _, event := range events {
				if err := HandleMIDIEvent(event, handler); err != nil {
					errch <- err
					return
				}
			}
		}
	}()

	return errch
}

// HandleMIDIEvent handles a MIDI event with the provided handler.
func HandleMIDIEvent(event portmidi.Event, handler EventHandler) error {
	switch event.Status {
	case NoteStatus:
		note, err := MIDINote(event)
		if err != nil {
			return err
		}
		if err := handler.Play(note); err != nil {
			return err
		}
	case CCStatus:
		ctrl, err := MIDICtrl(event)
		if err != nil {
			return err
		}
		if err := handler.Control(ctrl); err != nil {
			return err
		}
	}
	return nil
}
