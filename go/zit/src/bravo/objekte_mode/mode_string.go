// Code generated by "stringer -type=Mode"; DO NOT EDIT.

package objekte_mode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ModeEmpty-0]
	_ = x[ModeAddToBestandsaufnahme-2]
	_ = x[ModeUpdateTai-4]
	_ = x[ModeSchwanz-8]
}

const (
	_Mode_name_0 = "ModeEmpty"
	_Mode_name_1 = "ModeAddToBestandsaufnahme"
	_Mode_name_2 = "ModeUpdateTai"
	_Mode_name_3 = "ModeSchwanz"
)

func (i Mode) String() string {
	switch {
	case i == 0:
		return _Mode_name_0
	case i == 2:
		return _Mode_name_1
	case i == 4:
		return _Mode_name_2
	case i == 8:
		return _Mode_name_3
	default:
		return "Mode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
