package metadatei

type readerState int

const (
	readerStateEmpty = readerState(iota)
	readerStateFirstBoundary
	readerStateSecondBoundary
)
