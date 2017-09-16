package sc

import (
	"encoding/binary"
	"io"
)

// Variant provides a way to create synthdef presets.
type Variant struct {
	Name               string    `json:"name,omitempty"`
	InitialParamValues []float32 `json:"initialParamValues"`
}

// Write writes a variant to an io.Writer.
func (variant *Variant) Write(w io.Writer) error {
	if err := newPstring(variant.Name).Write(w); err != nil {
		return err
	}
	for _, v := range variant.InitialParamValues {
		if err := binary.Write(w, byteOrder, v); err != nil {
			return err
		}
	}
	return nil
}

// readVariant read a Variant from an io.Reader
func readVariant(r io.Reader, numParams int32) (*Variant, error) {
	name, err := readPstring(r)
	if err != nil {
		return nil, err
	}
	paramValues := make([]float32, numParams)
	for i := 0; int32(i) < numParams; i++ {
		if err := binary.Read(r, byteOrder, &paramValues[i]); err != nil {
			return nil, err
		}
	}
	v := Variant{name.String(), paramValues}
	return &v, nil
}
