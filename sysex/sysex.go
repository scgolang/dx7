// package sysex parses MIDI sysex messages created by a DX7.
// This package would not be possible without the incredible collection
// of DX7 resources maintained by Dave Benson
// https://homepages.abdn.ac.uk/mth192/pages/html/dx7.html
// as well as the folks at blitter.com
// http://www.blitter.com/~russtopia/MIDI/~jglatt/tech/midispec/sysex.htm
package sysex

import (
	"fmt"
	"io"
)

const (
	yamahaManufacturerID = 0x43
	headerLength         = 6
)

// Sysex defines a MIDI sysex message.
type Sysex struct {
	Substatus    int       `json:"substatus"`
	Channel      int       `json:"channel"`
	FormatNumber int       `json:"format_number"`
	ByteCount    int16     `json:"byte_count"`
	Data         *BulkDump `json:"data"`
}

// New parses a sysex message from an io.Reader.
func New(r io.Reader) (*Sysex, error) {
	hdr := make([]byte, headerLength)
	if _, err := r.Read(hdr); err != nil {
		return nil, err
	}

	if hdr[1] != yamahaManufacturerID {
		return nil, fmt.Errorf("Manufacturer is not Yamaha: %X", hdr[1])
	}

	var (
		syx = &Sysex{
			Substatus:    getSubstatus(hdr[2]),
			Channel:      getChannel(hdr[2]),
			FormatNumber: midiMask(hdr[3]),
			ByteCount:    getByteCount(hdr[4], hdr[5]),
		}
		data = make([]byte, syx.ByteCount)
	)

	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	if int16(n) != syx.ByteCount {
		return nil, fmt.Errorf("only read %d data bytes", n)
	}

	bulkDump, err := NewBulkDump(data)
	if err != nil {
		return nil, err
	}
	syx.Data = bulkDump

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
	return (int16(ms) << 7) + int16(ls)
}
