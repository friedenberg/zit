package workspace_config_blobs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
)

const (
	TypeV0 = builtin_types.WorkspaceConfigTypeTomlV0
)

type (
	Blob interface {
		config_mutable_blobs.Blob
		GetDefaultQueryGroup() string
	}
)

type TypeWithBlob = *ids.TypeWithObject[*Blob]

var typedCoders = map[string]interfaces.Coder[TypeWithBlob]{
	TypeV0: blobV0Coder{},
}

var Coder = triple_hyphen_io.Coder[TypeWithBlob]{
	Metadata: ids.TypedMetadataCoder[*Blob]{},
	Blob:     ids.TypedCoders[*Blob](typedCoders),
}

func DecodeFromFile(
	object TypeWithBlob,
	path string,
) (err error) {
	var reader io.Reader

	{
		var file *os.File

		if file, err = files.OpenExclusiveReadOnly(path); err != nil {
			if !errors.IsNotExist(err) {
				err = errors.Wrap(err)
			}

			return
		} else {
			defer errors.DeferredCloser(&err, file)

			reader = file
		}
	}

	if _, err = Coder.DecodeFrom(object, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
