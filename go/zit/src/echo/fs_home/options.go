package fs_home

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

type Options struct {
	BasePath             string
	Debug                debug.Options
	DryRun               bool
	PermitNoZitDirectory bool
	store_fs             string
}

func (o *Options) Validate() (err error) {
	if o.store_fs, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type ReadOptions struct {
	*age.Age
	CompressionType immutable_config.CompressionType

	io.Reader
}

type FileReadOptions struct {
	*age.Age
	CompressionType immutable_config.CompressionType
	Path            string
}

type WriteOptions struct {
	*age.Age
	CompressionType immutable_config.CompressionType

	io.Writer
}

type MoveOptions struct {
	*age.Age
	CompressionType immutable_config.CompressionType

	TempDir                   string
	ErrorOnAttemptedOverwrite bool
	LockFile                  bool
	FinalPath                 string
	GenerateFinalPathFromSha  bool
}
