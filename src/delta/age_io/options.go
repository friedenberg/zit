package age_io

import (
	"io"

	"github.com/friedenberg/zit/src/charlie/age"
)

type ReadOptions struct {
	age.Age
	UseZip bool

	io.Reader
}

type FileReadOptions struct {
	age.Age
	UseZip bool
	Path   string
}

type WriteOptions struct {
	age.Age
	UseZip bool

	io.Writer
}

type MoveOptions struct {
	age.Age
	UseZip bool

	TempDir                   string
	ErrorOnAttemptedOverwrite bool
	LockFile                  bool
	FinalPath                 string
	GenerateFinalPathFromSha  bool
}
