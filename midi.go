package main

import (
	"github.com/rakyll/portmidi"
)

// MIDIHandler
type MIDIHandler interface {
	Play(Note) error
	Control(CC) error
}

const (
	NoteEvent = iota
	ControlEvent
	NoteStatus = int64(144)
	CCStatus   = int64(176)
)

// Note is a MIDI note event.
type Note struct {
	Note, Velocity int
}

// NewNote creates a new MIDI note event.
// It panics if the event status does not indicate a MIDI note.
func NewNote(e portmidi.Event) Note {
	if e.Status != NoteStatus {
		panic("e is not a note event")
	}
	return Note{int(e.Data1), int(e.Data2)}
}

// CC is a MIDI control change event.
type CC struct {
	Num, Value int
}

// NewCC create a new MIDI control change event.
// It panics if the event status does not indicate a control change.
func NewCC(e portmidi.Event) CC {
	if e.Status != CCStatus {
		panic("e is not a control change event")
	}
	return CC{int(e.Data1), int(e.Data2)}
}
