// Code generated by "stringer -type=boundaryReaderState"; DO NOT EDIT.

package ohio

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[boundaryReaderStateEmpty-0]
	_ = x[boundaryReaderStateNeedsBoundary-1]
	_ = x[boundaryReaderStateOnlyContent-2]
	_ = x[boundaryReaderStatePartialBoundaryInBuffer-3]
	_ = x[boundaryReaderStateCompleteBoundaryInBuffer-4]
}

const _boundaryReaderState_name = "boundaryReaderStateEmptyboundaryReaderStateNeedsBoundaryboundaryReaderStateOnlyContentboundaryReaderStatePartialBoundaryInBufferboundaryReaderStateCompleteBoundaryInBuffer"

var _boundaryReaderState_index = [...]uint8{0, 24, 56, 86, 128, 171}

func (i boundaryReaderState) String() string {
	if i < 0 || i >= boundaryReaderState(len(_boundaryReaderState_index)-1) {
		return "boundaryReaderState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _boundaryReaderState_name[_boundaryReaderState_index[i]:_boundaryReaderState_index[i+1]]
}