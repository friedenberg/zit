// Code generated by "stringer -type=RemoteType"; DO NOT EDIT.

package repo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RemoteTypeUnspecified-0]
	_ = x[RemoteTypeNativeDotenvXDG-1]
	_ = x[RemoteTypeSocketUnix-2]
	_ = x[RemoteTypePortHTTP-3]
	_ = x[RemoteTypeStdioLocal-4]
	_ = x[RemoteTypeStdioSSH-5]
}

const _RemoteType_name = "RemoteTypeUnspecifiedRemoteTypeNativeDotenvXDGRemoteTypeSocketUnixRemoteTypePortHTTPRemoteTypeStdioLocalRemoteTypeStdioSSH"

var _RemoteType_index = [...]uint8{0, 21, 46, 66, 84, 104, 122}

func (i RemoteType) String() string {
	if i < 0 || i >= RemoteType(len(_RemoteType_index)-1) {
		return "RemoteType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RemoteType_name[_RemoteType_index[i]:_RemoteType_index[i+1]]
}