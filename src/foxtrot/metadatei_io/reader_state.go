package metadatei_io

type readerState int

const (
	readerStateEmpty = readerState(iota)
	readerStateFirstBoundary
	readerStateSecondBoundary
)
