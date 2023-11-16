package ohio_boundary_reader

//go:generate stringer -type=boundaryReaderState
type boundaryReaderState int

const (
	boundaryReaderStateEmpty = boundaryReaderState(iota)
	boundaryReaderStateNeedsBoundary
	boundaryReaderStateOnlyContent
	boundaryReaderStatePartialBoundaryInBuffer
	boundaryReaderStateCompleteBoundaryInBuffer
)
