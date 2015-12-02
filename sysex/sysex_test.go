package main

import "testing"

func TestSubStatus(t *testing.T) {
	for _, pair := range []struct {
		In  byte
		Out int
	}{
		{0x02, 0},
		{0x08, 0},
		{0x10, 1},
		{0x1A, 1},
		{0x20, 2},
		{0x2B, 2},
		{0x30, 3},
		{0x3C, 3},
		{0x40, 4},
		{0x4C, 4},
		{0x50, 5},
		{0x5C, 5},
		{0x60, 6},
		{0x6C, 6},
		{0x70, 7},
		{0x7C, 7},
		{0x8D, 0},
		{0xA0, 2},
		{0xB1, 3},
		{0xC2, 4},
		{0xD3, 5},
		{0xE4, 6},
		{0xF5, 7},
	} {
		if expected, got := pair.Out, getSubstatus(pair.In); expected != got {
			t.Fatalf("Expected %d got %d", expected, got)
		}
	}
}

func TestChannel(t *testing.T) {
	for _, pair := range []struct {
		In  byte
		Out int
	}{
		{0x02, 2},
		{0x08, 8},
		{0x10, 0},
		{0x1A, 10},
		{0x20, 0},
		{0x2B, 11},
		{0x30, 0},
		{0x3C, 12},
		{0x40, 0},
		{0x4C, 12},
		{0x50, 0},
		{0x5C, 12},
		{0x60, 0},
		{0x6C, 12},
		{0x70, 0},
		{0x7C, 12},
		{0x8D, 13},
		{0xA0, 0},
		{0xB1, 1},
		{0xC2, 2},
		{0xD3, 3},
		{0xE4, 4},
		{0xF5, 5},
	} {
		if expected, got := pair.Out, getChannel(pair.In); expected != got {
			t.Fatalf("Expected %d got %d", expected, got)
		}
	}
}

func TestByteCount(t *testing.T) {
	for _, val := range []struct {
		MS  byte
		LS  byte
		Out int16
	}{
		{0x01, 0x00, int16(256)},
		// TODO: more test cases
	} {
		if expected, got := val.Out, getByteCount(val.MS, val.LS); expected != got {
			t.Fatalf("Expected %d got %d", expected, got)
		}
	}
}
