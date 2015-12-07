package main

import "fmt"

const (
	bulkDumpLength            = 4096
	numOps                    = 6
	opDataLength              = 17
	kbdLevelScalingDataLength = 4
)

// Common errors.
var (
	ErrInvalidLength = fmt.Errorf("bulk dump must be %d bytes in length", bulkDumpLength)
)

// BulkDump contains 155 parameters for a DX7 voice.
type BulkDump []*Op

// NewBulkDump creates a new BulkDump from a byte slice.
func NewBulkDump(data []byte) (BulkDump, error) {
	if len(data) != bulkDumpLength {
		return nil, ErrInvalidLength
	}
	ops := make([]*Op, numOps)
	for i, j := numOps, 0; i > 0; i-- {
		ops[i-1] = NewOp(data[j*opDataLength : (j*opDataLength)+opDataLength])
		j++
	}
	return BulkDump(ops), nil
}

// Op contains all the parameters for a single operator.
type Op struct {
	Amp                    EG              `json:"amp"`
	KbdLevelScaling        KbdLevelScaling `json:"kbd_level_scaling"`
	KbdRateScaling         int8            `json:"kbd_rate_scaling"`
	OutputLevel            int8            `json:"output_level"`
	KbdVelocitySensitivity int8            `json:"kbd_velocity_sensitivity"`
	AmpModSensitivity      int8            `json:"amp_mod_sensitivity"`
	Oscillator             Oscillator      `json:"oscillator"`
}

// NewOp creates a new Op from a byte slice.
func NewOp(data []byte) *Op {
	return &Op{
		Amp:                    NewEG(data[0:8]),
		KbdLevelScaling:        NewKbdLevelScaling(data[8:12]),
		KbdRateScaling:         int8(data[12] & 0x07),
		KbdVelocitySensitivity: getKbdVelSens(data[13]),
		AmpModSensitivity:      int8(data[13] & 0x03),
		OutputLevel:            int8(data[14]),
		Oscillator: Oscillator{
			Mode:       int8(data[15] & 0x01),
			FreqCoarse: getFreqCoarse(data[15]),
			FreqFine:   int8(data[16]),
			Detune:     getOscDetune(data[12]),
		},
	}
}

// Oscillator contains the oscillator parameters of a DX7 voice.
type Oscillator struct {
	Mode       int8 `json:"mode"`
	FreqCoarse int8 `json:"freq_coarse"`
	FreqFine   int8 `json:"freq_fine"`
	Detune     int8 `json:"detune"`
}

// KbdLevelScaling contains all the keyboard level scaling parameters
// for a DX7 voice.
type KbdLevelScaling struct {
	Breakpoint int8 `json:"breakpoint"`
	Ldepth     int8 `json:"ldepth"`
	Rdepth     int8 `json:"rdepth"`
	Lcurve     int8 `json:"lcurve"`
	Rcurve     int8 `json:"rcurve"`
}

// NewKbdLevelScaling creates a new KbdLevelScaling from a byte slice.
func NewKbdLevelScaling(data []byte) KbdLevelScaling {
	return KbdLevelScaling{
		Breakpoint: int8(data[0]),
		Ldepth:     int8(data[1]),
		Rdepth:     int8(data[2]),
		Lcurve:     int8(data[3] & 0x03),
		Rcurve:     getRcurve(data[3]),
	}
}

// getRcurve gets the R Curve for Keyboard Level Scaling
func getRcurve(b byte) int8 {
	return int8((0x0C & b) >> 2)
}

// getOscDetune gets Oscillator Detune from a byte.
func getOscDetune(b byte) int8 {
	return int8((0x78 & b) >> 3)
}

// getFreqCoarse gets the freq coarse parameter from a byte.
func getFreqCoarse(b byte) int8 {
	return int8((0x3E & b) >> 1)
}

// getKbdVelSens gets the keyboard velocity sensitivity parameter
// from a byte.
func getKbdVelSens(b byte) int8 {
	return int8((0x1C & b) >> 2)
}

// EG contains all the parameters of an envelope generator for
// a DX7 voice.
type EG struct {
	L1 int8 `json:"l1"`
	L2 int8 `json:"l2"`
	L3 int8 `json:"l3"`
	L4 int8 `json:"l4"`
	R1 int8 `json:"r1"`
	R2 int8 `json:"r2"`
	R3 int8 `json:"r3"`
	R4 int8 `json:"r4"`
}

// NewEG creates a new EG from a byte slice.
// A runtime panic occurs if data is less than 8 bytes long.
func NewEG(data []byte) EG {
	return EG{
		R1: int8(data[0]),
		R2: int8(data[1]),
		R3: int8(data[2]),
		R4: int8(data[3]),
		L1: int8(data[4]),
		L2: int8(data[5]),
		L3: int8(data[6]),
		L4: int8(data[7]),
	}
}
