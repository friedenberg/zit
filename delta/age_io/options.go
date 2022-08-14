package age_io

import (
	"io"

	"github.com/friedenberg/zit/charlie/age"
)

type ReadOptions struct {
	age.Age
	UseZip bool

	io.Reader
}

type WriteOptions struct {
	age.Age
	UseZip bool

	io.Writer
	LockFile bool
}
