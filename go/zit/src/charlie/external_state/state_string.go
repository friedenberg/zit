// Code generated by "stringer -type=State"; DO NOT EDIT.

package external_state

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Tracked-1]
	_ = x[Untracked-2]
	_ = x[Recognized-3]
	_ = x[Conflicted-4]
	_ = x[Error-5]
}

const _State_name = "UnknownTrackedUntrackedRecognizedConflictedError"

var _State_index = [...]uint8{0, 7, 14, 23, 33, 43, 48}

func (i State) String() string {
	if i < 0 || i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}
