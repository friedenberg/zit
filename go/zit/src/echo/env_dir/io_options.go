package env_dir

import (
	"io"
)

type ReadOptions struct {
	Config
	io.Reader
}

type FileReadOptions struct {
	Config
	Path string
}

type WriteOptions struct {
	Config
	io.Writer
}

type MoveOptions struct {
	Config
	TemporaryFS
	ErrorOnAttemptedOverwrite bool
	FinalPath                 string
	GenerateFinalPathFromSha  bool
}
