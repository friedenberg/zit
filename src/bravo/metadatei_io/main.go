package metadatei_io

import (
	"io"
)

const (
	Boundary = "---"
)

type MetadateiWriterTo interface {
	io.WriterTo
	HasMetadateiContent() bool
}
