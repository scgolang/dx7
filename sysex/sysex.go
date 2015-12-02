// package sysex parses MIDI sysex messages created by a DX7.
// This package would not be possible without the incredible collection
// of DX7 resources maintained by Dave Benson
// https://homepages.abdn.ac.uk/mth192/pages/html/dx7.html
// as well as the folks at blitter.com
// http://www.blitter.com/~russtopia/MIDI/~jglatt/tech/midispec/sysex.htm
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	ManufacturerID = 0x43

	headerLength = 6
)

// Sysex defines a MIDI sysex message.
type Sysex struct {
	Substatus    int   `json:"substatus"`
	Channel      int   `json:"channel"`
	FormatNumber int   `json:"format_number"`
	ByteCount    int16 `json:"byte_count"`
	data         []byte
}

// New parses a sysex message from an io.Reader.
func New(r io.Reader) (*Sysex, error) {
	hdr := make([]byte, headerLength)
	if _, err := r.Read(hdr); err != nil {
		return nil, err
	}

	if hdr[1] != ManufacturerID {
		return nil, fmt.Errorf("Manufacturer is not Yamaha: %X", hdr[1])
	}

	syx := &Sysex{
		Substatus:    getSubstatus(hdr[2]),
		Channel:      getChannel(hdr[2]),
		FormatNumber: midiMask(hdr[3]),
		ByteCount:    getByteCount(hdr[4], hdr[5]),
	}

	syx.data = make([]byte, syx.ByteCount)

	if _, err := r.Read(syx.data); err != nil {
		return nil, err
	}

	return syx, nil
}

// getSubstatus gets the substatus value from a byte.
// The structure of the byte is 0sssnnnn.
func getSubstatus(b byte) int {
	return int((b & 0x70) >> 4)
}

// getChannel gets the channel value from the 3rd byte of
// a DX7 sysex message.
func getChannel(b byte) int {
	return int(b & 0x0F)
}

// midiMask masks a byte since MIDI data bytes cannot have
// the high bit set.
func midiMask(b byte) int {
	return int(b & 0x7F)
}

// getByteCount gets the count of data bytes.
func getByteCount(ms, ls byte) int16 {
	return (int16(ms) << 8) + int16(ls)
}

func main() {
	dir, err := os.Open("syx")
	if err != nil {
		log.Fatal(err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		r, err := os.Open(path.Join(dir.Name(), file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		syx, err := New(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(r.Name() + " ")
		if err := json.NewEncoder(os.Stdout).Encode(syx); err != nil {
			log.Fatal(err)
		}
	}
}
