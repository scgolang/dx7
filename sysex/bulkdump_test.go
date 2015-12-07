package sysex

import "testing"

func TestGetRcurve(t *testing.T) {
	for _, tc := range []struct {
		Input  byte
		Output int8
	}{
		{0x03, int8(0)},
		{0x04, int8(1)},
		{0x05, int8(1)},
		{0x06, int8(1)},
		{0x07, int8(1)},
		{0x08, int8(2)},
		{0x09, int8(2)},
		{0x0A, int8(2)},
		{0x0B, int8(2)},
		{0x0C, int8(3)},
		{0x0D, int8(3)},
		{0x0E, int8(3)},
		{0x0F, int8(3)},
		{0x10, int8(0)},
		{0x11, int8(0)},
		{0x12, int8(0)},
		{0x13, int8(0)},
	} {
		if expected, got := tc.Output, getRcurve(tc.Input); expected != got {
			t.Fatalf("Expected %d, got %d", expected, got)
		}
	}
}

func TestGetOscDetune(t *testing.T) {
	for _, tc := range []struct {
		Input  byte
		Output int8
	}{
		{0x07, int8(0)},
		{0x38, int8(7)},
		{0x70, int8(14)},
	} {
		if expected, got := tc.Output, getOscDetune(tc.Input); expected != got {
			t.Fatalf("Expected %d, got %d", expected, got)
		}
	}
}
