package sysex

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
type BulkDump struct {
	// Ops is the list of operators that define the voice.
	Ops []*Op `json:"ops"`

	// PitchEG is the pitch envelope generator.
	PitchEG EG `json:"pitch_eg"`

	// Algorithm determines the modulation routing for
	// the 6 operators.
	Algorithm int8 `json:"algorithm"`

	// OscKeySync
	OscKeySync int8 `json:"osc_key_sync"`

	// Feedback adjusts the amount of feedback for
	// algorithms that use feedback.
	Feedback int8 `json:"feedback"`

	// LFO controls the the low frequency oscillator.
	LFO LFO `json:"lfo"`

	// Transpose transposes a particular voice up or down
	// on the keyboard.
	Transpose int8 `json:"transpose"`

	// Name is the voice's name.
	Name string `json:"name"`
}

// NewBulkDump creates a new BulkDump from a byte slice.
func NewBulkDump(data []byte) (*BulkDump, error) {
	if len(data) != bulkDumpLength {
		return nil, ErrInvalidLength
	}
	ops := make([]*Op, numOps)
	for i, j := numOps, 0; i > 0; i-- {
		ops[i-1] = NewOp(data[j*opDataLength : (j*opDataLength)+opDataLength])
		j++
	}

	offset := opDataLength * numOps

	return &BulkDump{
		Ops:        ops,
		PitchEG:    NewEG(data[offset : offset+8]),
		Algorithm:  int8(data[offset+8] & 0x1F),
		OscKeySync: int8(data[offset+9] & 0x08),
		Feedback:   int8(data[offset+9] & 0x07),
		LFO:        NewLFO(data[offset+10 : offset+15]),
		Transpose:  int8(data[offset+15]),
		Name:       string(data[offset+16 : offset+26]),
	}, nil
}

// Op contains all the parameters for a single operator.
type Op struct {
	// AmpEG is the amplitude envelope generator.
	AmpEG EG `json:"amp"`

	// KbdLevelScaling provides a way to scale the output
	// level depending on what key is pressed.
	KbdLevelScaling KbdLevelScaling `json:"kbd_level_scaling"`

	// KbdRateScaling provides a way to shorten the
	// rates of the amplitude envelope generators as higher
	// notes are played.
	KbdRateScaling int8 `json:"kbd_rate_scaling"`

	// OutputLevel controls the output level of the operator.
	OutputLevel int8 `json:"output_level"`

	// KdbVelocitySensitivity controls how sensitive an operator
	// is to velocity (how much affect velocity will have on
	// the output level of the operator).
	KbdVelocitySensitivity int8 `json:"kbd_velocity_sensitivity"`

	// AmpModSensitivity not really sure what this does (briansorahan).
	AmpModSensitivity int8 `json:"amp_mod_sensitivity"`

	// Oscillator corresponds to the oscillator section on
	// the DX7.
	Oscillator Oscillator `json:"oscillator"`
}

// NewOp creates a new Op from a byte slice.
// A runtime panic occurs if data is less than 17 bytes long.
func NewOp(data []byte) *Op {
	return &Op{
		AmpEG:                  NewEG(data[0:8]),
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
	// Mode controls whether or not the keyboard maps to
	// oscillator frequency (0=tracking, 1=fixed).
	Mode int8 `json:"mode"`

	// FreqCoarse is a coarse adjustment of frequency.
	// In tracking mode 0 corresponds to a frequency ratio of 0.5x
	// and 99 corresponds to 31x.
	// In fixed mode 0 corresponds to 1Hz and 99 to 1000Hz.
	FreqCoarse int8 `json:"freq_coarse"`

	// FreqFine is a fine adjustment of frequency.
	FreqFine int8 `json:"freq_fine"`

	// Detune offers a frequency adjustment which is finer than FreqFine.
	Detune int8 `json:"detune"`
}

// KbdLevelScaling contains all the keyboard level scaling parameters
// for a DX7 voice.
// Keyboard level scaling lets you change the output level of an operator
// based on the key that is pressed.
// Each operator can be programmed to have any of 4 curves on either side
// of an adjustable breakpoint.
// The scaling can be used to make the tone and/or volume change as
// you move to different octaves.
type KbdLevelScaling struct {
	// Breakpoint is the break point in the level-scaling curve.
	// 0 corresponds to 1 1/3 octaves below the lowest note on
	// the keyboard (A-1), and 99 corresponds to 2 octaves above
	// the highest note on the keybard.
	Breakpoint int8 `json:"breakpoint"`

	// Ldepth controls the curve depth on the left side of
	// the breakpoint.
	Ldepth int8 `json:"ldepth"`

	// Rdepth controls the curve depth on the right side of
	// the breakpoint.
	Rdepth int8 `json:"rdepth"`

	// Lcurve controls the curve on the left side of the
	// breakpoint. The curve can have 4 different shapes:
	// negative linear, negative exponential, positive
	// exponential, and positive linear.
	Lcurve int8 `json:"lcurve"`

	// Rcurve controls the curve on the right side of the
	// breakpoint. The curve can have 4 different shapes:
	// negative linear, negative exponential, positive
	// exponential, and positive linear.
	Rcurve int8 `json:"rcurve"`
}

// NewKbdLevelScaling creates a new KbdLevelScaling from a byte slice.
// A runtime panic occurs if data is less than 4 bytes long.
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

// LFO contains the parameters for a DX7 LFO.
type LFO struct {
	// Speed adjusts LFO speed (0 is slowest, 99 is fastest)
	Speed int8 `json:"speed"`

	// Delay between when a key is pressed and modulation kicks in
	Delay int8 `json:"delay"`

	// PMD is Pitch Modulation Depth
	PMD int8 `json:"pmd"`

	// AMD is Amplitude Modulation Depth
	AMD int8 `json:"amd"`

	// PMSensitivity adjusts the sensitivity of individual
	// voices to LFO modulation
	PMSensitivity int8 `json:"pm_sensitivity"`

	// Wave 0=triangle, 1=sawdown, 2=sawup, 3=square, 4=sine, 5=s+h
	Wave int8 `json:"wave"`

	// Sync causes the LFO to restart every time a key is pressed
	Sync int8 `json:"sync"`
}

// NewLFO creates a new LFO from a byte slice.
// A runtime panic occurs if data is less than 5 bytes long.
func NewLFO(data []byte) LFO {
	return LFO{
		Speed:         int8(data[0]),
		Delay:         int8(data[1]),
		PMD:           int8(data[2]),
		AMD:           int8(data[3]),
		PMSensitivity: getPMSensitivity(data[4]),
		Wave:          getWave(data[4]),
		Sync:          int8(data[4] & 0x01),
	}
}

// getPMSensitivity gets the PMSensitivity parameter of an LFO.
func getPMSensitivity(b byte) int8 {
	return int8((b & 0x70) >> 4)
}

// getWave gets the wave parameter of an LFO.
func getWave(b byte) int8 {
	return int8((0x0E & b) >> 1)
}
