package standort

import (
	"io"

	"code.linenisgreat.com/zit/src/delta/age"
	"code.linenisgreat.com/zit/src/delta/angeboren"
)

type ReadOptions struct {
	*age.Age
	CompressionType angeboren.CompressionType

	io.Reader
}

type FileReadOptions struct {
	*age.Age
	CompressionType angeboren.CompressionType
	Path            string
}

type WriteOptions struct {
	*age.Age
	CompressionType angeboren.CompressionType

	io.Writer
}

type MoveOptions struct {
	*age.Age
	CompressionType angeboren.CompressionType

	TempDir                   string
	ErrorOnAttemptedOverwrite bool
	LockFile                  bool
	FinalPath                 string
	GenerateFinalPathFromSha  bool
}
