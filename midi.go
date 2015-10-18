package main

import (
	"fmt"
	"io"

	"github.com/rakyll/portmidi"
)

const (
	// NoteStatus is the status value of a MIDI Note message.
	NoteStatus = int64(144)

	// CCStatus is the status value of a MIDI CC message.
	CCStatus = int64(176)
)

// MIDINote creates a new MIDI note event.
// It panics if the event status does not indicate a MIDI note.
func MIDINote(e portmidi.Event) Note {
	if e.Status != NoteStatus {
		panic("e is not a note event")
	}
	return Note{int(e.Data1), int(e.Data2)}
}

// MIDICC create a new MIDI control change event.
// It panics if the event status does not indicate a control change.
func MIDICC(e portmidi.Event) Ctrl {
	if e.Status != CCStatus {
		panic("e is not a control change event")
	}
	return Ctrl{int(e.Data1), int(e.Data2)}
}

// PrintMidiDevices prints a list of portmidi devices on stdout.
func PrintMidiDevices(w io.Writer) {
	midiDevices := portmidi.CountDevices()
	// fmt.Printf("%d midi devices:\n", midiDevices)

	fmt.Fprintln(w, "| ID | Interface |         Name         | Input | Output |")
	fmt.Fprintln(w, "|----|-----------|----------------------|-------|--------|")

	row := "| %-2d | %-9s | %-20s | %-5t | %-6t |\n"
	for i := 0; i < midiDevices; i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceId(i))
		fmt.Printf(row, i, info.Interface, info.Name, info.IsInputAvailable, info.IsOutputAvailable)
	}
}

// MidiListen listens for MIDI events and handles them with
// a MIDIHandler.
func MidiListen(midiDeviceID portmidi.DeviceId, handler EventHandler) error {
	midiInput, err := portmidi.NewInputStream(midiDeviceID, midiBufferSize)
	if midiInput != nil {
		defer midiInput.Close()
	}
	if err != nil {
		return err
	}

	for event := range midiInput.Listen() {
		switch event.Status {
		case NoteStatus:
			if err := handler.Play(MIDINote(event)); err != nil {
				return err
			}
		}
	}

	return nil
}
