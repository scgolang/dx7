package main

import "fmt"

// EventHandler
type EventHandler interface {
	Play(*Note) error
	Control(*Ctrl) error
}

// Note is a MIDI note event.
type Note struct {
	Note, Velocity int
}

func (n Note) String() string {
	return fmt.Sprintf("note=%d velocity=%d", n.Note, n.Velocity)
}

// Ctrl is a control event.
type Ctrl struct {
	Num, Value int
}

func (c Ctrl) String() string {
	return fmt.Sprintf("num=%d value=%d", c.Num, c.Value)
}
