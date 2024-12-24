package dir_layout

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
)

type Options struct {
	BasePath             string
	PermitNoZitDirectory bool
	MakeXDGDirectories   bool
}

func (o Options) GetReadOnlyBlobStorePath() string {
	return os.Getenv("ZIT_READ_ONLY_BLOB_STORE_PATH")
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

	TemporaryFS
	ErrorOnAttemptedOverwrite bool
	LockFile                  bool
	FinalPath                 string
	GenerateFinalPathFromSha  bool
}
