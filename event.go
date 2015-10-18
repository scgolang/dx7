package main

// EventHandler
type EventHandler interface {
	Play(Note) error
	Control(Ctrl) error
}

// Note is a MIDI note event.
type Note struct {
	Note, Velocity int
}

// Ctrl is a control event.
type Ctrl struct {
	Num, Value int
}
